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
	ShowAddress   bool

	AddressFont    font.Face
	StartAddress   uint16
	BlockIncrement uint16
}

func (gc GridConfig) WithMem(startAddr, blockIncrement uint16) GridConfig {
	gc.StartAddress = startAddr
	gc.BlockIncrement = blockIncrement
	gc.ShowAddress = true
	return gc
}

type Highlight struct {
	BlockX int
	BlockY int
	Color  color.RGBA
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
		labelWidth = (font.MeasureString(cfg.AddressFont, "ffff ") + fixed.I(1)/2).Round()
	}

	// Adjust final image size to fit address labels
	sw := w*scale + gridSpacingX + labelWidth
	sh := h*scale + gridSpacingY
	dr := image.Rect(0, 0, sw, sh)
	img := image.NewRGBA(dr)

	// Fill background
	for y := 0; y < sh; y++ {
		for x := 0; x < sw; x++ {
			img.SetRGBA(x+labelWidth, y, cfg.FillColor)
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
					img.SetRGBA(dstX+dx+labelWidth, dstY+dy, col)
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
						img.SetRGBA(x+labelWidth, y, cfg.GridColor)
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
						img.SetRGBA(x+labelWidth, y, cfg.GridColor)
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
				text := fmt.Sprintf("%04X", address)
				drawer.Dot = fixed.Point26_6{
					X: fixed.I(2),
					Y: fixed.I(yPix),
				}
				drawer.DrawString(text)
				lastY = yPix
			}
			address += cfg.BlockIncrement * uint16(w/blockSize)
		}
	}
	return widget.Image{
		Src:   paint.NewImageOp(img),
		Scale: gtx.Metric.PxPerDp,
	}.Layout(gtx)
}
