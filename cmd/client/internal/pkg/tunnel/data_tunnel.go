package tunnel

import (
	"io"
	"log"
	"net"

	"github.com/3vilive/expit/cmd/client/internal/pkg/config"
	"github.com/3vilive/expit/pkg/common"
	"github.com/3vilive/expit/pkg/pb"
	protocol "github.com/3vilive/expit/pkg/protocol/v2"
	"github.com/pkg/errors"
)

type DataTunnel struct {
	serverAddr string
	clientId   string
	conn       net.Conn
}

func NewDataTunnel(server, clientId string) (*DataTunnel, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial server")
	}

	dataTunnel := &DataTunnel{
		serverAddr: server,
		clientId:   clientId,
		conn:       conn,
	}

	if err := dataTunnel.Handshake(); err != nil {
		return nil, errors.Wrap(err, "failed to handshake")
	}

	return dataTunnel, nil
}

func (t *DataTunnel) Handshake() error {
	err := protocol.WriteMessage(
		t.conn,
		protocol.NewHandshakeMsg("Darwin", common.VERSION, t.clientId, pb.Handshake_DATA),
	)
	if err != nil {
		return errors.Wrap(err, "write handshake message error")
	}

	replyMsg, err := protocol.ReadMessage(t.conn)
	if err != nil {
		return errors.Wrap(err, "read handshake reply message error")
	}

	if replyMsg.Type != pb.Message_HANDSHAKE_REPLY {
		return errors.Errorf("unexpected message type during handshake: %s", replyMsg.Type)
	}

	handshakeReply := replyMsg.GetHandshakeReply()
	if len(handshakeReply.Error) > 0 {
		return errors.Errorf("handshake error: %s", handshakeReply.Error)
	}

	return nil
}

func (t *DataTunnel) Start() {
	for {
		msg, err := protocol.ReadMessage(t.conn)
		if err != nil {
			log.Printf("[data tunnel] read message error: %s\n", err)
			return
		}

		switch msg.Type {
		case pb.Message_START_FORWARD:

			localApp := config.LocalAppAddr

			log.Printf("[data tunnel] connect to local app %s\n", localApp)
			localAppConn, err := net.Dial("tcp", localApp)

			if err != nil {
				log.Printf("[data tunnel] failed to dial local app: %s\n", err)
				return
			}

			if err := protocol.WriteMessage(t.conn, protocol.NewStartForwardMsgReply("")); err != nil {
				log.Printf("[data tunnel] failed to write forward reply message: %s\n", err)
				localAppConn.Close()
				return
			}

			go t.Forward(localAppConn)
			return
		}
	}
}

func (t *DataTunnel) Forward(localAppConn net.Conn) {
	go io.Copy(localAppConn, t.conn)
	go io.Copy(t.conn, localAppConn)
}
