package ethiopiancalendar

import (
	"errors"
	"fmt"
	"strings"
)

type EtDate struct {
	Year  int
	Month int
	Day   int
}

var monthNames = []string{"", "Meskerem", "Tikimt", "Hidar", "Tahsas", "Tir", "Yekatit", "Megabit", "Miazia", "Genbot", "Sene", "Hamle", "Nehase", "Pagume"}

const jdOffset = 1723856

func IsLeap(year int) bool {
	if year < 0 {
		return false
	}

	return (year % 4) == 3
}

func DaysInMonth(year, month int) int {
	if month < 1 || month > 13 {
		return 0
	}

	if month == 13 {
		if IsLeap(year) {
			return 6
		}
		return 5
	}
	return 30
}

func (d EtDate) Validate() error {
	if d.Year <= 0 {
		return errors.New("year must be positive")
	}

	if d.Month < 1 || d.Month > 13 {
		return errors.New("month must be between 1 and 13")
	}

	maxDay := DaysInMonth(d.Year, d.Month)

	if d.Day < 1 || d.Day > maxDay {
		return errors.New("day out of range for month")
	}
	return nil
}

func (d EtDate) ToJDN() (int, error) {
	if err := d.Validate(); err != nil {
		return 0, err
	}

	y := d.Year
	m := d.Month
	day := d.Day

	return jdOffset + 365 + 365*(y-1) + (y / 4) + 30*m + day - 31, nil
}

func JDNToEt(jdn int) (EtDate, error) {
	if jdn < jdOffset {
		return EtDate{}, errors.New("jdn before Ethiopian epoch")
	}

	fixed := jdn - jdOffset
	q := fixed / 1461
	r := fixed % 1461

	year := 4*q + (r / 365) + 1

	n := (r % 365) + 365*(r/1460)

	month := (n / 30) + 1

	day := (n % 30) + 1
	d := EtDate{Year: year, Month: month, Day: day}

	if err := d.Validate(); err != nil {
		return EtDate{}, err
	}

	return d, nil
}

func GregorianToJDN(year, month, day int) (int, error) {
	if year < 0 {
		return 0, errors.New("no year 0 in Gregorian")
	}

	a := (14 - month) / 12

	y := year + 4800 - a

	m := month + 12*a - 3

	jdn := day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045

	return jdn, nil
}

func JDNToGregorian(jdn int) (year, month, day int, err error) {
	a := jdn + 32044
	b := (4*a + 3) / 146097
	c := a - (146097*b)/4
	d := (4*c + 3) / 1461
	e := c - (1461*d)/4
	m := (5*e + 2) / 153
	day = e - (153*m+2)/5 + 1
	month = m + 3 - 12*(m/10)
	year = 100*b + d - 4800 + m/10
	if year <= 0 {
		return 0, 0, 0, errors.New("invalid Gregorian year")
	}

	return year, month, day, nil
}

func (d EtDate) ToGregorian() (int, int, int, error) {
	jdn, err := d.ToJDN()
	if err != nil {
		return 0, 0, 0, err
	}
	return JDNToGregorian(jdn)
}

func FromGregorian(year, month, day int) (EtDate, error) {
	jdn, err := GregorianToJDN(year, month, day)
	if err != nil {
		return EtDate{}, err
	}
	return JDNToEt(jdn)
}

func (d EtDate) Format(layout string) string {
	str := strings.ReplaceAll(layout, "YYYY", fmt.Sprintf("%04d", d.Year))
	str = strings.ReplaceAll(str, "MM", fmt.Sprintf("%02d", d.Month))
	str = strings.ReplaceAll(str, "DD", fmt.Sprintf("%02d", d.Day))
	str = strings.ReplaceAll(str, "Month", monthNames[d.Month])
	return str
}

func (d EtDate) AddDays(days int) (EtDate, error) {
	jdn, err := d.ToJDN()
	if err != nil {
		return EtDate{}, err
	}
	return JDNToEt(jdn + days)
}

func (d EtDate) AddMonths(months int) EtDate {
	y := d.Year + months/13
	m := d.Month + (months % 13)
	if m > 13 {
		m -= 13
		y++
	} else if m < 1 {
		m += 13
		y--
	}

	newDate := EtDate{Year: y, Month: m, Day: d.Day}

	maxDay := DaysInMonth(y, m)
	if newDate.Day > maxDay {
		newDate.Day = maxDay
	}

	return newDate
}

func (d EtDate) AddYears(years int) EtDate {
	newDate := EtDate{Year: d.Year + years, Month: d.Month, Day: d.Day}

	if newDate.Month == 13 {
		maxDay := DaysInMonth(newDate.Year, newDate.Month)
		if newDate.Day > maxDay {
			newDate.Day = maxDay
		}
	}
	return newDate
}
