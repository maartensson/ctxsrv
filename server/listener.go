package server

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/coreos/go-systemd/v22/activation"
)

func ListenTCP(port int) (net.Listener, error) {
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}
	return ln, nil
}

func ActivationListener(opts ...ListenerOption) (net.Listener, error) {
	cfg := defaultListenerConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	lis, err := activation.Listeners()
	if err != nil {
		return nil, fmt.Errorf("failed to get systemd socket listeners: %w", err)
	}

	if len(lis) == 0 {
		return nil, fmt.Errorf("no systemd sockets found")
	}

	if cfg.port != nil {
		for _, ln := range lis {
			if addr, ok := ln.Addr().(*net.TCPAddr); ok && addr.Port == *cfg.port {
				if cfg.tlsConfig != nil {
					return tls.NewListener(ln, cfg.tlsConfig), nil
				}
				return ln, nil
			}
		}
		return nil, fmt.Errorf("no socket listeners available for port: %d", *cfg.port)
	}

	if cfg.tlsConfig != nil {
		return tls.NewListener(lis[0], cfg.tlsConfig), nil
	}

	return lis[0], nil
}

type ListenerOption func(*listenerConfig)

type listenerConfig struct {
	port      *int
	tlsConfig *tls.Config
}

func defaultListenerConfig() listenerConfig {
	return listenerConfig{}
}

func WithPort(port int) ListenerOption {
	return func(c *listenerConfig) {
		c.port = &port
	}
}

func WithTLS(cfg *tls.Config) ListenerOption {
	return func(c *listenerConfig) {
		c.tlsConfig = cfg
	}
}
