package repositories

import (
	"attendance-system/models"
	"attendance-system/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: *NewBaseRepository(),
	}
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

func (r *UserRepository) FindAll(filters []Filter, search string, page, limit int) ([]models.User, *Pagination, error) {
	var users []models.User

	query := r.DB.Preload("Employee.Department")
	query = r.ApplyFilters(query, filters)
	query = r.ApplySearch(query, search, []string{"username", "email"})

	pagination, err := r.Paginate(query, page, limit, &users)
	if err != nil {
		return nil, nil, r.HandleError(err)
	}

	return users, pagination, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Employee.Department").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("user not found")
		}
		return nil, r.HandleError(err)
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Employee.Department").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("user not found")
		}
		return nil, r.HandleError(err)
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Employee.Department").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("user not found")
		}
		return nil, r.HandleError(err)
	}
	return &user, nil
}

// FindByEmployeeID finds a user by their employee ID
func (r *UserRepository) FindByEmployeeID(employeeID string) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Employee.Department").Where("employee_id = ?", employeeID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("user not found for this employee")
		}
		return nil, r.HandleError(err)
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return r.HandleError(err)
	}
	return nil
}

// UpdateByEmployeeID updates a user by their employee ID
func (r *UserRepository) UpdateByEmployeeID(employeeID string, updates map[string]interface{}) error {
	result := r.DB.Model(&models.User{}).Where("employee_id = ?", employeeID).Updates(updates)
	if result.Error != nil {
		return r.HandleError(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("user not found for this employee")
	}
	return nil
}

func (r *UserRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return r.HandleError(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("user not found")
	}
	return nil
}

// DeleteByEmployeeID deletes a user by their employee ID
func (r *UserRepository) DeleteByEmployeeID(employeeID string) error {
	result := r.DB.Where("employee_id = ?", employeeID).Delete(&models.User{})
	if result.Error != nil {
		return r.HandleError(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("user not found for this employee")
	}
	return nil
}

func (r *UserRepository) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return r.DB.Model(&models.User{}).Where("id = ?", userID).Update("last_login", &now).Error
}

func (r *UserRepository) CountByRole(role string) (int64, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("role = ?", role).Count(&count).Error
	if err != nil {
		return 0, r.HandleError(err)
	}
	return count, nil
}

// FindUsersByEmployeeIDs finds multiple users by their employee IDs
func (r *UserRepository) FindUsersByEmployeeIDs(employeeIDs []string) ([]models.User, error) {
	var users []models.User
	err := r.DB.Preload("Employee.Department").Where("employee_id IN ?", employeeIDs).Find(&users).Error
	if err != nil {
		return nil, r.HandleError(err)
	}
	return users, nil
}

// CheckEmployeeIDExists checks if an employee ID is already assigned to any user
func (r *UserRepository) CheckEmployeeIDExists(employeeID string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("employee_id = ?", employeeID).Count(&count).Error
	if err != nil {
		return false, r.HandleError(err)
	}
	return count > 0, nil
}