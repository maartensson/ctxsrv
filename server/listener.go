package server

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/coreos/go-systemd/v22/activation"
)

func ActivationListener(port int, fallbackTCP bool) (net.Listener, error) {
	lis, err := activation.Listeners()
	if err != nil {
		return nil, fmt.Errorf("failed to get systemd socket listeners: %w", err)
	}

	for _, ln := range lis {
		if addr, ok := ln.Addr().(*net.TCPAddr); ok && addr.Port == port {
			return ln, nil
		}
	}

	if fallbackTCP {
		return ListenTCP(port)
	}

	return nil, fmt.Errorf("no socket listeners available for port: %d", port)
}

func ActivationListenerTLS(port int, tlsConfig *tls.Config, fallbackTCP bool) (net.Listener, error) {
	ln, err := ActivationListener(port, fallbackTCP)
	if err != nil {
		return nil, fmt.Errorf("failed to get systemd socket listeners: %w", err)
	}

	if tlsConfig == nil {
		return ln, nil
	}

	return tls.NewListener(ln, tlsConfig), nil
}

func ListenTCP(port int) (net.Listener, error) {
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}
	return ln, nil
}
