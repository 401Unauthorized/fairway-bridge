package Launch_Monitors

import (
	"Fairway_Bridge/Shared"
)

// LaunchMonitorController is an interface that defines the methods for a launch monitor controller.
type LaunchMonitorController interface {
	Connect() error
	Close() error
	LaunchShot()
	SetOnShotCallback(f Shared.ShotHandlerFunc)
}
