// +build ignore

package main

import (
	"flag"
	"math/rand"
	"os"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgimg"
)

var (
	sameX = flag.Bool("same.x", false, "use same x-range fo all facets")
	sameY = flag.Bool("same.y", false, "use same y-range fo all facets")
)

func main() {
	flag.Parse()

	fp, err := plot.NewFacetedPlot([]string{"AB", "CD", "EF"}, []string{"Gamma", "Delta"})
	if err != nil {
		panic(err)
	}
	fp.SameX = *sameX
	fp.SameY = *sameY
	fp.Plots[0][0].Title.Text = "Example Faceted Plot"
	fp.Plots[0][0].X.Label.Text = "X-Axis Label"
	fp.Plots[0][0].Y.Label.Text = "Y-Axis Label"

	pltr, _ := plotter.NewScatter(randomPoints(10))
	fp.Add(0, 0, pltr)
	pltr, _ = plotter.NewScatter(randomPoints(15))
	fp.Add(1, 0, pltr)
	pltr, _ = plotter.NewScatter(randomPoints(20))
	fp.Add(2, 0, pltr)

	pltr, _ = plotter.NewScatter(randomPoints(25))
	fp.Add(0, 1, pltr)
	pltr, _ = plotter.NewScatter(randomPoints(15))
	fp.Add(1, 1, pltr)
	pltr, _ = plotter.NewScatter(randomPoints(5))
	fp.Add(2, 1, pltr)

	pngcanvas := vgimg.PngCanvas{Canvas: vgimg.New(8*vg.Inch, 6*vg.Inch)}
	fp.Draw(draw.New(pngcanvas))
	file, err := os.Create("faceted.png")
	if err != nil {
		panic(err)
	}
	_, err = pngcanvas.WriteTo(file)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}

}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}

// randomTriples returns some random x, y, z triples
// with some interesting kind of trend.
func randomTriples(n int) plotter.XYZs {
	data := make(plotter.XYZs, n)
	for i := range data {
		if i == 0 {
			data[i].X = rand.Float64()
		} else {
			data[i].X = data[i-1].X + 2*rand.Float64()
		}
		data[i].Y = data[i].X + 10*rand.Float64()
		data[i].Z = data[i].X
	}
	return data
}
