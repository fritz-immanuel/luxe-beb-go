package routes

import (
	"luxe-beb-go/src/app/businessweb"

	"github.com/gin-gonic/gin"

	"luxe-beb-go/library/data"
	"luxe-beb-go/library/notif"

	"github.com/jmoiron/sqlx"
)

// RegisterWebRoutes  is a function to register all WEB Routes in the projectbase
func RegisterWebRoutes(db *sqlx.DB, dataManager *data.Manager, slackNotifier *notif.SlackNotifier, router *gin.Engine) {
	v1 := router.Group("/web/v1")
	{
		businessweb.RegisterRoutes(db, dataManager, slackNotifier, router, v1)
	}
}
