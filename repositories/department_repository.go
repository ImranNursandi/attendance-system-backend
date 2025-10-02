package repositories

import (
	"attendance-system/models"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

func (r *DepartmentRepository) FindAll() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Find(&departments).Error
	return departments, err
}

func (r *DepartmentRepository) FindByID(id int) (models.Department, error) {
	var department models.Department
	err := r.db.Where("id = ?", id).First(&department).Error
	return department, err
}

func (r *DepartmentRepository) Update(department *models.Department) error {
	return r.db.Save(department).Error
}

func (r *DepartmentRepository) Delete(id int) error {
	return r.db.Where("id = ?", id).Delete(&models.Department{}).Error
}