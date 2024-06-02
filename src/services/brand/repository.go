package brand

import (
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

// Repository is the contract between Repository and usecase
type Repository interface {
	FindAll(*gin.Context, models.FindAllBrandParams) ([]*models.Brand, *types.Error)
	Find(*gin.Context, string) (*models.Brand, *types.Error)
	Create(*gin.Context, *models.Brand) (*models.Brand, *types.Error)
	Update(*gin.Context, *models.Brand) (*models.Brand, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Brand, *types.Error)
}
