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

func (game Game) Run() {
	game.MainLoop()
}

func (game Game) MainLoop() {
	game.ui.Paint()
mainLoop:
	for {
		if stop := game.ui.Tick(); stop{
			break mainLoop
		}
	}
}

func (game Game) Close() {
	game.ui.Close()
}
