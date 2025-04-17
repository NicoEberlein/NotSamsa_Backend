package http

import (
	"github.com/gin-gonic/gin"
	"math"
)

type Page[T any] struct {
	Items       []T         `json:"items"`
	PageDetails PageDetails `json:"pageDetails"`
}

type PageDetails struct {
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
}

func Pageate[T any](items []T, page int, limit int) Page[T] {

	var itemSlice []T

	start := (page - 1) * limit
	end := start + limit
	if start > len(items) {
		itemSlice = make([]T, 0)
	}
	if end > len(items) {
		end = len(items)
	}
	itemSlice = items[start:end]

	return Page[T]{
		Items: itemSlice,
		PageDetails: PageDetails{
			TotalItems:  len(itemSlice),
			TotalPages:  int(math.Ceil(float64(len(items)) / float64(limit))),
			CurrentPage: page,
			PageSize:    limit,
		},
	}
}

func InitialRouteSetup(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.MaxMultipartMemory = 128 << 20
}

func SetupUserRoutes(router *gin.Engine, userHandler *UserHandler) {

	userGroup := router.Group("/users")
	userGroup.Use(Authenticatior())

	userGroup.GET("/", userHandler.GetAllUsersHandler)
	userGroup.DELETE("/:userId", userHandler.DeleteUserHandler)
	userGroup.GET("/:userId", userHandler.GetUserHandler)
	userGroup.PUT("/:userId", userHandler.PutUserHandler)
	userGroup.POST("/changePassword", userHandler.ChangePasswordHandler)

}

func SetupCollectionRoutes(router *gin.Engine, collectionHandler *CollectionHandler) {

	collectionsGroup := router.Group("/collections")
	collectionsGroup.Use(Authenticatior())

	collectionsGroup.GET("/:collectionId", collectionHandler.GetCollection)
	collectionsGroup.POST("/", collectionHandler.PostImageCollection)
	collectionsGroup.DELETE("/:collectionId", collectionHandler.DeleteCollection)
	collectionsGroup.GET("/", Authenticatior(), collectionHandler.GetAllCollectionsOfUser)

}

func SetupAuthRoutes(router *gin.Engine, authHandler *AuthHandler) {

	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)

}

func SetupImageRoutes(router *gin.Engine, imageHandler *ImageHandler) {

	imageGroup := router.Group("")
	imageGroup.Use(Authenticatior())

	router.POST("collection/:collectionId/upload", imageHandler.UploadImage)
	imageGroup.GET("/collection/:collectionId/images", imageHandler.GetImagesByCollection)
	imageGroup.GET("/images/:imageId/download", imageHandler.DownloadImage)
	imageGroup.DELETE("/images/:imageId", imageHandler.DeleteImage)

}
