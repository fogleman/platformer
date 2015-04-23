package gfx

import (
	"encoding/csv"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"
)

func ParseInts(items []string) []int {
	result := make([]int, len(items))
	for i, item := range items {
		f, _ := strconv.ParseInt(item, 0, 0)
		result[i] = int(f)
	}
	return result
}

func LoadCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return csv.NewReader(file).ReadAll()
}

func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}

func ImageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}
