package usecase

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"luxe-beb-go/library/types"
	"luxe-beb-go/src/services/category"

	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	validator "gopkg.in/go-playground/validator.v9"
)

type CategoryUsecase struct {
	categoryRepo   category.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewCategoryUsecase(db *sqlx.DB, categoryRepo category.Repository) category.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &CategoryUsecase{
		categoryRepo:   categoryRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *CategoryUsecase) FindAll(ctx *gin.Context, params models.FindAllCategoryParams) ([]*models.Category, *types.Error) {
	result, err := u.categoryRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".CategoryUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *CategoryUsecase) Find(ctx *gin.Context, id string) (*models.Category, *types.Error) {
	result, err := u.categoryRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".CategoryUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *CategoryUsecase) Count(ctx *gin.Context, params models.FindAllCategoryParams) (int, *types.Error) {
	result, err := u.categoryRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".CategoryUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *CategoryUsecase) Create(ctx *gin.Context, obj models.Category) (*models.Category, *types.Error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errValidation := validate.Struct(obj)
	if errValidation != nil {
		return nil, &types.Error{
			Path:       ".CategoryUsecase->Create()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data := models.Category{
		ID:       uuid.New().String(),
		Name:     obj.Name,
		StatusID: models.DEFAULT_STATUS_ID,
	}

	result, err := u.categoryRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".CategoryUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *CategoryUsecase) Update(ctx *gin.Context, id string, obj models.Category) (*models.Category, *types.Error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errValidation := validate.Struct(obj)
	if errValidation != nil {
		return nil, &types.Error{
			Path:       ".CategoryUsecase->Update()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data, err := u.categoryRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".CategoryUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name

	result, err := u.categoryRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".CategoryUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *CategoryUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.categoryRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".CategoryUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *CategoryUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Category, *types.Error) {
	result, err := u.categoryRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".CategoryUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}
