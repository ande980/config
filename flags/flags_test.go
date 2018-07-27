package flags

import (
	"flag"
	"testing"
	"time"
)

type Config struct {
	A      string `flag:"Z" usage:"The string value of A"`
	B      bool
	C      time.Duration
	Server *Server
}

type Server struct {
	Addr string
}

func TestFlags(t *testing.T) {
	cfg := &Config{
		A: "test",
		B: true,
		C: time.Second,
		Server: &Server{
			Addr: ":8080",
		},
	}

	f := New()
	if err := f.parse(cfg, "--server-addr", ":80", "-c", "3h"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	f.VisitAll(func(f *flag.Flag) {
		if f.Name == "z" {
			if f.Usage != "The string value of A" {
				t.Errorf("%s: expected '%s', got '%s'", f.Name, "The string value of A", f.Usage)
			}
		} else {
			if f.Usage != "" {
				t.Errorf("%s: expected '%s', got '%s'", f.Name, "", f.Usage)
			}
		}
	})

	if !cfg.B {
		t.Errorf("expected %t, got %t", true, cfg.B)
	}

	if cfg.Server.Addr != ":80" {
		t.Errorf("expected '%s', got '%s'", ":80", cfg.Server.Addr)
	}

	if cfg.C != time.Hour*3 {
		t.Errorf("expected '%s', got '%s'", "3h", cfg.C)
	}
}
