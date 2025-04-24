package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/service"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/storage/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PostUserModel struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type PutUserModel struct {
	Mail string `json:"mail"`
}

type ChangePasswordModel struct {
	OldPassword string
	NewPassword string
}

func setupTestHandler() *Handler {
	gin.SetMode(gin.TestMode)
	userRepository := mock.NewUserRepository()
	handler := NewHandler(
		&service.UserService{userRepository},
		&service.CollectionService{mock.NewMockCollectionRepository(), userRepository},
		nil,
		gin.Default())

	handler.Router.POST("/register", handler.Register)
	handler.Router.POST("/login", handler.Login)

	authGroup := handler.Router.Group("")
	authGroup.Use(func(c *gin.Context) {
		c.Set("user", "00000000-0000-0000-0000-000000000000")
		c.Next()
	})

	authGroup.GET("/users", handler.GetAllUsersHandler)
	authGroup.POST("/users/changePassword", handler.ChangePasswordHandler)

	authGroup.DELETE("/users/:userId", handler.DeleteUserHandler)
	authGroup.GET("/users/:userId", handler.GetUserHandler)
	authGroup.PUT("/users/:userId", handler.PutUserHandler)
	authGroup.GET("/users/me", handler.GetMeUser)

	return handler
}

func TestRegisterUser(t *testing.T) {
	handler := setupTestHandler()
	router := handler.Router

	user := PostUserModel{
		Mail:     "test@mail.com",
		Password: "testpassword123",
	}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var data map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &data)

	assert.NoError(t, err)

	_, exists := data["id"]
	assert.True(t, exists)

}

func TestLoginUser(t *testing.T) {

	handler := setupTestHandler()
	router := handler.Router

	user := PostUserModel{
		Mail:     "test@mail.com",
		Password: "testpassword123",
	}

	_, err := handler.UserService.Create(context.Background(), &domain.User{Mail: user.Mail, Password: user.Password})

	assert.NoError(t, err)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &data)

	assert.NoError(t, err)

	_, exists := data["token"]
	assert.True(t, exists)

	assert.Equal(t, user.Mail, data["mail"])

}

func TestGetAllUsersHandler(t *testing.T) {
	handler := setupTestHandler()
	router := handler.Router

	handler.UserService.Create(context.Background(), &domain.User{Mail: "test0@mail.com", Password: "testpassword123"})
	handler.UserService.Create(context.Background(), &domain.User{Mail: "test1@mail.com", Password: "testpassword123"})
	handler.UserService.Create(context.Background(), &domain.User{Mail: "test2@mail.com", Password: "testpassword123"})

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &data)

	assert.NoError(t, err)
	assert.Len(t, data["items"], 3)
}

func TestChangePasswordHandler(t *testing.T) {
	handler := setupTestHandler()
	router := handler.Router

	_, err := handler.UserService.Create(context.Background(), &domain.User{Id: "00000000-0000-0000-0000-000000000000", Mail: "test0@mail.com", Password: "testpassword123"})
	assert.NoError(t, err)

	body, _ := json.Marshal(ChangePasswordModel{
		OldPassword: "testpassword123",
		NewPassword: "testpassword456",
	})

	req, _ := http.NewRequest("POST", "/users/changePassword", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	loginBody, _ := json.Marshal(PostUserModel{Mail: "test0@mail.com", Password: "testpassword456"})
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestGetUserHandler(t *testing.T) {
	handler := setupTestHandler()
	router := handler.Router

	id, err := handler.UserService.Create(context.Background(), &domain.User{Mail: "test0@mail.com", Password: "testpassword123"})
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/"+id, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)
	assert.Equal(t, id, data["id"])
}

func TestDeleteUserHandler(t *testing.T) {
	handler := setupTestHandler()
	router := handler.Router

	id, err := handler.UserService.Create(context.Background(), &domain.User{Mail: "test0@mail.com", Password: "testpassword123"})

	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/"+id, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("DELETE", "/users/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/users/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestPutUserHandler(t *testing.T) {
	handler := setupTestHandler()
	router := handler.Router

	id, err := handler.UserService.Create(context.Background(), &domain.User{Mail: "test0@mail.com", Password: "testpassword123"})

	assert.NoError(t, err)

	body, _ := json.Marshal(PutUserModel{Mail: "testnew@mail.com"})
	fmt.Println(string(body))

	req, _ := http.NewRequest("PUT", "/users/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/users/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]string

	err = json.Unmarshal(w.Body.Bytes(), &data)
	assert.NoError(t, err)

	fmt.Printf("%+v", data)

	assert.Equal(t, "testnew@mail.com", data["mail"])
}

func TestGetMeUser(t *testing.T) {

	handler := setupTestHandler()
	router := handler.Router

	_, err := handler.UserService.Create(context.Background(), &domain.User{Id: "00000000-0000-0000-0000-000000000000", Mail: "test0@mail.com", Password: "testpassword123"})

	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]string

	err = json.Unmarshal(w.Body.Bytes(), &data)

	assert.NoError(t, err)
	assert.Equal(t, data["mail"], "test0@mail.com")
}
