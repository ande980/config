package yaml

import (
	"strings"
	"testing"
	"time"
)

type Config struct {
	A      string `yaml:"abcd" flag:"Z" usage:"The string value of A"`
	B      bool
	C      time.Duration `yaml:"-"`
	Server *Server
}

type Server struct {
	Addr string
}

func TestYAML(t *testing.T) {
	cfg := &Config{
		A: "test",
		B: true,
		C: time.Second,
		Server: &Server{
			Addr: ":8080",
		},
	}

	r := strings.NewReader(`abcd: princes of the universe
b: false
c: 3h
server:
  addr: :9090
`)
	p := WithReader(r)
	if err := p.Parse(cfg); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if cfg.A != "princes of the universe" {
		t.Errorf("expected '%s', got '%s'", "princes of the universe", cfg.A)
	}

	if cfg.B {
		t.Errorf("expected %t, got %t", false, cfg.B)
	}

	if cfg.C != time.Second {
		t.Errorf("expected %s, got %s", time.Second, cfg.C)
	}

	if cfg.Server.Addr != ":9090" {
		t.Errorf("expected '%s', got '%s'", ":9090", cfg.Server.Addr)
	}
}

func TestNoFile(t *testing.T) {
	cfg := &Config{
		A: "test",
		B: true,
		C: time.Second,
		Server: &Server{
			Addr: ":8080",
		},
	}

	p := New()
	if err := p.Parse(cfg); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
