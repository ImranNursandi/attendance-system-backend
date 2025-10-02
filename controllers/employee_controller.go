package controllers

import (
	"attendance-system/models"
	"attendance-system/repositories"
	"attendance-system/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	repo *repositories.EmployeeRepository
}

func NewEmployeeController(repo *repositories.EmployeeRepository) *EmployeeController {
	return &EmployeeController{repo: repo}
}

func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
	var employee models.Employee
	if err := ctx.ShouldBindJSON(&employee); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	if err := c.repo.Create(&employee); err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to create employee")
		return
	}

	utils.SuccessJSON(ctx, http.StatusCreated, "Employee created successfully", employee)
}

func (c *EmployeeController) GetEmployees(ctx *gin.Context) {
	employees, err := c.repo.FindAll()
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to fetch employees")
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employees fetched successfully", employees)
}

func (c *EmployeeController) GetEmployee(ctx *gin.Context) {
	id := ctx.Param("id")
	employee, err := c.repo.FindByID(id)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Employee not found")
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employee fetched successfully", employee)
}

func (c *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	id := ctx.Param("id")
	var employee models.Employee
	if err := ctx.ShouldBindJSON(&employee); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	existingEmployee, err := c.repo.FindByID(id)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusNotFound, "Employee not found")
		return
	}

	employee.ID = existingEmployee.ID
	employee.UpdatedAt = time.Now()

	if err := c.repo.Update(&employee); err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to update employee")
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employee updated successfully", employee)
}

func (c *EmployeeController) DeleteEmployee(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.repo.Delete(id); err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to delete employee")
		return
	}

	utils.SuccessJSON(ctx, http.StatusOK, "Employee deleted successfully", nil)
}