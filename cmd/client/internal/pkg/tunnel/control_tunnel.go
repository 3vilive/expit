package tunnel

import (
	"log"
	"net"
	"time"

	"github.com/3vilive/expit/pkg/pb"
	protocol "github.com/3vilive/expit/pkg/protocol/v2"
	"github.com/pkg/errors"
)

type ControlTunnel struct {
	serverAddr     string
	dataServerAddr string
	clientId       string
	conn           net.Conn
}

func NewControlTunnel(server string) (*ControlTunnel, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, errors.Wrap(err, "dial server error")
	}

	tunnel := &ControlTunnel{
		serverAddr: server,
		conn:       conn,
	}

	if err := tunnel.Handshake(); err != nil {
		return nil, errors.Wrap(err, "handshake error")
	}

	return tunnel, nil
}

func (t *ControlTunnel) Handshake() error {
	handshakeMsg := protocol.NewHandshakeMsg("Drawin", "test", "", pb.Handshake_CONTROL)
	if err := protocol.WriteMessage(t.conn, handshakeMsg); err != nil {
		return errors.Wrap(err, "write handshake message error")
	}

	handshakeReply, err := protocol.ReadMessage(t.conn)
	if err != nil {
		return errors.Wrap(err, "read handshake reply error")
	}

	if handshakeReply.Type != pb.Message_HANDSHAKE_REPLY {
		return errors.New("unexpected message during handshake stage")
	}

	replyData := handshakeReply.GetHandshakeReply()
	if len(replyData.Error) != 0 {
		return errors.Errorf("handshake error: %s", replyData.Error)
	}

	t.dataServerAddr = replyData.DataServerAddr
	t.clientId = replyData.ClientId

	return nil
}

func (t *ControlTunnel) Start() {
	closed := make(chan bool)
	errChan := make(chan error)
	msgChan := make(chan *pb.Message)
	heartbeatInterval := time.NewTicker(5 * time.Second)

	go func() {
		for {
			msg, err := protocol.ReadMessage(t.conn)
			if err != nil {
				errChan <- err
				return
			}

			msgChan <- msg
		}
	}()

	for {
		select {
		case err := <-errChan:
			log.Printf("[control tunnel] error occurred: %s\n", err)
			close(closed)
			return

		case <-closed:
			return

		case <-heartbeatInterval.C:
			if err := protocol.WriteMessage(t.conn, protocol.NewHeartbeatMsg()); err != nil {
				log.Printf("[control tunnel] write heartbeat message error: %s\n", err)
				errChan <- err
			}

		case msg := <-msgChan:
			log.Printf("[control tunnel] received message: %s\n", msg)
			switch msg.Type {
			case pb.Message_REQ_DATA_TUNNEL:
				t.HandleReqDataTunnel(msg.GetReqDataTunnel())
				log.Printf("[control tunnel] spawn data tunnel")
			}
		}
	}
}

func (t *ControlTunnel) HandleReqDataTunnel(msg *pb.RequestDataTunnel) {
	if msg.Number > 20 {
		msg.Number = 20
	}

	for i := 0; i < int(msg.Number); i += 1 {
		go func() {
			tunnel, err := NewDataTunnel(t.serverAddr, t.clientId)
			if err != nil {
				log.Printf("[control tunnel] failed to create data tunnel: %s\n", err)
				return
			}

			log.Printf("[control tunnel] start running")
			tunnel.Start()
		}()
	}
}
