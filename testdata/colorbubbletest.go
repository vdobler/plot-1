// +build ignore

package main

import (
	"math"

	"github.com/gonum/plot"
	"github.com/gonum/plot/palette"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

func main() {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Color Bubbles"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	grid := plotter.NewGrid()
	p.Add(grid)

	xyzw := make(plotter.XYZWs, 10)
	for i, _ := range xyzw {
		fi := float64(i)
		xyzw[i].X = fi
		xyzw[i].Y = fi * fi / 10
		xyzw[i].Z = math.Log(fi + 1)
		xyzw[i].W = fi
	}

	pal := palette.Rainbow(30, palette.Red, palette.Blue, 1, 1, 1)
	b, err := plotter.NewColorBubbles(xyzw, 1*vg.Millimeter, 7*vg.Millimeter, pal)
	if err != nil {
		panic(err)
	}
	p.Add(b)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "colorbubbles.png"); err != nil {
		panic(err)
	}
}
