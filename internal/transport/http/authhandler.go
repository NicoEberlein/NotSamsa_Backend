package http

import (
	"errors"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type AuthHandler struct {
	UserService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		UserService: userService,
	}
}

type LoginRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Mail  string `json:"mail"`
	Token string `json:"token"`
}

func (authHandler *AuthHandler) Login(c *gin.Context) {

	var loginRequest LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := authHandler.UserService.Authenticate(c, loginRequest.Mail, loginRequest.Password)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	user, err := authHandler.UserService.FindByMail(c, loginRequest.Mail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if result {
		token, err := createToken(user.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		response := LoginResponse{
			Mail:  loginRequest.Mail,
			Token: token,
		}
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
	}

}

func (authHandler *AuthHandler) Register(c *gin.Context) {

	var user LoginRequest
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newUser *domain.User = domain.NewUser(user.Mail, user.Password)

	id, err := authHandler.UserService.Create(c, newUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
