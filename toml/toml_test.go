package toml

import (
	"strings"
	"testing"
	"time"
)

type Config struct {
	A      string `toml:"abcd" flag:"Z" usage:"The string value of A"`
	B      bool
	C      time.Duration `toml:"-"`
	Server *Server
}

type Server struct {
	Addr string
}

func TestTOML(t *testing.T) {
	cfg := &Config{
		A: "test",
		B: true,
		C: time.Second,
		Server: &Server{
			Addr: ":8080",
		},
	}

	r := strings.NewReader(`abcd = "princes of the universe"
	B = false
	C = "3h"
	[Server]
	Addr = ":9090"`)
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
