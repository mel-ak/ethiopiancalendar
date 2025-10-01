package ethiopiancalendar

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// EtDate represents a date in the Ethiopian Calendar.
type EtDate struct {
	Year  int
	Month int
	Day   int
}

var monthNames = []string{"", "Meskerem", "Tikimt", "Hidar", "Tahsas", "Tir", "Yekatit", "Megabit", "Miazia", "Genbot", "Sene", "Hamle", "Nehase", "Pagume"}

const jdOffset = 1724221 // JDN for 1/1/1 EC (1 Mäskäräm 1), approximately 8/27/8 CE

// IsLeap checks if the given Ethiopian year is a leap year.
func IsLeap(year int) bool {
	if year < 0 {
		return false
	}
	return (year % 4) == 3
}

// DaysInMonth returns the number of days in the specified Ethiopian month and year.
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

// Validate checks if the EtDate is valid.
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

// ToJDN converts an Ethiopian date to Julian Day Number.
func (d EtDate) ToJDN() (int, error) {
	if err := d.Validate(); err != nil {
		return 0, err
	}
	y := d.Year
	m := d.Month
	day := d.Day
	return jdOffset + 365*(y-1) + (y / 4) + 30*(m-1) + day - 1, nil
}

// JDNToEt converts a Julian Day Number to an Ethiopian Calendar date.
func JDNToEt(jdn int) (EtDate, error) {
	if jdn < jdOffset {
		return EtDate{}, errors.New("jdn before Ethiopian epoch")
	}

	// Calculate days since the Ethiopian epoch
	fixed := jdn - jdOffset

	// Estimate year: account for 365 days per year + leap days (1 every 4 years)
	year := fixed / 365
	yearStartJDN, err := (EtDate{Year: year, Month: 1, Day: 1}).ToJDN()
	if err != nil {
		return EtDate{}, err
	}
	if jdn < yearStartJDN {
		year--
		yearStartJDN, err = (EtDate{Year: year, Month: 1, Day: 1}).ToJDN()
		if err != nil {
			return EtDate{}, err
		}
	}

	// Calculate days since the start of the Ethiopian year
	daysSinceYearStart := jdn - yearStartJDN
	if daysSinceYearStart < 0 {
		return EtDate{}, errors.New("invalid day calculation")
	}

	// Calculate month and day
	month := daysSinceYearStart/30 + 1
	day := daysSinceYearStart%30 + 1
	if month > 13 {
		month = 1
		year++
		day = daysSinceYearStart - 12*30 + 1
	}

	d := EtDate{Year: year, Month: int(month), Day: int(day)}
	if err := d.Validate(); err != nil {
		return EtDate{}, err
	}
	return d, nil
}

// GregorianToJDN converts a Gregorian date to Julian Day Number.
func GregorianToJDN(year, month, day int) (int, error) {
	if year == 0 {
		return 0, errors.New("no year 0 in Gregorian")
	}
	if month < 1 || month > 12 {
		return 0, errors.New("month must be between 1 and 12")
	}
	if day < 1 {
		return 0, errors.New("day must be positive")
	}
	// Basic validation for day of month
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if month == 2 && (year%4 == 0 && (year%100 != 0 || year%400 == 0)) {
		daysInMonth[1] = 29
	}
	if day > daysInMonth[month-1] {
		return 0, errors.New("day is out of range for the given month")
	}

	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3
	jdn := day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
	return jdn, nil
}

// JDNToGregorian converts a Julian Day Number to a Gregorian date.
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

// ToGregorian converts an Ethiopian Calendar date to a Gregorian date.
func (d EtDate) ToGregorian() (int, int, int, error) {
	jdn, err := d.ToJDN()
	if err != nil {
		return 0, 0, 0, err
	}
	return JDNToGregorian(jdn)
}

// FromGregorian converts a Gregorian date to an Ethiopian Calendar date.
func FromGregorian(year, month, day int) (EtDate, error) {
	// Validate input using time package for precise Gregorian date validation
	_, err := time.Parse("2006-01-02", fmt.Sprintf("%04d-%02d-%02d", year, month, day))
	if err != nil {
		return EtDate{}, errors.New("invalid Gregorian date: " + err.Error())
	}

	jdn, err := GregorianToJDN(year, month, day)
	if err != nil {
		return EtDate{}, err
	}
	return JDNToEt(jdn)
}

// Format formats the Ethiopian date according to the specified layout.
func (d EtDate) Format(layout string) string {
	str := strings.ReplaceAll(layout, "YYYY", fmt.Sprintf("%04d", d.Year))
	str = strings.ReplaceAll(str, "MM", fmt.Sprintf("%02d", d.Month))
	str = strings.ReplaceAll(str, "DD", fmt.Sprintf("%02d", d.Day))
	str = strings.ReplaceAll(str, "Month", monthNames[d.Month])
	return str
}

// AddDays adds or subtracts the specified number of days to the Ethiopian date.
func (d EtDate) AddDays(days int) (EtDate, error) {
	jdn, err := d.ToJDN()
	if err != nil {
		return EtDate{}, err
	}
	return JDNToEt(jdn + days)
}

// AddMonths adds or subtracts the specified number of months to the Ethiopian date.
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

// AddYears adds or subtracts the specified number of years to the Ethiopian date.
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
