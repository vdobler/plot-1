// +build ignore

package main

import (
	"math/rand"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

func main() {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Constraint Autoscaling"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	p.X.Constraint.Min.Higher = -35
	p.X.Constraint.Min.Lower = -25

	p.X.Constraint.Max.Higher = 42.5
	p.X.Constraint.Max.Lower = 42.5

	p.Y.Constraint.Min.Lower = 3.8
	p.Y.Constraint.Max.Higher = 4.1

	grid := plotter.NewGrid()
	p.Add(grid)

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
