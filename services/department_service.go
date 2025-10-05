package services

import (
	"attendance-system/models"
	"attendance-system/repositories"
	"time"
)

type DepartmentService struct {
	departmentRepo *repositories.DepartmentRepository
}

// DepartmentRepo getter untuk akses dari controller
func (s *DepartmentService) DepartmentRepo() *repositories.DepartmentRepository {
	return s.departmentRepo
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		departmentRepo: repositories.NewDepartmentRepository(),
	}
}

func (s *DepartmentService) CreateDepartment(req models.DepartmentRequest) (*models.Department, error) {
	department := &models.Department{
		Name:              req.Name,
		Description:       req.Description,
		MaxClockIn:        req.MaxClockIn,
		MaxClockOut:       req.MaxClockOut,
		LateTolerance:     req.LateTolerance,
		EarlyLeavePenalty: req.EarlyLeavePenalty,
		Status:            req.Status,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if department.Status == "" {
		department.Status = "active"
	}
	if department.LateTolerance == 0 {
		department.LateTolerance = 15 // Default 15 minutes
	}
	if department.EarlyLeavePenalty == 0 {
		department.EarlyLeavePenalty = 30 // Default 30 minutes
	}

	if err := s.departmentRepo.Create(department); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) GetAllDepartments(filters []repositories.Filter, search string, page, limit int) ([]models.Department, *repositories.Pagination, error) {
	return s.departmentRepo.FindAll(filters, search, page, limit)
}

func (s *DepartmentService) GetDepartmentByID(id uint) (*models.Department, error) {
	return s.departmentRepo.FindByID(id)
}

func (s *DepartmentService) UpdateDepartment(id uint, req models.DepartmentRequest) (*models.Department, error) {
	department, err := s.departmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	department.Name = req.Name
	department.Description = req.Description
	department.MaxClockIn = req.MaxClockIn
	department.MaxClockOut = req.MaxClockOut
	department.LateTolerance = req.LateTolerance
	department.EarlyLeavePenalty = req.EarlyLeavePenalty
	department.Status = req.Status
	department.UpdatedAt = time.Now()

	if err := s.departmentRepo.Update(department); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) DeleteDepartment(id uint) error {
	return s.departmentRepo.Delete(id)
}

func (s *DepartmentService) GetDepartmentWithEmployees(id uint) (*models.Department, error) {
	return s.departmentRepo.GetDepartmentWithEmployees(id)
}

func (s *DepartmentService) GetActiveDepartments() ([]models.Department, error) {
	return s.departmentRepo.GetActiveDepartments()
}