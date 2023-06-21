package v2

import "github.com/3vilive/expit/pkg/pb"

func NewHandshakeMsg(os, version, clientId string, tunnelType pb.Handshake_TunnelType) *pb.Message {
	return &pb.Message{
		Type: pb.Message_HANDSHAKE,
		Data: &pb.Message_Handshake{
			Handshake: &pb.Handshake{
				Os:         os,
				Version:    version,
				TunnelType: tunnelType,
				ClientId:   clientId,
			},
		},
	}
}

func NewHandshakeReplyMsg(err, dataServerAddr, clientId string) *pb.Message {
	return &pb.Message{
		Type: pb.Message_HANDSHAKE_REPLY,
		Data: &pb.Message_HandshakeReply{
			HandshakeReply: &pb.HandshakeReply{
				Error:          err,
				DataServerAddr: dataServerAddr,
				ClientId:       clientId,
			},
		},
	}
}

func NewHeartbeatMsg() *pb.Message {
	return &pb.Message{
		Type: pb.Message_HEARTBEAT,
		Data: &pb.Message_Heartbeat{
			Heartbeat: &pb.Heartbeat{},
		},
	}
}

func NewRequestDataTunnelMsg(n int) *pb.Message {
	return &pb.Message{
		Type: pb.Message_REQ_DATA_TUNNEL,
		Data: &pb.Message_ReqDataTunnel{
			ReqDataTunnel: &pb.RequestDataTunnel{
				Number: int32(n),
			},
		},
	}
}

func NewStartForwardMsg(client string) *pb.Message {
	return &pb.Message{
		Type: pb.Message_START_FORWARD,
		Data: &pb.Message_StartForward{
			StartForward: &pb.StartForward{
				ClientRemoteAddr: client,
			},
		},
	}
}

func NewStartForwardMsgReply(err string) *pb.Message {
	return &pb.Message{
		Type: pb.Message_START_FORWARD_REPLY,
		Data: &pb.Message_StartForwardReply{
			StartForwardReply: &pb.StartForwardReply{
				Error: err,
			},
		},
	}
}
