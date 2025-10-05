package repositories

import (
	"attendance-system/models"
	"attendance-system/utils"
	"errors"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	BaseRepository
}

func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{
		BaseRepository: *NewBaseRepository(),
	}
}

func (r *DepartmentRepository) Create(department *models.Department) error {
	if err := r.DB.Create(department).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

func (r *DepartmentRepository) FindAll(filters []Filter, search string, page, limit int) ([]models.Department, *Pagination, error) {
	var departments []models.Department

	query := r.DB.Model(&models.Department{})
	query = r.ApplyFilters(query, filters)
	query = r.ApplySearch(query, search, []string{"name", "description"})

	pagination, err := r.Paginate(query, page, limit, &departments)
	if err != nil {
		return nil, nil, r.HandleError(err)
	}

	return departments, pagination, nil
}

func (r *DepartmentRepository) FindByID(id uint) (*models.Department, error) {
	var department models.Department
	err := r.DB.First(&department, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("department not found")
		}
		return nil, r.HandleError(err)
	}
	return &department, nil
}

func (r *DepartmentRepository) Update(department *models.Department) error {
	if err := r.DB.Save(department).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

func (r *DepartmentRepository) Delete(id uint) error {
	// Check if department has employees
	var employeeCount int64
	r.DB.Model(&models.Employee{}).Where("department_id = ?", id).Count(&employeeCount)
	if employeeCount > 0 {
		return utils.NewConflictError("cannot delete department with existing employees")
	}

	result := r.DB.Delete(&models.Department{}, id)
	if result.Error != nil {
		return r.HandleError(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("department not found")
	}
	return nil
}

func (r *DepartmentRepository) GetDepartmentWithEmployees(id uint) (*models.Department, error) {
	var department models.Department
	err := r.DB.Preload("Employees").First(&department, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("department not found")
		}
		return nil, r.HandleError(err)
	}
	return &department, nil
}

func (r *DepartmentRepository) GetActiveDepartments() ([]models.Department, error) {
	var departments []models.Department
	err := r.DB.Where("status = ?", "active").Find(&departments).Error
	if err != nil {
		return nil, r.HandleError(err)
	}
	return departments, nil
}