package scheduler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		return nextDayRule(now, date, parts)
	case "y":
		return nextYearRule(now, date)
	case "w":
		return nextWeekRule(now, date, parts)
	case "m":
		return nextMonthRule(now, date, parts)
	default:
		return "", errors.New("invalid repeat rule")
	}
}

func afterNow(a, b time.Time) bool {
	return a.After(b)
}

func nextDayRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid d format")
	}

	n, err := strconv.Atoi(parts[1])
	if err != nil || n <= 0 || n > 400 {
		return "", errors.New("invalid day interval")
	}

	for {
		date = date.AddDate(0, 0, n)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(layout), nil
}

func nextYearRule(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(layout), nil
}

func nextWeekRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid w format")
	}

	daysStr := strings.Split(parts[1], ",")

	validDays := make(map[time.Weekday]bool)

	for _, d := range daysStr {
		n, err := strconv.Atoi(d)
		if err != nil || n < 1 || n > 7 {
			return "", errors.New("invalid weekday")
		}

		// 1=Mon → time.Monday
		wd := time.Weekday(n % 7) // 7 → Sunday
		validDays[wd] = true
	}

	for {
		date = date.AddDate(0, 0, 1)

		if validDays[date.Weekday()] && afterNow(date, now) {
			return date.Format(layout), nil
		}
	}
}

func daysInMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func nextMonthRule(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("invalid m format")
	}

	daysPart := parts[1]
	var monthsPart string
	if len(parts) == 3 {
		monthsPart = parts[2]
	}

	// дни
	daysStr := strings.Split(daysPart, ",")
	var days []int

	for _, d := range daysStr {
		n, err := strconv.Atoi(d)
		if err != nil || (n < -2 || n == 0 || n > 31) {
			return "", errors.New("invalid day")
		}
		days = append(days, n)
	}

	// месяцы
	monthsMap := map[int]bool{}
	if monthsPart != "" {
		for _, m := range strings.Split(monthsPart, ",") {
			n, err := strconv.Atoi(m)
			if err != nil || n < 1 || n > 12 {
				return "", errors.New("invalid month")
			}
			monthsMap[n] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)

		// ф по месяцам
		if len(monthsMap) > 0 && !monthsMap[int(date.Month())] {
			continue
		}

		lastDay := daysInMonth(date)

		for _, d := range days {
			var target int

			switch d {
			case -1:
				target = lastDay
			case -2:
				target = lastDay - 1
			default:
				target = d
			}

			if target == date.Day() && afterNow(date, now) {
				return date.Format(layout), nil
			}
		}
	}
}
