package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	localError "musicAPI/internal/err"
	"musicAPI/internal/libs/parsers"
	"musicAPI/internal/models"
	"strconv"
)

type Service interface {
	AddNew(ctx context.Context, title *models.Title) error
	Delete(ctx context.Context, title *models.Title) error
	Edit(ctx context.Context, song *models.Song) error
	GetCouplets(ctx context.Context, title *models.Title, page int, limit int) (string, error)
	GetSongsByGroupsAndRelease(ctx context.Context, filter *models.Filter, page int, limit int) ([]models.Song, error)
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
		songs.POST("/couplets", h.Couplets)
		songs.POST("filter-by-group-and-date", h.GetSongs)
	}

	return engine
}
func (h *Handler) Delete(c *gin.Context) {
	var title models.Title
	if err := c.ShouldBindJSON(&title); err != nil {
		h.log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(h.ctx, &title); err != nil {
		h.log.Error(err.Error())
		if errors.Is(err, localError.ErrNotFound) {
			c.JSON(404, gin.H{})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Delete success")
	c.JSON(204, gin.H{})
}

func (h *Handler) Add(c *gin.Context) {
	var title models.Title
	if err := c.ShouldBindJSON(&title); err != nil {
		h.log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AddNew(h.ctx, &title); err != nil {
		h.log.Error(err.Error())
		if errors.Is(err, localError.ErrAlreadyExist) {
			c.JSON(409, gin.H{})
			return
		}
		if errors.Is(err, localError.ErrNotFound) {
			c.JSON(404, gin.H{}) // cannot find in API
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("AddNew success")
	c.JSON(201, gin.H{})
}

func (h *Handler) Edit(c *gin.Context) {
	song := &models.Song{}
	if err := c.ShouldBindJSON(song); err != nil {
		h.log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if song.Info == nil || song.Title.Song == "" ||
		song.Title.Group == "" || !parsers.IsValidDate(song.Info.ReleaseDate) {
		h.log.Error("Edit: empty required field or not corrected value")
		c.JSON(400, gin.H{})
		return
	}

	if err := h.service.Edit(h.ctx, song); err != nil {
		h.log.Error(err.Error())
		if errors.Is(err, localError.ErrNotFound) {
			c.JSON(404, gin.H{})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Edit success")
	c.JSON(204, gin.H{})
}

func (h *Handler) Couplets(c *gin.Context) {
	title := &models.Title{}
	if err := c.ShouldBindJSON(title); err != nil {
		h.log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if title.Group == "" || title.Song == "" {
		h.log.Error("Couplets: group and song are required")
		c.JSON(400, gin.H{})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		h.log.Error("Couplets: group and song are required")
		c.JSON(400, gin.H{})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	if err != nil || limit < 0 {
		h.log.Error("Couplets: group and song are required")
		c.JSON(400, gin.H{})
		return
	}

	couplets, err := h.service.GetCouplets(h.ctx, title, page, limit)
	if err != nil {
		h.log.Error(err.Error())
		if errors.Is(err, localError.ErrNotFound) {
			c.JSON(404, gin.H{})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("GetCouplets success")
	c.JSON(200, gin.H{"couplets": couplets})
}

func (h *Handler) GetSongs(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		h.log.Error("GetSongs: not correct page or limit")
		c.JSON(400, gin.H{})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	if err != nil || limit < -1 {
		h.log.Error("GetSongs: not correct page or limit")
		c.JSON(400, gin.H{})
		return
	}

	filter := &models.Filter{}
	if err := c.ShouldBindJSON(filter); err != nil {
		h.log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if filter.Per == nil {
		h.log.Error("GetSongs: not correct period")
		c.JSON(400, gin.H{})
		return
	}
	if filter.Per.Start != "" && !parsers.IsValidDate(filter.Per.Start) {
		h.log.Error("GetSongs: not correct date")
		c.JSON(400, gin.H{})
		return
	}
	if filter.Per.End != "" && !parsers.IsValidDate(filter.Per.End) {
		h.log.Error("GetSongs: not correct date")
		c.JSON(400, gin.H{})
		return
	}

	songs, err := h.service.GetSongsByGroupsAndRelease(h.ctx, filter, page, limit)
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("GetSongsByGroupsAndRelease success")
	c.JSON(200, gin.H{
		"data": songs,
	})
}
