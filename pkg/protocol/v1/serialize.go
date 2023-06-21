package v1

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"

	"io"

	"github.com/pkg/errors"
)

const (
	_SizeOfHeader = 4
)

func (m *Message) Encode() ([]byte, error) {
	body, err := json.Marshal(m.Body)
	if err != nil {
		return nil, err
	}
	bodySize := len(body)

	// data: [header(6 bytes), body]
	// header: [message version(2bytes), message type(2 bytes), body size(2 bytes)]
	var buffer bytes.Buffer
	// write header
	binary.Write(&buffer, binary.LittleEndian, uint16(MessageVersion))
	binary.Write(&buffer, binary.LittleEndian, uint16(m.Type))   // write message type
	binary.Write(&buffer, binary.LittleEndian, uint16(bodySize)) // write size of body

	// write body
	buffer.Write(body)

	return buffer.Bytes(), nil
}

func ReadMessage(r io.Reader) (*Message, error) {
	var version uint16
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return nil, errors.Wrap(err, "read message version error")
	}

	if version != MessageVersion {
		return nil, errors.Errorf("mismatch version, got %d want %d", version, MessageVersion)
	}

	var msgType MessageType
	if err := binary.Read(r, binary.LittleEndian, &msgType); err != nil {
		return nil, errors.Wrap(err, "read message type error")
	}

	var bodySize uint16
	if err := binary.Read(r, binary.LittleEndian, &bodySize); err != nil {
		return nil, errors.Wrap(err, "read message body size error")
	}

	var bodyBuffer = make([]byte, bodySize)
	if err := binary.Read(r, binary.LittleEndian, &bodyBuffer); err != nil {
		return nil, errors.Wrap(err, "read message body error")
	}
	log.Printf("bodyBuffer: %s\n", string(bodyBuffer))

	var (
		message      Message
		unmarshalErr error
	)

	switch msgType {
	case TypeHandshake:
		var body HandshakeData
		unmarshalErr = json.Unmarshal(bodyBuffer, &body)
		message.Body = body
	case TypeHandshakeReply:
		var body HandshakeReplyData
		unmarshalErr = json.Unmarshal(bodyBuffer, &body)
		message.Body = body
	case TypeHeartbeat:
		var body HeartbeatData
		unmarshalErr = json.Unmarshal(bodyBuffer, &body)
		message.Body = body
	default:
		return nil, errors.Errorf("unexpected message type %s", msgType)
	}

	if unmarshalErr != nil {
		return nil, errors.Wrap(unmarshalErr, "unmarshal body buffer error")
	}

	message.Type = msgType

	return &message, nil
}

func WriteMessage(w io.Writer, msg *Message) error {
	encodeMsg, err := msg.Encode()
	if err != nil {
		return err
	}

	_, err = w.Write(encodeMsg)
	return err
}
