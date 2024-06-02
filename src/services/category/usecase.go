package category

import (
	"luxe-beb-go/library/types"
	"luxe-beb-go/models"

	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllCategoryParams) ([]*models.Category, *types.Error)
	Find(*gin.Context, string) (*models.Category, *types.Error)
	Count(*gin.Context, models.FindAllCategoryParams) (int, *types.Error)
	Create(*gin.Context, models.Category) (*models.Category, *types.Error)
	Update(*gin.Context, string, models.Category) (*models.Category, *types.Error)

	FindStatus(*gin.Context) ([]*models.Status, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Category, *types.Error)
}
