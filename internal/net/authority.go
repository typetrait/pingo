package net

import (
	"context"
	"fmt"
	"net"
)

type Authority interface {
	Connect(ctx context.Context) error
	RequestMatch(ctx context.Context)
}

// ---

type ServerAuthority struct {
}

func NewServerAuthority() *ServerAuthority {
	return &ServerAuthority{}
}

func (sa *ServerAuthority) Connect(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", ":7777")
	if err != nil {
		return fmt.Errorf("resolving game server address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("connecting to game server (%s): %w", addr, err)
	}

	_, err = conn.Write([]byte("bytes"))
	if err != nil {
		return err
	}

	return nil
}

func (sa *ServerAuthority) RequestMatch(ctx context.Context) {

}
