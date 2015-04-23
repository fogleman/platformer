package gfx

type TileVertex struct {
	X, Y, U, V float32
}

type Tile [6]TileVertex

type SheetItem struct {
	X, Y, W, H int
}

type Sheet struct {
	Texture *Texture
	Items   map[string]SheetItem
}

func NewSheet(texture *Texture, csvPath string) (*Sheet, error) {
	rows, err := LoadCSV(csvPath)
	if err != nil {
		return nil, err
	}
	items := make(map[string]SheetItem)
	for _, row := range rows {
		x := ParseInts(row[1:])
		items[row[0]] = SheetItem{x[0], x[1], x[2], x[3]}
	}
	return &Sheet{texture, items}, nil
}

func NewSheetFromFile(unit int, pngPath, csvPath string) (*Sheet, error) {
	texture, err := NewTextureFromFile(unit, pngPath)
	if err != nil {
		return nil, err
	}
	sheet, err := NewSheet(texture, csvPath)
	if err != nil {
		texture.Delete()
		return nil, err
	}
	return sheet, nil
}

func (sheet *Sheet) Tile(name string, x, y int) Tile {
	item := sheet.Items[name]
	x0 := float32(x)
	y0 := float32(y)
	x1 := float32(x + item.W)
	y1 := float32(y + item.H)
	u0 := float32(item.X) / float32(sheet.Texture.Width)
	v1 := float32(item.Y) / float32(sheet.Texture.Height)
	u1 := float32(item.X+item.W) / float32(sheet.Texture.Width)
	v0 := float32(item.Y+item.H) / float32(sheet.Texture.Height)
	return Tile{
		TileVertex{x0, y0, u0, v0},
		TileVertex{x1, y0, u1, v0},
		TileVertex{x0, y1, u0, v1},
		TileVertex{x1, y0, u1, v0},
		TileVertex{x1, y1, u1, v1},
		TileVertex{x0, y1, u0, v1},
	}
}
