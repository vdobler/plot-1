// +build ignore

package main

import (
	"flag"
	"image/color"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

var (
	ticks = flag.Int("ticks", 6, "Allow at most `ticks` major ticks")
)

func main() {

	flag.Parse()
	plot.DateTimeAxisMaxNoTicks = *ticks

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Date/Time-Axis"
	p.X.Label.Text = "Start"
	p.Y.Label.Text = "End"
	p.BackgroundColor = color.Gray16{0xdddd}

	p.X.ReferenceTime = time.Date(2003, 5, 11, 23, 11, 9, 0, time.UTC)
	p.X.Tick.Marker = plot.DateTimeTicks{}
	p.Y.ReferenceTime = p.X.ReferenceTime
	p.Y.Tick.Marker = plot.DateTimeTicks{}

	grid := plotter.Grid{
		Vertical:   draw.LineStyle{Color: color.White, Width: 2},
		Horizontal: draw.LineStyle{Color: color.White, Width: 2},
	}
	p.Add(&grid)

	xy := make(plotter.XYs, 4)
	xy[0].X = p.X.TimeToFloat(time.Date(2003, 5, 12, 8, 45, 17, 0, time.UTC))
	xy[0].Y = p.Y.TimeToFloat(time.Date(2002, 2, 12, 0, 4, 17, 0, time.UTC))

	xy[1].X = p.X.TimeToFloat(time.Date(2003, 5, 12, 16, 0, 39, 0, time.UTC))
	xy[1].Y = p.Y.TimeToFloat(time.Date(2003, 4, 28, 12, 5, 3, 0, time.UTC))

	xy[2].X = p.X.TimeToFloat(time.Date(2003, 5, 11, 23, 11, 9, 0, time.UTC))
	xy[2].Y = p.Y.TimeToFloat(time.Date(2004, 10, 3, 8, 45, 55, 0, time.UTC))

	xy[3].X = p.X.TimeToFloat(time.Date(2003, 5, 12, 19, 54, 42, 0, time.UTC))
	xy[3].Y = p.Y.TimeToFloat(time.Date(2005, 12, 20, 23, 9, 1, 0, time.UTC))

	b, err := plotter.NewScatter(xy)
	if err != nil {
		panic(err)
	}
	p.Add(b)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "datetime.png"); err != nil {
		panic(err)
	}
}
