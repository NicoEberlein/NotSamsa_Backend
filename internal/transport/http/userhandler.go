package http

import (
	"errors"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserPutRestModel struct {
	Mail string `json:"mail"`
}

type UserChangePasswordRestModel struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (h *Handler) GetAllUsersHandler(c *gin.Context) {

	page, err0 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "5"))

	if err0 != nil || err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse page or limit parameter"})
		return
	}

	users, err := h.UserService.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, Pageate(users, page, limit))
}

func (h *Handler) GetUserHandler(c *gin.Context) {
	id := c.Param("userId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}

	user, err := h.UserService.FindById(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)

}

func (h *Handler) DeleteUserHandler(c *gin.Context) {
	id := c.Param("userId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be provided"})
	}

	err := h.UserService.Delete(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) PutUserHandler(c *gin.Context) {
	var user UserPutRestModel

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	id := c.Param("userId")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be provided"})
	}

	err := h.UserService.UpdateUserDetails(c, &domain.User{
		Id:   id,
		Mail: user.Mail,
	})

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) ChangePasswordHandler(c *gin.Context) {
	var model UserChangePasswordRestModel
	if err := c.ShouldBind(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	id := c.GetString("user")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be provided"})
	}

	err := h.UserService.UpdatePassword(c, id, model.OldPassword, model.NewPassword)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
