package repositories

import (
	"attendance-system/models"
	"attendance-system/utils"
	"errors"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	BaseRepository
}

func NewEmployeeRepository() *EmployeeRepository {
	return &EmployeeRepository{
		BaseRepository: *NewBaseRepository(),
	}
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
	if err := r.DB.Create(employee).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

func (r *EmployeeRepository) FindAll(filters []Filter, search string, page, limit int) ([]models.Employee, *Pagination, error) {
	var employees []models.Employee

	query := r.DB.Preload("Department")
	query = r.ApplyFilters(query, filters)
	query = r.ApplySearch(query, search, []string{"name", "employee_id", "position"})

	pagination, err := r.Paginate(query, page, limit, &employees)
	if err != nil {
		return nil, nil, r.HandleError(err)
	}

	return employees, pagination, nil
}

func (r *EmployeeRepository) FindByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.DB.Preload("Department").First(&employee, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("employee not found")
		}
		return nil, r.HandleError(err)
	}
	return &employee, nil
}

func (r *EmployeeRepository) FindByEmail(email string) (*models.Employee, error) {
	var employee models.Employee
	err := r.DB.Where("email = ?", email).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) FindByEmployeeID(employeeID string) (*models.Employee, error) {
	var employee models.Employee
	err := r.DB.Preload("Department").Where("employee_id = ?", employeeID).First(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("employee not found")
		}
		return nil, r.HandleError(err)
	}
	return &employee, nil
}

func (r *EmployeeRepository) FindByDepartment(departmentID uint) ([]models.Employee, error) {
	var employees []models.Employee
	err := r.DB.Preload("Department").Where("department_id = ?", departmentID).Find(&employees).Error
	if err != nil {
		return nil, r.HandleError(err)
	}
	return employees, nil
}

func (r *EmployeeRepository) Update(employee *models.Employee) error {
	if err := r.DB.Save(employee).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

func (r *EmployeeRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Employee{}, id)
	if result.Error != nil {
		return r.HandleError(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("employee not found")
	}
	return nil
}

func (r *EmployeeRepository) GetActiveEmployeesCount() (int64, error) {
	var count int64
	err := r.DB.Model(&models.Employee{}).Where("status = ?", "active").Count(&count).Error
	if err != nil {
		return 0, r.HandleError(err)
	}
	return count, nil
}

func (r *EmployeeRepository) SearchEmployees(query string, limit int) ([]models.Employee, error) {
	var employees []models.Employee
	searchPattern := "%" + query + "%"
	
	err := r.DB.Preload("Department").
		Where("name LIKE ? OR employee_id LIKE ? OR email LIKE ?", searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Find(&employees).Error
	if err != nil {
		return nil, r.HandleError(err)
	}
	return employees, nil
}

// FindAllEmployeeIDs gets all employee IDs for ID generation
func (r *EmployeeRepository) FindAllEmployeeIDs() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.DB.Select("employee_id").Find(&employees).Error
	if err != nil {
		return nil, r.HandleError(err)
	}
	return employees, nil
}

func (r *EmployeeRepository) FindLatestEmployee() (*models.Employee, error) {
	var employee models.Employee
	err := r.DB.Order("employee_id DESC").First(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No employees yet
		}
		return nil, r.HandleError(err)
	}
	return &employee, nil
}