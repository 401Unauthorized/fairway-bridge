package Virtual

import (
	"Fairway_Bridge/Shared"
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

// LaunchMonitor simulates a launch monitor.
type LaunchMonitor struct {
	log        *zap.SugaredLogger
	onShot     Shared.ShotHandlerFunc
	stopSignal chan struct{}
}

// NewLaunchMonitor initializes a virtual launch monitor.
func NewLaunchMonitor(logger *zap.Logger) *LaunchMonitor {
	log := logger.With(zap.String("component", "LAUNCH_MONITOR"), zap.String("type", "VIRTUAL")).Sugar()

	return &LaunchMonitor{
		log:        log,
		stopSignal: make(chan struct{}),
	}
}

// Connect simulates connecting to a launch monitor.
func (lm *LaunchMonitor) Connect() error {
	lm.log.Infof("connecting to virtual launch monitor...")
	time.Sleep(1 * time.Second) // Simulate connection delay
	lm.log.Infof("✅ virtual launch monitor is ready!")
	return nil
}

// Close simulates disconnecting from the launch monitor.
func (lm *LaunchMonitor) Close() error {
	lm.log.Infof("closing virtual launch monitor connection...")
	close(lm.stopSignal)        // Signal the goroutine to stop
	time.Sleep(1 * time.Second) // Simulate delay
	lm.log.Infof("virtual launch monitor connection closed.")
	return nil
}

// SetOnShotCallback sets the callback function to be called when a shot is launched.
func (lm *LaunchMonitor) SetOnShotCallback(f Shared.ShotHandlerFunc) {
	lm.onShot = f
}

// RandomShot simulates launching a shot and calls the onShot callback.
func (lm *LaunchMonitor) RandomShot() {
	lm.log.Infof("🏌️ simulating shot launch...")
	time.Sleep(1 * time.Second) // Simulate processing time

	// Generate random shot data
	ballData := Shared.StandardizedBallData{
		Speed:         rand.Float64()*50 + 120,
		SpinAxis:      rand.Float64()*10 - 5,
		TotalSpin:     rand.Float64()*4000 + 2000,
		HLA:           rand.Float64()*5 - 2.5,
		VLA:           rand.Float64()*10 + 10,
		CarryDistance: rand.Float64()*100 + 150,
	}

	clubData := Shared.StandardizedClubData{
		Speed:                rand.Float64()*20 + 90,
		AngleOfAttack:        rand.Float64()*5 - 2.5,
		FaceToTarget:         rand.Float64()*5 - 2.5,
		Lie:                  rand.Float64()*5 - 2.5,
		Loft:                 rand.Float64()*5 + 8,
		Path:                 rand.Float64()*5 - 2.5,
		SpeedAtImpact:        rand.Float64()*20 + 90,
		VerticalFaceImpact:   rand.Float64()*2 - 1,
		HorizontalFaceImpact: rand.Float64()*2 - 1,
		ClosureRate:          rand.Float64()*5 - 2.5,
	}

	lm.log.Infof("⛳️ shot launched! %v %v", ballData, clubData)

	// Call the callback function
	if lm.onShot != nil {
		lm.onShot(ballData, clubData, Shared.ShotDataOptions{
			ClubType:         "Driver",
			ContainsClubData: true,
			ContainsBallData: true,
		})
	}
}

// promptInput prompts the user for input and returns the value as a float64.
func promptInput(prompt string, defaultValue float64) float64 {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (default: %.1f): ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	if input == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("Invalid input, using default value.")
		return defaultValue
	}
	return value
}

// LaunchShot simulates launching a shot with a specified distance and fixed launch angle.
func (lm *LaunchMonitor) LaunchShot() {
	go func() {
		for {
			select {
			case <-lm.stopSignal:
				lm.log.Infof("stopping prompt input loop.")
				return
			default:
				distance := promptInput("Enter target distance (yards)", 150.0)

				lm.log.Infof("🏌️ simulating shot launch...")

				// Fixed launch angle (in degrees)
				launchAngle := 25.0
				theta := launchAngle * (math.Pi / 180.0)

				// Convert distance from yards to meters (1 yard = 0.9144 m)
				distanceMeters := distance * 0.9144
				g := 9.81 // gravitational acceleration (m/s²)

				// Ideal projectile motion formula: distance = (v² * sin(2θ)) / g
				// Solve for v (m/s): v = sqrt(distanceMeters * g / sin(2θ))
				vMPS := math.Sqrt(distanceMeters * g / math.Sin(2*theta))

				// Convert velocity from m/s to mph (1 m/s ≈ 2.23694 mph or divide by 0.44704)
				computedBallSpeed := vMPS / 0.44704

				// Correction factor to adjust for simulator physics (tweak this value as needed)
				correctionFactor := 1.2 // Lowered from 1.5 to reduce overshooting
				ballSpeed := computedBallSpeed * correctionFactor

				// Calculate club head speed as a fraction of ball speed
				clubHeadSpeed := ballSpeed * 0.95

				// Clamp club head speed to maximum allowed value (200 mph)
				if clubHeadSpeed > 200 {
					clubHeadSpeed = 200
				}

				// Other default shot parameters
				spinAxis := 3.0        // Draw shot
				totalSpin := 4000.0    // Total spin in RPM
				launchDirection := 0.0 // No horizontal deviation
				clubAngleFace := 0.0   // Club face square
				clubAnglePath := 0.0   // Neutral club path

				// Generate random shot data
				ballData := Shared.StandardizedBallData{
					Speed:     ballSpeed,
					SpinAxis:  spinAxis,
					TotalSpin: totalSpin,
					HLA:       launchDirection,
					VLA:       launchAngle,
				}

				clubData := Shared.StandardizedClubData{
					Speed:         clubHeadSpeed,
					AngleOfAttack: clubAnglePath,
					FaceToTarget:  clubAngleFace,
				}

				lm.log.Infof("⛳️ shot launched! %v %v", ballData, clubData)

				// Call the callback function
				if lm.onShot != nil {
					lm.onShot(ballData, clubData, Shared.ShotDataOptions{
						ClubType:         "7Iron",
						ContainsClubData: true,
						ContainsBallData: true,
					})
				}
			}
		}
	}()
}
