// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import "time"

// interval is a time interval used during rounding dates/times.
type interval int

const (
	second interval = iota
	minute
	hour
	day
	week
	month
	year
)

func (i interval) String() string {
	return []string{"sec", "min", "hour", "day", "week", "month", "year"}[i]
}

// seconds contains the average amount of seconds for each interval.
var seconds = []int{
	1,
	60,
	3600,
	24 * 3600,
	7 * 24 * 3600,
	30.4375 * 24 * 3600,
	365.25 * 24 * 3600,
}

// timeDelta contains the information for generating ticks.
type timeDelta struct {
	count   int
	unit    interval
	format  string // format for each tick
	first   string // format for first tick
	instant bool
	minors  int
}

func round0(x, n int) int {
	return n * (x / n) // TODO check negative
}

func round1(x, n int) int {
	return n*((x-1)/n) + 1 // TODO check negative
}

// RoundDown returns t rounded down to a full d.
func (d timeDelta) RoundDown(t time.Time) time.Time {
	Y, M, D := t.Date()
	h, m, s := t.Hour(), t.Minute(), t.Second()

	switch d.unit {
	case year:
		Y = round0(Y, d.count)
		M, D, h, m, s = 1, 1, 0, 0, 0
	case month:
		M = time.Month(round1(int(M), d.count))
		D, h, m, s = 1, 0, 0, 0
	case week:
		panic("unimplemented")
	case day:
		D = round1(D, d.count)
		h, m, s = 0, 0, 0
	case hour:
		h = round0(h, d.count)
		m, s = 0, 0
	case minute:
		m = round0(m, d.count)
		s = 0
	case second:
		s = round0(s, d.count)
	default:
		panic("ooops")
	}

	return time.Date(Y, M, D, h, m, s, 0, t.Location())
}

// timeDeltas constins suitable date/time axis tick spacings.
// TODO: think about it much more.
var timeDeltas = []timeDelta{
	{1, second, "15:04:05", "2 Jan 2006", true, 2},
	{2, second, "15:04:05", "2 Jan 2006", true, 2},
	{5, second, "15:04:05", "2 Jan 2006", true, 5},
	{10, second, "15:04:05", "2 Jan 2006", true, 2},
	{20, second, "15:04:05", "2 Jan 2006", true, 4},
	{30, second, "15:04:05", "2 Jan 2006", true, 3},
	{1, minute, "15:04", "2 Jan 2006", true, 4},
	{2, minute, "15:04", "2 Jan 2006", true, 4},
	{5, minute, "15:04", "2 Jan 2006", true, 5},
	{10, minute, "15:04", "2 Jan 2006", true, 2},
	{15, minute, "15:04", "2 Jan 2006", true, 3},
	{20, minute, "15:04", "2 Jan 2006", true, 4},
	{30, minute, "15:04", "2 Jan 2006", true, 3},
	{1, hour, "15h", "2 Jan 2006", true, 4},
	{2, hour, "15h", "2 Jan 2006", true, 4},
	{3, hour, "15h", "2 Jan 2006", true, 3},
	{4, hour, "15h", "2 Jan 2006", true, 4},
	{6, hour, "15h", "2 Jan 2006", true, 3},
	{12, hour, "15:04:00", "2 Jan 2006", true, 3},
	{1, day, "02.01.", "2006", false, 4},
	{1, week, "Week N", "2006", false, 7},
	{1, month, "Jan 2006", "", false, 2},
	{2, month, "Jan 2006", "", false, 2},
	{3, month, "Jan 2006", "", false, 3},
	{6, month, "Jan 2006", "", false, 3},
	{1, year, "2006", "", false, 4},
	{2, year, "2006", "", false, 2},
	{5, year, "2006", "", false, 5},
	{10, year, "2006", "", false, 2},
	{20, year, "2006", "", false, 4},
	{50, year, "2006", "", false, 5},
	{100, year, "2006", "", false, 4},
}

// suitableDelta returns the first timeDelta such that opt*timeDelta > rng
func suitableDelta(rng int64, max int) (timeDelta, int, time.Duration) {
	// TODO: handle edge cases
	for _, td := range timeDeltas {
		full := int64(td.count * seconds[td.unit])
		dur := time.Duration(full) * time.Second
		n := int(rng / int64(full))
		if int64(max)*full >= rng {
			return td, n, dur
		}
	}
	return timeDeltas[len(timeDeltas)-1], 2, 30 // TODO
}

// DateTimeAxisMaxNoTicks controls the maximal number of major ticks
// drawn on a dat/time-axis.
var DateTimeAxisMaxNoTicks int = 5

// Ticks returns Ticks in a specified range
func (tt DateTimeTicks) Ticks(a Axis) (ticks []Tick) {
	rng := int64(a.Max - a.Min)
	delta, _, dur := suitableDelta(rng, DateTimeAxisMaxNoTicks+1)

	minorDur := dur / time.Duration(delta.minors)
	t := a.FloatToTime(a.Min)

	firstMajor := true
	v := a.Min
	for v < a.Max {
		t = delta.RoundDown(t)
		// Add delta.minors-1 minor ticks.
		tm := t.Add(minorDur)
		for i := 1; i < delta.minors; i++ {
			if vm := a.TimeToFloat(tm); vm >= a.Min && vm <= a.Max { // Todo: slag
				ticks = append(ticks, Tick{Value: vm})
			}
			tm = tm.Add(minorDur)
		}

		v = a.TimeToFloat(t)
		if v >= a.Min && v <= a.Max { // Todo: slag
			tick := Tick{
				Value: v,
				Label: t.Format(delta.format),
			}
			if firstMajor && delta.first != "" {
				// TODO: limit okay? other formats?
				tick.Label += "\n" + t.Format(delta.first)
			}
			ticks = append(ticks, tick)
			firstMajor = false
		}

		t = t.Add(dur).Add(dur / 20)
	}

	return ticks
}
