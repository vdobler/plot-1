// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"fmt"
	"image/color"
	"math"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// FacetedPlots collects several individual Plots into one facted plot.
// The title, the x-axis label and the y-axis label are taken form the
// plot in the grid position (0,0). TODO: this is stupid.
type FacetedPlot struct {
	RowLabels, ColLabels []string
	Plots                [][]*Plot
	SameY, SameX         bool
}

// NewFacetedPlot returns a new faceted plot consisting of a grid of
// len(columnLables) x len(rowLabels) individual plots.
func NewFacetedPlot(columnLabels, rowLabels []string) (*FacetedPlot, error) {
	fp := &FacetedPlot{
		Plots:     make([][]*Plot, len(columnLabels)),
		RowLabels: rowLabels,
		ColLabels: columnLabels,
	}
	var err error
	for c := 0; c < len(columnLabels); c++ {
		fp.Plots[c] = make([]*Plot, len(rowLabels))
		for r := 0; r < len(rowLabels); r++ {
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
}

// Draw draws the faceted plot to c.
func (f *FacetedPlot) Draw(canvas draw.Canvas) {
	rows, cols := len(f.RowLabels), len(f.ColLabels)
	f.mergeAxis()

	// Draw title and axis labels determined by plot at (0,0).
	p00 := f.Plots[0][0]
	if p00.Title.Text != "" {
		canvas.FillText(p00.Title.TextStyle, canvas.Center().X, canvas.Max.Y, -0.5, -1, p00.Title.Text)
		canvas.Max.Y -= p00.Title.Height(p00.Title.Text) - p00.Title.Font.Extents().Descent
		canvas.Max.Y -= p00.Title.Padding
	}
	x := canvas.Min.X
	if p00.Y.Label.Text != "" {
		x += p00.Y.Label.Height(p00.Y.Label.Text)
		canvas.Push()
		canvas.Rotate(math.Pi / 2)
		canvas.FillText(p00.Y.Label.TextStyle, canvas.Center().Y, -x, -0.5, 0, p00.Y.Label.Text)
		canvas.Pop()
		x += -p00.Y.Label.Font.Extents().Descent
		p00.Y.Label.Text = ""
	}
	canvas.Min.X += x
	y := canvas.Min.Y
	if p00.X.Label.Text != "" {
		y -= p00.X.Label.Font.Extents().Descent
		canvas.FillText(p00.X.Label.TextStyle, canvas.Center().X, y, -0.5, 0, p00.X.Label.Text)
		y += p00.X.Label.Height(p00.X.Label.Text)
		p00.X.Label.Text = ""
	}
	canvas.Min.Y += y

	// Determine each plot size and train axis.
	ywidths, ymaxwidth := f.yAxisWidths()
	xheights, xmaxheight := f.xAxisHeights()
	rowLabelWidth := 5 * vg.Millimeter  // TODO: make configurable
	colLabelHeight := 5 * vg.Millimeter // TODO: make configurable
	gridSep := 2 * vg.Millimeter        // TODO: make configurable

	fwidth := (canvas.Rectangle.Size().X - ymaxwidth - rowLabelWidth) / vg.Length(cols)
	println("fwidth", fwidth)
	fheight := (canvas.Rectangle.Size().Y - xmaxheight - colLabelHeight) / vg.Length(rows)
	f.trainAxis(draw.NewCanvas(canvas, fwidth, fheight))

	// Draw the axis.
	for c := 0; c < cols; c++ {
		plt := f.Plots[c][0]
		pleft := ymaxwidth + vg.Length(c)*fwidth
		pright := canvas.Size().X - pleft - fwidth + gridSep
		box := canvas.Crop(pleft, xmaxheight-xheights[c], -pright, 0)
		tmp := plt.Y
		plt.HideY()
		ha := horizontalAxis{plt.X}
		ha.draw(box)
		plt.Y = tmp
	}
	for r := 0; r < rows; r++ {
		plt := f.Plots[0][r]
		pbot := xmaxheight + vg.Length(r)*fheight
		ptop := canvas.Size().Y - pbot - fheight + gridSep
		box := canvas.Crop(ymaxwidth-ywidths[r], pbot, 0, -ptop)
		tmp := plt.X
		plt.HideX()
		va := verticalAxis{plt.Y}
		va.draw(box)
		plt.X = tmp
	}
	canvas.Min.X += ymaxwidth
	canvas.Min.Y += xmaxheight

	// Draw column and row labels.
	miny, maxy := canvas.Max.Y-colLabelHeight, canvas.Max.Y
	for c := 0; c < cols; c++ {
		minx := canvas.Min.X + vg.Length(c)*fwidth
		maxx := minx + fwidth - gridSep
		fmt.Println(minx, canvas.Max.Y-colLabelHeight, maxx-canvas.Max.X, 0)
		box := draw.Rectangle{
			Min: draw.Point{X: minx, Y: miny},
			Max: draw.Point{X: maxx, Y: maxy},
		}
		canvas.SetColor(color.Gray16{0xaaaa}) // TODO: make configurable
		canvas.Fill(box.Path())
		canvas.FillText(p00.X.Label.TextStyle, (minx+maxx)/2, (miny+maxy)/2, -0.5, -0.5, f.ColLabels[c])
	}
	minx, maxx := canvas.Max.X-rowLabelWidth, canvas.Max.X
	for r := 0; r < rows; r++ {
		miny := canvas.Min.Y + vg.Length(r)*fheight
		maxy := miny + fheight - gridSep
		box := draw.Rectangle{
			Min: draw.Point{X: minx, Y: miny},
			Max: draw.Point{X: maxx, Y: maxy},
		}
		canvas.SetColor(color.Gray16{0xaaaa}) // TODO: make configurable
		canvas.Fill(box.Path())
		canvas.Push()
		canvas.Rotate(-math.Pi / 2)
		canvas.FillText(p00.Y.Label.TextStyle, -(miny+maxy)/2, (minx+maxx)/2, -0.5, -0.5, f.RowLabels[r])
		canvas.Pop()
	}
	canvas.Max.Y -= colLabelHeight
	canvas.Max.X -= rowLabelWidth

	// Draw the plain plots (with all axis and titles turned off).
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			println(c, r)
			p := f.Plots[c][r]

			minx := vg.Length(c) * fwidth
			miny := vg.Length(r) * fheight
			maxx := canvas.Size().X - minx - fwidth + gridSep
			maxy := canvas.Size().Y - miny - fheight + gridSep
			pc := canvas.Crop(minx, miny, -maxx, -maxy)
			p.BackgroundColor = nil
			p.drawBackground(pc)
			for _, data := range p.plotters {
				data.Plot(pc, p)
			}

		}
	}
}

func (p *Plot) drawBackground(c draw.Canvas) {
	trX, trY := p.Transforms(&c)

	c.SetColor(color.Gray16{0xeeee}) // TODO: make configurable
	c.Fill(c.Rectangle.Path())

	gls := draw.LineStyle{
		Color: color.White,
		Width: 2,
	}

	for _, tk := range p.X.Tick.Marker.Ticks(p.X.Min, p.X.Max) {
		if tk.IsMinor() {
			continue
		}
		x := trX(tk.Value)
		c.StrokeLine2(gls, x, c.Min.Y, x, c.Min.Y+c.Size().Y)
	}

	for _, tk := range p.Y.Tick.Marker.Ticks(p.Y.Min, p.Y.Max) {
		if tk.IsMinor() {
			continue
		}
		y := trY(tk.Value)
		c.StrokeLine2(gls, c.Min.X, y, c.Min.X+c.Size().X, y)
	}

}

// mergeAxis makes sure all plots of a row share a common y-axis and all plots
// of a column share a common x-axis.
// If SameY (SameX) is true than all plots share the same y-axis (x-axis).
// All axis are sanitized during homogenisation.
func (f *FacetedPlot) mergeAxis() {
	rows, cols := len(f.RowLabels), len(f.ColLabels)

	// Y-axis homogenisation.
	for r := 0; r < rows; r++ {
		// Find min and max of data ranges of y-axis of row r.
		f.Plots[0][r].Y.sanitizeRange()
		ymin, ymax := f.Plots[0][r].Y.Min, f.Plots[0][r].Y.Max
		println("Row", r, "Col", 0, "current Y-range ", ymin, ymax)
		for c := 1; c < cols; c++ {
			f.Plots[c][r].Y.sanitizeRange()
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
		f.Plots[c][0].X.sanitizeRange()
		xmin, xmax := f.Plots[c][0].X.Min, f.Plots[c][0].X.Max
		println("Row", 0, "Col", c, "current X-range ", xmin, xmax)
		for r := 1; r < rows; r++ {
			f.Plots[c][r].X.sanitizeRange()
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

// trainAxis expands the axis so that the glyphs fit.
// Each axis of each plot is trained individualy and the axis are
// merged again.
func (f *FacetedPlot) trainAxis(canvas draw.Canvas) {
	rows, cols := len(f.RowLabels), len(f.ColLabels)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			fmt.Printf("Training axis of facet (%d,%d)\n", c, r)
			f.Plots[c][r].trainAxis(canvas)
		}
	}
	f.mergeAxis()
}

func (f *FacetedPlot) yAxisWidths() (width []vg.Length, max vg.Length) {
	width = make([]vg.Length, len(f.Plots[0]))
	for r := range width {
		a := verticalAxis{f.Plots[0][r].Y}
		w := a.size()
		width[r] = w
		if w > max {
			max = w
		}
	}
	return width, max
}

func (f *FacetedPlot) xAxisHeights() (height []vg.Length, max vg.Length) {
	height = make([]vg.Length, len(f.Plots))
	for c := range height {
		a := horizontalAxis{f.Plots[c][0].X}
		h := a.size()
		height[c] = h
		if h > max {
			max = h
		}
	}
	return height, max
}
