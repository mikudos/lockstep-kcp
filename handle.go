package lockstep

import (
	"log"

	"github.com/xtaci/kcp-go/v5"
)

type (
	IHandler interface {
		HandleSession(*kcp.UDPSession)
	}
)

func (s *Server) HandleSession(session *kcp.UDPSession) {
	session.SetNoDelay(1, 10, 2, 1)
	session.SetStreamMode(true)
	session.SetWindowSize(4096, 4096)
	session.SetReadBuffer(4 * 1024 * 1024)
	session.SetWriteBuffer(4 * 1024 * 1024)
	session.SetACKNoDelay(true)
	buf := make([]byte, 4096)
	for {
		n, err := session.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		n, err = session.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
