package GoPro7

import (
	"Fairway_Bridge/Shared"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	goproIP   = "10.5.5.9"
	mediaPort = "8080"
)

// Camera encapsulates camera control parameters and state.
type Camera struct {
	log             *zap.SugaredLogger
	videoDir        string
	autoStopSeconds int
	overrideVideo   bool
	networkIP       string
	videoCounter    int
	manualStopChan  chan struct{}
	shutdownChan    chan struct{}
}

// NewCamera creates a new Camera instance, ensuring the video directory exists
// and setting up a custom HTTP client if a network IP is provided.
func NewCamera(log *zap.SugaredLogger, config Shared.Config) *Camera {
	cam := &Camera{
		videoDir:        config.Camera.VideoDir,
		autoStopSeconds: config.Camera.AutoStopSeconds,
		overrideVideo:   config.Camera.OverrideVideo,
		networkIP:       config.Camera.NetworkIP,
		manualStopChan:  make(chan struct{}),
		shutdownChan:    make(chan struct{}),
		log:             log,
	}
	cam.ensureVideoDir()
	if config.Camera.NetworkIP != "" {
		cam.setCustomHTTPClient()
	}
	return cam
}

// ensureVideoDir creates the video storage directory if it doesn't exist.
func (c *Camera) ensureVideoDir() {
	if _, err := os.Stat(c.videoDir); os.IsNotExist(err) {
		if err := os.MkdirAll(c.videoDir, 0755); err != nil {
			log.Fatalf("creating video directory: %v", err)
		}
	}
}

// setCustomHTTPClient creates an HTTP client that forces outbound connections
// to use the specified local IP address.
func (c *Camera) setCustomHTTPClient() {
	localAddr, err := net.ResolveTCPAddr("tcp", c.networkIP+":0")
	if err != nil {
		c.log.Fatalf("resolving local address %s: %v", c.networkIP, err)
	}
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			LocalAddr: localAddr,
			Timeout:   30 * time.Second,
		}).DialContext,
	}
	http.DefaultClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
	c.log.Infof("using local IP address %s for outbound connections", c.networkIP)
}

// sendCommand performs an HTTP GET to the given URL.
func (c *Camera) sendCommand(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("sendCommand GET error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sendCommand non-OK status: %s", resp.Status)
	}
	return nil
}

// startRecording sends the command to start recording.
func (c *Camera) startRecording() error {
	c.log.Info("starting recording...")
	url := fmt.Sprintf("http://%s/gp/gpControl/command/shutter?p=1", goproIP)
	return c.sendCommand(url)
}

// stopRecording sends the command to stop recording.
func (c *Camera) stopRecording() error {
	c.log.Info("stopping recording...")
	url := fmt.Sprintf("http://%s/gp/gpControl/command/shutter?p=0", goproIP)
	return c.sendCommand(url)
}

// deleteLastMedia instructs the GoPro to delete its last recorded media.
func (c *Camera) deleteLastMedia() error {
	c.log.Info("deleting last recorded video...")
	url := fmt.Sprintf("http://%s/gp/gpControl/command/storage/delete/last", goproIP)
	return c.sendCommand(url)
}

// getLatestVideoFile fetches the latest recorded video's folder and filename.
func (c *Camera) getLatestVideoFile() (folder, filename string, err error) {
	url := fmt.Sprintf("http://%s:%s/gp/gpMediaList", goproIP, mediaPort)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("getLatestVideoFile GET error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("getLatestVideoFile read error: %w", err)
	}

	// Remove carriage returns that might cause JSON unmarshal errors.
	body = []byte(strings.ReplaceAll(string(body), "\r", ""))

	var mediaList struct {
		Media []struct {
			D  string `json:"d"`
			Fs []struct {
				N string `json:"n"`
			} `json:"fs"`
		} `json:"media"`
	}

	if err := json.Unmarshal(body, &mediaList); err != nil {
		return "", "", fmt.Errorf("getLatestVideoFile unmarshal error: %w", err)
	}

	if len(mediaList.Media) == 0 || len(mediaList.Media[0].Fs) == 0 {
		return "", "", fmt.Errorf("no media files found")
	}

	latestFolder := mediaList.Media[0].D
	latestFile := mediaList.Media[0].Fs[len(mediaList.Media[0].Fs)-1].N

	return latestFolder, latestFile, nil
}

// downloadVideo downloads a video file from the GoPro.
func (c *Camera) downloadVideo(folder, filename string) error {
	downloadURL := fmt.Sprintf("http://%s/videos/DCIM/%s/%s", goproIP, folder, filename)
	var localFilePath string
	if c.overrideVideo {
		localFilePath = filepath.Join(c.videoDir, "video.mp4")
	} else {
		c.videoCounter++
		localFilePath = filepath.Join(c.videoDir, fmt.Sprintf("video_%03d.mp4", c.videoCounter))
	}

	c.log.Infof("downloading video: %s -> %s", filename, localFilePath)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("downloadVideo GET error: %w", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("downloadVideo create file error: %w", err)
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("downloadVideo copy error: %w", err)
	}

	c.log.Info("download complete:", localFilePath)
	return nil
}

func (c *Camera) Connect() error {
	c.log.Info("connecting to camera...")
	c.log.Info("✅ camera is ready!")
	return nil
}

func (c *Camera) SaveLastRecording() error {
	folder, filename, err := c.getLatestVideoFile()
	if err != nil {
		return fmt.Errorf("fetching latest video: %w", err)
	} else {
		if err := c.downloadVideo(folder, filename); err != nil {
			return fmt.Errorf("downloading video: %w", err)
		}
	}
	return nil
}

// LoopCapture starts the main recording loop. It continuously records sessions until a shutdown is signaled.
// Each session either stops automatically after autoStopSeconds (discarding the video) or is manually stopped
// (saving the recorded video). A shutdown signal during a session stops recording and saves the current video.
func (c *Camera) LoopCapture() error {
	go func() {
		for {
			// Check if a shutdown has been signaled before starting a new session.
			select {
			case <-c.shutdownChan:
				c.log.Info("shutdown signal received. Exiting run loop.")
			default:
			}

			c.log.Infof("starting new %d-second recording session...", c.autoStopSeconds)
			if err := c.startRecording(); err != nil {
				c.log.Errorf("starting recording: %v", err)
				time.Sleep(3 * time.Second)
				continue
			}

			// Wait for auto-stop timeout, a manual stop, or a shutdown during the session.
			select {
			case <-time.After(time.Duration(c.autoStopSeconds) * time.Second):
				c.log.Info("auto-stop triggered: stopping recording and discarding video.")
				if err := c.stopRecording(); err != nil {
					c.log.Errorf("stopping recording: %v", err)
				}
				time.Sleep(2 * time.Second)
				if err := c.deleteLastMedia(); err != nil {
					c.log.Errorf("deleting last media: %v", err)
				}
			case <-c.manualStopChan:
				c.log.Info("manual stop triggered: stopping recording and saving video.")
				if err := c.stopRecording(); err != nil {
					c.log.Errorf("stopping recording: %v", err)
				}
			case <-c.shutdownChan:
				c.log.Info("shutdown signal received during recording session.")
				if err := c.stopRecording(); err != nil {
					c.log.Infof("stopping recording during shutdown: %v", err)
				}
				c.log.Info("exiting recording loop due to shutdown.")
			}
			// Short delay before starting the next session.
			time.Sleep(3 * time.Second)
		}
	}()

	return nil
}

// StartCapture signals the camera to start the current recording session.
func (c *Camera) StartCapture() error {
	c.log.Infof("starting new recording session...")
	if err := c.startRecording(); err != nil {
		return fmt.Errorf("error starting recording: %w", err)
	}
	return nil
}

// StopCapture signals the camera to stop the current recording session.
func (c *Camera) StopCapture() error {
	select {
	case c.manualStopChan <- struct{}{}:
	default:
		c.log.Warnf("no active recording session to stop.")
		return fmt.Errorf("no active recording session to stop")
	}
	return nil
}

// DeleteLastRecording sends a command to delete the last recorded video.
func (c *Camera) DeleteLastRecording() error {
	err := c.deleteLastMedia()
	if err != nil {
		return err
	}
	return nil
}

// Shutdown signals the camera to stop recording entirely.
func (c *Camera) Shutdown() error {
	close(c.shutdownChan)
	return nil
}
