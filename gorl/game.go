package gorl

type Game struct {
	ui *UI
}

func NewGame() (*Game, error) {
	ui, err := NewUI()
	if err != nil {
		return nil, err
	}
	game := &Game{ui}
	return game, nil
}

func (game *Game) Run() {
	game.MainLoop()
}

func (game *Game) MainLoop() {
	game.ui.Paint()
mainLoop:
	for {
		game.ui.Tick()
		panic(game.ui.State)
		if game.ui.State == StateClosed {
			break mainLoop
		}
	}
}

func (game *Game) Close() {
	game.ui.Close()
}
