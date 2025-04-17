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
	ImageService      *service.ImageService
	ImageHandler      *http.ImageHandler
	AuthHandler       *http.AuthHandler
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
		ImageService: &service.ImageService{
			ImageRepository: gormstore.NewImageRepository(db),
			S3:              sss,
		},
	}

	Config.UserHandler = http.NewUserHandler(Config.UserService)
	Config.CollectionHandler = http.NewCollectionHandler(Config.CollectionService)
	Config.AuthHandler = http.NewAuthHandler(Config.UserService)
	Config.ImageHandler = http.NewImageHandler(Config.ImageService)
	Config.Router = gin.Default()

	http.InitialRouteSetup(Config.Router)
	http.SetupUserRoutes(Config.Router, Config.UserHandler)
	http.SetupCollectionRoutes(Config.Router, Config.CollectionHandler)
	http.SetupAuthRoutes(Config.Router, Config.AuthHandler)
	http.SetupImageRoutes(Config.Router, Config.ImageHandler)

	if err := Config.Router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
