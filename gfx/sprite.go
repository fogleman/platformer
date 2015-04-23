package gfx

type Sprite struct {
	Sheet            *Sheet
	Name             string
	X, Y             float64
	AnchorX, AnchorY float64
	ScaleX, ScaleY   float64
	FlipX, FlipY     bool
	Rotation         float64
}

func NewSprite(sheet *Sheet, name string) *Sprite {
	s := Sprite{}
	s.Sheet = sheet
	s.Name = name
	s.ScaleX = 1
	s.ScaleY = 1
	return &s
}

func (s *Sprite) SetPosition(x, y float64) {
	s.X = x
	s.Y = y
}

func (s *Sprite) SetAnchor(x, y float64) {
	s.AnchorX = x
	s.AnchorY = y
}

func (s *Sprite) SetScale(x, y float64) {
	s.ScaleX = x
	s.ScaleY = y
}

func (s *Sprite) Matrix() Matrix {
	item := s.Sheet.Items[s.Name]
	ax, ay := s.AnchorX-0.5, s.AnchorY-0.5
	sx, sy := s.ScaleX, s.ScaleY
	if s.FlipX {
		sx *= -1
	}
	if s.FlipY {
		sy *= -1
	}
	m := Identity()
	m = m.Scale(Vector{sx * float64(item.W) / 2, sy * float64(item.H) / 2, 1})
	m = m.Translate(Vector{-ax * float64(item.W), -ay * float64(item.H), 0})
	m = m.Rotate(Vector{0, 0, 1}, s.Rotation)
	m = m.Translate(Vector{s.X, s.Y, 0})
	return m
}

func (s *Sprite) Tile() Tile {
	return s.Sheet.TransformedTile(s.Name, s.Matrix())
}
