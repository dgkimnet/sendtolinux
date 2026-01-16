package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sendtolinux/internal/dbussvc"
	"sendtolinux/internal/httpserver"

	"github.com/godbus/dbus/v5"
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("connect session bus: %v", err)
	}
	defer conn.Close()

	reply, err := conn.RequestName(dbussvc.ServiceName, dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatalf("request name: %v", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatalf("name %s already taken", dbussvc.ServiceName)
	}

	svc := dbussvc.New(conn)
	if err := conn.Export(svc, dbus.ObjectPath(dbussvc.ObjectPath), dbussvc.InterfaceName); err != nil {
		log.Fatalf("export service: %v", err)
	}

	srv, err := httpserver.Start(svc)
	if err != nil {
		log.Fatalf("start http server: %v", err)
	}

	log.Printf("D-Bus service running as %s", dbussvc.ServiceName)
	if os.Getenv("STL_EMIT_TEST") == "1" {
		if err := svc.EmitTestSignal(); err != nil {
			log.Printf("emit test signal: %v", err)
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("shutting down")
	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http shutdown: %v", err)
		}
		cancel()
	}
}
