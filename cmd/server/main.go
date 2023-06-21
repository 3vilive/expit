package main

import (
	"log"
	"net"

	"github.com/3vilive/expit/cmd/server/internal/pkg/manager"
	"github.com/3vilive/expit/cmd/server/internal/pkg/tunnel"
	"github.com/3vilive/expit/pkg/pb"
	protocol "github.com/3vilive/expit/pkg/protocol/v2"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:20000")
	if err != nil {
		log.Printf("resolve tcp addr error: %s\n", err.Error())
		return
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Printf("listen error: %s\n", err.Error())
		return
	}

	log.Printf("listen on %s, start accept connections\n", addr)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("accept connection error: %s\n", err.Error())
			continue
		}

		go handleNewConnection(connection)
	}
}

func handleNewConnection(conn net.Conn) {
	// handshake stage
	handshakeMsg, err := protocol.ReadMessage(conn)
	if err != nil {
		log.Printf("read message error: %s\n", err)
		return
	}

	if handshakeMsg.Type != pb.Message_HANDSHAKE {
		log.Printf("unexpected message in handshake stage: %s\n", handshakeMsg.Type)
		return
	}

	handshakeData := handshakeMsg.GetHandshake()
	switch handshakeData.TunnelType {
	case pb.Handshake_CONTROL:
		handleControlTunnel(conn)
	case pb.Handshake_DATA:
		handleDataTunnel(conn, handshakeData)
	default:
		log.Printf("unknown register tunnel type: %s\n", handshakeData.TunnelType)
		return
	}
}

func handleControlTunnel(conn net.Conn) {
	remoteAddr := conn.RemoteAddr()
	log.Printf("[control tunnel %s] handle control tunnel\n", remoteAddr)

	controlTunnel, err := tunnel.NewControlTunnel(conn)
	if err != nil {
		log.Printf("[control tunnel %s] run control tunnel error: %s\n", remoteAddr, err)
		return
	}

	if err := controlTunnel.DoHandshakeReply(); err != nil {
		log.Printf("[control tunnel %s] write handshake replay msg error: %s\n", remoteAddr, err)
		return
	}

	go controlTunnel.StartHandleMsg()
	go controlTunnel.StartHandleDataServer()

	tunnelManager := manager.GetControlTunnelManager()
	tunnelManager.AddControlTunnel(controlTunnel.Id(), controlTunnel)
}

func handleDataTunnel(conn net.Conn, handshake *pb.Handshake) {
	remoteAddr := conn.RemoteAddr().String()

	dataTunnel := tunnel.NewDataTunnel(handshake.ClientId, conn)
	if err := dataTunnel.DoHandShakeReply(); err != nil {
		log.Printf("[data tunnel - %s] write handshake reply error: %s\n", remoteAddr, err)
		conn.Close()
		return
	}

	// go dataTunnel.Start()
	tunnelManager := manager.GetControlTunnelManager()
	controlTunnel := tunnelManager.GetControlTunnelById(handshake.ClientId)
	if controlTunnel != nil {
		controlTunnel.PutDataTunnel(dataTunnel)
	} else {
		log.Printf("unexpected nil control tunnel, client id %s\n", handshake.ClientId)
	}
}
