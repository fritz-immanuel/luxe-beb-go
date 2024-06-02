package usecase

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"luxe-beb-go/library/types"
	"luxe-beb-go/src/services/brand"

	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	validator "gopkg.in/go-playground/validator.v9"
)

type BrandUsecase struct {
	brandRepo      brand.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewBrandUsecase(db *sqlx.DB, brandRepo brand.Repository) brand.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &BrandUsecase{
		brandRepo:      brandRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *BrandUsecase) FindAll(ctx *gin.Context, params models.FindAllBrandParams) ([]*models.Brand, *types.Error) {
	result, err := u.brandRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".BrandUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BrandUsecase) Find(ctx *gin.Context, id string) (*models.Brand, *types.Error) {
	result, err := u.brandRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BrandUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BrandUsecase) Count(ctx *gin.Context, params models.FindAllBrandParams) (int, *types.Error) {
	result, err := u.brandRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".BrandUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *BrandUsecase) Create(ctx *gin.Context, obj models.Brand) (*models.Brand, *types.Error) {
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
			Path:       ".BrandUsecase->Create()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data := models.Brand{
		ID:       uuid.New().String(),
		Name:     obj.Name,
		StatusID: models.DEFAULT_STATUS_ID,
	}

	result, err := u.brandRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".BrandUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BrandUsecase) Update(ctx *gin.Context, id string, obj models.Brand) (*models.Brand, *types.Error) {
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
			Path:       ".BrandUsecase->Update()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data, err := u.brandRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BrandUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name

	result, err := u.brandRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".BrandUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *BrandUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.brandRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".BrandUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BrandUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Brand, *types.Error) {
	result, err := u.brandRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".BrandUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}
