package common

import (
	"time"
)

func BeginOfDay(t time.Time) time.Time {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return t
}

func BeginOfNextDay(t time.Time) time.Time {
	t = BeginOfDay(t).Add(24 * time.Hour)
	return t
}

func SameDay(t1, t2 time.Time) bool {
	d1, m1, y1 := t1.Date()
	d2, m2, y2 := t2.Date()
	return d1 == d2 && m1 == m2 && y1 == y2
}

func SameWeek(t1, t2 time.Time) bool {
	w1, y1 := t1.ISOWeek()
	w2, y2 := t2.ISOWeek()
	return w1 == w2 && y1 == y2
}

func SameMonth(t1, t2 time.Time) bool {
	_, m1, y1 := t1.Date()
	_, m2, y2 := t2.Date()
	return m1 == m2 && y1 == y2
}

func GetTimeEndShift(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 19, 0, 0, 0, time.Local)
}

func WeekStart(year, week int) time.Time {
	// start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}
