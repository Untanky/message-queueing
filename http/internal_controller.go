package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	queueing "message-queueing"
	"net/http"
	"strconv"
)

type internalController struct {
	repository queueing.Repository
	storage    queueing.BlockStorage
}

type GetManifestResponse struct {
	Blobs []string `json:"files"`
}

func (controller *internalController) getManifest(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK, GetManifestResponse{
			Blobs: []string{"abc"},
		},
	)
}

func (controller *internalController) getFile(ctx *gin.Context) {
	w := ctx.Writer

	header := w.Header()
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Length", strconv.FormatInt(controller.storage.Length(), 10))
	w.WriteHeader(http.StatusOK)
	controller.storage.WriteTo(w)
}

func (controller *internalController) addMessage(ctx *gin.Context) {
	message := new(queueing.QueueMessage)
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusBadRequest,
			Msg:    "could not parse JSON",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	err = controller.repository.Create(ctx, message)
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusInternalServerError,
			Msg:    "could not add to repository",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (controller *internalController) getMessage(ctx *gin.Context) {
	id, ok := ctx.Params.Get("messageID")
	if !ok {
		httpErr := HttpError{
			Status: http.StatusBadRequest,
			Msg:    "could not find messageID",
			Base:   nil,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	messageID, err := uuid.Parse(id)

	message, err := controller.repository.GetByID(ctx, messageID)
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusInternalServerError,
			Msg:    "could not add to repository",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	ctx.JSON(http.StatusOK, message)
}

func (controller *internalController) updateMessage(ctx *gin.Context) {
	message := new(queueing.QueueMessage)
	err := ctx.ShouldBindJSON(message)
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusBadRequest,
			Msg:    "could not parse JSON",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	err = controller.repository.Update(ctx, message)
	if err != nil {
		httpErr := HttpError{
			Status: http.StatusInternalServerError,
			Msg:    "could not add to repository",
			Base:   err,
		}
		ctx.Error(httpErr)
		ctx.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
