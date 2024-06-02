package repository

import (
	"fmt"
	"net/http"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

type CategoryRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewCategoryRepository(repository data.GenericStorage, statusRepository data.GenericStorage) CategoryRepository {
	return CategoryRepository{repository: repository, statusRepository: statusRepository}
}

func (s CategoryRepository) FindAll(ctx *gin.Context, params models.FindAllCategoryParams) ([]*models.Category, *types.Error) {
	data := []*models.Category{}
	bulks := []*models.CategoryBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(" AND %s", params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(" AND categories.%s", params.FindAllParams.StatusID)
	}

	if params.FindAllParams.SortBy != "" {
		where += fmt.Sprintf(" ORDER BY %s", params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where += ` LIMIT :limit OFFSET :offset`
	}

	query := fmt.Sprintf(`
  SELECT
    categories.id, categories.name,
    categories.status_id, status.name status_name
  FROM categories
  JOIN status ON status.id = categories.status_id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":     params.FindAllParams.Size,
		"offset":    ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id": params.FindAllParams.StatusID,
	})

	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.Category{
			ID:       v.ID,
			Name:     v.Name,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}

		data = append(data, obj)
	}

	return data, nil
}

func (s CategoryRepository) Find(ctx *gin.Context, id string) (*models.Category, *types.Error) {
	result := models.Category{}
	bulks := []*models.CategoryBulk{}
	var err error

	query := `
  SELECT
    categories.id, categories.name,
    categories.status_id, status.name status_name
  FROM categories
  JOIN status ON status.id = categories.status_id
  WHERE categories.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Category{
			ID:       v.ID,
			Name:     v.Name,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".CategoryStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s CategoryRepository) Create(ctx *gin.Context, obj *models.Category) (*models.Category, *types.Error) {
	data := models.Category{}
	result, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	lastID, _ := (*result).LastInsertId()
	err = s.repository.FindByID(ctx, &data, lastID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s CategoryRepository) Update(ctx *gin.Context, obj *models.Category) (*models.Category, *types.Error) {
	data := models.Category{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s CategoryRepository) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	businessStatus := []*models.Status{}

	err := s.statusRepository.Where(ctx, &businessStatus, "1=1", map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->FindStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return businessStatus, nil
}

func (s CategoryRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Category, *types.Error) {
	data := models.Category{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".CategoryStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
