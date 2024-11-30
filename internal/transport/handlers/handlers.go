package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"musicAPI/internal/models"
	"strconv"
)

type Service interface {
	AddNew(title *models.Title) error
	Delete(title *models.Title) error
	Edit(song *models.Song) error
	GetCouplet(title *models.Title, startInd int, endInd int) (string, error)
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
		songs.GET("filter-by-group-and-date", h.Search)
	}

	return engine
}
func (h *Handler) Delete(c *gin.Context) {
	var title models.Title
	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(&title); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "deleted"})
}

func (h *Handler) Add(c *gin.Context) {
	var title models.Title
	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddNew(&title); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "added"})
}

func (h *Handler) Edit(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Edit(&song); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, nil)
}

func (h *Handler) Couplets(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))

	if group == "" || song == "" {
		c.JSON(400, gin.H{"error": "group and song are required"})
		return
	}

	couplets, err := h.service.GetCouplet(&models.Title{Group: group, Song: song}, page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"couplets": couplets})
}

func (h *Handler) Search(c *gin.Context) {
	// Implementation for search endpoint
}