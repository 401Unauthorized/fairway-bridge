package Simulators

import "Fairway_Bridge/Shared"

// SimulatorController is an interface that defines the methods required for a simulator controller.
type SimulatorController interface {
	Connect() error
	LaunchShot(ballData Shared.StandardizedBallData, clubData Shared.StandardizedClubData, shotDataOptions Shared.ShotDataOptions) error
	Close() error
}
