// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"errors"
	"image/color"
	"math"

	"github.com/gonum/plot"
	"github.com/gonum/plot/palette"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// ColorBubbles implements the Plotter interface, drawing
// a bubble plot of x, y, z, w quadruples where the z value
// determines the area of the bubble and w its color.
type ColorBubbles struct {
	XYZWs

	// MinRadius and MaxRadius give the minimum
	// and maximum bubble radius respectively.
	// The radii of each bubble is interpolated linearly
	// between these two values.
	MinRadius, MaxRadius vg.Length

	// MinZ and MaxZ are the minimum and
	// maximum Z values from the data.
	MinZ, MaxZ float64

	// Palette is the palette to use as color scale.
	// TODO: Should go to color axis
	Palette palette.Palette

	MinW, MaxW float64
}

// NewColorBubbles creates as new bubble plot plotter for
// the given data, with a minimum and maximum
// bubble radius.
func NewColorBubbles(xyzw XYZWer, min, max vg.Length, pal palette.Palette) (*ColorBubbles, error) {
	cpy, err := CopyXYZWs(xyzw)
	if err != nil {
		return nil, err
	}
	if min > max {
		return nil, errors.New("Min bubble radius is greater than the max radius")
	}
	minz, maxz := cpy[0].Z, cpy[0].Z
	minw, maxw := cpy[0].W, cpy[0].W
	for _, d := range cpy {
		minz, maxz = math.Min(minz, d.Z), math.Max(maxz, d.Z)
		minw, maxw = math.Min(minw, d.W), math.Max(maxw, d.W)
	}
	return &ColorBubbles{
		XYZWs:     cpy,
		Palette:   pal,
		MinRadius: min,
		MaxRadius: max,
		MinZ:      minz,
		MaxZ:      maxz,
		MinW:      minw,
		MaxW:      maxw,
	}, nil
}

// Plot implements the Plot method of the plot.Plotter interface.
func (cb *ColorBubbles) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)

	colors := cb.Palette.Colors()

	for _, d := range cb.XYZWs {
		x := trX(d.X)
		y := trY(d.Y)
		if !c.Contains(draw.Point{x, y}) {
			continue
		}

		rad := cb.radius(d.Z)
		col := cb.color(colors, d.W)

		// Draw a circle with radius rad centered at x, y in color col.
		c.SetColor(col)
		var p vg.Path
		p.Move(x+rad, y)
		p.Arc(x, y, rad, 0, 2*math.Pi)
		p.Close()
		c.Fill(p)
	}
}

// radius returns the radius of a bubble by linear interpolation.
func (cb *ColorBubbles) radius(z float64) vg.Length {
	rng := cb.MaxRadius - cb.MinRadius
	if cb.MaxZ == cb.MinZ {
		return rng/4 + cb.MinRadius
	}
	a := (z - cb.MinZ) / (cb.MaxZ - cb.MinZ)
	return vg.Length(a*a)*rng + cb.MinRadius
}

func (cb *ColorBubbles) color(colors []color.Color, w float64) color.Color {
	rng := len(colors) - 1
	if cb.MaxW == cb.MinW {
		return colors[rng/2]
	}
	d := float64(rng) * (w - cb.MinW) / (cb.MaxW - cb.MinW)
	return colors[int(d)]
}

// DataRange implements the DataRange method
// of the plot.DataRanger interface.
func (cb *ColorBubbles) DataRange() (xmin, xmax, ymin, ymax float64) {
	return XYRange(cb.XYZWs)
}

// GlyphBoxes implements the GlyphBoxes method
// of the plot.GlyphBoxer interface.
func (cb *ColorBubbles) GlyphBoxes(plt *plot.Plot) []plot.GlyphBox {
	boxes := make([]plot.GlyphBox, len(cb.XYZWs))
	for i, d := range cb.XYZWs {
		boxes[i].X = plt.X.Norm(d.X)
		boxes[i].Y = plt.Y.Norm(d.Y)
		r := cb.radius(d.Z)
		boxes[i].Rectangle = draw.Rectangle{
			Min: draw.Point{-r, -r},
			Max: draw.Point{+r, +r},
		}
	}
	return boxes
}
