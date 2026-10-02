package transports

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"db-viewer/internal/engine/entities"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestSSHTunnel(t *testing.T) {
	_, key, _ := ed25519.GenerateKey(rand.Reader)
	signer, _ := ssh.NewSignerFromKey(key)
	cfg := &ssh.ServerConfig{PasswordCallback: func(_ ssh.ConnMetadata, _ []byte) (*ssh.Permissions, error) { return nil, nil }}
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
