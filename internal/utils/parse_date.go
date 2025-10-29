package utils

import (
	"fmt"
	"time"
)

func ParseYearMonth(s string) (time.Time, error) {
	var y, m int

	if len(s) == 7 && s[4] == '-' { // YYYY-MM
		if _, err := fmt.Sscanf(s, "%d-%d", &y, &m); err != nil {
			return time.Time{}, err
		}
	} else if len(s) == 7 && s[2] == '-' { // MM-YYYY
		if _, err := fmt.Sscanf(s, "%d-%d", &m, &y); err != nil {
			return time.Time{}, err
		}
	} else {
		return time.Time{}, fmt.Errorf("bad year-month format, use YYYY-MM or MM-YYYY")
	}

	if m < 1 || m > 12 {
		return time.Time{}, fmt.Errorf("month out of range (1..12)")
	}

	return time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC), nil
}

func YmString(t time.Time) string { return t.Format("2006-01") }

func MonthsOverlapInclusive(aStart time.Time, aEnd *time.Time, bStart time.Time, bEnd *time.Time) int {
	start := aStart
	if bStart.After(start) {
		start = bStart
	}

	var end time.Time
	if aEnd != nil && bEnd != nil {
		if aEnd.Before(*bEnd) {
			end = *aEnd
		} else {
			end = *bEnd
		}
	} else if aEnd != nil {
		end = *aEnd
	} else if bEnd != nil {
		end = *bEnd
	} else {
		return 0 // open-open shouldn't happen for query period
	}
	if end.Before(start) {
		return 0
	}
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	return years*12 + months + 1
}
