package category

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"

	"luxe-beb-go/library/helpers"
	"luxe-beb-go/middleware"
	"luxe-beb-go/models"
	"luxe-beb-go/src/services/category"
	"luxe-beb-go/src/services/category/repository"
	"luxe-beb-go/src/services/category/usecase"

	"github.com/gin-gonic/gin"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/http/response"
	"luxe-beb-go/library/notif"
	"luxe-beb-go/library/types"
)

var ()

type CategoryHandler struct {
	CategoryUsecase category.Usecase
	dataManager     *data.Manager
	Result          gin.H
	Status          int
	notifier        *notif.SlackNotifier
}

func (h CategoryHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, slackNotifier *notif.SlackNotifier, router *gin.Engine, v *gin.RouterGroup) {
	categoryRepo := repository.NewCategoryRepository(
		data.NewMySQLStorage(db, "categories", models.Category{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uCategory := usecase.NewCategoryUsecase(db, &categoryRepo)

	base := &CategoryHandler{CategoryUsecase: uCategory, dataManager: dataManager, notifier: slackNotifier}

	rs := v.Group("/categories")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)
		rs.PUT("/status", middleware.Auth, base.UpdateStatus)
	}

	status := v.Group("/statuses")
	{
		status.GET("/categories", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *CategoryHandler) FindAll(c *gin.Context) {
	var params models.FindAllCategoryParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.CategoryUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.CategoryUsecase.Count(c, params)
	if err != nil {
		err.Path = ".CategoryHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Success", StatusCode: http.StatusOK, Message: "Data shown successfuly", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *CategoryHandler) Find(c *gin.Context) {
	id := c.Param("id")

	result, err := h.CategoryUsecase.Find(c, id)
	if err != nil {
		err.Path = ".CategoryHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, h.notifier, "Category not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data shown successfuly", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.Category
	var data *models.Category

	obj.Name = c.PostForm("Name")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.CategoryUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".CategoryHandler->Create()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data created successfuly", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.Category
	var data *models.Category

	id := c.Param("id")

	obj.Name = c.PostForm("Name")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.CategoryUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".CategoryHandler->Update()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Category successfuly updated", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *CategoryHandler) FindStatus(c *gin.Context) {
	datas, err := h.CategoryUsecase.FindStatus(c)
	if err != nil {
		response.Error(c, h.notifier, err.Message, http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data successfuly shown", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *CategoryHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Category

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".CategoryHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.CategoryUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".CategoryHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Status update success", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
