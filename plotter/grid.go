// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

var (
	// DefaultBackgroundColor is the default background color behind
	// grid lines.
	DefaultGridBackgroundColor = color.Gray{0xee}

	// DefaultGridLineStyle is the default style for major grid lines.
	DefaultGridLineStyle = draw.LineStyle{
		Color: color.Gray{0xff},
		Width: vg.Points(1.5),
	}

	// DefaultMinorGridLineStyle is the default style for mainor grid lines.
	DefaultMinorGridLineStyle = draw.LineStyle{
		Color: color.Gray{0xff},
		Width: vg.Points(0.75),
	}
)

// Grid implements the plot.Plotter interface, drawing
// a set of grid lines at the major and minor tick marks
// after filling the plot area background.
type Grid struct {
	// Background is the background color of the plot area.
	Background color.Color

	// Vertical is the style of the vertical lines at major ticks.
	Vertical draw.LineStyle

	// Horizontal is the style of the horizontal lines at major ticks.
	Horizontal draw.LineStyle

	// MinorVertical is the style of the vertical lines at minor ticks.
	MinorVertical draw.LineStyle

	// MinorHorizontal is the style of the horizontal lines at minor ticks.
	MinorHorizontal draw.LineStyle
}

// NewGrid returns a new grid with both vertical and
// horizontal lines using the default grid line style.
func NewGrid() *Grid {
	return &Grid{
		Background:      DefaultGridBackgroundColor,
		Vertical:        DefaultGridLineStyle,
		Horizontal:      DefaultGridLineStyle,
		MinorVertical:   DefaultMinorGridLineStyle,
		MinorHorizontal: DefaultMinorGridLineStyle,
	}
}

// Plot implements the plot.Plotter interface.
func (g *Grid) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)

	// Draw background below grid lines.
	if g.Background != nil {
		c.SetColor(g.Background)
		c.Fill(c.Rectangle.Path())
	}

	// Minor grid lines first.
	if g.MinorVertical.Color != nil && g.MinorVertical.Width > 0 {
		for _, tk := range plt.X.Tick.Marker.Ticks(plt.X) {
			if !tk.IsMinor() {
				continue
			}
			x := trX(tk.Value)
			c.StrokeLine2(g.MinorVertical, x, c.Min.Y, x, c.Min.Y+c.Size().Y)
		}
	}
	if g.MinorHorizontal.Color != nil && g.MinorHorizontal.Width > 0 {
		for _, tk := range plt.Y.Tick.Marker.Ticks(plt.Y) {
			if !tk.IsMinor() {
				continue
			}
			y := trY(tk.Value)
			c.StrokeLine2(g.MinorHorizontal, c.Min.X, y, c.Min.X+c.Size().X, y)
		}
	}

	// The major grid lines are drawn after (over) the minor ones.
	if g.Vertical.Color != nil && g.Vertical.Width > 0 {
		for _, tk := range plt.X.Tick.Marker.Ticks(plt.X) {
			if tk.IsMinor() {
				continue
			}
			x := trX(tk.Value)
			c.StrokeLine2(g.Vertical, x, c.Min.Y, x, c.Min.Y+c.Size().Y)
		}
	}
	if g.Horizontal.Color != nil && g.Horizontal.Width > 0 {
		for _, tk := range plt.Y.Tick.Marker.Ticks(plt.Y) {
			if tk.IsMinor() {
				continue
			}
			y := trY(tk.Value)
			c.StrokeLine2(g.Horizontal, c.Min.X, y, c.Min.X+c.Size().X, y)
		}
	}
}
