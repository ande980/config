package env

import (
	"os"
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

func TestEnv(t *testing.T) {
	cfg := &Config{
		A: "test",
		B: true,
		C: time.Second,
		Server: &Server{
			Addr: ":8080",
		},
	}

	os.Setenv("BUBBLES_A", "not test")
	os.Setenv("BUBBLES_B", "false")
	os.Setenv("C", "3h")
	os.Setenv("BUBBLES_SERVER_ADDR", ":9090")

	p := WithPrefix("bubbles")
	if err := p.Parse(cfg); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if cfg.A != "not test" {
		t.Errorf("expected '%s', got '%s'", "not test", cfg.A)
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
