package log

import (
	"fmt"
	"math"
	"time"
)

type unit struct {
	name      string
	threshold float64
}

var units = []unit{
	unit{"ns", 1000},
	unit{"μs", 1000},
	unit{"ms", 1000},
	unit{"s", 60},
	unit{"m", 60},
	unit{"h", 24},
	unit{"d", 7},
	unit{"w", 0},
}

func Convert(d time.Duration) string {
	n := float64(d.Nanoseconds())
	var u unit
	for _, u = range units {
		if u.threshold == 0 || n < u.threshold {
			break
		}
		n /= u.threshold
	}
	return fmt.Sprintf("%.1f%s", math.Floor(n*10)/10, u.name)
}
