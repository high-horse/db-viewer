package transports

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"db-lens/internal/engine/entities"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

func TestSSHTunnel(t *testing.T) {
	_, key, _ := ed25519.GenerateKey(rand.Reader)
	signer, _ := ssh.NewSignerFromKey(key)
	cfg := &ssh.ServerConfig{PasswordCallback: func(meta ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		if meta.User() != "tester" || string(password) != "secret" {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, nil
	}}
	cfg.PublicKeyCallback = func(meta ssh.ConnMetadata, public ssh.PublicKey) (*ssh.Permissions, error) {
		if meta.User() == "tester" && bytes.Equal(public.Marshal(), signer.PublicKey().Marshal()) {
			return nil, nil
		}
		return nil, fmt.Errorf("invalid key")
	}
	cfg.AddHostKey(signer)
	server, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.Mkdir(filepath.Join(home, ".ssh"), 0700); err != nil {
		t.Fatal(err)
	}
	known := filepath.Join(home, ".ssh", "known_hosts")
	if err := os.WriteFile(known, []byte(knownhosts.Line([]string{knownhosts.Normalize(server.Addr().String())}, signer.PublicKey())+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			raw, err := server.Accept()
			if err != nil {
				return
			}
			go func() {
				conn, channels, requests, err := ssh.NewServerConn(raw, cfg)
				if err != nil {
					raw.Close()
					return
				}
				defer conn.Close()
				go ssh.DiscardRequests(requests)
				for request := range channels {
					var target struct {
						Host       string
						Port       uint32
						Origin     string
						OriginPort uint32
					}
					_ = ssh.Unmarshal(request.ExtraData(), &target)
					if request.ChannelType() != "direct-tcpip" || target.Host != "database.internal" || target.Port != 5432 {
						request.Reject(ssh.ConnectionFailed, "wrong target")
						continue
					}
					channel, reqs, err := request.Accept()
					if err != nil {
						continue
					}
					go ssh.DiscardRequests(reqs)
					go func() { defer channel.Close(); _, _ = io.Copy(channel, channel) }()
				}
			}()
		}
	}()
	host, portString, _ := net.SplitHostPort(server.Addr().String())
	port, _ := strconv.Atoi(portString)
	config := entities.SSHConfig{Host: host, Port: port, Username: "tester", AuthMethod: "password", Password: "secret"}
	if err := TestSSH(context.Background(), config); err != nil {
		t.Fatalf("SSH test failed: %v", err)
	}
	spaced := config
	spaced.Host = " " + config.Host + " "
	spaced.Username = " tester "
	if err := TestSSH(context.Background(), spaced); err != nil {
		t.Fatalf("whitespace normalization failed: %v", err)
	}
	invalid := config
	invalid.Password = "wrong"
	if err := TestSSH(context.Background(), invalid); err == nil {
		t.Fatal("SSH test accepted invalid credentials")
	}
	cancelled, stop := context.WithCancel(context.Background())
	stop()
	if err := TestSSH(cancelled, config); err == nil {
		t.Fatal("SSH test ignored cancellation")
	}
	keyBlock, err := ssh.MarshalPrivateKey(key, "test identity")
	if err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(home, ".ssh", "id_ed25519")
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(keyBlock), 0600); err != nil {
		t.Fatal(err)
	}
	config.Password = ""
	t.Setenv("SSH_AUTH_SOCK", "")
	if err := TestSSH(context.Background(), config); err != nil {
		t.Fatalf("default local key failed: %v", err)
	}
	customPath := filepath.Join(home, ".ssh", "custom_key")
	if err := os.Rename(keyPath, customPath); err != nil {
		t.Fatal(err)
	}
	explicit := config
	explicit.AuthMethod, explicit.PrivateKey = "private_key", customPath
	if err := TestSSH(context.Background(), explicit); err != nil {
		t.Fatalf("explicit key failed: %v", err)
	}
	socket := filepath.Join(home, "agent.sock")
	agentListener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	defer agentListener.Close()
	keyring := agent.NewKeyring()
	if err := keyring.Add(agent.AddedKey{PrivateKey: key}); err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			conn, err := agentListener.Accept()
			if err != nil {
				return
			}
			go func() { defer conn.Close(); _ = agent.ServeAgent(keyring, conn) }()
		}
	}()
	t.Setenv("SSH_AUTH_SOCK", socket)
	if err := TestSSH(context.Background(), config); err != nil {
		t.Fatalf("agent authentication failed: %v", err)
	}
	tunnel := NewSSHConfig("database.internal", 5432, config)
	ctx, cancel := context.WithCancel(context.Background())
	if err := tunnel.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	cancel()
	defer tunnel.Close()
	local, err := net.DialTimeout("tcp", tunnel.Address(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer local.Close()
	_ = local.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := local.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 5)
	if _, err := io.ReadFull(local, buffer); err != nil {
		t.Fatal(err)
	}
	if string(buffer) != "hello" {
		t.Fatalf("got %q", buffer)
	}
	if err := tunnel.Close(); err != nil {
		t.Fatal(err)
	}
	if tunnel.Address() != "" {
		t.Fatal("listener still present")
	}
	if err := tunnel.Close(); err != nil {
		t.Fatal(err)
	}
	_, other, _ := ed25519.GenerateKey(rand.Reader)
	otherSigner, _ := ssh.NewSignerFromKey(other)
	if err := os.WriteFile(known, []byte(knownhosts.Line([]string{knownhosts.Normalize(server.Addr().String())}, otherSigner.PublicKey())+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := TestSSH(context.Background(), config); err == nil {
		t.Fatal("SSH test accepted an untrusted host key")
	}
	if err := NewSSHConfig("database.internal", 5432, config).Connect(context.Background()); err == nil {
		t.Fatal("accepted an untrusted host key")
	}
}

func TestTransportSelection(t *testing.T) {
	for _, dialect := range []string{"mysql", "pgx", "mongodb"} {
		transport, err := ForConfig(entities.ConnectionConfig{Type: dialect, SSHConfig: &entities.SSHConfig{}})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := transport.(*SSH); !ok {
			t.Fatal("SSH transport not selected")
		}
	}
	if _, err := ForConfig(entities.ConnectionConfig{Type: "sqlite", SSHConfig: &entities.SSHConfig{}}); err == nil {
		t.Fatal("SQLite accepted SSH")
	}
	if got := NewDirect("::1", 5432).Address(); got != "[::1]:5432" {
		t.Fatal(got)
	}
}
