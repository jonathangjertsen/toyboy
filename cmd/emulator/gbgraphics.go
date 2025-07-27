package main

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

var fonts = map[string]font.Face{
	"Basic": basicfont.Face7x13,
}

type Highlight struct {
	BlockX    int
	BlockY    int
	Color     color.RGBA
	Text      string
	TextColor color.RGBA
	Font      font.Face
}

var DefaultGridConfig = ConfigGraphics{
	ShowGrid:       true,
	BlockSize:      8,
	ShowAddress:    false,
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
	cfg ConfigGraphics,
	highlights []Highlight,
) D {
	blockSize := cfg.BlockSize
	gridCols := (w / blockSize) - 1
	gridRows := (h / blockSize) - 1
	gridT := gui.Config.GUI.Graphics.GridThickness

	// Prepare a fast lookup for highlights
	highlightMap := make(map[[2]int]color.RGBA)
	for _, hl := range highlights {
		highlightMap[[2]int{hl.BlockX, hl.BlockY}] = hl.Color
	}

	gridSpacingX := 0
	gridSpacingY := 0
	if cfg.ShowGrid && !gui.Config.GUI.Graphics.Overlay {
		gridSpacingX = gridCols * gridT
		gridSpacingY = gridRows * gridT
	}

	// Compute width for text labels
	face, ok := fonts[gui.Config.GUI.Graphics.Font]
	if !ok {
		face = basicfont.Face7x13
	}
	labelWidth := 0
	if cfg.ShowAddress {
		labelWidth = (font.MeasureString(face, "ffffh ") + fixed.I(1)/2).Round()
	}
	labelHeight := 0
	if cfg.ShowOffsets {
		labelHeight = face.Metrics().Height.Ceil() + 2
	}

	// Adjust final image size to fit address labels
	sw := w*cfg.Scale + gridSpacingX + labelWidth
	sh := h*cfg.Scale + gridSpacingY + labelHeight
	dr := image.Rect(0, 0, sw, sh)
	img := image.NewRGBA(dr)

	// Fill background
	for y := 0; y < sh; y++ {
		for x := 0; x < sw; x++ {
			img.SetRGBA(x+labelWidth, y+labelHeight, gui.Config.GUI.Graphics.FillColor)
		}
	}

	// Draw image pixels
	for y := 0; y < h; y++ {
		gridOffsetY := 0
		if cfg.ShowGrid && !gui.Config.GUI.Graphics.Overlay {
			gridOffsetY = (y / blockSize) * gridT
		}
		for x := 0; x < w; x++ {
			gridOffsetX := 0
			if cfg.ShowGrid && !gui.Config.GUI.Graphics.Overlay {
				gridOffsetX = (x / blockSize) * gridT
			}
			i := y*w + x

			col := fb[i].RGBA()

			// Blend highlight color if block is highlighted
			if hlCol, ok := highlightMap[[2]int{(x / blockSize), (y / blockSize)}]; ok {
				col = blendRGBA(col, hlCol)
			}

			dstX := x*cfg.Scale + gridOffsetX
			dstY := y*cfg.Scale + gridOffsetY

			for dy := 0; dy < cfg.Scale; dy++ {
				for dx := 0; dx < cfg.Scale; dx++ {
					img.SetRGBA(dstX+dx+labelWidth, dstY+dy+labelHeight, col)
				}
			}
		}
	}

	if cfg.ShowGrid {
		// Vertical dashed grid lines
		for gx := 1; gx <= gridCols; gx++ {
			baseX := gx * blockSize * cfg.Scale
			if !gui.Config.GUI.Graphics.Overlay {
				baseX += (gx - 1) * gridT
			}
			for t := 0; t < gridT; t++ {
				x := baseX + t
				if x >= sw {
					continue
				}
				for y := 0; y < sh; y++ {
					useDash := ((y / cfg.Scale) % gui.Config.GUI.Graphics.DashLen) < (gui.Config.GUI.Graphics.DashLen / 2)
					if useDash {
						img.SetRGBA(x+labelWidth, y+labelHeight, gui.Config.GUI.Graphics.GridColor)
					}
				}
			}
		}

		// Horizontal dashed grid lines
		for gy := 1; gy <= gridRows; gy++ {
			baseY := gy * blockSize * cfg.Scale
			if !gui.Config.GUI.Graphics.Overlay {
				baseY += (gy - 1) * gridT
			}
			for t := 0; t < gridT; t++ {
				y := baseY + t
				if y >= sh {
					continue
				}
				for x := 0; x < sw; x++ {
					useDash := ((x / cfg.Scale) % gui.Config.GUI.Graphics.DashLen) < (gui.Config.GUI.Graphics.DashLen / 2)
					if useDash {
						img.SetRGBA(x+labelWidth, y+labelHeight, gui.Config.GUI.Graphics.GridColor)
					}
				}
			}
		}
	}
	if cfg.ShowAddress {
		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.Black),
			Face: face,
		}
		fontHeight := face.Metrics().Height.Ceil()
		margin := 4
		lastY := -9999 // ensure first label is always drawn

		address := cfg.StartAddress
		for row := 0; row < h/blockSize; row++ {
			yPix := row*blockSize*cfg.Scale + row*gridT + fontHeight

			// Avoid overlap with previous label
			if yPix-lastY >= fontHeight+margin {
				var text string
				if cfg.DecimalAddress {
					text = address.Dec()
				} else {
					text = address.Hex()
				}
				drawer.Dot = fixed.Point26_6{
					X: fixed.I(2),
					Y: fixed.I(yPix),
				}
				drawer.DrawString(text)
				lastY = yPix
			}
			if cfg.LineIncrement > 0 {
				address += cfg.LineIncrement * model.Addr(blockSize)
			} else {
				address += cfg.BlockIncrement * model.Addr(w/blockSize)
			}
		}
	}
	if cfg.ShowOffsets {
		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.Black),
			Face: face,
		}
		fontWidth := font.MeasureString(face, "00h").Ceil()
		for col := 0; col < w/blockSize; col++ {
			xPix := col*blockSize*cfg.Scale + col*gridT + labelWidth
			var text string
			offset := model.Addr(col) * cfg.BlockIncrement
			if cfg.DecimalAddress {
				text = fmt.Sprintf("%02d ", offset)
			} else {
				text = fmt.Sprintf("%02Xh", offset)
			}
			textX := xPix + (blockSize*cfg.Scale-fontWidth)/2
			if textX < labelWidth {
				textX = labelWidth
			}
			drawer.Dot = fixed.P(textX, face.Metrics().Ascent.Ceil())
			drawer.DrawString(text)
		}
	}
	// Render text inside highlighted blocks
	for _, hl := range highlights {
		if hl.Text == "" || hl.Font == nil {
			continue
		}

		textX := hl.BlockX*blockSize*cfg.Scale + hl.BlockX*gridT + labelWidth
		textY := hl.BlockY*blockSize*cfg.Scale + hl.BlockY*gridT + labelHeight
		blockW := blockSize * cfg.Scale
		blockH := blockSize * cfg.Scale

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

func tiledata(vram []model.Data8) []model.Color {
	tileData := vram[:0x1800]
	tiles := make([]model.Tile, len(tileData)/16)
	for i := range tiles {
		tiles[i] = model.DecodeTile(tileData[i*16 : (i+1)*16])
	}
	return placetiles(tiles, 24, 16)
}

func tilemap(vram []model.Data8, addr model.Addr, signedAddressing bool) []model.Color {
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
					col := tile[rowInTile][colInTile].Color()
					fb[(tileRow*8+rowInTile)*(8*w)+tileCol*8+colInTile] = col
				}
			}
		}
	}
	return fb
}

func oambuffer(vram []model.Data8, buf model.OAMBuffer) []model.Color {
	tiles := make([]model.Tile, 10)
	for slot := range buf.Level {
		obj := buf.Buffer[slot]
		tileIndex := int(obj.TileIndex)
		tile := model.DecodeTile(vram[16*tileIndex : 16*(tileIndex+1)])
		tiles[slot] = tile
	}
	return placetiles(tiles, 10, 1)
}

func oam(vram []model.Data8, oam []model.Data8) []model.Color {
	tiles := make([]model.Tile, 40)
	for i := 0; i < 40; i++ {
		obj := model.DecodeObject(oam[i*4 : (i+1)*4])
		tileIndex := int(obj.TileIndex)
		tile := model.DecodeTile(vram[16*tileIndex : 16*(tileIndex+1)])
		tiles[i] = tile
	}
	return placetiles(tiles, 10, 4)
}
