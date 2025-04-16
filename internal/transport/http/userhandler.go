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

type UserCreateRestModel struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type UserPutRestModel struct {
	Mail string `json:"mail"`
}

type UserChangePasswordRestModel struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (userHandler *UserHandler) PostUserHandler(c *gin.Context) {

	var user UserCreateRestModel
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newUser *domain.User = domain.NewUser(user.Mail, user.Password)

	id, err := userHandler.UserService.Create(c, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (userHandler *UserHandler) GetAllUsersHandler(c *gin.Context) {

	page, err0 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "5"))

	if err0 != nil || err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse page or limit parameter"})
		return
	}

	users, err := userHandler.UserService.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, Pageate(users, page, limit))
}

func (userHandler *UserHandler) GetUserHandler(c *gin.Context) {
	id := c.Param("userId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}

	user, err := userHandler.UserService.FindById(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)

}

func (userHandler *UserHandler) DeleteUserHandler(c *gin.Context) {
	id := c.Param("userId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be provided"})
	}

	err := userHandler.UserService.Delete(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (userHandler *UserHandler) PutUserHandler(c *gin.Context) {
	var user UserPutRestModel

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	id := c.Param("userId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be provided"})
	}

	err := userHandler.UserService.UpdateUserDetails(c, &domain.User{
		Id:   id,
		Mail: user.Mail,
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (userHandler *UserHandler) ChangePasswordHandler(c *gin.Context) {
	var model UserChangePasswordRestModel
	if err := c.ShouldBind(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	id := c.Param("userId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be provided"})
	}

	err := userHandler.UserService.UpdatePassword(c, id, model.OldPassword, model.NewPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
