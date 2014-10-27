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
	game.ui.MainLoop()
}

func (game Game) Close() {
	game.ui.Close()
}
