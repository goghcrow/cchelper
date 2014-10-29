package link

import (
	"errors"
	"time"
)

var SessionNotExistError = errors.New("Session does not exist")

// Notice:返回session需要先检查是否nil,然后检查是否关闭
func (server *Server) Session(sid uint64) (*Session, error) {
	if session, ok := server.sessions[sid]; ok {
		return session, nil
	} else {
		return nil, SessionNotExistError
	}
}

func (server *Server) Send(sid uint64, message Message) error {
	if session, ok := server.sessions[sid]; ok {
		if !session.IsClosed() {
			return session.Send(message)
		}
		return SendToClosedError
	}
	return SessionNotExistError
}

func (server *Server) SendPacket(sid uint64, packet []byte) error {
	if session, ok := server.sessions[sid]; ok {
		if !session.IsClosed() {
			return session.SendPacket(packet)
		}
		return SendToClosedError
	}
	return SessionNotExistError

}

func (server *Server) TrySend(sid uint64, message Message, timeout time.Duration) error {
	if session, ok := server.sessions[sid]; ok {
		if !session.IsClosed() {
			return session.TrySend(message, timeout)
		}
		return SendToClosedError
	}
	return SessionNotExistError
}

func (server *Server) TrySendPacket(sid uint64, packet []byte, timeout time.Duration) error {
	if session, ok := server.sessions[sid]; ok {
		if !session.IsClosed() {
			return session.TrySendPacket(packet, timeout)
		}
		return SendToClosedError
	}
	return SessionNotExistError
}
