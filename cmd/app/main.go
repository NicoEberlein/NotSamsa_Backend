package main

import (
	"fmt"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/app/service"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/app/storage/gorm"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/app/transport/http"
	"github.com/NicoEberlein/NotSamsa_Backend/internal/storage/s3"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"log"
)

type GlobalConfig struct {
	DB                *gorm.DB
	S3                *minio.Client
	UserService       *service.UserService
	CollectionService *service.CollectionService
	ImageService      *service.ImageService
	Handler           *http.Handler
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
	}

	imageRepository := gormstore.NewImageRepository(Config.DB)

	Config.ImageService = &service.ImageService{
		ImageRepository: imageRepository,
		S3:              sss,
	}

	userRepository := gormstore.NewUserRepository(Config.DB)

	Config.UserService = &service.UserService{
		UserRepository: userRepository,
	}

	Config.CollectionService = &service.CollectionService{
		CollectionRepository: gormstore.NewImageCollectionRepository(db),
		UserRepository:       userRepository,
		ImageRepository:      imageRepository,
		S3:                   Config.S3,
	}

	Config.Router = gin.Default()

	Config.Handler = http.NewHandler(
		Config.UserService,
		Config.CollectionService,
		Config.ImageService,
		Config.Router,
	)

	Config.Handler.InitialRouteSetup()
	Config.Handler.SetupRoutes()

	if err := Config.Router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
