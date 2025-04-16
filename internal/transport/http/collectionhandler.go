package http

import (
	"errors"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type CollectionHandler struct {
	CollectionService *service.ImageCollectionService
}

func NewCollectionHandler(collectionService *service.ImageCollectionService) *CollectionHandler {
	return &CollectionHandler{
		CollectionService: collectionService,
	}
}

type PostCollectionModel struct {
	OwnerId string `json:"ownerId"`
	Name    string `json:"name"`
}

func (collectionHandler *CollectionHandler) PostImageCollection(c *gin.Context) {

	var model PostCollectionModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	collection := &domain.ImageCollection{
		OwnerId: model.OwnerId,
		Name:    model.Name,
	}

	id, err := collectionHandler.CollectionService.Create(c, collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (collectionHandler *CollectionHandler) GetAllCollectionsOfUser(c *gin.Context) {

	userId := c.Param("userId")
	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
	}

	page, err0 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err0 != nil || err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err0.Error()})
		return
	}

	collections, err := collectionHandler.CollectionService.FindByUser(c, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Pageate(collections, page, limit))
}

func (collectionHandler *CollectionHandler) DeleteCollection(c *gin.Context) {
	collectionId := c.Param("collectionId")
	if len(collectionId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collectionId is required"})
	}

	err := collectionHandler.CollectionService.Delete(c, collectionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}

func (collectionHandler *CollectionHandler) GetCollection(c *gin.Context) {

	id := c.Param("collectionId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}
	collection, err := collectionHandler.CollectionService.FindById(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, collection)

}
