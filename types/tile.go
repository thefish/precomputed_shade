package types

type Tile struct {
	Char        string
	Name        string
	BlocksPass  bool
	BlocksSight bool
	Visible     bool
}

func NewWall() *Tile {
	return &Tile{
		"#",
		"Wall",
		true,
		true,
		false,
	}
}

func NewFloor() *Tile {
	return &Tile{
		".",
		"Floor",
		false,
		false,
		false,
	}
}
