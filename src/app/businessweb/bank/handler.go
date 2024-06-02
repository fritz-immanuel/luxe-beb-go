package bank

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"

	"luxe-beb-go/library/helpers"
	"luxe-beb-go/middleware"
	"luxe-beb-go/models"
	"luxe-beb-go/src/services/bank"
	"luxe-beb-go/src/services/bank/repository"
	"luxe-beb-go/src/services/bank/usecase"

	"github.com/gin-gonic/gin"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/http/response"
	"luxe-beb-go/library/notif"
	"luxe-beb-go/library/types"
)

var (
	strToDateFormat      = "2006-01-02"
	strToTimestampFormat = "2006-01-02 15:04:05"
)

type BankHandler struct {
	BankUsecase bank.Usecase
	dataManager *data.Manager
	Result      gin.H
	Status      int
	notifier    *notif.SlackNotifier
}

func (h BankHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, slackNotifier *notif.SlackNotifier, router *gin.Engine, v *gin.RouterGroup) {
	bankRepo := repository.NewBankRepository(
		data.NewMySQLStorage(db, "banks", models.Bank{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uBank := usecase.NewBankUsecase(db, &bankRepo)

	base := &BankHandler{BankUsecase: uBank, dataManager: dataManager, notifier: slackNotifier}

	rs := v.Group("/banks")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)
		rs.PUT("/status", middleware.Auth, base.UpdateStatus)
	}

	status := v.Group("/statuses")
	{
		status.GET("/banks", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *BankHandler) FindAll(c *gin.Context) {
	var params models.FindAllBankParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.BankUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.BankUsecase.Count(c, params)
	if err != nil {
		err.Path = ".BankHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Bank Berhasil Ditampilkan", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *BankHandler) Find(c *gin.Context) {
	id := c.Param("id")

	result, err := h.BankUsecase.Find(c, id)
	if err != nil {
		err.Path = ".BankHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, h.notifier, "Bank not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Bank Berhasil Ditampilkan", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BankHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.Bank
	var data *models.Bank

	obj.Name = c.PostForm("Name")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.BankUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".BankHandler->Create()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Bank Berhasil Ditambahkan", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BankHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.Bank
	var data *models.Bank

	id := c.Param("id")

	obj.Name = c.PostForm("Name")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.BankUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}
		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".BankHandler->Update()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Bank Berhasil Diperbarui", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *BankHandler) FindStatus(c *gin.Context) {
	datas, err := h.BankUsecase.FindStatus(c)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}
	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data Bank Status Berhasil Ditampilkan", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *BankHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Bank

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".BankHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.BankUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".BankHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Status Bank Berhasil Diperbarui", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
