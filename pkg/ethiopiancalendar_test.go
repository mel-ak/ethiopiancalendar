package ethiopiancalendar

import "testing"

func TestIsLeap(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{2015, true},
		{2016, false},
		{0, false},
	}

	for _, tt := range tests {
		if got := IsLeap(tt.year); got != tt.want {
			t.Errorf("IsLeap(%d) = %v, want %v", tt.year, got, tt.want)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	if DaysInMonth(2015, 13) != 6 {
		t.Error("Expected 6 days in Pagume 2015")
	}
	if DaysInMonth(2016, 13) != 5 {
		t.Error("Expected 5 days in Pagume 2016")
	}
	if DaysInMonth(2017, 5) != 30 {
		t.Error("Expected 30 days in month 5")
	}
}

func TestConversion(t *testing.T) {
	et := EtDate{Year: 2016, Month: 1, Day: 1}
	gy, gm, gd, err := et.ToGregorian()

	if err != nil {
		t.Error(err)
	}

	if gy != 2023 || gm != 9 || gd != 12 {
		t.Errorf("Expected 2023-09-12, got %d-%d-%d", gy, gm, gd)
	}
}

func TestConversionBack(t *testing.T) {
	gy, gm, gd := 2023, 9, 12
	et, err := FromGregorian(gy, gm, gd)

	if err != nil {
		t.Error(err)
	}

	if et.Year != 2016 || et.Month != 1 || et.Day != 1 {
		t.Errorf("Expected 2016-01-01, got %d-%d-%d", et.Year, et.Month, et.Day)
	}
}

func TestAddDays(t *testing.T) {
	et := EtDate{Year: 2016, Month: 1, Day: 1}
	future, err := et.AddDays(10)

	if err != nil {
		t.Error(err)
	}

	if future.Year != 2016 || future.Month != 1 || future.Day != 11 {
		t.Errorf("Expected 2016-01-11, got %d-%d-%d", future.Year, future.Month, future.Day)
	}
}
