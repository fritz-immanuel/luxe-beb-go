package repository

import (
	"fmt"
	"net/http"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

type ProductRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewProductRepository(repository data.GenericStorage, statusRepository data.GenericStorage) ProductRepository {
	return ProductRepository{repository: repository, statusRepository: statusRepository}
}

func (s ProductRepository) FindAll(ctx *gin.Context, params models.FindAllProductParams) ([]*models.Product, *types.Error) {
	data := []*models.Product{}
	bulks := []*models.ProductBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where += fmt.Sprintf(` AND %s`, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND products.%s`, params.FindAllParams.StatusID)
	}

	if params.Code != "" {
		where += ` AND products.code = :code`
	}

	if params.Name != "" {
		where += ` AND products.name LIKE ":name%%"`
	}

	if params.BrandID != "" {
		where += ` AND products.brand_id = :brand_id`
	}

	if params.CategoryID != "" {
		where += ` AND products.category_id = :category_id`
	}

	if params.FindAllParams.SortBy != "" {
		where += fmt.Sprintf(` ORDER BY %s`, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where += ` LIMIT :limit OFFSET :offset`
	}

	query := fmt.Sprintf(`
  SELECT
    products.id, products.code, products.name, products.price,
    products.brand_id, products.category_id, products.description,
    products.status_id, status.name status_name
  FROM products
  JOIN product_status ON product_status.id = products.status_id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":       params.FindAllParams.Size,
		"offset":      ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id":   params.FindAllParams.StatusID,
		"code":        params.Code,
		"name":        params.Name,
		"brand_id":    params.BrandID,
		"category_id": params.CategoryID,
	})

	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.Product{
			ID:          v.ID,
			Code:        v.Code,
			Name:        v.Name,
			Price:       v.Price,
			BrandID:     v.BrandID,
			CategoryID:  v.CategoryID,
			Description: v.Description,
			StatusID:    v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}

		data = append(data, obj)
	}

	return data, nil
}

func (s ProductRepository) Find(ctx *gin.Context, id string) (*models.Product, *types.Error) {
	result := models.Product{}
	bulks := []*models.ProductBulk{}
	var err error

	query := `
  SELECT
    products.id, products.code, products.name, products.price,
    products.brand_id, products.category_id, products.description,
    products.status_id, status.name status_name
  FROM products
  JOIN product_status ON product_status.id = products.status_id
  WHERE products.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Product{
			ID:          v.ID,
			Code:        v.Code,
			Name:        v.Name,
			Price:       v.Price,
			BrandID:     v.BrandID,
			CategoryID:  v.CategoryID,
			Description: v.Description,
			StatusID:    v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".ProductStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s ProductRepository) Create(ctx *gin.Context, obj *models.Product) (*models.Product, *types.Error) {
	data := models.Product{}
	result, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Create()",
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
			Path:       ".ProductStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ProductRepository) Update(ctx *gin.Context, obj *models.Product) (*models.Product, *types.Error) {
	data := models.Product{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ProductRepository) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	businessStatus := []*models.Status{}

	err := s.statusRepository.Where(ctx, &businessStatus, "1=1", map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->FindStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return businessStatus, nil
}

func (s ProductRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Product, *types.Error) {
	data := models.Product{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
