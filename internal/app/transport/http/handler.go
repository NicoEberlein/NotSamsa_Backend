package http

import (
	"github.com/NicoEberlein/NotSamsa_Backend/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"math"
	"time"
)

type Handler struct {
	UserService       *service.UserService
	CollectionService *service.CollectionService
	ImageService      *service.ImageService
	Router            *gin.Engine
	RestClient        *resty.Client
}

func NewHandler(
	userService *service.UserService,
	collectionService *service.CollectionService,
	imageService *service.ImageService,
	router *gin.Engine) *Handler {

	microserviceClient := resty.New().
		SetTimeout(3 * time.Second).
		SetBaseURL("http://localhost:8081")

	return &Handler{
		UserService:       userService,
		CollectionService: collectionService,
		ImageService:      imageService,
		Router:            router,
		RestClient:        microserviceClient,
	}
}

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
