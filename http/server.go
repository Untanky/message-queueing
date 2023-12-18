package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	queueing "message-queueing"
	"net/http"
	"time"
)

func NewServer(service queueing.Service, storage queueing.BlockStorage, repository queueing.Repository) http.Handler {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(logger())

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

func logger() gin.HandlerFunc {
	log := slog.Default().With("service", "http-server")

	return func(ctx *gin.Context) {
		start := time.Now()

		requestID := ctx.Request.Header.Get("request-id")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx.Set("request-id", requestID)
		ctx.Header("request-id", requestID)

		ctx.Next()

		duration := time.Now().Sub(start)

		attr := []slog.Attr{
			slog.Int("status", ctx.Writer.Status()),
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.Duration("latency", duration),
			slog.String("host", ctx.ClientIP()),
			slog.Int("bytes", ctx.Writer.Size()),
		}
		if ctx.Writer.Status() < 400 {
			log.InfoContext(ctx, "request succeeded", "request", attr, "requestID", requestID)
		} else {
			log.ErrorContext(ctx, "request failed", "request", attr, "requestID", requestID)
		}
	}
}
