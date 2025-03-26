package main

import (
	"Fairway_Bridge/Cameras"
	"Fairway_Bridge/HTTP"
	"Fairway_Bridge/Router"
	"Fairway_Bridge/Shared"
	"Fairway_Bridge/Storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse CLI flags into the configuration.
	config := Shared.ParseFlags()

	// Create a logger.
	logger, logFile, err := Shared.NewLogger(config)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Print out the configuration.
	config.PrintConfig(logger)

	// Create a file for storing shots.
	storage, err := Storage.NewFileStorage(logger, config)
	if err != nil {
		logger.Sugar().Fatalf("Failed to create FileStorage: %v", err)
	}

	// Start the Camera Controller
	cam, err := Cameras.NewCamera(logger, config)
	if err != nil {
		logger.Sugar().Fatalf("Failed to create camera: %v", err)
	}

	// Connect to the Camera
	err = cam.Connect()
	if err != nil {
		logger.Sugar().Fatalf("Failed to connect to camera: %v", err)
	}

	// Create an API Server & Host UI pages.
	err = HTTP.Serve(logger, config, cam)
	if err != nil {
		logger.Sugar().Fatalf("Failed to create API server: %v", err)
	}

	// Connect to the Simulator & Launch Monitor
	launchMonitor, simulator, err := Router.LaunchMonitorToSimulator(logger, storage, cam, config)
	if err != nil {
		logger.Sugar().Fatalf("Failed to start the system router: %v", err)
	}

	logger.Sugar().Infof("Servers are running. Press Ctrl+C to exit.")

	// Block main goroutine until an interrupt signal is received.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	logger.Sugar().Infof("Shutdown signal received. Exiting...")

	// Shutdown procedures.
	if cam != nil {
		err = cam.Shutdown()
		if err != nil {
			logger.Sugar().Fatalf("Failed to shutdown camera: %v", err)
		}
	}
	if launchMonitor != nil {
		err = launchMonitor.Close()
		if err != nil {
			logger.Sugar().Fatalf("Failed to shutdown launch monitor: %v", err)
		}
	}
	if simulator != nil {
		err = simulator.Close()
		if err != nil {
			logger.Sugar().Fatalf("Failed to shutdown simulator: %v", err)
		}
	}

	if logger != nil {
		err := logger.Sync()
		if err != nil {
			logger.Sugar().Fatalf("Failed to shutdown logger: %v", err)
		}
	}

	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			logger.Sugar().Fatalf("Failed to close log file: %v", err)
		}
	}

}
