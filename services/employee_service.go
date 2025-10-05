package services

import (
	"attendance-system/models"
	"attendance-system/repositories"
	"attendance-system/utils"
	"time"
	"strings"

	"strconv"
	"fmt"
)

type EmployeeService struct {
	employeeRepo *repositories.EmployeeRepository
	userRepo     *repositories.UserRepository
}

func NewEmployeeService() *EmployeeService {
	return &EmployeeService{
		employeeRepo: repositories.NewEmployeeRepository(),
		userRepo:     repositories.NewUserRepository(),
	}
}

func (s *EmployeeService) CreateEmployee(req models.EmployeeRequest) (*models.Employee, error) {
	// Generate employee ID if not provided
	employeeID := req.EmployeeID
	if employeeID == "" {
		var err error
		employeeID, err = s.generateEmployeeID()
		if err != nil {
			return nil, err
		}
	} else {
		// Check if provided employee ID already exists
		existingEmployee, _ := s.employeeRepo.FindByEmployeeID(employeeID)
		if existingEmployee != nil {
			return nil, utils.NewConflictError("Employee ID Already Exists")
		}
	}

	// Check if user email already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, utils.NewConflictError("User Account Already Exists With This Email")
	}

	employee := &models.Employee{
		EmployeeID:   employeeID,
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
		Phone:        req.Phone,
		Address:      req.Address,
		Position:     req.Position,
		Status:       req.Status,
		JoinDate:     req.JoinDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if employee.Status == "" {
		employee.Status = "active"
	}
	if employee.JoinDate.IsZero() {
		employee.JoinDate = time.Now()
	}

	// Create employee
	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, err
	}

	// Create user account for the employee
	if err := s.createUserForEmployee(employee, req.Email); err != nil {
		// If user creation fails, delete the employee to maintain consistency
		s.employeeRepo.Delete(employee.ID)
		return nil, err
	}

	return employee, nil
}

func (s *EmployeeService) generateEmployeeID() (string, error) {
	// Get all employee IDs to find the highest number
	employees, err := s.employeeRepo.FindAllEmployeeIDs()
	if err != nil {
		return "", err
	}

	maxID := 0
	for _, emp := range employees {
		// Extract numeric part from EMP001, EMP002, etc.
		if strings.HasPrefix(emp.EmployeeID, "EMP") {
			if idNum, err := strconv.Atoi(emp.EmployeeID[3:]); err == nil {
				if idNum > maxID {
					maxID = idNum
				}
			}
		}
	}

	nextID := maxID + 1
	return fmt.Sprintf("EMP%03d", nextID), nil
}

func (s *EmployeeService) createUserForEmployee(employee *models.Employee, email string) error {
	// Generate username from email
	username := strings.Split(email, "@")[0]
	
	// Check if username already exists, if so, append employee ID
	existingUserByUsername, _ := s.userRepo.FindByUsername(username)
	if existingUserByUsername != nil {
		username = username + "_" + employee.EmployeeID
	}

	// Create user
	user := &models.User{
		Username:   username,
		Email:      email,
		Password:   "Welcome123", // Default password
		Role:       "employee",   // Default role for employees
		EmployeeID: &employee.EmployeeID,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *EmployeeService) GetAllEmployees(filters []repositories.Filter, search string, page, limit int) ([]models.Employee, *repositories.Pagination, error) {
	employees, pagination, err := s.employeeRepo.FindAll(filters, search, page, limit)
	if err != nil {
		return nil, nil, err
	}

	return employees, pagination, nil
}

func (s *EmployeeService) GetEmployeeByID(id uint) (*models.Employee, error) {
	return s.employeeRepo.FindByID(id)
}

func (s *EmployeeService) GetEmployeeByEmployeeID(employeeID string) (*models.Employee, error) {
	return s.employeeRepo.FindByEmployeeID(employeeID)
}

func (s *EmployeeService) UpdateEmployee(id uint, req models.EmployeeRequest) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if employee ID is being changed and if it already exists
	if employee.EmployeeID != req.EmployeeID {
		existingEmployee, _ := s.employeeRepo.FindByEmployeeID(req.EmployeeID)
		if existingEmployee != nil {
			return nil, utils.NewConflictError("Employee ID Already Exists")
		}
		
		// IMPORTANT: Update employee_id in user table if employee ID changes
		if err := s.updateUserEmployeeID(employee.EmployeeID, req.EmployeeID); err != nil {
			return nil, err
		}
	}

	employee.EmployeeID = req.EmployeeID
	employee.DepartmentID = req.DepartmentID
	employee.Name = req.Name
	employee.Phone = req.Phone
	employee.Address = req.Address
	employee.Position = req.Position
	employee.Status = req.Status
	employee.JoinDate = req.JoinDate
	employee.UpdatedAt = time.Now()

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, err
	}

	return employee, nil
}

// updateUserEmployeeID updates the employee_id in the user table
func (s *EmployeeService) updateUserEmployeeID(oldEmployeeID, newEmployeeID string) error {
	// Check if user exists with the old employee ID
	_, err := s.userRepo.FindByEmployeeID(oldEmployeeID)
	if err != nil {
		// If no user found, it's not an error - just return
		return nil
	}
	
	// Update the employee_id
	updates := map[string]interface{}{
		"employee_id": newEmployeeID,
		"updated_at":  time.Now(),
	}
	
	return s.userRepo.UpdateByEmployeeID(oldEmployeeID, updates)
}

func (s *EmployeeService) DeleteEmployee(id uint) error {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Also delete associated user account
	if employee.EmployeeID != "" {
		if err := s.userRepo.DeleteByEmployeeID(employee.EmployeeID); err != nil {
			// Log the error but continue with employee deletion
			// In production, use a proper logger instead of println
			println("Warning: Failed to delete user account for employee:", err.Error())
		}
	}

	return s.employeeRepo.Delete(id)
}

func (s *EmployeeService) GetEmployeesByDepartment(departmentID uint) ([]models.Employee, error) {
	return s.employeeRepo.FindByDepartment(departmentID)
}

func (s *EmployeeService) GetActiveEmployeesCount() (int64, error) {
	return s.employeeRepo.GetActiveEmployeesCount()
}

func (s *EmployeeService) SearchEmployees(query string, limit int) ([]models.Employee, error) {
	return s.employeeRepo.SearchEmployees(query, limit)
}

// GetEmployeeWithUserEmail gets employee details along with user email
func (s *EmployeeService) GetEmployeeWithUserEmail(employeeID string) (*models.EmployeeWithUserResponse, error) {
	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, err
	}

	// Find associated user
	user, err := s.userRepo.FindByEmployeeID(employeeID)
	if err != nil {
		// If no user found, return employee data without email
		response := &models.EmployeeWithUserResponse{
			Employee:  employee.ToResponse(),
			UserEmail: "",
		}
		return response, nil
	}

	response := &models.EmployeeWithUserResponse{
		Employee:  employee.ToResponse(),
		UserEmail: user.Email,
	}

	return response, nil
}

// GetEmployeeWithUser gets complete employee and user data
func (s *EmployeeService) GetEmployeeWithUser(employeeID string) (*models.Employee, *models.User, error) {
	employee, err := s.employeeRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.FindByEmployeeID(employeeID)
	if err != nil {
		return employee, nil, nil // Return employee even if no user found
	}

	return employee, user, nil
}

// UpdateEmployeeStatus updates only the employee status
func (s *EmployeeService) UpdateEmployeeStatus(id uint, status string) error {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return err
	}

	employee.Status = status
	employee.UpdatedAt = time.Now()

	return s.employeeRepo.Update(employee)
}