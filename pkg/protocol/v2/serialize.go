package v2

import (
	"encoding/binary"
	"io"

	"github.com/3vilive/expit/pkg/pb"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func ReadMessage(r io.Reader) (*pb.Message, error) {
	var bodySize uint16
	if err := binary.Read(r, binary.LittleEndian, &bodySize); err != nil {
		return nil, errors.Wrap(err, "read message body size error")
	}

	var bodyBuffer = make([]byte, bodySize)
	if err := binary.Read(r, binary.LittleEndian, &bodyBuffer); err != nil {
		return nil, errors.Wrap(err, "read message body error")
	}

	msg := &pb.Message{}
	if err := proto.Unmarshal(bodyBuffer, msg); err != nil {
		return nil, errors.Wrap(err, "unmarshal message body error")
	}

	return msg, nil
}

func WriteMessage(w io.Writer, msg *pb.Message) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "marshal message body error")
	}

	if err := binary.Write(w, binary.LittleEndian, uint16(len(b))); err != nil {
		return errors.Wrap(err, "write message body size error")
	}
	if err := binary.Write(w, binary.LittleEndian, b); err != nil {
		return errors.Wrap(err, "write message body error")
	}

	return nil
}
