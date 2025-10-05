package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/ananthvk/wiz-lights-kde-night-light/internal/dbusclient"
	"github.com/ananthvk/wiz-lights-kde-night-light/internal/light"
	"github.com/godbus/dbus/v5"
)

// TODO: Edge case: When the lamp is turned off, but the monitor has not yet updated `switchedOn`, if the night light changes temperature,
// the application sends a packet to change the temperature, thereby switching on the light again

const (
	defaultBulbIP   = "192.168.1.140"
	defaultBulbPort = "38899"
	monitorInterval = time.Second * 15
	statusTimeout   = time.Second * 5
)

var switchedOn atomic.Bool

type Application struct {
	conn      net.PacketConn
	lightAddr *net.UDPAddr
}

func (app *Application) signalHandler(signal *dbus.Signal) {
	temperature, err := dbusclient.GetCurrentTemperature(signal)
	if err != nil {
		slog.Error("signal error", "error", err)
	} else {
		if switchedOn.Load() {
			err = light.ChangeLightTemperature(app.conn, app.lightAddr, temperature)
			if err != nil {
				slog.Error("udp packet send failed", "error", err)
			}
			slog.Info("changed temperature of light", "temperature", temperature)
		} else {
			slog.Info("lamp switched off, not changing temperature", "temperature", temperature)
		}
	}
}

func MonitorLampStatus(conn net.PacketConn, addr *net.UDPAddr, interval time.Duration) {
	slog.Info("monitoring lamp")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		status, err := light.CheckStatus(conn, addr, statusTimeout)
		if err != nil {
			slog.Error("error while fetching lamp status", "error", err)
			continue
		}
		if status.Error != nil {
			slog.Error("lamp returned error", "error", status.Error)
			continue
		}
		if status.Result != nil {
			slog.Info("fetched lamp status", "method", *status.Method, "switchedOn", *(*status.Result).State, "scene", *(*status.Result).SceneID, "temp", *(*status.Result).Temp)
			switchedOn.Store(*(*status.Result).State)
		}
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
	defer connection.Close()

	// Open another udp socket to check for led status
	statusConnection, err := net.ListenPacket("udp", "")
	if err != nil {
		slog.Error("udp socket open failed", "error", err)
		os.Exit(1)
	}
	defer statusConnection.Close()

	slog.Info("created udp socket")
	lightAddr, err := net.ResolveUDPAddr("udp", *bulbIp+":"+*bulbPort)
	if err != nil {
		slog.Error("could not resolve light ip", "error", err)
		os.Exit(1)
	}

	go MonitorLampStatus(statusConnection, lightAddr, monitorInterval)

	app := &Application{
		conn:      connection,
		lightAddr: lightAddr,
	}
	defer conn.Close()
	conn.SetSignalHandler(app.signalHandler)
	slog.Info("starting run loop")
	conn.RunLoop(context.Background())
}
