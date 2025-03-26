package Cameras

import (
	"Fairway_Bridge/Cameras/GoPro7"
	"Fairway_Bridge/Cameras/Virtual"
	"Fairway_Bridge/Shared"
	"fmt"
	"go.uber.org/zap"
)

// NewCamera creates a new camera controller based on the configuration.
func NewCamera(logger *zap.Logger, config Shared.Config) (CameraController, error) {
	log := logger.With(zap.String("component", "CAMERA")).Sugar()
	var cam CameraController
	switch config.Camera.Name {
	case "GOPRO7":
		cam = GoPro7.NewCamera(log, config)
	case "VIRTUAL":
		cam = Virtual.NewCamera(log, config)
	default:
		return nil, fmt.Errorf("camera %s is not supported", config.Camera.Name)
	}
	return cam, nil
}
