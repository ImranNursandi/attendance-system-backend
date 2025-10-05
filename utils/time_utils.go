package utils

import (
	"fmt"
	"time"
)

// CheckClockInLate checks if clock in is late and returns late minutes
func CheckClockInLate(clockInTime time.Time, maxClockIn string, tolerance int) (bool, int) {
	// Parse max clock in time
	maxTime, err := time.Parse("15:04:05", maxClockIn)
	if err != nil {
		return false, 0
	}

	// Create comparable times
	clockIn := time.Date(0, 1, 1, clockInTime.Hour(), clockInTime.Minute(), clockInTime.Second(), 0, time.UTC)
	max := time.Date(0, 1, 1, maxTime.Hour(), maxTime.Minute(), maxTime.Second(), 0, time.UTC)

	// Add tolerance
	maxWithTolerance := max.Add(time.Duration(tolerance) * time.Minute)

	// Check if late
	if clockIn.After(maxWithTolerance) {
		lateMinutes := int(clockIn.Sub(maxWithTolerance).Minutes())
		return true, lateMinutes
	}

	return false, 0
}

// CheckClockOutEarly checks if clock out is early and returns early minutes
func CheckClockOutEarly(clockOutTime time.Time, maxClockOut string, penaltyThreshold int) (bool, int) {
	// Parse max clock out time
	minTime, err := time.Parse("15:04:05", maxClockOut)
	if err != nil {
		return false, 0
	}

	// Create comparable times
	clockOut := time.Date(0, 1, 1, clockOutTime.Hour(), clockOutTime.Minute(), clockOutTime.Second(), 0, time.UTC)
	min := time.Date(0, 1, 1, minTime.Hour(), minTime.Minute(), minTime.Second(), 0, time.UTC)

	// Subtract penalty threshold
	minWithThreshold := min.Add(-time.Duration(penaltyThreshold) * time.Minute)

	// Check if early leave
	if clockOut.Before(minWithThreshold) {
		earlyMinutes := int(minWithThreshold.Sub(clockOut).Minutes())
		return true, earlyMinutes
	}

	return false, 0
}

// CalculateWorkHours calculates work hours between clock in and clock out
func CalculateWorkHours(clockIn, clockOut time.Time) float64 {
	duration := clockOut.Sub(clockIn)
	return duration.Hours()
}

// FormatDuration formats duration to human readable string
func FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}

// GetStartOfDay returns the start of the day for the given time
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay returns the end of the day for the given time
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// IsWeekend checks if the given time is a weekend
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// GetBusinessDays calculates business days between two dates
func GetBusinessDays(start, end time.Time) int {
	businessDays := 0
	current := start

	for current.Before(end) || current.Equal(end) {
		if !IsWeekend(current) {
			businessDays++
		}
		current = current.AddDate(0, 0, 1)
	}

	return businessDays
}

// ParseDateRange parses start and end date strings
func ParseDateRange(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Adjust end date to end of day
	endDate = GetEndOfDay(endDate)

	return startDate, endDate, nil
}