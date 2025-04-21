package http

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitialRouteSetup() {
	h.Router.Use(gin.Recovery())
	h.Router.MaxMultipartMemory = 128 << 20
}

func (h *Handler) SetupRoutes() {

	h.Router.POST("/login", h.Login)
	h.Router.POST("/register", h.Register)

	authGroup := h.Router.Group("")
	authGroup.Use(h.Authenticator())

	authGroup.GET("/users", h.GetAllUsersHandler)
	authGroup.POST("/users/changePassword", h.ChangePasswordHandler)

	authGroup.DELETE("/users/:userId", h.DeleteUserHandler)
	authGroup.GET("/users/:userId", h.GetUserHandler)
	authGroup.PUT("/users/:userId", h.PutUserHandler)

	authGroup.POST("/collections", h.PostImageCollection)
	authGroup.GET("/collections", h.GetAllCollectionsOfUser)

	// Check ownership
	authGroup.DELETE("/collections/:collectionId", h.CheckCollectionOwnership(false), h.DeleteCollection)

	mustBeParticipantGroup := authGroup.Group("")
	mustBeParticipantGroup.Use(h.CheckCollectionOwnership(true))

	// Check ownership or participant
	mustBeParticipantGroup.GET("/collections/:collectionId", h.GetCollection)

	mustBeParticipantGroup.POST("collection/:collectionId/images", h.UploadImage)
	mustBeParticipantGroup.GET("/collection/:collectionId/images", h.GetImagesByCollection)
	mustBeParticipantGroup.GET("/collection/:collectionId/images/:imageId", h.DownloadImage)
	mustBeParticipantGroup.DELETE("/collection/:collectionId/images/:imageId", h.DeleteImage)

}
