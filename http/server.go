package http

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	queueing "message-queueing"
	"net/http"
	"strconv"
	"time"
)

func NewServer(service queueing.Service) http.Handler {
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

	api := router.Group("/api/v1")
	queueAPI := api.Group("/queues/:queueID")
	queueAPI.POST("/messages", controller.postMessage)
	queueAPI.GET("/messages/available", controller.getAvailableMessages)
	queueAPI.POST("/messages/:messageID/acknowledge", controller.postAcknowledgeMessage)

	return router
}

type QueueMessageController struct {
	service queueing.Service
}

type HttpError struct {
	Status int    `json:"status"`
	Msg    string `json:"error"`
	Base   error  `json:"description"`
}

func (err HttpError) Unwrap() error {
	return err.Base
}

func (err HttpError) Error() string {
	return err.Base.Error()
}

type InputMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

type MessageReceipt struct {
	MessageID uuid.UUID `json:"messageID"`
	DataHash  []byte    `json:"dataHash"`
	Timestamp uint64    `json:"timestamp"`
}

func (controller *QueueMessageController) postMessage(ctx *gin.Context) {
	var body InputMessage
	if err := ctx.ShouldBindJSON(&body); err != nil {
		httpErr := HttpError{
			Status: http.StatusBadRequest,
			Msg:    "bad request",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	messageID := uuid.New()
	f := false
	now := time.Now().UnixMilli()
	hash := md5.New()
	hash.Write(body.Data)
	dataHash := hash.Sum(nil)
	message := queueing.QueueMessage{
		MessageID:     messageID[:],
		Data:          body.Data,
		DataHash:      dataHash,
		Attributes:    body.Attributes,
		Timestamp:     &now,
		LastRetrieved: nil,
		Acknowledged:  &f,
	}

	err := controller.service.Enqueue(ctx, &message)
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusInternalServerError,
			Msg:    "internal server error",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	ctx.JSON(
		http.StatusCreated, MessageReceipt{
			MessageID: messageID,
			DataHash:  dataHash,
			Timestamp: uint64(now),
		},
	)
}

func (controller *QueueMessageController) getAvailableMessages(ctx *gin.Context) {
	messageCount := int64(10)
	var err error
	if ctx.Query("messageCount") != "" {
		messageCount, err = strconv.ParseInt(ctx.Query("messageCount"), 10, 64)
	}
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusBadRequest,
			Msg:    "bad request",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	messages := make([]*queueing.QueueMessage, messageCount)
	actualCount, err := controller.service.Retrieve(ctx, messages)
	if err != nil && !errors.Is(err, queueing.NextMessageUnavailableError) {
		httpErr := HttpError{
			Status: http.StatusInternalServerError,
			Msg:    "internal server error",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	ctx.JSON(http.StatusOK, messages[:actualCount])
}

func (controller *QueueMessageController) postAcknowledgeMessage(ctx *gin.Context) {
	messageID, err := uuid.Parse(ctx.Param("messageID"))
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusBadRequest,
			Msg:    "internal server error",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	err = controller.service.Acknowledge(ctx, messageID)
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusInternalServerError,
			Msg:    "internal server error",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
