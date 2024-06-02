package product

import (
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllProductParams) ([]*models.Product, *types.Error)
	Find(*gin.Context, string) (*models.Product, *types.Error)
	Count(*gin.Context, models.FindAllProductParams) (int, *types.Error)
	Create(*gin.Context, models.Product) (*models.Product, *types.Error)
	Update(*gin.Context, string, models.Product) (*models.Product, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Product, *types.Error)
}
