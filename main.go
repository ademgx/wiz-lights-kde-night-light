package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"

	"github.com/ananthvk/wiz-lights-kde-night-light/internal/dbusclient"
	"github.com/ananthvk/wiz-lights-kde-night-light/internal/light"
	"github.com/godbus/dbus/v5"
)

const (
	defaultBulbIP   = "192.168.1.140"
	defaultBulbPort = "38899"
)

type Application struct {
	conn      net.PacketConn
	lightAddr *net.UDPAddr
}

func (app *Application) signalHandler(signal *dbus.Signal) {
	temperature, err := dbusclient.GetCurrentTemperature(signal)
	if err != nil {
		slog.Error("signal error", "error", err)
	} else {
		err = light.ChangeLightTemperature(app.conn, app.lightAddr, temperature)
		if err != nil {
			slog.Error("udp packet send failed", "error", err)
		}
		slog.Info("changed temperature of light", "temperature", temperature)
	}
}

func main() {
	bulbIp := flag.String("bulb-ip", defaultBulbIP, "ip address of the bulb")
	bulbPort := flag.String("bulb-port", defaultBulbPort, "port of the bulb")
	flag.Parse()

	conn, err := dbusclient.NewConnection()
	if err != nil {
		slog.Error("error while creating dbus connection", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to dbus")
	connection, err := net.ListenPacket("udp", "")
	if err != nil {
		slog.Error("udp socket open failed", "error", err)
		os.Exit(1)
	}
	slog.Info("created udp socket")
	lightAddr, err := net.ResolveUDPAddr("udp", *bulbIp+":"+*bulbPort)
	if err != nil {
		slog.Error("could not resolve light ip", "error", err)
		os.Exit(1)
	}
	app := &Application{
		conn:      connection,
		lightAddr: lightAddr,
	}
	defer conn.Close()
	conn.SetSignalHandler(app.signalHandler)
	slog.Info("starting run loop")
	conn.RunLoop(context.Background())
}
