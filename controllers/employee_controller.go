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

type EmployeeController struct {
	employeeService *services.EmployeeService
}

func NewEmployeeController() *EmployeeController {
	return &EmployeeController{
		employeeService: services.NewEmployeeService(),
	}
}

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Create a new employee with the provided data
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body models.EmployeeRequest true "Employee data"
// @Success 201 {object} utils.Response{data=models.EmployeeResponse}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /employees [post]
func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
	var req models.EmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	employee, err := c.employeeService.CreateEmployee(req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusCreated, "Employee created successfully", employee.ToResponse())
}

// GetAllEmployees godoc
// @Summary Get all employees
// @Description Get paginated list of employees with optional filtering and search
// @Tags employees
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param department_id query int false "Filter by department ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} utils.Response{data=[]models.EmployeeResponse}
// @Failure 500 {object} utils.Response
// @Router /employees [get]
func (c *EmployeeController) GetAllEmployees(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")

	var filters []repositories.Filter
	if departmentID := ctx.Query("department_id"); departmentID != "" {
		if deptID, err := strconv.ParseUint(departmentID, 10, 32); err == nil {
			filters = append(filters, repositories.Filter{Field: "department_id", Value: deptID})
		}
	}
	if status := ctx.Query("status"); status != "" {
		filters = append(filters, repositories.Filter{Field: "status", Value: status})
	}

	employees, pagination, err := c.employeeService.GetAllEmployees(filters, search, page, limit)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Convert to response objects
	var employeeResponses []models.EmployeeResponse
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, employee.ToResponse())
	}

	response := map[string]interface{}{
		"employees":  employeeResponses,
		"pagination": pagination,
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employees retrieved successfully", response)
}

// GetEmployeeByID godoc
// @Summary Get employee by ID
// @Description Get employee details by ID
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} utils.Response{data=models.EmployeeResponse}
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /employees/{id} [get]
func (c *EmployeeController) GetEmployeeByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	employee, err := c.employeeService.GetEmployeeByID(uint(id))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employee retrieved successfully", employee.ToResponse())
}

// GetEmployeeWithUser godoc
// @Summary Get employee with user details
// @Description Get employee details along with user account information
// @Tags employees
// @Accept json
// @Produce json
// @Param employee_id path string true "Employee ID"
// @Success 200 {object} utils.Response{data=models.EmployeeWithUserResponse}
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /employees/{employee_id}/with-user [get]
func (c *EmployeeController) GetEmployeeWithUser(ctx *gin.Context) {
	employeeID := ctx.Param("employee_id")
	if employeeID == "" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Employee ID is required")
		return
	}

	response, err := c.employeeService.GetEmployeeWithUserEmail(employeeID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employee with user details retrieved successfully", response)
}

// UpdateEmployee godoc
// @Summary Update employee
// @Description Update employee details
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body models.EmployeeRequest true "Employee data"
// @Success 200 {object} utils.Response{data=models.EmployeeResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /employees/{id} [put]
func (c *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	var req models.EmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	employee, err := c.employeeService.UpdateEmployee(uint(id), req)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employee updated successfully", employee.ToResponse())
}

// DeleteEmployee godoc
// @Summary Delete employee
// @Description Delete employee by ID
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /employees/{id} [delete]
func (c *EmployeeController) DeleteEmployee(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	if err := c.employeeService.DeleteEmployee(uint(id)); err != nil {
		utils.HandleError(ctx, err)
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employee deleted successfully", nil)
}

// GetEmployeesByDepartment godoc
// @Summary Get employees by department
// @Description Get all employees in a specific department
// @Tags employees
// @Accept json
// @Produce json
// @Param department_id path int true "Department ID"
// @Success 200 {object} utils.Response{data=[]models.EmployeeResponse}
// @Failure 500 {object} utils.Response
// @Router /employees/department/{department_id} [get]
func (c *EmployeeController) GetEmployeesByDepartment(ctx *gin.Context) {
	departmentID, err := strconv.ParseUint(ctx.Param("department_id"), 10, 32)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid department ID")
		return
	}

	employees, err := c.employeeService.GetEmployeesByDepartment(uint(departmentID))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	var employeeResponses []models.EmployeeResponse
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, employee.ToResponse())
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employees retrieved successfully", employeeResponses)
}

// SearchEmployees godoc
// @Summary Search employees
// @Description Search employees by name, employee ID, or position
// @Tags employees
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Maximum results" default(10)
// @Success 200 {object} utils.Response{data=[]models.EmployeeResponse}
// @Failure 500 {object} utils.Response
// @Router /employees/search [get]
func (c *EmployeeController) SearchEmployees(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Search query is required")
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50
	}

	employees, err := c.employeeService.SearchEmployees(query, limit)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	var employeeResponses []models.EmployeeResponse
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, employee.ToResponse())
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employees search completed", employeeResponses)
}