package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"

	"luxe-beb-go/library/helpers"
	"luxe-beb-go/middleware"
	"luxe-beb-go/models"
	"luxe-beb-go/src/services/product"
	"luxe-beb-go/src/services/product/repository"
	"luxe-beb-go/src/services/product/usecase"

	"github.com/gin-gonic/gin"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/http/response"
	"luxe-beb-go/library/notif"
	"luxe-beb-go/library/types"
)

var ()

type ProductHandler struct {
	ProductUsecase product.Usecase
	dataManager    *data.Manager
	Result         gin.H
	Status         int
	notifier       *notif.SlackNotifier
}

func (h ProductHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, slackNotifier *notif.SlackNotifier, router *gin.Engine, v *gin.RouterGroup) {
	productRepo := repository.NewProductRepository(
		data.NewMySQLStorage(db, "products", models.Product{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "product_status", models.Status{}, data.MysqlConfig{}),
	)

	uProduct := usecase.NewProductUsecase(db, &productRepo)

	base := &ProductHandler{ProductUsecase: uProduct, dataManager: dataManager, notifier: slackNotifier}

	rs := v.Group("/products")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)
		rs.PUT("/status", middleware.Auth, base.UpdateStatus)
	}

	status := v.Group("/statuses")
	{
		status.GET("/products", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *ProductHandler) FindAll(c *gin.Context) {
	var params models.FindAllProductParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.ProductUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.ProductUsecase.Count(c, params)
	if err != nil {
		err.Path = ".ProductHandler->FindAll()" + err.Path
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

func (h *ProductHandler) Find(c *gin.Context) {
	id := c.Param("id")

	result, err := h.ProductUsecase.Find(c, id)
	if err != nil {
		err.Path = ".ProductHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, h.notifier, "Product not found", http.StatusUnprocessableEntity, *err)
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

func (h *ProductHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.Product
	var data *models.Product

	obj.Name = c.PostForm("Name")
	obj.Price, _ = strconv.ParseFloat(c.PostForm("Price"), 64)
	obj.BrandID = c.PostForm("BrandID")
	obj.CategoryID = c.PostForm("CategoryID")
	obj.Description = c.PostForm("Description")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ProductUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".ProductHandler->Create()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Data created successfuly", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ProductHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.Product
	var data *models.Product

	id := c.Param("id")

	obj.Name = c.PostForm("Name")
	obj.Price, _ = strconv.ParseFloat(c.PostForm("Price"), 64)
	obj.BrandID = c.PostForm("BrandID")
	obj.CategoryID = c.PostForm("CategoryID")
	obj.Description = c.PostForm("Description")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ProductUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".ProductHandler->Update()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Product successfuly updated", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ProductHandler) FindStatus(c *gin.Context) {
	datas, err := h.ProductUsecase.FindStatus(c)
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

func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Product

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".ProductHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.ProductUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".ProductHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Success", StatusCode: http.StatusOK, Message: "Status update success", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
