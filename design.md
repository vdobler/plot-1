Proof of Concept for a Redisign of Package Plot
===============================================

Most of the examples (which are single file package main commands)
have command line parameters which allow to play with the feature
to demonstrate.


Faceted Plots
-------------

Faceted plots are produced by a new type `FacetedPlot` which takes care
of layouting the different plots, homogenising and drawing axis and the
row/column lables.

Facetting of the data is not done by package plot: Each facet has to be
drawn individually through the `Add` method of `FacetedPlot`

TODO:
 - Single row (or single column) faceted plots have a redundant column
   or row) label which should be suppressed.
 - Make layout and colors configurable.
 - Manualy setting x- and y-axis range cannot be done (only autorangeing).
 - Maybe: Completely free axis (each plots has x and y).
 - Error handling and edge cases (e.g. missing facets).
 - The grey background with the white gridlines is fixed (not customizable)
   as the gridlines are needed as there are factes without any axis.
   This is not flexible enough.

Change:
 - New type `FacetedPlot` with some methods.

Example: plot/testdata/factedtest.go


Full Autorangeing of X and Y Axis
----------------------------------

The old version did pad the drawing area with white space to allow drawing
the glyphs unclipped. This is a sensible and easy to implement but has some
drawbacks:
 - The axis do not range over the whole drawng area which I dislike and
   prevents e.g. a grey background (like ggplot2) or an axis-box
   (like gnuplot).
 - Homogenising the axis across rows or columns in a faceted plot is not
   possible (at least I didn't see how to do it).

The redisgn really increases the Min and Max of the axis albait in a dead
ugly way (expanding iteratively until convergence). But I think this could be
done in one step with a bit more math.

The API and Semantic for DataRanger and GlyphBoxer is unchanged.

Changes:
 - The type `Normalizer` needs to provide the inverse of the Normalize method
   to allow the calculation of the needed axis expansion.
 - The padX and padY functions are no longer needed and have been deleted.

Example: plot/testdate/bubbletest.go


Date/Time Axis
--------------

The packages plot and plotter work with float64 so date/times have to
converted to float64. This is done by specifying a reference time
an interpreting the float64s as seconds from this reference time.
This allows a resolution of nanoseconds if a whole day has to be
covered and a resolution of microseconds for a whole cenutury which should
be okay for plotting, even if the reference time is choosen far from
the actual data range.

The reference times time zone is used to determine the time zone the
ticks are displayed in.

Changes:
 - New field `ReferenceTime` in Axis.
 - New methods `TimeToFloat` and `FloatToTime` for Axis.
 - API change in `Ticker` interface: Use `Ticks(a Axis) []Tick` instead of the
   old `Ticks(min, max float64) []Tick` to allow the tick method access the
   reference time. 
 - New type `DateTimeTicker` which implements the Ticker interface and produces
   sensible date time ticks

TODO:
 - Maybe rename to TimeToFloat64 and Float64ToTime.
 - Handle edge cases (e.g. ranges less than 1 second, or more than several
   hundred of years).
 - The `DateTimeTicker` is not internationalized, it produces english month
   names with german date layouts which is a bit strange.  Maybe export
   `timeDeltas` or make the slice of timeDeltas a parameter to `DateTimeTicker`
   so that users may provide their own set of deltas and date layouts (still
   would produce englisch month names).

Examples: plot/testdata/timetest.go  and timetest2.go


Constrained Autoscaling
-----------------------
 
Constraining the range an end of an axis may be autosclaed to is done by
a new field `Constraint` in Axis.

Changes:
 - New field (with subfields) `Constraint` in type Axis.

TOOD:
 - Setting up an Axis manually (i.e. without `makeAxis`) became much more
   complicated. Maybe make Constraint a pointer with nil meaning "unconstrained
   autoscaling". 

Example: plot/testdata/constrained-as.go


Color Bubble Plotter
--------------------

A new plotter to create bubble plot with bubble area nd bubble color controled
by two individual values.  Currently just a toy but implemented to test the
planed color and size scales (axis).

Example: plot/testdata/colorbubbletest.go


Next Steps / Ideas
------------------

### Additional scales for color, gylph size, glyph type and line style

Let's assume the follwing scenario: The user creates a faceted plot using
ColorBubbles to display quadruples <x,y,z,w> with x and y maped to the normal
x and y axis, using z for the bubble radius and displaying w as the color
of the bubbles.
In this scenario all facets should use the same (maximal) color range and
the same palette. Also all facets should use the same z to radius
transformation.
And most probably the user would like to draw a color scale with tics so
he can read of w-values from color. Most probably he would want a discrete
(or nominal) "bubble size scale" where bubble sizes are ticked and labeled.

I cannot see how this is doable with a Thumbnailer.

To solve this:
 - Add new Axis to each Plot: Color, Size, Line and Glyph with Color and Size
   continuous axis and Line and Glyp categorical / nominal axis.
 - Expand DataRanger to provide not only x and y (min,max)-pairs but (min,max)
   for all ranges used in this plot:  Scatter would return only x and y, while
   Bubble returns x, y and size, ColorBubbles return x, y, size and color
   and HeatMap returns x, y and color.
 - Autoscale all axis/scales actually used by the different plotters
 - In faceted plots: Homogenize Color, Size, Line and Glyph over all factes.
 - Allow plots to draw scales for color, size, glyph and line.

In ggplot2 terms a plotter would be a 'geom'. The x and y axis are already
there and can be used unchanged adding new scales seems straight forward
and works nice in ggplot2.

TODO:
 - Using Scale for other axis/sclaes than x and y might be sensible, e.g.
   to display a logarithmic color scale.
 - How to combine glyph and line scale (these are typically the same)?


### Expansion of axis to next minor or major ticks

This is soemthing I always liked very much in gnuplot. With more exposure to
ggplot2 I do not find it that necessary any more; but why not provide it.

This needs massive changes in tick generation.



Initial ideas, thoughts and remarks
-----------------------------------

### Additional Scale Types

A plot.Plotter may produce more than x and y axis and this should
be reflected in a homogenious way. The following types of scales
should be considered:
 - Normal x and y scales (orthogonal)
 - A color scale
 - A size scale
 - Polar r, ϑ scales
 - Discrete, symbol scales

A plotter might produce/need several of these scales to read of the
values displayed: Think of a bubble plot with variable fill color
which may represent (x,y,r,c) quadruples of floats.

Currently DataRanger just returns [xmin,xmax] and [ymin,ymax].
It should return a complete set of scales covered, e.g.
    map[Scale]struct{min, max float64}

When a plot is rendered all scale data ranges from all Plotters can
be collected and all scales can be ranged and drawn.

In a faceted plot the x and y scales can be ranged individually or free.
A color or a size scale would be global.
Both is easy to implement with the enhanced DataRanger.


### Autoscaling, Ticks, Visual Extent and Dodging of Glyphs

The DataRange may be used to range/train the scales. But this should be
just the basic training of the scales. After determination of which
data range is needed (or wanted) the actual range typically needs to be
expanded for the following reasons:
 - Add space around data (ggplot2 uses 5% of the data range on both sides).
 - Expand range to next full minor or major tick mark.
 - Expand range to accomodate for large glyphs (large bubbles).
 - Expand range to accomodate for dodged glyphes (dodged bars or boxplots).
The first two work in data units, the last two typically work in pixel
units. This suggests that in addition to the Normalizer of a Axis
the inverse is needed too:
 - Normalizer: Convert data units to [0,1]
 - Denormalizer: Convert pixel units to data units
Extreme case: Think of dodged boxplots on a log scale where 3 different
boxes have to be painted at 10, 100 and 1000, each box 10 pixels wide,
boxes seperated by 2 pixels so that the rightmost glyph might be painted
at a data value of 1800 = Denorm(Pixel(Norm(1000))+5+2+10).

               |
           |  ###  |
          ### ###  |
          ### ### ###
           |  ### ###
               |   |
                   |

    -----------+-------
             1000

The current design accomodates for the glyph space by calling GlyphBoxes.
This allows to add appropriate padding but does not allow to expand the
x or y axis. 

This is a threefold problem:
   
    data range [min,max]  <--->  normalized axis coordinates [0,1]  <---->  pixel
                           
                    --- Normalizer -->                       --- simple linear -->

                   <-- DeNormalizer ---                      <-- simple linear ---


How to handle (r,theta) coordinates: A large glyph at (r0,theta0) might
fall well into the drawing area, if not, only the range of the r-scale needs to
be increasd on the max-side. Uhoh, that's ugly...

Maybe the following would work:
 1. Increase data range according to autoscale rules:
    Either 5% on both sides or to next tick, or whatever
 2. See if all GlybBoxes fall into drawing area, if not:
 3. Estimate how much larger which range should be and repeat.
As 3 may increase the range this might switch to different tick spacings
and thus to completely different ranges.
Does this process converge? If yes fast enough?
For linear scales this estimation should be perfectly fine, for log
and polar this is just some math. Unfortunately this cannot be done with the
current impl: Only an instance of a Normalizer interface is there

Conclusion: plot.Normalizer must get a second method DeNormalize.


### Tic Generation

To expand the axis range to the next minor or major tic the Ticker must
provide this information, e.g. with an additional method in Ticker
    OutsideTicks(min, max) [4]Tick

One obstacle in tick generation is the available space. The happy case
of a data range for [2,7] is trivial, but assume a data range of [1000,9000]:
Naviely using steps of 1000 and expanding to the next major tick produces
the ticks:
    0  1000 2000 3000 4000 5000 6000 7000 8000 9000 10000
This might be harmless in a full plot, but as ticks in a faceted plot with
each panel 200 pixels wide the ticks would overlap; in such a case
ticks of 0, 5000 and 10000 would be much nicer or even 0, 5k, 10k.
Too long and thus overlapping tick labels is a problem for tick values
much smaler than 1 (e.g. "0.0007"), and date/time ticks (e.g. "Apr 2015")


Whether the tick labels are properly seperated or do overlap is dependent on:
 - Number of ticks
 - Length of each tick (number of characters)
 - Font size used in tick printing

This is not an issue for manualy set ticks (ConstantTicks) but for
automatically generated ones.
If rotated (45° or 90° or anything) tick labels are to be supported this
problem becomes much less dominant, but I dislike rotated tick labels.

Summary: Lots of long labels are a non-issue if the axis is long (lots
of pixels) and the font is small or if the labels are rotated at least 30°.

Maybe a hint like "don't do more than 5 ticks" to DefaultTicks would be enough.


### Tics, Grid and Tic Labels

The current generation of a grid is done by the special plotter.Grid. This
makes it impossible to turn on grids from package plot itself as this would
create an import cycle.

The convention of "unlabled ticks are minor ticks, labled ones are major ticks"
is nice but inflexible: In a faceted plot you cannot have an inner axis with
major ticks without labels.  The labels just occupy precious space in faceted
plots.

Having neither a grid nor ticks makes faceted plots ugly to unusable.
ggplot2 doesn't user inner axis (neither ticks, nor labels) whcih works
fine due to the gray background and the default grid.

Adding a grid during plot creation (outside of package plot) doesn't work
as the plotter.Grid relies on the ticks and a turned of axis has no ticks.

Proposed solution: Decouple axis drawing from tick/label generation:
Instead of a method HideX have a bool HideX per Axis. (Or even a bitset
which would allow hiding major ticks, minor ticks, labels and the axis itself
individualy).



