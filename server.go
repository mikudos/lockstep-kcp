package lockstep

import (
	"log"

	"github.com/mikudos/lockstep-kcp/scene"
	"github.com/xtaci/kcp-go/v5"
)

type (
	// IServer IServer
	IServer interface {
		RegisterScene(string, scene.IScene)
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
		sceneMap map[string]scene.IScene
	}
)

// New generate lockstep server
func New(option *Option) IServer {
	if option == nil {
		option = &Option{"0.0.0.0:8090", nil, 0, 0}
	}
	return &Server{
		option:   *option,
		sceneMap: make(map[string]scene.IScene),
	}
}

// RegisterScene RegisterScene
func (s *Server) RegisterScene(key string, scene scene.IScene) {
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
		go s.LoopScene()
		for {
			session, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go s.HandleSession(session)
		}
	} else {
		log.Fatal(err)
	}
}

// LoopScene LoopScene
func (s *Server) LoopScene() {
	for _, s := range s.sceneMap {
		s.Frame()
	}
}
