package controllers

import (
	"attendance-system/services"
	"attendance-system/utils"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController() *ReportController {
	return &ReportController{
		reportService: services.NewReportService(),
	}
}

// GenerateAttendanceReport godoc
// @Summary Generate attendance report
// @Description Generate detailed attendance report for a period
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param department_id query int false "Filter by department ID"
// @Success 200 {object} utils.Response{data=[]services.AttendanceReport}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /reports/attendance [get]
func (c *ReportController) GenerateAttendanceReport(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	departmentID, _ := strconv.ParseUint(ctx.Query("department_id"), 10, 32)

	if startDate == "" || endDate == "" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	reports, err := c.reportService.GenerateAttendanceReport(startDate, endDate, uint(departmentID))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Attendance report generated successfully", reports)
}

// GenerateSummaryReport godoc
// @Summary Generate summary report
// @Description Generate summary attendance report for a period
// @Tags reports
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param department_id query int false "Filter by department ID"
// @Success 200 {object} utils.Response{data=services.SummaryReport}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /reports/summary [get]
func (c *ReportController) GenerateSummaryReport(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	departmentID, _ := strconv.ParseUint(ctx.Query("department_id"), 10, 32)

	if startDate == "" || endDate == "" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	summary, err := c.reportService.GenerateSummaryReport(startDate, endDate, uint(departmentID))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Summary report generated successfully", summary)
}

// GenerateDepartmentReport godoc
// @Summary Generate department report
// @Description Generate monthly department attendance report
// @Tags reports
// @Accept json
// @Produce json
// @Param department_id path int true "Department ID"
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Success 200 {object} utils.Response{data=map[string]interface{}}
// @Failure 500 {object} utils.Response
// @Router /reports/department/{department_id} [get]
func (c *ReportController) GenerateDepartmentReport(ctx *gin.Context) {
	departmentID, err := strconv.ParseUint(ctx.Param("department_id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid department ID")
		return
	}

	month, _ := strconv.Atoi(ctx.Query("month"))
	year, _ := strconv.Atoi(ctx.Query("year"))

	report, err := c.reportService.GenerateDepartmentReport(uint(departmentID), month, year)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Department report generated successfully", report)
}

// ExportAttendanceReport godoc
// @Summary Export attendance report
// @Description Export attendance report in Excel or CSV format
// @Tags reports
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,text/csv
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param department_id query int false "Filter by department ID"
// @Param format query string false "Export format (excel or csv)" default(excel)
// @Success 200 {file} file "Exported file"
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /reports/export/attendance [get]
func (c *ReportController) ExportAttendanceReport(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	departmentID, _ := strconv.ParseUint(ctx.Query("department_id"), 10, 32)
	format := ctx.DefaultQuery("format", "excel")

	if startDate == "" || endDate == "" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	reports, err := c.reportService.GenerateAttendanceReport(startDate, endDate, uint(departmentID))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	filename := fmt.Sprintf("attendance-report-%s-to-%s", startDate, endDate)

	switch format {
	case "csv":
		err = c.exportAttendanceCSV(ctx, reports, filename)
	default:
		err = c.exportAttendanceExcel(ctx, reports, filename)
	}

	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
}

// ExportSummaryReport godoc
// @Summary Export summary report
// @Description Export summary report in Excel or CSV format
// @Tags reports
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,text/csv
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param department_id query int false "Filter by department ID"
// @Param format query string false "Export format (excel or csv)" default(excel)
// @Success 200 {file} file "Exported file"
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /reports/export/summary [get]
func (c *ReportController) ExportSummaryReport(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	departmentID, _ := strconv.ParseUint(ctx.Query("department_id"), 10, 32)
	format := ctx.DefaultQuery("format", "excel")

	if startDate == "" || endDate == "" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Start date and end date are required")
		return
	}

	summary, err := c.reportService.GenerateSummaryReport(startDate, endDate, uint(departmentID))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	filename := fmt.Sprintf("summary-report-%s-to-%s", startDate, endDate)

	switch format {
	case "csv":
		err = c.exportSummaryCSV(ctx, summary, filename)
	default:
		err = c.exportSummaryExcel(ctx, summary, filename)
	}

	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
}

// ExportDepartmentReport godoc
// @Summary Export department report
// @Description Export department report in Excel or CSV format
// @Tags reports
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,text/csv
// @Param department_id path int true "Department ID"
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Param format query string false "Export format (excel or csv)" default(excel)
// @Success 200 {file} file "Exported file"
// @Failure 500 {object} utils.Response
// @Router /reports/export/department/{department_id} [get]
func (c *ReportController) ExportDepartmentReport(ctx *gin.Context) {
	departmentID, err := strconv.ParseUint(ctx.Param("department_id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid department ID")
		return
	}

	month, _ := strconv.Atoi(ctx.Query("month"))
	year, _ := strconv.Atoi(ctx.Query("year"))
	format := ctx.DefaultQuery("format", "excel")

	report, err := c.reportService.GenerateDepartmentReport(uint(departmentID), month, year)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	filename := fmt.Sprintf("department-report-%d-%d-%d", departmentID, month, year)

	switch format {
	case "csv":
		err = c.exportDepartmentCSV(ctx, report, filename)
	default:
		err = c.exportDepartmentExcel(ctx, report, filename)
	}

	if err != nil {
		utils.HandleError(ctx, err)
		return
	}
}

// Excel Export Methods
func (c *ReportController) exportAttendanceExcel(ctx *gin.Context, reports []services.AttendanceReport, filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	// Create a new sheet
	sheetName := "Attendance Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Set headers
	headers := []string{"Employee ID", "Employee Name", "Department", "Date", "Clock In", "Clock Out", "Work Hours", "Status", "Late Minutes", "Early Minutes"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Set data
	for i, report := range reports {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), report.EmployeeID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.EmployeeName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), report.Department)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), report.Date)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), report.ClockIn.Format("2006-01-02 15:04:05"))
		
		if !report.ClockOut.IsZero() {
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), report.ClockOut.Format("2006-01-02 15:04:05"))
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "Not Clocked Out")
		}
		
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("%.2f", report.WorkHours))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), report.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), report.LateMinutes)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), report.EarlyMinutes)
	}

	// Set active sheet and apply styling
	f.SetActiveSheet(index)
	
	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 15)
	f.SetColWidth(sheetName, "B", "B", 25)
	f.SetColWidth(sheetName, "C", "C", 20)
	f.SetColWidth(sheetName, "D", "F", 18)
	f.SetColWidth(sheetName, "G", "J", 15)

	// Style headers
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2c5aa0"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", "J1", style)

	// Set response headers
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")

	return f.Write(ctx.Writer)
}

func (c *ReportController) exportSummaryExcel(ctx *gin.Context, summary *services.SummaryReport, filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Summary Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Set summary data
	f.SetCellValue(sheetName, "A1", "Summary Report")
	f.SetCellValue(sheetName, "A2", "Period")
	f.SetCellValue(sheetName, "B2", summary.Period)
	
	f.SetCellValue(sheetName, "A3", "Total Employees")
	f.SetCellValue(sheetName, "B3", summary.TotalEmployees)
	
	f.SetCellValue(sheetName, "A4", "Total Present")
	f.SetCellValue(sheetName, "B4", summary.TotalPresent)
	
	f.SetCellValue(sheetName, "A5", "Total Late")
	f.SetCellValue(sheetName, "B5", summary.TotalLate)
	
	f.SetCellValue(sheetName, "A6", "Total Absent")
	f.SetCellValue(sheetName, "B6", summary.TotalAbsent)
	
	f.SetCellValue(sheetName, "A7", "Total Work Hours")
	f.SetCellValue(sheetName, "B7", summary.TotalWorkHours)
	
	f.SetCellValue(sheetName, "A8", "Average Work Hours")
	f.SetCellValue(sheetName, "B8", summary.AverageWorkHours)

	// Apply styling
	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "B", "B", 25)

	// Style headers
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16},
	})
	f.SetCellStyle(sheetName, "A1", "B1", headerStyle)

	labelStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle(sheetName, "A2", "A8", labelStyle)

	f.SetActiveSheet(index)

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")

	return f.Write(ctx.Writer)
}

func (c *ReportController) exportDepartmentExcel(ctx *gin.Context, report map[string]interface{}, filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Department Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Safely extract department report data
	deptReportInterface, exists := report["department_report"]
	if !exists {
		return fmt.Errorf("department_report not found in response")
	}

	// Convert to the actual struct type
	deptReportMap, ok := deptReportInterface.(map[string]interface{})
	if !ok {
		// If it's already a struct, we need to handle it differently
		// For now, let's use a type assertion with proper error handling
		return fmt.Errorf("unexpected type for department_report: %T", deptReportInterface)
	}

	departmentName, _ := deptReportMap["department_name"].(string)
	totalEmployees, _ := deptReportMap["total_employees"].(float64) // JSON numbers are float64
	period := report["period"].(string)

	// Set department header
	f.SetCellValue(sheetName, "A1", "Department Report")
	f.SetCellValue(sheetName, "A2", "Department")
	f.SetCellValue(sheetName, "B2", departmentName)
	f.SetCellValue(sheetName, "A3", "Period")
	f.SetCellValue(sheetName, "B3", period)
	f.SetCellValue(sheetName, "A4", "Total Employees")
	f.SetCellValue(sheetName, "B4", int(totalEmployees))

	// Set summary data
	summaryInterface, exists := deptReportMap["summary"]
	if !exists {
		return fmt.Errorf("summary not found in department_report")
	}

	summary, ok := summaryInterface.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for summary: %T", summaryInterface)
	}

	f.SetCellValue(sheetName, "A6", "Summary Statistics")
	f.SetCellValue(sheetName, "A7", "Total Present")
	f.SetCellValue(sheetName, "B7", summary["total_present"])
	f.SetCellValue(sheetName, "A8", "Total Late")
	f.SetCellValue(sheetName, "B8", summary["total_late"])
	f.SetCellValue(sheetName, "A9", "Total Absent")
	f.SetCellValue(sheetName, "B9", summary["total_absent"])
	f.SetCellValue(sheetName, "A10", "Attendance Rate")
	f.SetCellValue(sheetName, "B10", fmt.Sprintf("%.2f%%", summary["attendance_rate"]))
	f.SetCellValue(sheetName, "A11", "Average Work Hours")
	f.SetCellValue(sheetName, "B11", fmt.Sprintf("%.2f", summary["average_work_hours"]))

	// Employee statistics
	f.SetCellValue(sheetName, "A13", "Employee Statistics")
	employeeStatsInterface, exists := deptReportMap["employee_stats"]
	if !exists {
		return fmt.Errorf("employee_stats not found in department_report")
	}

	employeeStats, ok := employeeStatsInterface.([]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for employee_stats: %T", employeeStatsInterface)
	}

	headers := []string{"Employee ID", "Employee Name", "Present Days", "Late Days", "Absent Days", "Avg Work Hours"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 14)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, empStatInterface := range employeeStats {
		empStat, ok := empStatInterface.(map[string]interface{})
		if !ok {
			continue // Skip invalid entries
		}

		statsInterface, exists := empStat["stats"]
		if !exists {
			continue
		}

		stats, ok := statsInterface.(map[string]interface{})
		if !ok {
			continue
		}

		row := i + 15
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), empStat["employee_id"])
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), empStat["employee_name"])
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), stats["total_present"])
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), stats["total_late"])
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), stats["total_absent"])
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), stats["avg_work_hours"])
	}

	// Apply styling
	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "B", "B", 25)
	f.SetColWidth(sheetName, "C", "F", 15)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 16},
	})
	f.SetCellStyle(sheetName, "A1", "A1", headerStyle)

	sectionStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
	})
	f.SetCellStyle(sheetName, "A6", "A6", sectionStyle)
	f.SetCellStyle(sheetName, "A13", "A13", sectionStyle)

	tableHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2c5aa0"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A14", "F14", tableHeaderStyle)

	f.SetActiveSheet(index)

	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")

	return f.Write(ctx.Writer)
}

// CSV Export Methods
func (c *ReportController) exportAttendanceCSV(ctx *gin.Context, reports []services.AttendanceReport, filename string) error {
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")

	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	// Write headers
	headers := []string{"Employee ID", "Employee Name", "Department", "Date", "Clock In", "Clock Out", "Work Hours", "Status", "Late Minutes", "Early Minutes"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data
	for _, report := range reports {
		clockOut := ""
		if !report.ClockOut.IsZero() {
			clockOut = report.ClockOut.Format("2006-01-02 15:04:05")
		} else {
			clockOut = "Not Clocked Out"
		}

		record := []string{
			report.EmployeeID,
			report.EmployeeName,
			report.Department,
			report.Date,
			report.ClockIn.Format("2006-01-02 15:04:05"),
			clockOut,
			fmt.Sprintf("%.2f", report.WorkHours),
			report.Status,
			strconv.Itoa(report.LateMinutes),
			strconv.Itoa(report.EarlyMinutes),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func (c *ReportController) exportSummaryCSV(ctx *gin.Context, summary *services.SummaryReport, filename string) error {
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")

	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	// Write headers
	if err := writer.Write([]string{"Metric", "Value"}); err != nil {
		return err
	}

	// Write data
	records := [][]string{
		{"Period", summary.Period},
		{"Total Employees", strconv.FormatInt(summary.TotalEmployees, 10)},
		{"Total Present", strconv.FormatInt(summary.TotalPresent, 10)},
		{"Total Late", strconv.FormatInt(summary.TotalLate, 10)},
		{"Total Absent", strconv.FormatInt(summary.TotalAbsent, 10)},
		{"Total Work Hours", summary.TotalWorkHours},
		{"Average Work Hours", summary.AverageWorkHours},
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func (c *ReportController) exportDepartmentCSV(ctx *gin.Context, report map[string]interface{}, filename string) error {
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")

	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	// Safely extract department report data
	deptReportInterface, exists := report["department_report"]
	if !exists {
		return fmt.Errorf("department_report not found in response")
	}

	deptReport, ok := deptReportInterface.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for department_report: %T", deptReportInterface)
	}

	summaryInterface, exists := deptReport["summary"]
	if !exists {
		return fmt.Errorf("summary not found in department_report")
	}

	summary, ok := summaryInterface.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for summary: %T", summaryInterface)
	}

	employeeStatsInterface, exists := deptReport["employee_stats"]
	if !exists {
		return fmt.Errorf("employee_stats not found in department_report")
	}

	employeeStats, ok := employeeStatsInterface.([]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for employee_stats: %T", employeeStatsInterface)
	}

	// Write department header
	if err := writer.Write([]string{"Department Report"}); err != nil {
		return err
	}
	writer.Write([]string{"Department", deptReport["department_name"].(string)})
	writer.Write([]string{"Period", report["period"].(string)})
	writer.Write([]string{"Total Employees", fmt.Sprintf("%.0f", deptReport["total_employees"])})
	writer.Write([]string{}) // Empty line

	// Write summary
	writer.Write([]string{"Summary Statistics"})
	writer.Write([]string{"Total Present", fmt.Sprintf("%.0f", summary["total_present"])})
	writer.Write([]string{"Total Late", fmt.Sprintf("%.0f", summary["total_late"])})
	writer.Write([]string{"Total Absent", fmt.Sprintf("%.0f", summary["total_absent"])})
	writer.Write([]string{"Attendance Rate", fmt.Sprintf("%.2f%%", summary["attendance_rate"])})
	writer.Write([]string{"Average Work Hours", fmt.Sprintf("%.2f", summary["average_work_hours"])})
	writer.Write([]string{}) // Empty line

	// Write employee statistics header
	writer.Write([]string{"Employee Statistics"})
	writer.Write([]string{"Employee ID", "Employee Name", "Present Days", "Late Days", "Absent Days", "Avg Work Hours"})

	// Write employee data
	for _, empStatInterface := range employeeStats {
		empStat, ok := empStatInterface.(map[string]interface{})
		if !ok {
			continue // Skip invalid entries
		}

		statsInterface, exists := empStat["stats"]
		if !exists {
			continue
		}

		stats, ok := statsInterface.(map[string]interface{})
		if !ok {
			continue
		}

		record := []string{
			empStat["employee_id"].(string),
			empStat["employee_name"].(string),
			fmt.Sprintf("%.0f", stats["total_present"]),
			fmt.Sprintf("%.0f", stats["total_late"]),
			fmt.Sprintf("%.0f", stats["total_absent"]),
			fmt.Sprintf("%.2f", stats["avg_work_hours"]),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}