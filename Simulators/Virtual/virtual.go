package Virtual

import (
	"Fairway_Bridge/Shared"
	"go.uber.org/zap"
	"time"
)

// Simulator creates a virtual golf simulator.
type Simulator struct {
	log *zap.SugaredLogger
}

// NewSimulator initializes the virtual simulator.
func NewSimulator(logger *zap.Logger) *Simulator {
	log := logger.With(zap.String("component", "SIMULATOR"), zap.String("type", "VIRTUAL")).Sugar()
	return &Simulator{
		log: log,
	}
}

// Connect establishes a connection.
func (vs *Simulator) Connect() error {
	vs.log.Infof("connecting to virtual simulator...")
	time.Sleep(200 * time.Millisecond) // Simulate delay
	vs.log.Infof("✅ connected to virtual simulator!")
	return nil
}

// Close simulates closing the connection.
func (vs *Simulator) Close() error {
	vs.log.Infof("closing connection...")
	time.Sleep(1 * time.Second) // Simulate processing time
	vs.log.Infof("connection closed.")
	return nil
}

// LaunchShot simulates launching a golf shot.
func (vs *Simulator) LaunchShot(ballData Shared.StandardizedBallData, cludData Shared.StandardizedClubData, shotDataOptions Shared.ShotDataOptions) error {
	vs.log.Infof("🏌️ simulating shot launch...")
	time.Sleep(1 * time.Second) // Simulate processing time
	vs.log.Infof("⛳️ shot launched on virtual simulator! %v %v %v ", ballData, cludData, shotDataOptions)
	return nil
}
