package repositories

import (
	"attendance-system/models"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *EmployeeRepository) FindAll() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.db.Preload("Department").Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) FindByID(id string) (models.Employee, error) {
	var employee models.Employee
	err := r.db.Preload("Department").Where("id = ?", id).First(&employee).Error
	return employee, err
}

func (r *EmployeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *EmployeeRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Employee{}).Error
}