package usecase

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"luxe-beb-go/library/types"
	"luxe-beb-go/src/services/bank"

	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
	validator "gopkg.in/go-playground/validator.v9"
)

type BankUsecase struct {
	bankRepo       bank.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewBankUsecase(db *sqlx.DB, bankRepo bank.Repository) bank.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &BankUsecase{
		bankRepo:       bankRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *BankUsecase) FindAll(ctx *gin.Context, filterFindAllParams models.FindAllBankParams) ([]*models.Bank, *types.Error) {
	result, err := u.bankRepo.FindAll(ctx, filterFindAllParams)
	if err != nil {
		err.Path = ".BankUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BankUsecase) Find(ctx *gin.Context, id string) (*models.Bank, *types.Error) {
	result, err := u.bankRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BankUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BankUsecase) Count(ctx *gin.Context, filterFindAllParams models.FindAllBankParams) (int, *types.Error) {
	result, err := u.bankRepo.FindAll(ctx, filterFindAllParams)
	if err != nil {
		err.Path = ".BankUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *BankUsecase) Create(ctx *gin.Context, obj models.Bank) (*models.Bank, *types.Error) {
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
			Path:       ".BankUsecase->Create()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data := models.Bank{
		ID:       uuid.New().String(),
		Name:     obj.Name,
		StatusID: models.DEFAULT_STATUS_ID,
	}

	result, err := u.bankRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".BankUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BankUsecase) Update(ctx *gin.Context, id string, obj models.Bank) (*models.Bank, *types.Error) {
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
			Path:       ".BankUsecase->Update()",
			Message:    errValidation.Error(),
			Error:      errValidation,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "validation-error",
		}
	}

	data, err := u.bankRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".BankUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name

	result, err := u.bankRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".BankUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *BankUsecase) FindStatus(ctx *gin.Context) ([]*models.Status, *types.Error) {
	result, err := u.bankRepo.FindStatus(ctx)
	if err != nil {
		err.Path = ".BankUsecase->FindStatus()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *BankUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Bank, *types.Error) {
	result, err := u.bankRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".BankUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}
