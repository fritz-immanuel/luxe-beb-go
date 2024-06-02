package brand

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"

	"luxe-beb-go/library/helpers"
	"luxe-beb-go/middleware"
	"luxe-beb-go/models"
	"luxe-beb-go/src/services/brand"
	"luxe-beb-go/src/services/brand/repository"
	"luxe-beb-go/src/services/brand/usecase"

	"github.com/gin-gonic/gin"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/http/response"
	"luxe-beb-go/library/notif"
	"luxe-beb-go/library/types"
)

var ()

type BrandHandler struct {
	BrandUsecase brand.Usecase
	dataManager  *data.Manager
	Result       gin.H
	Status       int
	notifier     *notif.SlackNotifier
}

func (h BrandHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, slackNotifier *notif.SlackNotifier, router *gin.Engine, v *gin.RouterGroup) {
	brandRepo := repository.NewBrandRepository(
		data.NewMySQLStorage(db, "brands", models.Brand{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uBrand := usecase.NewBrandUsecase(db, &brandRepo)

	base := &BrandHandler{BrandUsecase: uBrand, dataManager: dataManager, notifier: slackNotifier}

	rs := v.Group("/brands")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)
		rs.PUT("/status", middleware.Auth, base.UpdateStatus)
	}

	status := v.Group("/statuses")
	{
		status.GET("/brands", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *BrandHandler) FindAll(c *gin.Context) {
	var params models.FindAllBrandParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.BrandUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.BrandUsecase.Count(c, params)
	if err != nil {
		err.Path = ".BrandHandler->FindAll()" + err.Path
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

func (h *BrandHandler) Find(c *gin.Context) {
	id := c.Param("id")

	result, err := h.BrandUsecase.Find(c, id)
	if err != nil {
		err.Path = ".BrandHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, h.notifier, "Brand not found", http.StatusUnprocessableEntity, *err)
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

func (h *BrandHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.Brand
	var data *models.Brand

	obj.Name = c.PostForm("Name")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.BrandUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".BrandHandler->Create()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data created successfuly", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BrandHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.Brand
	var data *models.Brand

	id := c.Param("id")

	obj.Name = c.PostForm("Name")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.BrandUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".BrandHandler->Update()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Brand successfuly updated", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BrandHandler) FindStatus(c *gin.Context) {
	datas, err := h.BrandUsecase.FindStatus(c)
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

func (h *BrandHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Brand

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".BrandHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.BrandUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".BrandHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Status update success", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
