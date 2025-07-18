package gui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/jonathangjertsen/toyboy/model"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type GridConfig struct {
	Show          bool
	Overlay       bool
	BlockSize     int
	GridColor     color.RGBA // RGBA grid line color
	FillColor     color.RGBA // RGBA fill/background
	DashLen       int
	GridThickness int

	ShowAddress    bool
	AddressFont    font.Face
	StartAddress   uint16
	BlockIncrement uint16
	LineIncrement  uint16
	DecimalAddress bool

	ShowOffsets bool
}

func (gc GridConfig) WithMem(startAddr uint16, blockIncrement uint16) GridConfig {
	gc.StartAddress = startAddr
	gc.BlockIncrement = blockIncrement
	gc.ShowAddress = true
	return gc
}

type Highlight struct {
	BlockX    int
	BlockY    int
	Color     color.RGBA
	Text      string
	TextColor color.RGBA
	Font      font.Face
}

var DefaultGridConfig = GridConfig{
	Show:           true,
	Overlay:        false,
	BlockSize:      8,
	GridColor:      color.RGBA{R: 136, G: 136, B: 136, A: 255},
	FillColor:      color.RGBA{R: 240, G: 240, B: 240, A: 255},
	DashLen:        4,
	GridThickness:  1,
	ShowAddress:    false,
	AddressFont:    basicfont.Face7x13,
	StartAddress:   0x8000,
	BlockIncrement: 16,
}

func blendRGBA(src, over color.RGBA) color.RGBA {
	a := float64(over.A) / 255
	invA := 1 - a
	return color.RGBA{
		R: uint8(float64(src.R)*invA + float64(over.R)*a),
		G: uint8(float64(src.G)*invA + float64(over.G)*a),
		B: uint8(float64(src.B)*invA + float64(over.B)*a),
		A: 255,
	}
}

func (gui *GUI) GBGraphics(
	gtx C,
	w, h int,
	fb []model.Color,
	scale int,
	cfg GridConfig,
	highlights []Highlight,
) D {
	blockSize := cfg.BlockSize
	gridCols := (w / blockSize) - 1
	gridRows := (h / blockSize) - 1
	gridT := cfg.GridThickness

	// Prepare a fast lookup for highlights
	highlightMap := make(map[[2]int]color.RGBA)
	for _, hl := range highlights {
		highlightMap[[2]int{hl.BlockX, hl.BlockY}] = hl.Color
	}

	gridSpacingX := 0
	gridSpacingY := 0
	if cfg.Show && !cfg.Overlay {
		gridSpacingX = gridCols * gridT
		gridSpacingY = gridRows * gridT
	}

	// Compute width for text labels
	labelWidth := 0
	if cfg.ShowAddress {
		labelWidth = (font.MeasureString(cfg.AddressFont, "ffffh ") + fixed.I(1)/2).Round()
	}
	labelHeight := 0
	if cfg.ShowOffsets {
		labelHeight = cfg.AddressFont.Metrics().Height.Ceil() + 2
	}

	// Adjust final image size to fit address labels
	sw := w*scale + gridSpacingX + labelWidth
	sh := h*scale + gridSpacingY + labelHeight
	dr := image.Rect(0, 0, sw, sh)
	img := image.NewRGBA(dr)

	// Fill background
	for y := 0; y < sh; y++ {
		for x := 0; x < sw; x++ {
			img.SetRGBA(x+labelWidth, y+labelHeight, cfg.FillColor)
		}
	}

	// Draw image pixels
	for y := 0; y < h; y++ {
		gridOffsetY := 0
		if cfg.Show && !cfg.Overlay {
			gridOffsetY = (y / blockSize) * gridT
		}
		for x := 0; x < w; x++ {
			gridOffsetX := 0
			if cfg.Show && !cfg.Overlay {
				gridOffsetX = (x / blockSize) * gridT
			}
			i := y*w + x

			col := fb[i].RGBA()

			// Blend highlight color if block is highlighted
			if hlCol, ok := highlightMap[[2]int{(x / blockSize), (y / blockSize)}]; ok {
				col = blendRGBA(col, hlCol)
			}

			dstX := x*scale + gridOffsetX
			dstY := y*scale + gridOffsetY

			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.SetRGBA(dstX+dx+labelWidth, dstY+dy+labelHeight, col)
				}
			}
		}
	}

	if cfg.Show {
		// Vertical dashed grid lines
		for gx := 1; gx <= gridCols; gx++ {
			baseX := gx * blockSize * scale
			if !cfg.Overlay {
				baseX += (gx - 1) * gridT
			}
			for t := 0; t < gridT; t++ {
				x := baseX + t
				if x >= sw {
					continue
				}
				for y := 0; y < sh; y++ {
					useDash := ((y / scale) % cfg.DashLen) < (cfg.DashLen / 2)
					if useDash {
						img.SetRGBA(x+labelWidth, y+labelHeight, cfg.GridColor)
					}
				}
			}
		}

		// Horizontal dashed grid lines
		for gy := 1; gy <= gridRows; gy++ {
			baseY := gy * blockSize * scale
			if !cfg.Overlay {
				baseY += (gy - 1) * gridT
			}
			for t := 0; t < gridT; t++ {
				y := baseY + t
				if y >= sh {
					continue
				}
				for x := 0; x < sw; x++ {
					useDash := ((x / scale) % cfg.DashLen) < (cfg.DashLen / 2)
					if useDash {
						img.SetRGBA(x+labelWidth, y+labelHeight, cfg.GridColor)
					}
				}
			}
		}
	}
	if cfg.ShowAddress {
		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.Black),
			Face: cfg.AddressFont,
		}
		fontHeight := cfg.AddressFont.Metrics().Height.Ceil()
		margin := 4
		lastY := -9999 // ensure first label is always drawn

		address := cfg.StartAddress
		for row := 0; row < h/blockSize; row++ {
			yPix := row*blockSize*scale + row*gridT + fontHeight

			// Avoid overlap with previous label
			if yPix-lastY >= fontHeight+margin {
				var text string
				if cfg.DecimalAddress {
					text = fmt.Sprintf("%04d", address)
				} else {
					text = fmt.Sprintf("%04Xh", address)
				}
				drawer.Dot = fixed.Point26_6{
					X: fixed.I(2),
					Y: fixed.I(yPix),
				}
				drawer.DrawString(text)
				lastY = yPix
			}
			if cfg.LineIncrement > 0 {
				address += cfg.LineIncrement * uint16(blockSize)
			} else {
				address += cfg.BlockIncrement * uint16(w/blockSize)
			}
		}
	}
	if cfg.ShowOffsets {
		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.Black),
			Face: cfg.AddressFont,
		}
		fontWidth := font.MeasureString(cfg.AddressFont, "00h").Ceil()
		for col := 0; col < w/blockSize; col++ {
			xPix := col*blockSize*scale + col*gridT + labelWidth
			var text string
			offset := uint16(col) * cfg.BlockIncrement
			if cfg.DecimalAddress {
				text = fmt.Sprintf("%02d ", offset)
			} else {
				text = fmt.Sprintf("%02Xh", offset)
			}
			textX := xPix + (blockSize*scale-fontWidth)/2
			if textX < labelWidth {
				textX = labelWidth
			}
			drawer.Dot = fixed.P(textX, cfg.AddressFont.Metrics().Ascent.Ceil())
			drawer.DrawString(text)
		}
	}
	// Render text inside highlighted blocks
	for _, hl := range highlights {
		if hl.Text == "" || hl.Font == nil {
			continue
		}

		textX := hl.BlockX*blockSize*scale + hl.BlockX*gridT + labelWidth
		textY := hl.BlockY*blockSize*scale + hl.BlockY*gridT + labelHeight
		blockW := blockSize * scale
		blockH := blockSize * scale

		textWidth := font.MeasureString(hl.Font, hl.Text).Ceil()
		textHeight := hl.Font.Metrics().Height.Ceil()

		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(hl.TextColor),
			Face: hl.Font,
			Dot: fixed.P(
				textX+(blockW-textWidth)/2,
				textY+(blockH-textHeight)/2+hl.Font.Metrics().Ascent.Ceil(),
			),
		}
		drawer.DrawString(hl.Text)
	}
	return widget.Image{
		Src:   paint.NewImageOp(img),
		Scale: gtx.Metric.PxPerDp,
	}.Layout(gtx)
}

func tiledata(vram []uint8) []model.Color {
	tileData := vram[:0x1800]
	tiles := make([]model.Tile, len(tileData)/16)
	for i := range tiles {
		tiles[i] = model.DecodeTile(tileData[i*16 : (i+1)*16])
	}
	return placetiles(tiles, 24, 16)
}

func tilemap(vram []uint8, addr uint16, signedAddressing bool) []model.Color {
	tileMap := vram[addr-model.AddrVRAMBegin : addr-model.AddrVRAMBegin+0x400]
	tiles := make([]model.Tile, len(tileMap))
	for i := range tiles {
		tileID := tileMap[i]
		var offset uint16
		if signedAddressing {
			offset = uint16(int32(0x1000) + 16*int32(int8(tileID)))
		} else {
			offset = 16 * uint16(tileID)
		}
		tile := vram[offset : offset+16]
		tiles[i] = model.DecodeTile(tile)
	}
	return placetiles(tiles, 32, 32)
}

func placetiles(tiles []model.Tile, w, h int) []model.Color {
	fb := make([]model.Color, h*w*8*8)
	for tileRow := range h {
		for tileCol := range w {
			tile := tiles[tileRow*w+tileCol]
			for rowInTile := range 8 {
				for colInTile := range 8 {
					col := tile[rowInTile][colInTile].Color
					fb[(tileRow*8+rowInTile)*(8*w)+tileCol*8+colInTile] = col
				}
			}
		}
	}
	return fb
}

func oambuffer(vram []uint8, buf model.OAMBuffer) []model.Color {
	tiles := make([]model.Tile, 10)
	for slot := range buf.Level {
		obj := buf.Buffer[slot]
		tileIndex := int(obj.TileIndex)
		tile := model.DecodeTile(vram[16*tileIndex : 16*(tileIndex+1)])
		tiles[slot] = tile
	}
	return placetiles(tiles, 10, 1)
}

func oam(vram []uint8, oam []uint8) []model.Color {
	fb := make([]model.Color, 8*8*10*4)
	objects := make([]model.Sprite, 40)
	for i := 0; i < 40; i += 4 {
		objects = append(objects, model.DecodeSprite(oam[i*4:(i+1)*4]))
	}
	for tileRow := range 4 {
		for tileCol := range 10 {
			obj := objects[tileRow*10+tileCol]
			tileIndex := int(obj.TileIndex)
			tile := model.DecodeTile(vram[16*tileIndex : 16*(tileIndex+1)])
			for rowInTile := range 8 {
				for colInTile := range 8 {
					col := tile[rowInTile][colInTile].Color
					fb[(tileRow*8+rowInTile)*(8*10)+tileCol*8+colInTile] = col
				}
			}
		}
	}
	return fb
}
