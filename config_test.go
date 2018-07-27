package config

import (
	"config/env"
	"config/json"
	"os"
	"strings"
	"testing"
	"time"
)

type Config struct {
	A      string `json:"abcd" flag:"Z" usage:"The string value of A"`
	B      bool
	C      time.Duration `json:"-"`
	Server *Server
}

type Server struct {
	Addr string
}

func TestConfig(t *testing.T) {
	cfg := &Config{
		A: "test",
		B: true,
		C: time.Second,
		Server: &Server{
			Addr: ":8080",
		},
	}

	os.Setenv("BUBBLES_A", "not test")
	os.Setenv("B", "false")
	os.Setenv("BUBBLES_C", "4h")
	os.Setenv("BUBBLES_SERVER_ADDR", ":9091")

	r := strings.NewReader(`{"abcd":"princes of the universe","B":true,"C":"3h","Server":{"Addr":":9090"}}`)

	providers = []Provider{
		json.WithReader(r),
		env.WithPrefix("bubbles"),
	}

	if err := Parse(cfg); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if cfg.A != "not test" {
		t.Errorf("expected '%s', got '%s'", "not test", cfg.A)
	}

	if !cfg.B {
		t.Errorf("expected %t, got %t", true, cfg.B)
	}

	if cfg.C != time.Hour*4 {
		t.Errorf("expected %s, got %s", time.Hour*4, cfg.C)
	}

	if cfg.Server.Addr != ":9091" {
		t.Errorf("expected '%s', got '%s'", ":9091", cfg.Server.Addr)
	}
}
