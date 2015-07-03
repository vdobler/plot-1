// +build ignore

package main

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

func main() {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Bubbles"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	grid := plotter.NewGrid()
	p.Add(grid)

	xyz := make(plotter.XYZs, 3)
	xyz[0].X, xyz[0].Y, xyz[0].Z = 2, 2, 10
	xyz[1].X, xyz[1].Y, xyz[1].Z = 3, 3, 1
	xyz[2].X, xyz[2].Y, xyz[2].Z = 4, 4, 10
	b, err := plotter.NewBubbles(xyz, 2*vg.Millimeter, 20*vg.Millimeter)
	if err != nil {
		panic(err)
	}
	p.Add(b)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "bubbles.png"); err != nil {
		panic(err)
	}
}
