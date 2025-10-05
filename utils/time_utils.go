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
	maxTime, err := time.Parse("15:04:05", maxClockOut)
	if err != nil {
		return false, 0
	}

	// Create comparable times
	clockOut := time.Date(0, 1, 1, clockOutTime.Hour(), clockOutTime.Minute(), clockOutTime.Second(), 0, time.UTC)
	max := time.Date(0, 1, 1, maxTime.Hour(), maxTime.Minute(), maxTime.Second(), 0, time.UTC)

	// Subtract penalty threshold to get the minimum allowed time
	minAllowedTime := max.Add(-time.Duration(penaltyThreshold) * time.Minute)

	// Check if early leave (clocking out before minimum allowed time)
	if clockOut.Before(minAllowedTime) {
		earlyMinutes := int(minAllowedTime.Sub(clockOut).Minutes())
		return true, earlyMinutes
	}

	return false, 0
}

// NEW: Check if clock out is after max clock out time
func CheckClockOutLate(clockOutTime time.Time, maxClockOut string) (bool, int) {
	maxTime, err := time.Parse("15:04:05", maxClockOut)
	if err != nil {
		return false, 0
	}

	clockOut := time.Date(0, 1, 1, clockOutTime.Hour(), clockOutTime.Minute(), clockOutTime.Second(), 0, time.UTC)
	max := time.Date(0, 1, 1, maxTime.Hour(), maxTime.Minute(), maxTime.Second(), 0, time.UTC)

	if clockOut.After(max) {
		overtimeMinutes := int(clockOut.Sub(max).Minutes())
		return true, overtimeMinutes
	}

	return false, 0
}

// CalculateWorkHours calculates work hours between clock in and clock out
func CalculateWorkHours(clockIn, clockOut time.Time) float64 {
	duration := clockOut.Sub(clockIn)
	return duration.Hours()
}

// CalculateWorkHoursDecimal calculates work hours with 2 decimal precision
func CalculateWorkHoursDecimal(clockIn, clockOut time.Time) float64 {
	duration := clockOut.Sub(clockIn)
	hours := duration.Hours()
	return float64(int(hours*100)) / 100 // Round to 2 decimal places
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

// FormatWorkHours formats work hours to decimal format
func FormatWorkHours(hours float64) string {
	return fmt.Sprintf("%.2f hours", hours)
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

// IsHoliday checks if the given date is a holiday (you can extend this with a database)
func IsHoliday(t time.Time) bool {
	// Simple implementation - extend with database lookup
	holidays := map[string]bool{
		"01-01": true, // New Year
		"12-25": true, // Christmas
		// Add more holidays as needed
	}
	
	key := fmt.Sprintf("%02d-%02d", t.Month(), t.Day())
	return holidays[key]
}

// IsWorkingDay checks if it's a working day (not weekend and not holiday)
func IsWorkingDay(t time.Time) bool {
	return !IsWeekend(t) && !IsHoliday(t)
}

// GetBusinessDays calculates business days between two dates
func GetBusinessDays(start, end time.Time) int {
	businessDays := 0
	current := GetStartOfDay(start)
	endDate := GetStartOfDay(end)

	for !current.After(endDate) {
		if IsWorkingDay(current) {
			businessDays++
		}
		current = current.AddDate(0, 0, 1)
	}

	return businessDays
}

// ParseDateRange parses start and end date strings
func ParseDateRange(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start date: %v", err)
		}
		startDate = GetStartOfDay(startDate)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end date: %v", err)
		}
		endDate = GetEndOfDay(endDate)
	}

	return startDate, endDate, nil
}

// NEW: Calculate punctuality status
func CalculatePunctualityStatus(clockIn time.Time, clockOut *time.Time, deptMaxClockIn, deptMaxClockOut string, lateTolerance, earlyLeavePenalty int) (isLate bool, lateMinutes int, isEarlyLeave bool, earlyMinutes int, punctuality string) {
	// Check clock-in punctuality
	isLate, lateMinutes = CheckClockInLate(clockIn, deptMaxClockIn, lateTolerance)
	
	// Check clock-out punctuality
	if clockOut != nil {
		isEarlyLeave, earlyMinutes = CheckClockOutEarly(*clockOut, deptMaxClockOut, earlyLeavePenalty)
	}

	// Determine overall punctuality
	switch {
	case isLate:
		punctuality = "late"
	case isEarlyLeave:
		punctuality = "early_leave"
	default:
		punctuality = "on_time"
	}

	return
}

// NEW: Validate time format
func IsValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04:05", timeStr)
	return err == nil
}

// NEW: Get current fiscal year
func GetCurrentFiscalYear() int {
	now := time.Now()
	year := now.Year()
	
	// Assuming fiscal year starts in April
	if now.Month() < time.April {
		year--
	}
	
	return year
}

// NEW: Get month name
func GetMonthName(month int) string {
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	
	if month < 1 || month > 12 {
		return "Unknown"
	}
	
	return months[month-1]
}

// NEW: Calculate age from birth date
func CalculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	
	// Adjust if birthday hasn't occurred this year
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	
	return age
}