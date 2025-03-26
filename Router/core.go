package Router

import (
	"Fairway_Bridge/Cameras"
	"Fairway_Bridge/Launch_Monitors"
	Garmin_R10 "Fairway_Bridge/Launch_Monitors/Garmin-R10"
	"Fairway_Bridge/Launch_Monitors/Virtual"
	"Fairway_Bridge/Shared"
	"Fairway_Bridge/Simulators"
	"Fairway_Bridge/Simulators/GSPro"
	Virtual2 "Fairway_Bridge/Simulators/Virtual"
	"Fairway_Bridge/Storage"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// LaunchMonitorToSimulator initializes the launch monitor and simulator based on the provided configuration.
func LaunchMonitorToSimulator(logger *zap.Logger, storage *Storage.FileStorage, camera Cameras.CameraController, config Shared.Config) (Launch_Monitors.LaunchMonitorController, Simulators.SimulatorController, error) {
	log := logger.With(zap.String("component", "ROUTER")).Sugar()

	var launchMonitor Launch_Monitors.LaunchMonitorController
	var simulator Simulators.SimulatorController

	// Create the simulator
	switch config.Simulator.Name {
	case "GSPRO":
		gsClient := GSPro.NewSimulator(config.Simulator.IPAddress, config.Simulator.Port, logger)
		if err := gsClient.Connect(); err != nil {
			return nil, nil, fmt.Errorf("failed to connect to GSPro: %v", err)
		}
		simulator = gsClient
	case "VIRTUAL":
		simClient := Virtual2.NewSimulator(logger)
		if err := simClient.Connect(); err != nil {
			return nil, nil, fmt.Errorf("failed to connect to Virtual: %v", err)
		}
		simulator = simClient
	default:
		return nil, nil, fmt.Errorf("simulator %s is not supported", config.Simulator.Name)
	}

	// Create the launch monitor
	switch config.LaunchMonitor.Name {
	case "R10":
		r10 := Garmin_R10.NewLaunchMonitor(config.Bridge.IPAddress, config.Bridge.Port, logger, config)
		if err := r10.Connect(); err != nil {
			return nil, nil, fmt.Errorf("failed to connect to R10: %v", err)
		}
		launchMonitor = r10
	case "VIRTUAL":
		virtual := Virtual.NewLaunchMonitor(logger)
		if err := virtual.Connect(); err != nil {
			return nil, nil, fmt.Errorf("failed to connect to R10: %v", err)
		}
		launchMonitor = virtual
	default:
		return nil, nil, fmt.Errorf("launch monitor %s is not supported", config.LaunchMonitor.Name)
	}

	// Create the shot callback function
	launchMonitor.SetOnShotCallback(func(standardBall Shared.StandardizedBallData, standardClub Shared.StandardizedClubData, shotDataOptions Shared.ShotDataOptions) {
		log.Infof("received shot callback from %s", config.LaunchMonitor.Name)

		// Use Controllable Modifiers
		modifiers := Shared.GetModifiers()

		// Apply Multiplier Adjustment for ball & club data
		adjustedBallOut := standardBall.ApplyAdjustment(modifiers.BallData)
		adjustedClubOut := standardClub.ApplyAdjustment(modifiers.ClubData)

		log.Infof("Ball Data Adjusted - Speed: %.2f -> %.2f, SpinAxis: %.2f -> %.2f, TotalSpin: %.2f -> %.2f, HLA: %.2f -> %.2f, VLA: %.2f -> %.2f",
			standardBall.Speed, adjustedBallOut.Speed,
			standardBall.SpinAxis, adjustedBallOut.SpinAxis,
			standardBall.TotalSpin, adjustedBallOut.TotalSpin,
			standardBall.HLA, adjustedBallOut.HLA,
			standardBall.VLA, adjustedBallOut.VLA,
		)

		log.Infof("Club Data Adjusted - Speed: %.2f -> %.2f, AngleOfAttack: %.2f -> %.2f, FaceToTarget: %.2f -> %.2f, Lie: %.2f -> %.2f, Loft: %.2f -> %.2f, Path: %.2f -> %.2f, SpeedAtImpact: %.2f -> %.2f, VerticalFaceImpact: %.2f -> %.2f, HorizontalFaceImpact: %.2f -> %.2f, ClosureRate: %.2f -> %.2f",
			standardClub.Speed, adjustedClubOut.Speed,
			standardClub.AngleOfAttack, adjustedClubOut.AngleOfAttack,
			standardClub.FaceToTarget, adjustedClubOut.FaceToTarget,
			standardClub.Lie, adjustedClubOut.Lie,
			standardClub.Loft, adjustedClubOut.Loft,
			standardClub.Path, adjustedClubOut.Path,
			standardClub.SpeedAtImpact, adjustedClubOut.SpeedAtImpact,
			standardClub.VerticalFaceImpact, adjustedClubOut.VerticalFaceImpact,
			standardClub.HorizontalFaceImpact, adjustedClubOut.HorizontalFaceImpact,
			standardClub.ClosureRate, adjustedClubOut.ClosureRate,
		)

		// Send the shot to the simulator
		if err := simulator.LaunchShot(adjustedBallOut, standardClub, shotDataOptions); err != nil {
			log.Errorf("launching shot via %s: %v", config.Simulator.Name, err)
		} else {
			log.Infof("✅ shot sent to simulator successfully!")
		}

		// Save the shot data
		if err := storage.SaveShot(time.Now(), standardBall, standardClub, adjustedBallOut, adjustedClubOut, shotDataOptions.ClubType); err != nil {
			log.Errorf("saving shot: %v", err)
		} else {
			log.Infof("✅ shot saved successfully!")
		}

		// Stop the camera and save the recording
		if camera != nil {
			if err := camera.StopCapture(); err != nil {
				log.Errorf("stopping camera: %v", err)
			} else {
				log.Infof("✅ camera stopped successfully!")
			}
			if err := camera.SaveLastRecording(); err != nil {
				log.Errorf("saving recording: %v", err)
			} else {
				log.Infof("✅ recording saved successfully!")
			}
		}
	})

	// Start user input loop if the launch monitor is virtual
	if config.LaunchMonitor.Name == "VIRTUAL" {
		launchMonitor.LaunchShot()
	}

	return launchMonitor, simulator, nil
}
