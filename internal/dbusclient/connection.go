package dbusclient

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	busName           = "com.ananthvk.nightlightmonitor"
	nightLightBusPath = "/org/kde/KWin/NightLight"
	signalChannelSize = 10
)

// Connection represents a connection to the DBUS. A handler must be set, which will be called whenever this connection
// gets a new signal
type Connection struct {
	conn          *dbus.Conn
	handler       func(signal *dbus.Signal)
	signalChannel chan *dbus.Signal
}

// SetSignalHandler sets a function as a handler. This function will be called whenever this connection receives a new
// signal
func (c *Connection) SetSignalHandler(handler func(signal *dbus.Signal)) {
	c.handler = handler
}

// Close closes the underlying dbus connection
func (c *Connection) Close() error {
	return c.conn.Close()
}

// NewConnection creates a new DBUS connection
func NewConnection() (*Connection, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("connect failed: %w", err)
	}
	if _, err := conn.RequestName(busName, dbus.NameFlagDoNotQueue); err != nil {
		return nil, fmt.Errorf("request connection name failed: %w", err)
	}
	if err = conn.AddMatchSignal(dbus.WithMatchPathNamespace(nightLightBusPath)); err != nil {
		return nil, fmt.Errorf("add match signal failed: %w", err)
	}
	e := make(chan *dbus.Signal, signalChannelSize)
	conn.Signal(e)
	connection := &Connection{conn: conn, signalChannel: e}
	return connection, nil
}

// RunLoop starts the main event loop for the D-Bus connection, processing incoming
// signals until the provided context is cancelled. This method blocks until the
// context is done. It should be run in a separate goroutine to avoid blocking.
// Each signal is handled in a separate goroutine by the handler function.
func (c *Connection) RunLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case signal := <-c.signalChannel:
			if c.handler == nil {
				panic("signal handler not set")
			}
			go c.handler(signal)
		}
	}
}
