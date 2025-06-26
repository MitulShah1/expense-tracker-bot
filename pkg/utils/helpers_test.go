package utils

import (
	"strconv"
	"testing"
	"time"
)

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		amount float64
		want   string
	}{
		{123.456, "₹123.46"},
		{0, "₹0.00"},
		{-10.5, "₹-10.50"},
		{9999999.99, "₹9999999.99"},
	}
	for _, tt := range tests {
		got := FormatCurrency(tt.amount)
		if got != tt.want {
			t.Errorf("FormatCurrency(%v) = %q, want %q", tt.amount, got, tt.want)
		}
	}
}

func TestFormatDate(t *testing.T) {
	d := time.Date(2024, 6, 1, 15, 4, 5, 0, time.UTC)
	want := "01 Jun 2024"
	if got := FormatDate(d); got != want {
		t.Errorf("FormatDate() = %q, want %q", got, want)
	}
}

func TestFormatDateTime(t *testing.T) {
	d := time.Date(2024, 6, 1, 15, 4, 5, 0, time.UTC)
	want := "01 Jun 2024 15:04:05"
	if got := FormatDateTime(d); got != want {
		t.Errorf("FormatDateTime() = %q, want %q", got, want)
	}
}

func TestFormatMonth(t *testing.T) {
	d := time.Date(2024, 6, 1, 15, 4, 5, 0, time.UTC)
	want := "2024-06"
	if got := FormatMonth(d); got != want {
		t.Errorf("FormatMonth() = %q, want %q", got, want)
	}
}

func TestCalculateAverage(t *testing.T) {
	tests := []struct {
		name    string
		numbers []float64
		want    float64
	}{
		{"empty", []float64{}, 0},
		{"single", []float64{10}, 10},
		{"multiple", []float64{1, 2, 3, 4}, 2.5},
		{"negative", []float64{-1, -2, -3}, -2},
	}
	for _, tt := range tests {
		got := CalculateAverage(tt.numbers)
		if got != tt.want {
			t.Errorf("%s: CalculateAverage(%v) = %v, want %v", tt.name, tt.numbers, got, tt.want)
		}
	}
}

func TestCalculateTotal(t *testing.T) {
	tests := []struct {
		name    string
		numbers []float64
		want    float64
	}{
		{"empty", []float64{}, 0},
		{"single", []float64{10}, 10},
		{"multiple", []float64{1, 2, 3, 4}, 10},
		{"negative", []float64{-1, -2, -3}, -6},
	}
	for _, tt := range tests {
		got := CalculateTotal(tt.numbers)
		if got != tt.want {
			t.Errorf("%s: CalculateTotal(%v) = %v, want %v", tt.name, tt.numbers, got, tt.want)
		}
	}
}

func TestGetMonthRange(t *testing.T) {
	start, end := GetMonthRange(2024, time.June)
	wantStart := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
	if !start.Equal(wantStart) {
		t.Errorf("GetMonthRange start = %v, want %v", start, wantStart)
	}
	if !end.Equal(wantEnd) {
		t.Errorf("GetMonthRange end = %v, want %v", end, wantEnd)
	}
}

func TestGetYearRange(t *testing.T) {
	start, end := GetYearRange(2024)
	wantStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
	if !start.Equal(wantStart) {
		t.Errorf("GetYearRange start = %v, want %v", start, wantStart)
	}
	if !end.Equal(wantEnd) {
		t.Errorf("GetYearRange end = %v, want %v", end, wantEnd)
	}
}

func TestGetCurrentMonthRange(t *testing.T) {
	start, end := GetCurrentMonthRange()
	now := time.Now().UTC()
	wantStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	wantEnd := wantStart.AddDate(0, 1, 0).Add(-time.Second)
	if !start.Equal(wantStart) {
		t.Errorf("GetCurrentMonthRange start = %v, want %v", start, wantStart)
	}
	if !end.Equal(wantEnd) {
		t.Errorf("GetCurrentMonthRange end = %v, want %v", end, wantEnd)
	}
}

func TestGetCurrentYearRange(t *testing.T) {
	start, end := GetCurrentYearRange()
	now := time.Now().UTC()
	wantStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := wantStart.AddDate(1, 0, 0).Add(-time.Second)
	if !start.Equal(wantStart) {
		t.Errorf("GetCurrentYearRange start = %v, want %v", start, wantStart)
	}
	if !end.Equal(wantEnd) {
		t.Errorf("GetCurrentYearRange end = %v, want %v", end, wantEnd)
	}
}

func TestIsValidAmount(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"123.45", true},
		{"0", true},
		{"-10.5", true},
		{"abc", false},
		{"", false},
		{"1e3", true},
	}
	for _, tt := range tests {
		got := IsValidAmount(tt.input)
		if got != tt.want {
			t.Errorf("IsValidAmount(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		input   string
		want    float64
		wantErr bool
	}{
		{"123.45", 123.45, false},
		{"0", 0, false},
		{"-10.5", -10.5, false},
		{"abc", 0, true},
		{"", 0, true},
		{"1e3", 1000, false},
	}
	for _, tt := range tests {
		got, err := FormatAmount(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("FormatAmount(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("FormatAmount(%q) = %v, want %v", tt.input, got, tt.want)
		}
		if !tt.wantErr && strconv.FormatFloat(got, 'f', -1, 64) != strconv.FormatFloat(tt.want, 'f', -1, 64) {
			t.Errorf("FormatAmount(%q) = %v, want %v (string compare)", tt.input, got, tt.want)
		}
	}
}
