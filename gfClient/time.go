package gfClient

// Used to parse grafana time specifications. These can take various forms:
//	 * relative: "now", "now-1h", "now-2d", "now-3w", "now-5M", "now-1y"
//   * human friendly boundary:
// 			From:"now/d" -> start of today
//			To:  "now/d" -> end of today
//			To:  "now/w" -> end of the week
//			To:  "now-1d/d" -> end of yesterday
//			When used as boundary, the same string will evaluate to a different time if used in 'From' or 'To'
//	 * absolute unix time: "142321234"
//
// The required behaviour is clearly documented in the unit tests, time_test.go.

type TimeRange struct {
	From string
	To   string
}

func NewTimeRange(from, to string) TimeRange {
	if from == "" {
		from = "now-1h"
	}
	if to == "" {
		to = "now"
	}
	return TimeRange{from, to}
}
