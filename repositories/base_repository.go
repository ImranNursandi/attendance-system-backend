package repositories

import (
	"attendance-system/config"
	"attendance-system/utils"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{
		DB: config.DB,
	}
}

type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type Filter struct {
	Field string
	Value interface{}
}

func (r *BaseRepository) Paginate(query *gorm.DB, page, limit int, result interface{}) (*Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Clone the query for counting
	countQuery := query.Session(&gorm.Session{NewDB: true})
	var total int64
	if err := countQuery.Model(result).Count(&total).Error; err != nil {
		return nil, err
	}

	err := query.Offset(offset).Limit(limit).Find(result).Error
	if err != nil {
		return nil, err
	}

	totalPage := (int(total) + limit - 1) / limit

	return &Pagination{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

func (r *BaseRepository) ApplyFilters(query *gorm.DB, filters []Filter) *gorm.DB {
	for _, filter := range filters {
		if filter.Value != nil && filter.Value != "" {
			if strings.Contains(filter.Field, ".") {
				// Relation filter
				parts := strings.Split(filter.Field, ".")
				if len(parts) == 2 {
					query = query.Joins(parts[0]).Where(parts[0]+"."+parts[1]+" = ?", filter.Value)
				}
			} else {
				// Direct field filter
				query = query.Where(filter.Field+" = ?", filter.Value)
			}
		}
	}
	return query
}

func (r *BaseRepository) ApplySearch(query *gorm.DB, search string, searchFields []string) *gorm.DB {
	if search != "" && len(searchFields) > 0 {
		searchPattern := "%" + strings.ToLower(search) + "%"
		conditions := []string{}
		args := []interface{}{}

		for _, field := range searchFields {
			conditions = append(conditions, "LOWER("+field+") LIKE ?")
			args = append(args, searchPattern)
		}

		if len(conditions) > 0 {
			query = query.Where(strings.Join(conditions, " OR "), args...)
		}
	}
	return query
}

func (r *BaseRepository) HandleError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.NewNotFoundError("record not found")
	}
	if strings.Contains(err.Error(), "Duplicate entry") {
		return utils.NewConflictError("duplicate entry")
	}
	if strings.Contains(err.Error(), "foreign key constraint") {
		return utils.NewBadRequestError("foreign key constraint violation")
	}
	return utils.NewInternalServerError(err.Error())
}