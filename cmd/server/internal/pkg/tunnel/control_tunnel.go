package tunnel

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/3vilive/expit/pkg/pb"
	protocol "github.com/3vilive/expit/pkg/protocol/v2"
	"github.com/pkg/errors"
)

type ControlTunnel struct {
	id             string
	conn           net.Conn
	dataServer     *net.TCPListener
	msgChan        chan *pb.Message
	errChan        chan error
	dataTunnelChan chan *DataTunnel
	closed         chan bool
}

func NewControlTunnel(conn net.Conn) (*ControlTunnel, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, errors.Wrap(err, "resolve tcp addr error")
	}
	dataServer, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "listen tcp server error")
	}

	controlTunnel := &ControlTunnel{
		id:             conn.RemoteAddr().String(),
		conn:           conn,
		dataServer:     dataServer,
		errChan:        make(chan error),
		dataTunnelChan: make(chan *DataTunnel, 20),
		closed:         make(chan bool),
	}

	return controlTunnel, nil
}

func (t *ControlTunnel) Id() string {
	return t.id
}

func (t *ControlTunnel) GetDataServerAddr() string {
	return t.dataServer.Addr().String()
}

func (t *ControlTunnel) DoHandshakeReply() error {
	// reply handshake
	dataServerLocalAddr := t.GetDataServerAddr()
	dataServerPort := strings.Split(dataServerLocalAddr, ":")[1]
	dataServerRemoteAddr := fmt.Sprintf("0.0.0.0:%s", dataServerPort) // TODO: remote addr

	return t.WriteMessage(protocol.NewHandshakeReplyMsg(
		"", dataServerRemoteAddr, t.id,
	))
}

func (t *ControlTunnel) StartHandleDataServer() {
	// pre request data tunnel
	t.RequestDataTunnel(10)

	dataServerAddr := t.dataServer.Addr().String()
	log.Printf("[data server - %s] start running\n", dataServerAddr)

	clientConnChan := make(chan net.Conn)

	go func() {
		for {
			clientConn, err := t.dataServer.Accept()
			if err != nil {
				log.Printf("[data server - %s]accept connection error: %s\n", dataServerAddr, err)
				t.errChan <- err
				return
			}

			clientConnChan <- clientConn
		}
	}()

	for {
		select {
		case <-t.closed:
			t.dataServer.Close()
			return

		case clientConn := <-clientConnChan:
			log.Printf("[data server - %s] accept connection %s\n", dataServerAddr, clientConn.RemoteAddr())
			// get data tunnel
			dataTunnel := t.GetDataTunnel()
			if dataTunnel == nil {
				log.Printf("[data server - %s] failed to get data tunnel, close client connection", dataServerAddr)
				clientConn.Close()
				return
			}
			// forward data
			go dataTunnel.StartForward(clientConn)
		}
	}
}

func (t *ControlTunnel) GetDataTunnel() *DataTunnel {
	select {
	case tunnel := <-t.dataTunnelChan:
		return tunnel
	default:
		t.RequestDataTunnel(10)
		timeout := time.NewTimer(5 * time.Second)

		select {
		case <-timeout.C:
			return nil
		case tunnel := <-t.dataTunnelChan:
			return tunnel
		}
	}
}

func (t *ControlTunnel) PutDataTunnel(tunnel *DataTunnel) {
	t.dataTunnelChan <- tunnel
}

func (t *ControlTunnel) StartHandleMsg() error {
	for {
		select {
		case msg := <-t.Messages():
			log.Printf("[control tunnel - %s] received msg: %v\n", t.id, msg)
		case err := <-t.errChan:
			close(t.closed)
			return err
		case <-t.closed:
			return nil
		}
	}
}

func (t *ControlTunnel) Messages() <-chan *pb.Message {
	if t.msgChan != nil {
		return t.msgChan
	}

	t.msgChan = make(chan *pb.Message)

	go func() {
		for {
			msg, err := protocol.ReadMessage(t.conn)
			if err != nil {
				t.errChan <- err
				return
			}

			t.msgChan <- msg
		}
	}()

	return t.msgChan
}

func (t *ControlTunnel) WriteMessage(msg *pb.Message) error {
	return protocol.WriteMessage(t.conn, msg)
}

func (t *ControlTunnel) RequestDataTunnel(n int) error {
	return t.WriteMessage(protocol.NewRequestDataTunnelMsg(n))
}
