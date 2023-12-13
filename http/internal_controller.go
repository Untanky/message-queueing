package http

import (
	"github.com/gin-gonic/gin"
	queueing "message-queueing"
	"net/http"
	"strconv"
)

type internalController struct {
	storage queueing.BlockStorage
}

type GetManifestResponse struct {
	Files []string `json:"files"`
}

func (controller *internalController) getManifest(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK, GetManifestResponse{
			Files: []string{"abc"},
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
