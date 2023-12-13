package http

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	queueing "message-queueing"
	"net/http"
)

func NewServer(service queueing.Service, storage queueing.BlockStorage) http.Handler {
	router := gin.Default()

	tls.NewListener(
		nil, &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	)

	router.Use(gin.Recovery())

	controller := QueueMessageController{
		service: service,
	}

	internalCtrlr := internalController{
		storage: storage,
	}

	api := router.Group("/api/v1")
	queueAPI := api.Group("/queues/:queueID")
	queueAPI.POST("/messages", controller.postMessage)
	queueAPI.GET("/messages/available", controller.getAvailableMessages)
	queueAPI.POST("/messages/:messageID/acknowledge", controller.postAcknowledgeMessage)

	internal := router.Group("/internal")
	internal.GET("/queue/:queueID/manifest", internalCtrlr.getManifest)
	internal.GET("/queue/:queueID/file/:fileID", internalCtrlr.getFile)
	//internal.POST("/queue/:queueID/messages")
	//internal.GET("/queue/:queueID/messages/available")
	//internal.GET("/queue/:queueID/messages/:messageID")
	//internal.POST("/queue/:queueID/messages/:messageID/acknowledge")

	return router
}
