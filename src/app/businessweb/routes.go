package businessweb

import (
	http_bank "luxe-beb-go/src/app/businessweb/bank"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/notif"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	bankHandler http_bank.BankHandler
)

func RegisterRoutes(db *sqlx.DB, dataManager *data.Manager, slackNotifier *notif.SlackNotifier, router *gin.Engine, v *gin.RouterGroup) {
	v1 := v.Group("")
	{
		bankHandler.RegisterAPI(db, dataManager, slackNotifier, router, v1)
	}
}
