package scene

type (
	RunCallback       func(IScene)
	FrameCallback     func(IScene)
	BroadcastCallback func(IScene)
)

type (
	// IScene IScene
	IScene interface {
		Run()
		Frame()
		BroadCast()
	}

	// Scene 游戏场景
	Scene struct {
		runCallback       RunCallback
		frameCallback     FrameCallback
		broadcastCallback BroadcastCallback
	}
)

// New 构造新的场景
func New(rCallback RunCallback, fCallback FrameCallback, bCallback BroadcastCallback) IScene {
	return &Scene{runCallback: rCallback, frameCallback: fCallback, broadcastCallback: bCallback}
}

// Run Run
func (sc *Scene) Run() {
	sc.runCallback(sc)
}

// Frame Frame
func (sc *Scene) Frame() {
	sc.frameCallback(sc)
	sc.BroadCast()
}

// BroadCast BroadCast
func (sc *Scene) BroadCast() {
	sc.broadcastCallback(sc)
}
