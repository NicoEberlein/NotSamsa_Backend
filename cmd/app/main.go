package main

import (
	"fmt"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/service"
	gormstore "github.com/NicoEberlein/NotSamsa_Backend/internal/storage/gorm"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/storage/s3"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"log"
)

type GlobalConfig struct {
	DB                *gorm.DB
	S3                *minio.Client
	UserService       *service.UserService
	UserHandler       *http.UserHandler
	CollectionService *service.ImageCollectionService
	CollectionHandler *http.CollectionHandler
	Router            *gin.Engine
}

var Config GlobalConfig

func main() {
	db := gormstore.GetDbConnection()
	database, _ := db.DB()
	pingErr := database.Ping()
	fmt.Println(pingErr)

	sss := s3.ConnectToS3()

	Config = GlobalConfig{
		DB: db,
		S3: sss,
		UserService: &service.UserService{
			UserRepository: gormstore.NewUserRepository(db),
		},
		CollectionService: &service.ImageCollectionService{
			ImageCollectionRepository: gormstore.NewImageCollectionRepository(db),
		},
	}

	Config.UserHandler = http.NewUserHandler(Config.UserService)
	Config.CollectionHandler = http.NewCollectionHandler(Config.CollectionService)
	Config.Router = gin.Default()

	Config.Router.POST("/user", Config.UserHandler.PostUserHandler)
	Config.Router.GET("/user", Config.UserHandler.GetAllUsersHandler)
	Config.Router.DELETE("/user/:userId", Config.UserHandler.DeleteUserHandler)
	Config.Router.GET("/user/:userId", Config.UserHandler.GetUserHandler)
	Config.Router.PUT("/user/:userId", Config.UserHandler.PutUserHandler)
	Config.Router.POST("/user/:userId/changePassword", Config.UserHandler.ChangePasswordHandler)

	Config.Router.GET("/collection/:collectionId", Config.CollectionHandler.GetCollection)
	Config.Router.GET("/user/:userId/collection", Config.CollectionHandler.GetAllCollectionsOfUser)
	Config.Router.POST("/collection", Config.CollectionHandler.PostImageCollection)
	Config.Router.DELETE("/collection/:collectionId", Config.CollectionHandler.DeleteCollection)

	if err := Config.Router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
