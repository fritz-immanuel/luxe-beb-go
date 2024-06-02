package category

import (
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

// Repository is the contract between Repository and usecase
type Repository interface {
	FindAll(*gin.Context, models.FindAllCategoryParams) ([]*models.Category, *types.Error)
	Find(*gin.Context, string) (*models.Category, *types.Error)
	Create(*gin.Context, *models.Category) (*models.Category, *types.Error)
	Update(*gin.Context, *models.Category) (*models.Category, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Category, *types.Error)
}
