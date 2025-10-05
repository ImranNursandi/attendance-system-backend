package controllers

import (
	"attendance-system/services"
	"attendance-system/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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