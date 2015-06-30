// +build ignore

package main

import (
	"image/color"
	"math/rand"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

func main() {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Constraint Autoscaling"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.BackgroundColor = color.Gray16{0xdddd}

	p.X.Constraint.Min.Higher = -35
	p.X.Constraint.Min.Lower = -25

	p.Y.Constraint.Min.Lower = 3.8
	p.Y.Constraint.Max.Higher = 4.1

	grid := plotter.Grid{
		Vertical:   draw.LineStyle{Color: color.White, Width: 2},
		Horizontal: draw.LineStyle{Color: color.White, Width: 2},
	}
	p.Add(&grid)

	xy := make(plotter.XYs, 25)
	for i := range xy {
		xy[i].X = rand.Float64()*100 - 50
		xy[i].Y = rand.Float64()/10 + 4
	}
	b, err := plotter.NewScatter(xy)
	if err != nil {
		panic(err)
	}
	p.Add(b)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "constrained-as.png"); err != nil {
		panic(err)
	}
}
