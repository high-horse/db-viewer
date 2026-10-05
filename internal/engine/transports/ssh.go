package transports

import (
	"context"
	"db-lens/internal/engine/entities"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSH forwards a loopback-only listener to the database through an SSH server.
// The tunnel lives until Close, independently of the connection request context.
type SSH struct {
	lifecycle sync.Mutex
	dbHost    string
	dbPort    int
	config    entities.SSHConfig
	mu        sync.Mutex
	client    *ssh.Client
	listener  net.Listener
	cancel    context.CancelFunc
	sockets   map[net.Conn]struct{}
	wg        sync.WaitGroup
}

func NewSSH(dbHost string, dbPort int, sshHost string, sshPort int, sshUser string, sshKeyPath string, sshPassword string) *SSH {
	method := "password"
	if sshKeyPath != "" {
		method = "private_key"
	}
	return NewSSHConfig(dbHost, dbPort, entities.SSHConfig{Host: sshHost, Port: sshPort, Username: sshUser, AuthMethod: method, PrivateKey: sshKeyPath, Password: sshPassword})
}

func NewSSHConfig(host string, port int, config entities.SSHConfig) *SSH {
	return &SSH{dbHost: host, dbPort: port, config: config}
}

func (s *SSH) Connect(ctx context.Context) error {
	s.lifecycle.Lock()
	defer s.lifecycle.Unlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client != nil {
		return nil
	}
	if s.dbPort < 1 || s.dbPort > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	client, err := dialSSH(ctx, s.config)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		_ = client.Close()
		return fmt.Errorf("listen for SSH tunnel: %w", err)
	}
	tunnelCtx, tunnelCancel := context.WithCancel(context.Background())
	s.client, s.listener, s.cancel = client, listener, tunnelCancel
	s.sockets = make(map[net.Conn]struct{})
	s.wg.Add(1)
	go s.serve(tunnelCtx, listener, client)
	return nil
}

// TestSSH checks authentication and host identity without creating a database tunnel.
func TestSSH(ctx context.Context, config entities.SSHConfig) error {
	client, err := dialSSH(ctx, config)
	if err != nil {
		return err
	}
	return client.Close()
}

func dialSSH(ctx context.Context, cfg entities.SSHConfig) (*ssh.Client, error) {
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Username = strings.TrimSpace(cfg.Username)
	if cfg.Host == "" || cfg.Username == "" {
		return nil, fmt.Errorf("SSH host and username are required")
	}
	if cfg.Port == 0 {
		cfg.Port = 22
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("SSH port must be between 1 and 65535")
	}
	var auth []ssh.AuthMethod
	var agentConn net.Conn
	defer func() {
		if agentConn != nil {
			_ = agentConn.Close()
		}
	}()
	switch cfg.AuthMethod {
	case "", "auto", "password":
		// Match the common passwordless terminal workflow: agent, then local keys.
		var signers []ssh.Signer
		if socket := os.Getenv("SSH_AUTH_SOCK"); socket != "" {
			agentCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			agentConn, _ = (&net.Dialer{}).DialContext(agentCtx, "unix", socket)
			cancel()
			if agentConn != nil {
				_ = agentConn.SetDeadline(time.Now().Add(3 * time.Second))
				signers, _ = agent.NewClient(agentConn).Signers()
				_ = agentConn.SetDeadline(time.Time{})
			}
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		for _, name := range []string{"id_ed25519", "id_ecdsa", "id_rsa"} {
			key, err := os.ReadFile(filepath.Join(home, ".ssh", name))
			if err != nil {
				continue
			}
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				signers = append(signers, signer)
			}
		}
		if len(signers) > 0 {
			auth = append(auth, ssh.PublicKeys(signers...))
		}
		if cfg.Password != "" {
			auth = append(auth, ssh.Password(cfg.Password))
		}
	case "private_key":
		key := []byte(cfg.PrivateKey)
		if !strings.Contains(cfg.PrivateKey, "-----BEGIN") {
			path := cfg.PrivateKey
			if strings.HasPrefix(path, "~/") {
				home, err := os.UserHomeDir()
				if err != nil {
					return nil, err
				}
				path = filepath.Join(home, path[2:])
			}
			var err error
			key, err = os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("read SSH private key: %w", err)
			}
		}
		var signer ssh.Signer
		var err error
		if cfg.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(cfg.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(key)
		}
		if err != nil {
			return nil, fmt.Errorf("parse SSH private key: %w", err)
		}
		auth = append(auth, ssh.PublicKeys(signer))
	default:
		return nil, fmt.Errorf("unsupported SSH authentication method %q", cfg.AuthMethod)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	verify, err := knownhosts.New(filepath.Join(home, ".ssh", "known_hosts"))
	if err != nil {
		return nil, fmt.Errorf("load SSH known_hosts (connect with OpenSSH first to trust the server): %w", err)
	}
	setup, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	address := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	raw, err := (&net.Dialer{}).DialContext(setup, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("dial SSH server: %w", err)
	}
	stop := context.AfterFunc(setup, func() { _ = raw.Close() })
	deadline, _ := setup.Deadline()
	_ = raw.SetDeadline(deadline)
	conn, channels, requests, err := ssh.NewClientConn(raw, address, &ssh.ClientConfig{User: cfg.Username, Auth: auth, HostKeyCallback: verify})
	stop()
	if err != nil {
		_ = raw.Close()
		if (cfg.AuthMethod == "" || cfg.AuthMethod == "auto" || cfg.AuthMethod == "password") && cfg.Password == "" {
			return nil, fmt.Errorf("SSH handshake using local keys/agent: %w (for a custom or encrypted key, choose Private Key or load it into ssh-agent)", err)
		}
		return nil, fmt.Errorf("SSH handshake: %w", err)
	}
	if err := setup.Err(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	_ = raw.SetDeadline(time.Time{})
	client := ssh.NewClient(conn, channels, requests)
	return client, nil
}

func (s *SSH) serve(ctx context.Context, listener net.Listener, client *ssh.Client) {
	defer s.wg.Done()
	for {
		local, err := listener.Accept()
		if err != nil {
			return
		}
		s.mu.Lock()
		if ctx.Err() != nil {
			s.mu.Unlock()
			_ = local.Close()
			return
		}
		s.sockets[local] = struct{}{}
		s.wg.Add(1)
		s.mu.Unlock()
		go s.forward(ctx, client, local)
	}
}

func (s *SSH) forward(ctx context.Context, client *ssh.Client, local net.Conn) {
	defer s.wg.Done()
	defer func() { _ = local.Close(); s.mu.Lock(); delete(s.sockets, local); s.mu.Unlock() }()
	host := s.dbHost
	if host == "" {
		host = "127.0.0.1"
	}
	dialCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	remote, err := client.DialContext(dialCtx, "tcp", net.JoinHostPort(host, strconv.Itoa(s.dbPort)))
	cancel()
	if err != nil {
		return
	}
	defer remote.Close()
	done := make(chan struct{})
	go func() { _, _ = io.Copy(remote, local); _ = remote.Close(); close(done) }()
	_, _ = io.Copy(local, remote)
	_ = local.Close()
	<-done
}

func (s *SSH) Close() error {
	s.lifecycle.Lock()
	defer s.lifecycle.Unlock()
	s.mu.Lock()
	if s.client == nil {
		s.mu.Unlock()
		return nil
	}
	s.cancel()
	_ = s.listener.Close()
	err := s.client.Close()
	for socket := range s.sockets {
		_ = socket.Close()
	}
	s.mu.Unlock()
	s.wg.Wait()
	s.mu.Lock()
	s.client, s.listener, s.cancel = nil, nil, nil
	s.mu.Unlock()
	return err
}

func (s *SSH) Address() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

func ForConfig(config entities.ConnectionConfig) (Transport, error) {
	if config.SSHConfig != nil {
		if config.Type != "mysql" && config.Type != "pgx" && config.Type != "mongodb" {
			return nil, fmt.Errorf("SSH tunnels are supported only for MySQL, PostgreSQL, and MongoDB")
		}
		if strings.HasPrefix(config.Host, "mongodb://") || strings.HasPrefix(config.Host, "mongodb+srv://") {
			return nil, fmt.Errorf("use a database host and port instead of a MongoDB URI with SSH")
		}
		return NewSSHConfig(config.Host, config.Port, *config.SSHConfig), nil
	}
	if config.SSHConfigID != nil {
		return nil, fmt.Errorf("SSH configuration details are missing")
	}
	return NewDirect(config.Host, config.Port), nil
}
