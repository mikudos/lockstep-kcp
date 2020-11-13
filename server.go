package lockstep

import (
	"log"

	"github.com/xtaci/kcp-go/v5"
)

type (
	// IServer IServer
	IServer interface {
		RegisterScene(string, IScene)
		Start(func() bool)
	}

	// Option Option
	Option struct {
		Addr         string
		Block        kcp.BlockCrypt
		DataShards   int
		ParityShards int
	}
	// Server Server
	Server struct {
		option   Option
		sceneMap map[string]IScene
	}
)

// New generate lockstep server
func New(option *Option) IServer {
	if option == nil {
		option = &Option{"0.0.0.0:8090", nil, 0, 0}
	}
	return &Server{
		option:   *option,
		sceneMap: make(map[string]IScene),
	}
}

// RegisterScene RegisterScene
func (s *Server) RegisterScene(key string, scene IScene) {
	s.sceneMap[key] = scene
	scene.Run()
}

// Start Start
func (s *Server) Start(fn func() bool) {
	if listener, err := kcp.ListenWithOptions(s.option.Addr, s.option.Block, s.option.DataShards, s.option.ParityShards); err == nil {
		// spin-up the client
		if ok := fn(); !ok {
			log.Fatal("condition function return false")
		}
		go s.Frame()
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handleEcho(s)
		}
	} else {
		log.Fatal(err)
	}
}

// Frame Frame
func (s *Server) Frame() {
	for _, s := range s.sceneMap {
		s.Frame()
	}
}

// handleEcho send back everything it received
func handleEcho(conn *kcp.UDPSession) {
	conn.SetNoDelay(1, 10, 2, 1)
	conn.SetStreamMode(true)
	conn.SetWindowSize(4096, 4096)
	conn.SetReadBuffer(4 * 1024 * 1024)
	conn.SetWriteBuffer(4 * 1024 * 1024)
	conn.SetACKNoDelay(true)
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
