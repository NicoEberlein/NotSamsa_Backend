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
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
}

type PostParticipantModel struct {
	Ids []string `json:"userIds"`
}

func (h *Handler) PostImageCollection(c *gin.Context) {

	var model PostCollectionModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userId := c.GetString("user")
	if len(userId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is required"})
		return
	}

	user, err := h.UserService.FindById(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	collection := &domain.Collection{
		Owner:       user,
		Name:        model.Name,
		Description: model.Description,
		Latitude:    model.Latitude,
		Longitude:   model.Longitude,
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
		page = 1
		limit = 10
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
		return
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

func (h *Handler) GetParticipants(c *gin.Context) {
	id := c.Param("collectionId")

	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}

	page, err0 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err0 != nil || err1 != nil {
		page = 1
		limit = 10
	}

	collection, err := h.CollectionService.FindById(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	c.JSON(http.StatusOK, Pageate(collection.Participants, page, limit))
}

func (h *Handler) AddParticipant(c *gin.Context) {

	collectionId := c.Param("collectionId")
	if len(collectionId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collectionId is required"})
		return
	}

	var participant PostParticipantModel

	if err := c.ShouldBindJSON(&participant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("participant: %+v", participant)

	for _, userId := range participant.Ids {
		fmt.Println(userId)
		err := h.CollectionService.AddParticipant(c, collectionId, userId)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})

}

func (h *Handler) DeleteParticipant(c *gin.Context) {

	collectionId := c.Param("collectionId")
	participantId := c.Param("participantId")

	if len(collectionId) == 0 || len(participantId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collectionId and participantId is required"})
		return
	}

	if err := h.CollectionService.DeleteParticipant(c, collectionId, participantId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) PatchCollection(c *gin.Context) {

	collectionId := c.Param("collectionId")
	fmt.Println("collectionId: ", collectionId)
	if len(collectionId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collectionId is required"})
		return
	}

	var postCollection PostCollectionModel
	if err := c.ShouldBindJSON(&postCollection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.CollectionService.Patch(c,
		&domain.Collection{
			Id:             collectionId,
			Name:           postCollection.Name,
			Description:    postCollection.Description,
			Latitude:       postCollection.Latitude,
			Longitude:      postCollection.Longitude,
			PreviewImageId: postCollection.PreviewImageId,
		}); err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}
