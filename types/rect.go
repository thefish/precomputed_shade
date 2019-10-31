package types

type Rect struct {
	X, Y, W, H int
}

func NewRect(x, y, w, h int) *Rect {
	return &Rect{x, y, w, h}
}

func (self *Rect) Intersects(other *Rect) bool {
	if self.X <= (other.X+other.W) &&
		(self.X+self.W) >= other.X &&
		self.Y <= (other.Y+other.Y) &&
		(self.Y+self.H) >= other.Y {
		return true
	}
	return false
}

func (r *Rect) InBounds(c Coords) bool {
	return c.X >= r.X && c.X <= (r.X+r.W-1) && c.Y >= r.Y && c.Y <= (r.Y+r.H-1)
}
