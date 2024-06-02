package bank

import (
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllBankParams) ([]*models.Bank, *types.Error)
	Find(*gin.Context, string) (*models.Bank, *types.Error)
	Count(*gin.Context, models.FindAllBankParams) (int, *types.Error)
	Create(*gin.Context, models.Bank) (*models.Bank, *types.Error)
	Update(*gin.Context, string, models.Bank) (*models.Bank, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Bank, *types.Error)
}
