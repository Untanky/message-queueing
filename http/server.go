package http

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	queueing "message-queueing"
	"net/http"
)

func NewServer(service queueing.Service, storage queueing.BlockStorage, repository queueing.Repository) http.Handler {
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
		storage:    storage,
		repository: repository,
	}

	api := router.Group("/api/v1")
	queueAPI := api.Group("/queues/:queueID")
	queueAPI.POST("/messages", controller.postMessage)
	queueAPI.GET("/messages/available", controller.getAvailableMessages)
	queueAPI.POST("/messages/:messageID/acknowledge", controller.postAcknowledgeMessage)

	internal := router.Group("/internal")
	internal.GET("/queues/:queueID/manifest", internalCtrlr.getManifest)
	internal.GET("/queues/:queueID/blob/:fileID", internalCtrlr.getFile)
	internal.POST("/queues/:queueID/messages", internalCtrlr.addMessage)
	internal.GET("/queues/:queueID/messages/:messageID", internalCtrlr.getMessage)
	internal.PUT("/queues/:queueID/messages/:messageID", internalCtrlr.updateMessage)

	return router
}
