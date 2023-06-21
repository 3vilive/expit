package v1

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMessageSerialize(t *testing.T) {
	var message = NewMessage(
		TypeHandshake,
		HandshakeData{
			OS:      "Darwin",
			Version: "0.0.1",
		},
	)

	encodeMsg, err := message.Encode()
	if err != nil {
		t.Errorf("encode message error: %s\n", err.Error())
		return
	}

	parsedMsg, err := ReadMessage(bytes.NewBuffer(encodeMsg))
	if err != nil {
		t.Errorf("parse message erorr: %s\n", err.Error())
		return
	}

	if parsedMsg.Type != message.Type {
		t.Errorf("unexpcted message type, want `%s` but got `%s`\n", message.Type, parsedMsg.Type)
		return
	}

	if !reflect.DeepEqual(parsedMsg.Body, message.Body) {
		t.Errorf("unexpcted message body, want `%v` but got `%v`\n", message.Body, parsedMsg.Body)
		return
	}
}
