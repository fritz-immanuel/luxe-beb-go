package user

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"

	"luxe-beb-go/library/helpers"
	"luxe-beb-go/middleware"
	"luxe-beb-go/models"
	"luxe-beb-go/src/services/user"
	"luxe-beb-go/src/services/user/repository"
	"luxe-beb-go/src/services/user/usecase"

	"github.com/gin-gonic/gin"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/http/response"
	"luxe-beb-go/library/notif"
	"luxe-beb-go/library/types"
)

var ()

type UserHandler struct {
	UserUsecase user.Usecase
	dataManager *data.Manager
	Result      gin.H
	Status      int
	notifier    *notif.SlackNotifier
}

func (h UserHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, slackNotifier *notif.SlackNotifier, router *gin.Engine, v *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(
		data.NewMySQLStorage(db, "users", models.User{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uUser := usecase.NewUserUsecase(db, &userRepo)

	base := &UserHandler{UserUsecase: uUser, dataManager: dataManager, notifier: slackNotifier}

	rs := v.Group("/users")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)
		rs.PUT("/status", middleware.Auth, base.UpdateStatus)

		rs.POST("auth/login", base.Login)

		//  TO DO
		// - Update Username & Password API
	}

	status := v.Group("/statuses")
	{
		status.GET("/users", middleware.AuthCheckIP, base.FindStatus)
	}
}

func (h *UserHandler) FindAll(c *gin.Context) {
	var params models.FindAllUserParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.UserUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.UserUsecase.Count(c, params)
	if err != nil {
		err.Path = ".UserHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data User Berhasil Ditampilkan", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *UserHandler) Find(c *gin.Context) {
	id := c.Param("id")

	result, err := h.UserUsecase.Find(c, id)
	if err != nil {
		err.Path = ".UserHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, h.notifier, "User not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data User Berhasil Ditampilkan", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) Create(c *gin.Context) {
	var err *types.Error
	var obj models.User
	var data *models.User

	obj.Name = c.PostForm("Name")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.UserUsecase.Create(c, obj)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->Create()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data User Berhasil Ditambahkan", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) Update(c *gin.Context) {
	var err *types.Error
	var obj models.User
	var data *models.User

	id := c.Param("id")

	obj.Name = c.PostForm("Name")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.UserUsecase.Update(c, id, obj)
		if err != nil {
			return err
		}
		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->Update()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data User Berhasil Diperbarui", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) FindStatus(c *gin.Context) {
	datas, err := h.UserUsecase.FindStatus(c)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, h.notifier, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}
	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Data User Status Berhasil Ditampilkan", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.User

	var ids []*models.IDNameTemplate

	newStatusID := c.PostForm("NewStatusID")

	errJson := json.Unmarshal([]byte(c.PostForm("ID")), &ids)
	if errJson != nil {
		err = &types.Error{
			Path:  ".UserHandler->UpdateStatus()",
			Error: errJson,
			Type:  "convert-error",
		}
		response.Error(c, h.notifier, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		for _, id := range ids {
			data, err = h.UserUsecase.UpdateStatus(c, id.ID, newStatusID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, h.notifier, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Status User Berhasil Diperbarui", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

// LOGIN
func (h *UserHandler) Login(c *gin.Context) {
	hash := md5.New()
	io.WriteString(hash, c.PostForm("Password"))

	username := c.PostForm("Username")
	password := fmt.Sprintf("%x", hash.Sum(nil))

	var params models.FindAllUserParams
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	params.Username = username
	params.Password = password
	params.FindAllParams.StatusID = "status_id = 1"

	datas, err := h.UserUsecase.Login(c, params)
	if err != nil {
		c.JSON(401, response.ErrorResponse{
			Code:    "LoginFailed",
			Status:  "Warning",
			Message: "Login Failed",
			Data: &response.DataError{
				Message: err.Message,
				Status:  401,
			},
		})
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Login Berhasil", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

// // //
