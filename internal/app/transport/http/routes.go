package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func (h *Handler) InitialRouteSetup() {
	h.Router.Use(gin.Recovery())
	h.Router.MaxMultipartMemory = 128 << 20

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	h.Router.Use(cors.New(config))
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
	authGroup.GET("/users/me", h.GetMeUser)

	authGroup.POST("/collections", h.PostImageCollection)
	authGroup.GET("/collections", h.GetAllCollectionsOfUser)

	mustBeOwnerGroup := authGroup.Group("")
	mustBeOwnerGroup.Use(h.CheckCollectionOwnership(false))

	// Check ownership
	mustBeOwnerGroup.DELETE("/collections/:collectionId", h.DeleteCollection)
	mustBeOwnerGroup.PATCH("/collections/:collectionId", h.PatchCollection) // todo implement
	mustBeOwnerGroup.POST("/collections/:collectionId/participants", h.AddParticipant)
	mustBeOwnerGroup.DELETE("/collections/:collectionId/participants/:participantId", h.DeleteParticipant)

	mustBeParticipantGroup := authGroup.Group("")
	mustBeParticipantGroup.Use(h.CheckCollectionOwnership(true))

	// Check ownership or participant
	mustBeParticipantGroup.GET("/collections/:collectionId", h.GetCollection)

	mustBeParticipantGroup.POST("collections/:collectionId/images", h.UploadImage)
	mustBeParticipantGroup.GET("/collections/:collectionId/images", h.GetImagesByCollection)
	mustBeParticipantGroup.GET("/collections/:collectionId/images/:imageId", h.CreateDownloadImage(false))
	h.Router.GET("/collections/:collectionId/previews/:imageId", h.CreateDownloadImage(true))
	mustBeParticipantGroup.DELETE("/collections/:collectionId/images/:imageId", h.DeleteImage)

	mustBeParticipantGroup.GET("/collections/:collectionId/participants", h.GetParticipants)

}
