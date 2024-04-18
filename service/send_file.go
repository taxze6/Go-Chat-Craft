package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func DownloadEmojiZip(c *gin.Context) {
	// Obtain the path of the zip file
	filePath := "./assets/download/fluentui_emoji_icon_data.zip"

	// open-file
	file, err := os.Open(filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file")
		return
	}
	defer file.Close()

	//Gets file information, including size
	fileInfo, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting file info")
		return
	}

	// Set the response header to tell the browser that this is a zip file
	c.Header("Content-Type", "application/zip")

	c.Header("Content-Disposition", "attachment; filename=emoji_icon_data.zip")

	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Send the zip file to the client
	c.File(filePath)
}
