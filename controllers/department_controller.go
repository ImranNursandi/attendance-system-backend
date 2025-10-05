package controllers

import (
	"attendance-system/models"
	"attendance-system/repositories"
	"attendance-system/services"
	"attendance-system/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct {
	departmentService *services.DepartmentService
}

func NewDepartmentController() *DepartmentController {
	return &DepartmentController{
		departmentService: services.NewDepartmentService(),
	}
}

// CreateDepartment godoc
// @Summary Create a new department
// @Description Create a new department with the provided data
// @Tags departments
// @Accept json
// @Produce json
// @Param department body models.DepartmentRequest true "Department data"
// @Success 201 {object} utils.Response{data=models.Department}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /departments [post]
func (c *DepartmentController) CreateDepartment(ctx *gin.Context) {
	var req models.DepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	department, err := c.departmentService.CreateDepartment(req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusCreated, "Department created successfully", department.ToResponse())
}

// GetAllDepartments godoc
// @Summary Get all departments
// @Description Get paginated list of departments with optional filtering and search
// @Tags departments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param status query string false "Filter by status"
// @Success 200 {object} utils.Response{data=[]models.DepartmentResponse}
// @Failure 500 {object} utils.Response
// @Router /departments [get]
func (c *DepartmentController) GetAllDepartments(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")

	var filters []repositories.Filter
	if status := ctx.Query("status"); status != "" {
		filters = append(filters, repositories.Filter{Field: "status", Value: status})
	}

	departments, pagination, err := c.departmentService.GetAllDepartments(filters, search, page, limit)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Convert to response objects
	var departmentResponses []models.DepartmentResponse
	for _, department := range departments {
		deptResp := department.ToResponse()
		// Get employee count for each department using the repository directly
		var employeeCount int64
		c.departmentService.DepartmentRepo().DB.Model(&models.Employee{}).
			Where("department_id = ?", department.ID).Count(&employeeCount)
		deptResp.EmployeeCount = int(employeeCount)
		departmentResponses = append(departmentResponses, deptResp)
	}

	response := map[string]interface{}{
		"departments": departmentResponses,
		"pagination":  pagination,
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Departments retrieved successfully", response)
}

// GetDepartmentByID godoc
// @Summary Get department by ID
// @Description Get department details by ID
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} utils.Response{data=models.DepartmentResponse}
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /departments/{id} [get]
func (c *DepartmentController) GetDepartmentByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid department ID")
		return
	}

	department, err := c.departmentService.GetDepartmentByID(uint(id))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Get employee count using the repository directly
	var employeeCount int64
	c.departmentService.DepartmentRepo().DB.Model(&models.Employee{}).
		Where("department_id = ?", department.ID).Count(&employeeCount)
	
	deptResp := department.ToResponse()
	deptResp.EmployeeCount = int(employeeCount)

	utils.SuccessJSON(ctx, http.StatusOK, "Department retrieved successfully", deptResp)
}

// UpdateDepartment godoc
// @Summary Update department
// @Description Update department details
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param department body models.DepartmentRequest true "Department data"
// @Success 200 {object} utils.Response{data=models.DepartmentResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /departments/{id} [put]
func (c *DepartmentController) UpdateDepartment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid department ID")
		return
	}

	var req models.DepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	department, err := c.departmentService.UpdateDepartment(uint(id), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Department updated successfully", department.ToResponse())
}

// DeleteDepartment godoc
// @Summary Delete department
// @Description Delete department by ID
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /departments/{id} [delete]
func (c *DepartmentController) DeleteDepartment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid department ID")
		return
	}

	if err := c.departmentService.DeleteDepartment(uint(id)); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Department deleted successfully", nil)
}

// GetActiveDepartments godoc
// @Summary Get active departments
// @Description Get list of all active departments
// @Tags departments
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.DepartmentResponse}
// @Failure 500 {object} utils.Response
// @Router /departments/active [get]
func (c *DepartmentController) GetActiveDepartments(ctx *gin.Context) {
	departments, err := c.departmentService.GetActiveDepartments()
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	var departmentResponses []models.DepartmentResponse
	for _, department := range departments {
		departmentResponses = append(departmentResponses, department.ToResponse())
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Active departments retrieved successfully", departmentResponses)
}