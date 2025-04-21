package http

import (
	"github.com/NicoEberlein/NotSamsa_Backend/internal/service"
	"github.com/gin-gonic/gin"
	"math"
)

type Handler struct {
	UserService       *service.UserService
	CollectionService *service.ImageCollectionService
	ImageService      *service.ImageService
	Router            *gin.Engine
}

func NewHandler(
	userService *service.UserService,
	collectionService *service.ImageCollectionService,
	imageService *service.ImageService,
	router *gin.Engine) *Handler {

	return &Handler{
		UserService:       userService,
		CollectionService: collectionService,
		ImageService:      imageService,
		Router:            router,
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
