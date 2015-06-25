// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"fmt"
	"math"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// FacetedPlots collects several individual Plots into one facted plot.
type FacetedPlot struct {
	RowLabels, ColLabels []string
	Plots                [][]*Plot
	SameY, SameX         bool
}

// NewFacetedPlot returns a new faceted plot consisting of len(rows) x len(cols)
// individual plots.
func NewFacetedPlot(cols, rows []string) (*FacetedPlot, error) {
	fp := &FacetedPlot{
		Plots:     make([][]*Plot, len(cols)),
		RowLabels: rows,
		ColLabels: cols,
	}
	var err error
	for c := 0; c < len(cols); c++ {
		fp.Plots[c] = make([]*Plot, len(rows))
		for r := 0; r < len(rows); r++ {
			fp.Plots[c][r], err = New()
			if err != nil {
				return fp, err
			}
		}
	}
	return fp, nil
}

// Add the Plotters ps to the facet in row and col.
func (f *FacetedPlot) Add(col, row int, ps ...Plotter) {
	f.Plots[col][row].Add(ps...)
	println("Added", len(ps), "plotters to facet ", col, row)
}

// Draw draws the faceted plot to c.
func (f *FacetedPlot) Draw(canvas draw.Canvas) {
	fmt.Printf("FacetedPlot.Draw to Rect %v\n", canvas.Rectangle)
	f.mergeAxis()

	rows, cols := len(f.RowLabels), len(f.ColLabels)

	// Hack, use Plot.Title as row/col label
	for r := 0; r < rows-1; r++ {
		for c := 0; c < cols-1; c++ {
			f.Plots[c][r].Title.Text = ""
		}
		f.Plots[cols-1][r].Title.Text = f.RowLabels[r]
	}
	for c := 0; c < cols-1; c++ {
		f.Plots[c][rows-1].Title.Text = f.ColLabels[c]
	}
	f.Plots[cols-1][rows-1].Title.Text = f.ColLabels[cols-1] + " / " + f.RowLabels[rows-1]

	// Turn off unused axis.
	for c := 1; c < cols; c++ {
		for r := 1; r < rows; r++ {
			f.Plots[c][r].HideAxes()
		}
		f.Plots[c][0].HideY()
	}
	for r := 1; r < rows; r++ {
		f.Plots[0][r].HideX()
	}

	// Determine individual plot sizes. BUG: first column must be broader
	// to allow y-axis, last column must be broader for facet label, first row
	// must be higher for facet label, last row must be higher for x-axis.
	fwidth := canvas.Rectangle.Size().X / vg.Length(cols)
	fheight := canvas.Rectangle.Size().Y / vg.Length(rows)
	fmt.Printf("Canvas.Rectangle.Size = %v\nPanels of size %.0f x %.0f \n",
		canvas.Rectangle.Size(), fwidth, fheight)
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			minx := canvas.Min.X + vg.Length(c)*fwidth
			miny := canvas.Min.Y + vg.Length(r)*fheight
			maxx := minx + fwidth - 5*vg.Millimeter
			maxy := miny + fheight - 5*vg.Millimeter
			println(minx, " --x-- ", maxx)
			pc := canvas.Crop(minx, miny, maxx-canvas.Max.X, maxy-canvas.Max.Y)
			fmt.Printf("pc = %v\n", pc.Rectangle)
			f.Plots[c][r].Draw(pc)
		}
	}
}

// mergeAxis makes sure all plots of a row share a common y-axis and all plots
// of a column share a common x-axis.
// If SameY (SameX) is true than all plots share the same y-axis (x-axis).
func (f *FacetedPlot) mergeAxis() {
	rows, cols := len(f.RowLabels), len(f.ColLabels)

	// Y-axis homogenisation.
	for r := 0; r < rows; r++ {
		// Find min and max of data ranges of y-axis of row r.
		ymin, ymax := f.Plots[0][r].Y.Min, f.Plots[0][r].Y.Max
		println("Row", r, "Col", 0, "current Y-range ", ymin, ymax)
		for c := 1; c < cols; c++ {
			ymin = math.Min(ymin, f.Plots[c][r].Y.Min)
			ymax = math.Max(ymax, f.Plots[c][r].Y.Max)
			println("Row", r, "Col", c, "current Y-range ", ymin, ymax)
		}
		// Spread min/max to whole row.
		for c := 0; c < cols; c++ {
			f.Plots[c][r].Y.Min, f.Plots[c][r].Y.Max = ymin, ymax
		}
	}
	if f.SameY {
		// Find global y-range.
		ymin, ymax := f.Plots[0][0].Y.Min, f.Plots[0][0].Y.Max
		for r := 1; r < rows; r++ {
			ymin = math.Min(ymin, f.Plots[0][r].Y.Min)
			ymax = math.Max(ymax, f.Plots[0][r].Y.Max)
		}
		// Spread global min/max to all plots.
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				f.Plots[c][r].Y.Min, f.Plots[c][r].Y.Max = ymin, ymax
			}
		}
	}

	// X-axis homogenisation.
	for c := 0; c < cols; c++ {
		// Find min and max of data ranges of x-axis of column c.
		xmin, xmax := f.Plots[c][0].X.Min, f.Plots[c][0].X.Max
		println("Row", 0, "Col", c, "current X-range ", xmin, xmax)
		for r := 1; r < rows; r++ {
			xmin = math.Min(xmin, f.Plots[c][r].X.Min)
			xmax = math.Max(xmax, f.Plots[c][r].X.Max)
			println("Row", r, "Col", c, "current X-range ", xmin, xmax)
		}
		// Spread min/max to whole column.
		for r := 0; r < rows; r++ {
			f.Plots[c][r].X.Min, f.Plots[c][r].X.Max = xmin, xmax
		}
	}
	if f.SameX {
		// Find global x-range.
		xmin, xmax := f.Plots[0][0].X.Min, f.Plots[0][0].X.Max
		for c := 1; c < cols; c++ {
			xmin = math.Min(xmin, f.Plots[c][0].X.Min)
			xmax = math.Max(xmax, f.Plots[c][0].X.Max)
		}
		// Spread global min/max to all plots.
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				f.Plots[c][r].X.Min, f.Plots[c][r].X.Max = xmin, xmax
			}
		}
	}
}
