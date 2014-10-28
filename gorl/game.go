package gorl

const (
	DefaultDungeonWidth = 100
	DefaultDungeonHeight = 100
)

type Game struct {
	ui *UI
	dungeons []*Dungeon
	currentDungeon *Dungeon
}

func NewGame() (*Game, error) {
	ui, err := NewUI()
	if err != nil {
		return nil, err
	}
	dungeons := make([]*Dungeon, 10)
	dungeons[0] = NewDungeon(DefaultDungeonWidth, DefaultDungeonHeight)
	game := &Game{ui, dungeons, nil}
	game.SetDungeon(dungeons[0])
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

func (game *Game) SetDungeon(d *Dungeon) {
	game.currentDungeon = d
	game.ui.PointCameraAt(d, d.Origin)
}
