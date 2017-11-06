package mt

import (
	"time"
	"regexp"
	"log"
	"errors"
	"strings"
	"strconv"
)

type TimeExpParser struct {
	base time.Time
	start *time.Time
	end *time.Time
}

func NewParser(start,end time.Time)*TimeExpParser{
	var p=&TimeExpParser{}
	p.SetAnyBase(start,end)
	return p
}

//time base and unit
const (
	year       = "y"
	year_alias = "a"

	quarter = "q"

	month       = "m"
	month_alias = "b"

	tenday       = "td"
	week         = "w"
	today        = "t"
	hour         = "h"
	minute       = "mi"
	minute_alias = "n"

	day    = "d"
	second = "s"
	cur    = "*"
)

// default start and end
const(
	start="$s"
	end="$e"
)

//symbols
const (
	add = "+"
	sub = "-"
)

const (
	symbol_pattern = "(\\"+add+"|\\"+sub+"){1}"
	number_pattern = "[\\d]*"
)

var (
	reg_sp  = regexp.MustCompile(symbol_pattern)
	reg_num = regexp.MustCompile(number_pattern)
)

func (this *TimeExpParser) SetBase(base time.Time) {
	this.base = base
	this.start=&base
	this.end=&base
}

func (this *TimeExpParser)SetAnyBase(start,end time.Time) {
	this.base = start
	this.start=&start
	this.end=&end
}

//time exp parser
func (this *TimeExpParser) Parse(texp string) (*time.Time, error) {

	texp = strings.ToLower(texp)

	symbols := reg_sp.FindAllString(texp, -1)

	vars := reg_sp.Split(texp, -1)

	varsLen := len(vars)

	if varsLen == 0 {
		log.Println(texp, "no validted vars")
		return nil, errors.New("time exp error")
	}

	startTime := this.getStartTime(vars[0])

	if varsLen == 1 {
		return &startTime, nil
	}

	symbolIndex := 0

	for _, v := range vars[1:] {
		startTime = delSymbol(startTime, symbols[symbolIndex], v)
		symbolIndex++
	}

	return &startTime, nil
}

func delSymbol(start time.Time, symbol, v string) time.Time {

	muStr := reg_num.FindAllString(v, -1)[0]
	unit := v[len(muStr):]

	if muStr == "" {
		muStr = "1"
	}
	mu, err := strconv.Atoi(muStr)

	if err != nil {
		log.Println(v, "mutil not number")
		mu = 1
	}

	if symbol == sub {
		mu = -mu
	}

	switch unit {
	case year:
		start = time.Date(start.Year()+mu, start.Month(), start.Day(), start.Hour(),
			start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
	case quarter:
		m := int(start.Month()) + mu*3
		start = time.Date(start.Year(), time.Month(m), start.Day(), start.Hour(),
			start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
	case month:
		m := int(start.Month()) + mu
		start = time.Date(start.Year(), time.Month(m), start.Day(), start.Hour(),
			start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
	case tenday:
		start = time.Date(start.Year(), start.Month(), start.Day()+mu*10, start.Hour(),
			start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
	case week:
		start = time.Date(start.Year(), start.Month(), start.Day()+mu*7, start.Hour(),
			start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
	case day:
		start = time.Date(start.Year(), start.Month(), start.Day()+mu, start.Hour(),
			start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
	case hour:
		start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour()+mu,
			start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
	case minute:
		start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(),
			start.Minute()+mu, start.Second(), start.Nanosecond(), time.UTC)
	case second:
		start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(),
			start.Minute(), start.Second()+mu, start.Nanosecond(), time.UTC)
	default:

	}

	return start
}

func (this *TimeExpParser) getStartTime(v string) (time.Time) {

	switch v {
	case year, year_alias:
		return time.Date(this.base.Year(), time.January, 1, 0, 0,
			0, 0, time.UTC)
	case quarter:
		return time.Date(this.base.Year(), getQuarterStart(this.base.Month()),
			1, 0, 0, 0, 0, time.UTC)
	case month, month_alias:
		return time.Date(this.base.Year(), this.base.Month(), 1, 0,
			0, 0, 0, time.UTC)
	case tenday:
		return time.Date(this.base.Year(), this.base.Month(), getTenDayStart(this.base.Day()),
			0, 0, 0, 0, time.UTC)
	case week:
		return time.Date(this.base.Year(), this.base.Month(),
			getWeekStart(this.base.Day(), this.base.Weekday()),
			0, 0, 0, 0, time.UTC)
	case today, day:
		return time.Date(this.base.Year(), this.base.Month(), this.base.Day(),
			0, 0, 0, 0, time.UTC)
	case hour:
		return time.Date(this.base.Year(), this.base.Month(), this.base.Day(),
			this.base.Hour(), 0, 0, 0, time.UTC)
	case minute, minute_alias:
		return time.Date(this.base.Year(), this.base.Month(), this.base.Day(),
			this.base.Hour(), this.base.Minute(), 0, 0, time.UTC)
	case second, cur:
		return time.Date(this.base.Year(), this.base.Month(), this.base.Day(),
			this.base.Hour(), this.base.Minute(), this.base.Second(), 0, time.UTC)
	case start:
		return time.Date(this.base.Year(), this.base.Month(), this.base.Day(),
			this.base.Hour(), this.base.Minute(), this.base.Second(), 0, time.UTC)
	case end:
		return time.Date(this.end.Year(), this.end.Month(), this.end.Day(),
			this.end.Hour(), this.end.Minute(), this.end.Second(), 0, time.UTC)
	default:
		return time.Date(this.base.Year(), this.base.Month(), this.base.Day(),
			this.base.Hour(), this.base.Minute(), this.base.Second(), 0, time.UTC)
	}
}

func getWeekStart(day int, weekDay time.Weekday) int {
	return day - int(weekDay-time.Monday)
}

func getTenDayStart(day int) int {
	switch {
	case day >= 1 && day <= 10:
		return 1
	case day >= 11 && day <= 20:
		return 11
	default:
		return 21
	}
}

func getQuarterStart(m time.Month) time.Month {
	switch m {
	case time.January, time.February, time.March:
		return time.January
	case time.April, time.May, time.June:
		return time.April
	case time.July, time.August, time.September:
		return time.July
	case time.October, time.November, time.December:
		return time.October
	default:
		return m
	}
}
