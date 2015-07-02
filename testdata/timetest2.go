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
	days  = flag.Int("days", 2, "Add `days` to reference time")
	hours = flag.Duration("hours", 17, "Add duration to reference time")
	zone  = flag.Int("zone", 4*3600, "Timezone offest in seconds")
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

	refTZ := time.FixedZone("Artificial", *zone)

	p.X.ReferenceTime = time.Date(2000, 1, 1, 0, 0, 0, 0, refTZ)
	p.X.Tick.Marker = plot.DateTimeTicks{}
	p.Y.ReferenceTime = p.X.ReferenceTime
	p.Y.Tick.Marker = plot.DateTimeTicks{}

	grid := plotter.Grid{
		Vertical:   draw.LineStyle{Color: color.White, Width: 2},
		Horizontal: draw.LineStyle{Color: color.White, Width: 2},
	}
	p.Add(&grid)

	t0 := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
	t1 := t0.Add(time.Duration(*days) * 24 * time.Hour).Add(*hours)

	xy := make(plotter.XYs, 2)
	xy[0].X = p.X.TimeToFloat(t0)
	xy[0].Y = p.Y.TimeToFloat(t0)

	xy[1].X = p.X.TimeToFloat(t1)
	xy[1].Y = p.Y.TimeToFloat(t1)
	b, err := plotter.NewScatter(xy)
	if err != nil {
		panic(err)
	}
	p.Add(b)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "timeaxis.png"); err != nil {
		panic(err)
	}
}
