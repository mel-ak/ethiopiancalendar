package main

import (
	"fmt"

	ethiopiancalendar "github.com/mel-ak/ethiopiancalendar/pkg"
)

func main() {
	date := ethiopiancalendar.EtDate{Year: 2015, Month: 1, Day: 1}

	fmt.Println(date)

	fmt.Println(ethiopiancalendar.DaysInMonth(2016, 13))
	fmt.Println(ethiopiancalendar.DaysInMonth(2015, 13))
	fmt.Println(ethiopiancalendar.DaysInMonth(2017, 5))

	date = ethiopiancalendar.EtDate{Year: 2017, Month: 2, Day: 6}
	if err := date.Validate(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("====================================================")
	etDate := ethiopiancalendar.EtDate{Year: 2016, Month: 1, Day: 1}

	gy, gm, gd, err := etDate.ToGregorian()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Gregorian: %d-%02d-%02d\n", gy, gm, gd)

	gYear, gMonth, gDay := 2023, 9, 12
	et, err := ethiopiancalendar.FromGregorian(gYear, gMonth, gDay)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Ethiopian: %d/%d/%d\n", et.Year, et.Month, et.Day)

	fmt.Println(date.Format("DD Month YYYY"))

	fmt.Println("====================================================")

	newDate, err := date.AddDays(10)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(newDate)

	monthAdded := date.AddMonths(3)
	fmt.Println(monthAdded)

	yearAdded := date.AddYears(1)
	fmt.Println(yearAdded)
}
