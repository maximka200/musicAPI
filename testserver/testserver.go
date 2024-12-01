package testserver

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func NewTestServer() *gin.Engine {
	router := gin.Default()

	router.GET("/info", func(c *gin.Context) {
		group := c.Query("group")
		song := c.Query("song")

		if group == "" || song == "" {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		response := gin.H{
			"releaseDate": "16.07.2006",
			"text": "Ooh baby, don't you know I suffer?\\nOoh" +
				"baby, can you hear me moan?\\nYou caught me under false pretenses\\n" +
				"How long before you let me go?\\n\\nOoh\\nYou set my soul alight\\" +
				"nOoh\\nYou set my soul alight",
			"link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}
		c.JSON(http.StatusOK, response)
	})

	return router
}

func main() {
	server := NewTestServer()

	log.Println()
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
