package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/app/service"
	mock2 "github.com/NicoEberlein/NotSamsa_Backend/internal/app/storage/mock"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupCollectionTestHandler() *Handler {
	gin.SetMode(gin.TestMode)
	userRepository := mock2.NewUserRepository()
	handler := NewHandler(
		&service.UserService{userRepository},
		&service.CollectionService{mock2.NewMockCollectionRepository(), userRepository},
		nil,
		gin.Default())

	authGroup := handler.Router.Group("")
	authGroup.Use(func(c *gin.Context) {
		c.Set("user", "00000000-0000-0000-0000-000000000000")
		c.Next()
	})

	authGroup.POST("/collections", handler.PostImageCollection)
	authGroup.GET("/collections", handler.GetAllCollectionsOfUser)
	authGroup.GET("/collections/:collectionId", handler.GetCollection)

	return handler
}

func TestGetAllCollectionsOfUser(t *testing.T) {

	handler := setupCollectionTestHandler()
	router := handler.Router

	_, err := handler.CollectionService.Create(context.Background(), &domain.Collection{
		OwnerId:     "00000000-0000-0000-0000-000000000000",
		Name:        "TestCollection",
		Description: "TestCollectionDescription",
	})

	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/collections", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &data)

	assert.NoError(t, err)
	assert.Len(t, data["items"], 1)
	assert.Equal(t, (data["items"].([]any))[0].(map[string]any)["name"], "TestCollection")
}

func TestPostCollection(t *testing.T) {

	handler := setupCollectionTestHandler()
	router := handler.Router

	_, err := handler.UserService.Create(context.Background(), &domain.User{Id: "00000000-0000-0000-0000-000000000000", Mail: "test@mail.com", Password: "testpassword123"})

	assert.NoError(t, err)

	body, _ := json.Marshal(PostCollectionModel{
		Name:        "test",
		Description: "test description",
	})

	req, _ := http.NewRequest("POST", "/collections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var data map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	req, _ = http.NewRequest(http.MethodGet, "/collections/"+data["id"], nil)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
