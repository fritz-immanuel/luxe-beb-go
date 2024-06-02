package usecase

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"luxe-beb-go/library/types"
	"luxe-beb-go/src/services/product"

	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	validator "gopkg.in/go-playground/validator.v9"
)

type ProductUsecase struct {
	productRepo    product.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewProductUsecase(db *sqlx.DB, productRepo product.Repository) product.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &ProductUsecase{
		productRepo:    productRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *ProductUsecase) FindAll(ctx *gin.Context, filterFindAllParams models.FindAllProductParams) ([]*models.Product, *types.Error) {
	result, err := u.productRepo.FindAll(ctx, filterFindAllParams)
	if err != nil {
		err.Path = ".ProductUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ProductUsecase) Find(ctx *gin.Context, id string) (*models.Product, *types.Error) {
	result, err := u.productRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ProductUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ProductUsecase) Count(ctx *gin.Context, filterFindAllParams models.FindAllProductParams) (int, *types.Error) {
	result, err := u.productRepo.FindAll(ctx, filterFindAllParams)
	if err != nil {
		err.Path = ".ProductUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *ProductUsecase) Create(ctx *gin.Context, obj models.Product) (*models.Product, *types.Error) {
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
			Path:       ".ProductUsecase->Create()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data := models.Product{
		ID:          uuid.New().String(),
		Code:        obj.Code,
		Name:        obj.Name,
		Price:       obj.Price,
		BrandID:     obj.BrandID,
		CategoryID:  obj.CategoryID,
		Description: obj.Description,
		StatusID:    models.DEFAULT_STATUS_ID,
	}

	result, err := u.productRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".ProductUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ProductUsecase) Update(ctx *gin.Context, id string, obj models.Product) (*models.Product, *types.Error) {
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
			Path:       ".ProductUsecase->Update()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data, err := u.productRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ProductUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Price = obj.Price
	data.BrandID = obj.BrandID
	data.CategoryID = obj.CategoryID
	data.Description = obj.Description

	result, err := u.productRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".ProductUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *ProductUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.productRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".ProductUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ProductUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Product, *types.Error) {
	result, err := u.productRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".ProductUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}
