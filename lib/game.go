package gorl

const (
	DefaultDungeonWidth  = 256
	DefaultDungeonHeight = 256
)

type Game struct {
	ui             *UI
	player         *Player
	dungeons       []*Dungeon
	currentDungeon *Dungeon
}

func NewGame() (*Game, error) {
	game := new(Game)

	dungeon := NewDungeon(DefaultDungeonWidth, DefaultDungeonHeight)
	game.dungeons = make([]*Dungeon, 10)
	game.dungeons[0] = dungeon

	game.player = NewPlayer()
	game.player.loc = dungeon.origin
	dungeon.AddMob(game.player)

	ui, err := NewUI()
	if err != nil {
		return nil, err
	}
	game.ui = ui
	ui.game = game
	game.SetDungeon(dungeon)

	return game, nil
}

type Movement struct {
	x, y int
}

func (c Coord) Plus(m Movement) Coord {
	c.x += m.x
	c.y += m.y
	return c
}

func (game *Game) Move(movement Movement) {
	dest := game.player.loc.Plus(movement)
	dest_tile := game.currentDungeon.Tile(dest.x, dest.y)
	if dest_tile.flags & FlagCrossable != 0 {
		game.player.Move(movement)
		game.ui.PointCameraAt(game.currentDungeon, game.player.loc)
	}
}

func (game *Game) Run() {
	game.MainLoop()
}

func (game *Game) MainLoop() {
	game.ui.Paint()
mainLoop:
	for {
		game.ui.Tick()
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
	game.ui.PointCameraAt(d, d.origin)
}
