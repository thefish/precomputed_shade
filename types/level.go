package types

type Level struct {
	*Rect
	Name   string
	Branch string
	Depth  int
	Tiles  []*Tile
}

func (l *Level) GetTile(coords Coords) *Tile {
	return l.Tiles[coords.Y*l.W+coords.X]
}

func (l *Level) GetTileByXY(x, y int) *Tile {
	return l.Tiles[y*l.W+x]
}

func (l *Level) SetTile(coords Coords, tile *Tile) {
	l.Tiles[coords.Y*l.W+coords.X] = tile
}

func (l *Level) SetTileByXY(x, y int, tile *Tile) {
	l.Tiles[y*l.W+x] = tile
}
