package HTTP

import (
	"Fairway_Bridge/Cameras"
	"Fairway_Bridge/Shared"
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	dataMutex sync.RWMutex
)

const (
	uploadDir = "Assets"
	imageFile = "stats.png"
)

// getModifiers handles GET /modifiers
func getModifiers(c *gin.Context) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	c.JSON(http.StatusOK, Shared.GetModifiers())
}

// updateModifiers handles PUT /modifiers
func updateModifiers(c *gin.Context) {
	var newModifiers Shared.ModifierData
	if err := c.ShouldBindJSON(&newModifiers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	dataMutex.Lock()
	defer dataMutex.Unlock()
	Shared.SetModifiers(newModifiers)
	c.Status(http.StatusOK)
}

// saveModifiers handles POST /modifiers/save
func saveModifiers(c *gin.Context) {
	// TODO: Implement saving logic for modifiers if needed.
	c.Status(http.StatusOK)
}

// getStatsImage handles GET /stats-image
func getStatsImage(c *gin.Context) {
	imagePath := filepath.Join(uploadDir, imageFile)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	c.File(imagePath)
}

// handleUploads handles POST /upload
func handleUploads(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse file"})
		return
	}

	imagePath := filepath.Join(uploadDir, imageFile)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Upload successful: %s", imagePath))
}

// LogResponse represents the structure of the log response
type LogResponse struct {
	Logs []string `json:"Logs"`
}

// readLastNLines reads the last n lines from a file.
func readLastNLines(filePath string, n int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return lines, nil
}

// getLogs handles GET /logs
func getLogs(logFile string) gin.HandlerFunc {
	return func(c *gin.Context) {
		lines, err := readLastNLines(logFile, 100)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read log file"})
			return
		}

		c.JSON(http.StatusOK, LogResponse{Logs: lines})
	}
}

// startCamera handles POST /camera/start
func startCamera(cam Cameras.CameraController) gin.HandlerFunc {
	return func(c *gin.Context) {
		delayStr := c.Query("delay")
		delay, err := strconv.Atoi(delayStr)
		if err != nil {
			delay = 0
		}

		time.Sleep(time.Duration(delay) * time.Millisecond)

		err = cam.StartCapture()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start camera"})
			return
		}

		c.Status(http.StatusOK)
	}
}

// stopCamera handles POST /camera/stop
func stopCamera(cam Cameras.CameraController) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := cam.StopCapture()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop camera"})
			return
		}

		c.Status(http.StatusOK)
	}
}

// saveCamera handles POST /camera/save
func saveCamera(cam Cameras.CameraController) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := cam.SaveLastRecording()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save camera data"})
			return
		}

		c.Status(http.StatusOK)
	}
}

// deleteCamera handles POST /camera/delete
func deleteCamera(cam Cameras.CameraController) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := cam.DeleteLastRecording()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete camera data"})
			return
		}

		c.Status(http.StatusOK)
	}
}

// Serve starts the HTTP server with the given logger, log buffer, IP address, and port
func Serve(logger *zap.Logger, config Shared.Config, cam Cameras.CameraController) error {
	log := logger.With(zap.String("component", "HTTP")).Sugar()

	// Ensure upload directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Serve static files (HTML, images)
	r.Static("assets", "Assets")

	// API Endpoints
	r.GET("/modifiers", getModifiers)
	r.PUT("/modifiers", updateModifiers)
	r.POST("/modifiers/save", saveModifiers)
	r.GET("/stats-image", getStatsImage)
	r.GET("/logs", getLogs(config.Bridge.LogFile))
	r.POST("/upload", handleUploads)

	r.POST("/camera/start", startCamera(cam))
	r.POST("/camera/stop", stopCamera(cam))
	r.POST("/camera/save", saveCamera(cam))
	r.POST("/camera/delete", deleteCamera(cam))

	// Serve HTML pages
	r.GET("/tv", func(c *gin.Context) {
		c.File("Assets/tv.html")
	})
	r.GET("/settings", func(c *gin.Context) {
		c.File("Assets/settings.html")
	})

	address := fmt.Sprintf("%s:%d", config.HTTP.IPAddress, config.HTTP.Port)
	log.Infof("Server running on http://%s", address)

	go func() {
		if err := r.Run(address); err != nil {
			log.Fatalf("HTTP Server Error: %v", err)
		}
	}()

	return nil
}
