package tunnel

import (
	"io"
	"log"
	"net"

	"github.com/3vilive/expit/pkg/pb"
	protocol "github.com/3vilive/expit/pkg/protocol/v2"
)

type DataTunnel struct {
	clientId string
	conn     net.Conn
}

func NewDataTunnel(clientId string, conn net.Conn) *DataTunnel {
	return &DataTunnel{
		clientId,
		conn,
	}
}

func (t *DataTunnel) DoHandShakeReply() error {
	return protocol.WriteMessage(t.conn, protocol.NewHandshakeReplyMsg(
		"", "", "",
	))
}

func (t *DataTunnel) Start() {
	// for {
	// 	msg, err := protocol.ReadMessage(t.conn)
	// 	if err != nil {
	// 		log.Printf("[data tunnel - %s] read message error: %s\n", t.conn.RemoteAddr(), err)
	// 		return
	// 	}

	// 	log.Printf("[data tunnel - %s] received message: %s\n", t.conn.RemoteAddr(), msg)
	// }
}

func (t *DataTunnel) StartForward(forwardConn net.Conn) {
	// send forward message
	err := protocol.WriteMessage(t.conn, protocol.NewStartForwardMsg(forwardConn.RemoteAddr().String()))
	if err != nil {
		log.Printf("[data tunnel] write start forward message error: %s\n", err)
		forwardConn.Close()
		return
	}

	replyMsg, err := protocol.ReadMessage(t.conn)
	if err != nil {
		log.Printf("[data tunnel] read start forward reply message error: %s\n", err)
		forwardConn.Close()
		return
	}

	if replyMsg.Type != pb.Message_START_FORWARD_REPLY {
		log.Printf("[data tunnel] unexpected reply msg: %s\n", replyMsg.Type)
		forwardConn.Close()
		return
	}

	forwardReply := replyMsg.GetStartForwardReply()
	if len(forwardReply.Error) > 0 {
		log.Printf("[data tunnel] start forward error: %s\n", forwardReply.Error)
		forwardConn.Close()
		return
	}

	// forward data
	go io.Copy(t.conn, forwardConn)
	go io.Copy(forwardConn, t.conn)
}
