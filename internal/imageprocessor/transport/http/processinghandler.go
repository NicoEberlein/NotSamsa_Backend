package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"image"
	"image/png"
	"io"
	"net/http"
)

type PostProcessingHandlerModel struct {
	GetFromPath string
	SaveToPath  string
}

func (h *Handler) StartProcessingHandler(c *gin.Context) {

	var body PostProcessingHandlerModel

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.S3.StatObject(c, "notsamsa", body.GetFromPath, minio.StatObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	go resizeAndUpload(h.S3, body, 600, 300)

	c.JSON(http.StatusAccepted, gin.H{"accepted": "image accepted"})
}

func resizeAndUpload(client *minio.Client, body PostProcessingHandlerModel, width int, height int) {
	obj, err := client.GetObject(context.Background(), "notsamsa", body.GetFromPath, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, obj); err != nil {
		fmt.Println(err)
	}

	im, _, err := image.Decode(&buf)
	if err != nil {
		fmt.Println(err)
	}

	var resizedBuf bytes.Buffer

	resizedImage := imaging.Fill(im, width, height, imaging.Center, imaging.Lanczos)
	if err = png.Encode(&resizedBuf, resizedImage); err != nil {
		fmt.Println(err)
	}

	if _, err = client.PutObject(
		context.Background(),
		"notsamsa",
		body.SaveToPath,
		&resizedBuf,
		int64(resizedBuf.Len()),
		minio.PutObjectOptions{
			ContentType: fmt.Sprintf("image/png"),
		},
	); err != nil {
		fmt.Println(err)
	}

}
