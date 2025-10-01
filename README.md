# Ethiopian Calendar Go Library

[![Go](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](https://github.com/mel-ak/ethiopian-calendar-go/actions)
[![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)](https://github.com/mel-ak/ethiopian-calendar-go)

A clean, idiomatic Go utility library for working with the Ethiopian (Ge'ez) calendar system. Supports bidirectional conversions between Ethiopian and Gregorian dates, leap year calculations, month day counts, date formatting, and basic arithmetic (add/subtract days, months, years).

The Ethiopian calendar has 13 months: 12 with 30 days each, and Pagumē (13th) with 5 or 6 days in leap years. It lags the Gregorian by ~7-8 years.

## Features

- **Date Representation**: `EtDate` struct with validation.
- **Leap Years**: Rule: `year % 4 == 3`.
- **Conversions**: Via Julian Day Number (JDN) for accuracy.
- **Formatting**: Custom layouts with month names (Amharic transliterated).
- **Arithmetic**: Add/subtract days (precise), months/years (with normalization).
- **No Dependencies**: Pure Go, lightweight.

## Installation

```bash
go get github.com/mel-ak/ethiopian-calendar-go
