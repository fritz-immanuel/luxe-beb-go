package repository

import (
	"fmt"
	"net/http"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

type BankRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewBankRepository(repository data.GenericStorage, statusRepository data.GenericStorage) BankRepository {
	return BankRepository{repository: repository, statusRepository: statusRepository}
}

func (s BankRepository) FindAll(ctx *gin.Context, params models.FindAllBankParams) ([]*models.Bank, *types.Error) {
	data := []*models.Bank{}
	bulks := []*models.BankBulk{}

	var err error

	where := `true`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where = fmt.Sprintf("%s AND banks.%s", where, params.FindAllParams.StatusID)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    banks.id, banks.name,
    banks.status_id, status.name status_name
  FROM banks
  JOIN status ON banks.status_id = status.id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":     params.FindAllParams.Size,
		"offset":    ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id": params.FindAllParams.StatusID,
	})

	if err != nil {
		return nil, &types.Error{
			Path:       ".BankStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.Bank{
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

func (s BankRepository) Find(ctx *gin.Context, id string) (*models.Bank, *types.Error) {
	result := models.Bank{}
	bulks := []*models.BankBulk{}
	var err error

	query := `
  SELECT
    banks.id, banks.name,
    banks.status_id, status.name status_name
  FROM banks
  JOIN status ON banks.status_id = status.id
  WHERE banks.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".BankStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Bank{
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
			Path:       ".BankStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s BankRepository) Create(ctx *gin.Context, obj *models.Bank) (*models.Bank, *types.Error) {
	data := models.Bank{}
	result, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BankStorage->Create()",
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
			Path:       ".BankStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s BankRepository) Update(ctx *gin.Context, obj *models.Bank) (*models.Bank, *types.Error) {
	data := models.Bank{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BankStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BankStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s BankRepository) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	businessStatus := []*models.Status{}

	err := s.statusRepository.Where(ctx, &businessStatus, "1=1", map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".BankStorage->FindStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return businessStatus, nil
}

func (s BankRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Bank, *types.Error) {
	data := models.Bank{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BankStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".BankStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
