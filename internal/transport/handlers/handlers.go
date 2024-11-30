package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"musicAPI/internal/models"
)

type Service interface {
	AddNew(title *models.Title) error
	Delete(title *models.Title) error
	Edit(song *models.Song) error
	GetCouplet(title *models.Title, startInd int, endInd int) ([]string, error)
	GetSongsByGroupsAndRelease(groups []string, period models.Period) ([]models.Song, error)
}

type Handler struct {
	service Service
	log     *slog.Logger
	ctx     context.Context
}

func NewHandler(log *slog.Logger, clientMethods Service, ctx context.Context) *Handler {
	return &Handler{
		clientMethods,
		log,
		ctx}
}

func (h *Handler) InitRouter() *gin.Engine {
	engine := gin.Default()

	songs := engine.Group("/songs")
	{
		songs.DELETE("/delete", h.Delete)
		songs.POST("/add", h.Add)
		songs.PATCH("/edit", h.Edit)
		songs.GET("/couplets", h.Couplets)
		songs.GET("/search", h.Search)
	}

	return engine
}

func (h *Handler) Delete(c *gin.Context) {
	panic("not impl")
}

func (h *Handler) Add(c *gin.Context) {
	panic("not impl")
}

func (h *Handler) Edit(c *gin.Context) {
	panic("not impl")
}

func (h *Handler) Couplets(c *gin.Context) {
	panic("not impl")
}

func (h *Handler) Search(c *gin.Context) {
	panic("not impl")
}
