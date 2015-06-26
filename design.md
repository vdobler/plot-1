Design document for more flexible axis
======================================


Faceted Plots
-------------

Currently a Plot has a single Add method to add different Plotters
to the plot.
To allow facted plots an additional indirection is needed: Instead of
adding a Plotter to a Plot the Plotters have to be added to a Panel
and the grid of Panels make up the whole Plot.
Something like
   // AddToPanel adds the Plotters ps to the panel (h,v) in a
   // faceted plot.
   func (p *Plot) AddToPanel(h, v int, ps ...Plotter)

   // FacetLabels sets the row and column labels in a faceted plot.
   func (p *Plot) FacetLabels(colLabels, rowLabels []string)


Additional Scale Types
----------------------

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


Autoscaling, Ticks, Visual Extent and Dodging of Glyphs
-------------------------------------------------------

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


Tic Generation
--------------

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


Tics, Grid and Tic Labels
-------------------------

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



