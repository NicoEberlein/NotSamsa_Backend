package http

import (
	"errors"
	"fmt"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostCollectionModel struct {
	Name string `json:"name"`
}

func (h *Handler) PostImageCollection(c *gin.Context) {

	var model PostCollectionModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userId := c.GetString("user")
	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is required"})
	}

	user, err := h.UserService.FindById(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	collection := &domain.Collection{
		Owner: user,
		Name:  model.Name,
	}

	id, err := h.CollectionService.Create(c, collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Collection: %+v", collection)

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) GetAllCollectionsOfUser(c *gin.Context) {

	userId := c.GetString("user")
	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
	}

	page, err0 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err0 != nil || err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err0.Error()})
		return
	}

	collections, err := h.CollectionService.FindByUser(c, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Pageate(collections, page, limit))
}

func (h *Handler) DeleteCollection(c *gin.Context) {
	collectionId := c.Param("collectionId")
	if len(collectionId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collectionId is required"})
	}

	err := h.CollectionService.Delete(c, collectionId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}

func (h *Handler) GetCollection(c *gin.Context) {

	id := c.Param("collectionId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}
	collection, err := h.CollectionService.FindById(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, collection)

}
