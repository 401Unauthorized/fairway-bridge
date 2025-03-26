package Virtual

import (
	"Fairway_Bridge/Shared"
	"go.uber.org/zap"
	"time"
)

// Camera encapsulates virtual control parameters and state.
type Camera struct {
	log             *zap.SugaredLogger
	videoDir        string
	autoStopSeconds int
	overrideVideo   bool
	networkIP       string
	config          Shared.Config
}

// NewCamera creates a new Camera instance.
func NewCamera(log *zap.SugaredLogger, config Shared.Config) *Camera {
	cam := &Camera{
		config:          config,
		videoDir:        config.Camera.VideoDir,
		autoStopSeconds: config.Camera.AutoStopSeconds,
		overrideVideo:   config.Camera.OverrideVideo,
		networkIP:       config.Camera.NetworkIP,
		log:             log,
	}
	return cam
}

// Connect establishes a connection to the camera.
func (c *Camera) Connect() error {
	c.log.Infof("connecting to camera...")
	time.Sleep(1 * time.Second) // Simulate connection delay
	c.log.Infof("✅ camera is ready!")
	return nil
}

// StartCapture starts the camera recording.
func (c *Camera) StartCapture() error {
	c.log.Infof("starting new recording session...")
	return nil
}

// StopCapture stops the camera recording.
func (c *Camera) StopCapture() error {
	c.log.Infof("stopping recording session.")
	return nil
}

// SaveLastRecording saves the last recording.
func (c *Camera) SaveLastRecording() error {
	c.log.Infof("saving last recording...")
	time.Sleep(1 * time.Second) // Simulate connection delay
	c.log.Infof("✅ recording saved!")
	return nil
}

// DeleteLastRecording deletes the last recording.
func (c *Camera) DeleteLastRecording() error {
	c.log.Infof("deleting last recording...")
	time.Sleep(1 * time.Second) // Simulate connection delay
	c.log.Infof("✅ recording deleted!")
	return nil
}

// LoopCapture starts a loop capture session.
func (c *Camera) LoopCapture() error {
	c.log.Infof("starting new %d-second recording session...", c.autoStopSeconds)
	return nil
}

// Shutdown handles the shutdown signal.
func (c *Camera) Shutdown() error {
	c.log.Infof("shutdown signal received.")
	return nil
}
