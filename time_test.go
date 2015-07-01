// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"testing"
	"time"
)

func TestRound0(t *testing.T) {
	for i, tc := range []struct{ x, n, w int }{
		{0, 5, 0},
		{1, 5, 0},
		{4, 5, 0},
		{5, 5, 5},
		{9, 5, 5},
		{10, 5, 10},
		{11, 5, 10},
		{14, 5, 10},
		{15, 5, 15},
		{17, 5, 15},
	} {
		if got := round0(tc.x, tc.n); got != tc.w {
			t.Errorf("%d: round0(%d,%d)=%d want %d", i, tc.x, tc.n, got, tc.w)
		}
	}

}

func TestRound1(t *testing.T) {
	for i, tc := range []struct{ x, n, w int }{
		{1, 3, 1},
		{2, 3, 1},
		{3, 3, 1},
		{4, 3, 4},
		{5, 3, 4},
		{6, 3, 4},
		{7, 3, 7},
		{8, 3, 7},
		{9, 3, 7},
		{10, 3, 10},
		{11, 3, 10},
		{12, 3, 10},
		{1, 6, 1},
		{6, 6, 1},
		{7, 6, 7},
		{12, 6, 7},
	} {
		if got := round1(tc.x, tc.n); got != tc.w {
			t.Errorf("%d: round1(%d,%d)=%d want %d", i, tc.x, tc.n, got, tc.w)
		}
	}
}

func TestSuitableDelta(t *testing.T) {
	for i, tc := range []struct {
		rng  int64
		max  int
		cnt  int
		unit interval
	}{
		{35, 3, 20, second},
		{35, 4, 10, second},
		{35, 7, 5, second},
		{6*3600 + 120, 3, 3, hour},
		{2*24*3600 - 1000, 4, 12, hour},
		{2*24*3600 + 1000, 4, 1, day},
		{2*24*3600 + 1000, 5, 12, hour},
		{180 * 24 * 3600, 5, 2, month},
		{180 * 24 * 3600, 6, 1, month},
	} {
		td, n, _ := suitableDelta(tc.rng, tc.max)
		if td.count != tc.cnt || td.unit != tc.unit || n > tc.max {
			t.Errorf("%d: suitableDelta(%d,%d)=%d*%s n=%d want %d*%s",
				i, tc.rng, tc.max, td.count, td.unit, n, tc.cnt, tc.unit)
		}

	}
}

func TestRoundDown(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	t0, err := time.Parse(layout, "2009-12-28 08:42:36")
	if err != nil {
		t.Fatal(err)
	}

	for i, tc := range []struct {
		td   timeDelta
		want string
	}{
		{timeDelta{count: 2, unit: second}, "2009-12-28 08:42:36"},
		{timeDelta{count: 5, unit: second}, "2009-12-28 08:42:35"},
		{timeDelta{count: 10, unit: second}, "2009-12-28 08:42:30"},
		{timeDelta{count: 20, unit: second}, "2009-12-28 08:42:20"},
		{timeDelta{count: 30, unit: second}, "2009-12-28 08:42:30"},
		{timeDelta{count: 1, unit: minute}, "2009-12-28 08:42:00"},
		{timeDelta{count: 2, unit: minute}, "2009-12-28 08:42:00"},
		{timeDelta{count: 5, unit: minute}, "2009-12-28 08:40:00"},
		{timeDelta{count: 10, unit: minute}, "2009-12-28 08:40:00"},
		{timeDelta{count: 20, unit: minute}, "2009-12-28 08:40:00"},
		{timeDelta{count: 30, unit: minute}, "2009-12-28 08:30:00"},
		{timeDelta{count: 1, unit: hour}, "2009-12-28 08:00:00"},
		{timeDelta{count: 2, unit: hour}, "2009-12-28 08:00:00"},
		{timeDelta{count: 3, unit: hour}, "2009-12-28 06:00:00"},
		{timeDelta{count: 4, unit: hour}, "2009-12-28 08:00:00"},
		{timeDelta{count: 6, unit: hour}, "2009-12-28 06:00:00"},
		{timeDelta{count: 12, unit: hour}, "2009-12-28 00:00:00"},
		{timeDelta{count: 1, unit: day}, "2009-12-28 00:00:00"},
		{timeDelta{count: 2, unit: day}, "2009-12-27 00:00:00"},
		{timeDelta{count: 5, unit: day}, "2009-12-26 00:00:00"},
		{timeDelta{count: 10, unit: day}, "2009-12-21 00:00:00"},
		{timeDelta{count: 15, unit: day}, "2009-12-16 00:00:00"},
		{timeDelta{count: 1, unit: month}, "2009-12-01 00:00:00"},
		{timeDelta{count: 2, unit: month}, "2009-11-01 00:00:00"},
		{timeDelta{count: 3, unit: month}, "2009-10-01 00:00:00"},
		{timeDelta{count: 4, unit: month}, "2009-09-01 00:00:00"},
		{timeDelta{count: 1, unit: year}, "2009-01-01 00:00:00"},
		{timeDelta{count: 2, unit: year}, "2008-01-01 00:00:00"},
		{timeDelta{count: 5, unit: year}, "2005-01-01 00:00:00"},
		{timeDelta{count: 10, unit: year}, "2000-01-01 00:00:00"},
	} {
		got := tc.td.RoundDown(t0)
		want, err := time.Parse(layout, tc.want)
		if err != nil {
			t.Errorf("%d: unexpected error %s", i, err)
			continue
		}
		if want != got {
			t.Errorf("%d: %s to %d %s = %s, want %s", i,
				t0.Format(layout), tc.td.count, tc.td.unit,
				got.Format(layout), tc.want)
		}
	}
}
