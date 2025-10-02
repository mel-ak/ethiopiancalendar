# Ethiopian Calendar Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/mel-ak/ethiopian-calendar-go.svg)](https://pkg.go.dev/github.com/mel-ak/ethiopian-calendar-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/mel-ak/ethiopian-calendar-go)](https://goreportcard.com/report/github.com/mel-ak/ethiopian-calendar-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](https://github.com/mel-ak/ethiopian-calendar-go/actions)
[![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)](https://github.com/mel-ak/ethiopian-calendar-go)

A clean, idiomatic Go utility library for working with the Ethiopian (Ge'ez) calendar system. Supports bidirectional conversions between Ethiopian and Gregorian dates, leap year calculations, month day counts, date formatting, and basic arithmetic operations.

The Ethiopian calendar has 13 months: 12 with 30 days each, and Pagumē (13th) with 5 or 6 days in leap years. It lags the Gregorian by ~7-8 years.

## Features

- **Date Representation**: `EtDate` struct with validation
- **Leap Year Detection**: Accurate calculation of Ethiopian leap years
- **Bidirectional Conversion**: Between Ethiopian and Gregorian calendars
- **Date Arithmetic**: Add/subtract days, months, and years
- **Flexible Formatting**: Custom date formatting with month names
- **REST API**: Built-in HTTP server with JSON endpoints
- **Zero Dependencies**: Pure Go implementation
- **Comprehensive Tests**: 95%+ test coverage

## Installation

```bash
go get github.com/mel-ak/ethiopian-calendar-go
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"
	"time"

	ethio "github.com/mel-ak/ethiopian-calendar-go"
)

func main() {
	// Get current date in Ethiopian calendar
	now := time.Now()
	etDate, _ := ethio.FromGregorian(now.Year(), int(now.Month()), now.Day())
	
	// Format and print Ethiopian date
	fmt.Printf("Today in Ethiopian: %s\n", etDate.Format("DD Month YYYY"))
	// Output: Today in Ethiopian: 21 Meskerem 2016

	// Convert back to Gregorian
	gy, gm, gd, _ := etDate.ToGregorian()
	fmt.Printf("Back to Gregorian: %d-%02d-%02d\n", gy, gm, gd)

	// Check if a year is a leap year
	year := 2015
	fmt.Printf("Is %d a leap year? %v\n", year, ethio.IsLeap(year))
	// Output: Is 2015 a leap year? true

	// Add 40 days to current date
	future, _ := etDate.AddDays(40)
	fmt.Printf("40 days from now: %s\n", future.Format("DD/MM/YYYY"))
}
```

## API Reference

### Types

#### EtDate
Represents a date in the Ethiopian calendar.

```go
type EtDate struct {
    Year  int
    Month int // 1-13 (1-12 are 30-day months, 13 is Pagume with 5 or 6 days)
    Day   int
}
```

### Functions

#### Date Creation and Validation

- `FromGregorian(year, month, day int) (EtDate, error)`: Creates an Ethiopian date from Gregorian date
- `(d EtDate) Validate() error`: Validates the Ethiopian date

#### Date Conversion

- `(d EtDate) ToGregorian() (int, int, int, error)`: Converts to Gregorian date
- `(d EtDate) ToJDN() (int, error)`: Converts to Julian Day Number
- `JDNToEt(jdn int) (EtDate, error)`: Creates Ethiopian date from JDN

#### Date Arithmetic

- `(d EtDate) AddDays(days int) (EtDate, error)`: Adds/subtracts days
- `(d EtDate) AddMonths(months int) (EtDate, error)`: Adds/subtracts months
- `(d EtDate) AddYears(years int) (EtDate, error)`: Adds/subtracts years

#### Formatting

- `(d EtDate) Format(layout string) string`: Formats the date using the specified layout
  - `YYYY`: 4-digit year (e.g., 2016)
  - `MM`: 2-digit month (01-13)
  - `DD`: 2-digit day (01-30, or 01-05/06 for Pagume)
  - `Month`: Full month name (e.g., "Meskerem")

#### Calendar Information

- `IsLeap(year int) bool`: Checks if a year is a leap year
- `DaysInMonth(year, month int) int`: Returns number of days in a month

## Web API

The package includes a REST API server that can be started as follows:

```go
package main

import (
	"log"
	"net/http"

	"github.com/mel-ak/ethiopian-calendar-go/api"
)

func main() {
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### API Endpoints

- `POST /api/convert`: Convert between Ethiopian and Gregorian dates
  ```json
  {
    "type": "etToGreg" | "gregToEt",
    "year": 2016,
    "month": 1,
    "day": 1
  }
  ```

- `POST /api/format`: Format an Ethiopian date
  ```json
  {
    "year": 2016,
    "month": 1,
    "day": 1,
    "layout": "DD Month YYYY"
  }
  ```

- `GET /api/leap?year=2015`: Check if a year is a leap year
- `GET /api/daysinmonth?year=2015&month=13`: Get days in a month

## Running Tests

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Ethiopian calendar calculations based on Julian Day Numbers
- inspired by various calendar conversion algorithms
