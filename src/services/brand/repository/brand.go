package repository

import (
	"fmt"
	"net/http"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

type BrandRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewBrandRepository(repository data.GenericStorage, statusRepository data.GenericStorage) BrandRepository {
	return BrandRepository{repository: repository, statusRepository: statusRepository}
}

func (s BrandRepository) FindAll(ctx *gin.Context, params models.FindAllBrandParams) ([]*models.Brand, *types.Error) {
	data := []*models.Brand{}
	bulks := []*models.BrandBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(" AND %s", params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(" AND brands.%s", params.FindAllParams.StatusID)
	}

	if params.FindAllParams.SortBy != "" {
		where += fmt.Sprintf(" ORDER BY %s", params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where += ` LIMIT :limit OFFSET :offset`
	}

	query := fmt.Sprintf(`
  SELECT
    brands.id, brands.name,
    brands.status_id, status.name status_name
  FROM brands
  JOIN status ON status.id = brands.status_id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":     params.FindAllParams.Size,
		"offset":    ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id": params.FindAllParams.StatusID,
	})

	if err != nil {
		return nil, &types.Error{
			Path:       ".BrandStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.Brand{
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

func (s BrandRepository) Find(ctx *gin.Context, id string) (*models.Brand, *types.Error) {
	result := models.Brand{}
	bulks := []*models.BrandBulk{}
	var err error

	query := `
  SELECT
    brands.id, brands.name,
    brands.status_id, status.name status_name
  FROM brands
  JOIN status ON status.id = brands.status_id
  WHERE brands.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".BrandStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Brand{
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
			Path:       ".BrandStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s BrandRepository) Create(ctx *gin.Context, obj *models.Brand) (*models.Brand, *types.Error) {
	data := models.Brand{}
	result, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BrandStorage->Create()",
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
			Path:       ".BrandStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s BrandRepository) Update(ctx *gin.Context, obj *models.Brand) (*models.Brand, *types.Error) {
	data := models.Brand{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BrandStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BrandStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s BrandRepository) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	businessStatus := []*models.Status{}

	err := s.statusRepository.Where(ctx, &businessStatus, "1=1", map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".BrandStorage->FindStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return businessStatus, nil
}

func (s BrandRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Brand, *types.Error) {
	data := models.Brand{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BrandStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BrandStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
