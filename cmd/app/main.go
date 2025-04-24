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
