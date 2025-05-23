package http

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/gin-gonic/gin"
	"image"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type CreatePreviewBody struct {
	GetFromPath string
	SaveToPath  string
	Width       int
	Height      int
}

func (h *Handler) UploadImage(c *gin.Context) {
	collectionId := c.Param("collectionId")
	if len(collectionId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection Id required"})
		return
	}
	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["images"]

	for _, file := range files {
		f, err := file.Open()
		defer f.Close()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var buf bytes.Buffer
		_, err = io.Copy(&buf, f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		_, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "im must be of format jpeg, png or gif"})
			return
		}

		imageModel := domain.NewImage(collectionId, format, file.Filename, int64(buf.Len()), time.Now(), &buf)

		if err = h.ImageService.Create(c, imageModel); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// Send create review request
		_, err = h.RestClient.R().
			SetContext(c).
			SetBody(&CreatePreviewBody{
				GetFromPath: imageModel.Path,
				SaveToPath:  strings.Replace(imageModel.Path, "images", "previews", 1),
				Width:       600,
				Height:      300,
			}).Post("/generate-preview")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "accepted"})
	}
}

func (h *Handler) CreateDownloadImage(preview bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		imageId := c.Param("imageId")
		if len(imageId) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image Id required"})
		}

		im, err := h.ImageService.FindById(c, imageId, preview)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := im.ImageBinary.Bytes()

		c.Header("Content-Disposition", "attachment; filename="+im.Name)
		c.Header("Content-Type", fmt.Sprintf("image/%s", im.Format))
		c.Header("Content-Length", fmt.Sprintf("%d", len(data)))

		_, err = c.Writer.Write(data)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		c.Writer.WriteHeader(http.StatusOK)

	}
}

func (h *Handler) DeleteImage(c *gin.Context) {

	imageId := c.Param("imageId")
	if len(imageId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image Id required"})
	}

	err := h.ImageService.Delete(c, imageId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})

}

func (h *Handler) GetImagesByCollection(c *gin.Context) {

	collectionId := c.Param("collectionId")
	page, err0 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err0 != nil || err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err0.Error()})
		return
	}

	if len(collectionId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection Id required"})
	}

	images, err := h.ImageService.FindByCollection(c, collectionId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, Pageate(images, page, limit))
}
