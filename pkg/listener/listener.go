package listener

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

type Listener struct {
	ID       string
	Protocol string
	Host     string
	Port     int
	Status   string
	Cancel   context.CancelFunc
}

func NewListener(protocol, host string, port int) *Listener {
	return &Listener{
		ID:       uuid.New().String(),
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Status:   "stopped",
	}
}

func (l *Listener) StartTLS(ctx context.Context, connHandler func(net.Conn), tlsConfig *tls.Config) error {
	listener, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", l.Host, l.Port), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening on %s", fmt.Sprintf("%s:%d", l.Host, l.Port))

	go func() {
		<-ctx.Done()
		listener.Close()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Printf("failed to accept connection: %v", err)
				continue
			}
		}
		go connHandler(conn)
	}
}
