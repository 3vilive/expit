package v1

const (
	MessageVersion = 1
)

type MessageType uint16

const (
	TypeHandshake      MessageType = iota
	TypeHandshakeReply MessageType = iota
	TypeHeartbeat      MessageType = iota
)

func (t MessageType) String() string {
	switch t {
	case TypeHandshake:
		return "Handshake"
	case TypeHandshakeReply:
		return "HandshakeReply"
	case TypeHeartbeat:
		return "Heartbeat"
	default:
		return "Need Implement String()"
	}
}

const (
	TunnelTypeControl = "control"
	TunnelTypeData    = "data"
)

type Message struct {
	Type MessageType
	Body interface{}
}

func NewMessage(t MessageType, body interface{}) *Message {
	return &Message{
		Type: t,
		Body: body,
	}
}

type HandshakeData struct {
	OS         string `json:"os"`          // client os
	Version    string `json:"version"`     // client version
	TunnelType string `json:"tunnel_type"` // tunnel type (control, data)
}

type HandshakeReplyData struct {
	Error string `json:"error"`
}

type HeartbeatData struct {
	
}