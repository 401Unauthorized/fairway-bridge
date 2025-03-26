package Shared

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
)

var Version = "0.1.0"

type Camera struct {
	Name            string
	VideoDir        string
	AutoStopSeconds int
	OverrideVideo   bool
	NetworkIP       string
}

type LaunchMonitor struct {
	Name string
}

type Simulator struct {
	Name      string
	IPAddress string
	Port      int
}

type Bridge struct {
	IPAddress string
	Port      int
	LogFile   string
	LogType   string
	ShotFile  string
	Version   string
}

type HTTP struct {
	IPAddress string
	Port      int
}

type Config struct {
	LaunchMonitor
	Simulator
	Bridge
	Camera
	HTTP
}

// ParseFlags parses command-line flags and returns a Config struct.
func ParseFlags() Config {
	launchMonitor := flag.String("launch-monitor", "", "Name of the launch monitor (required)")
	simulator := flag.String("simulator", "", "Name of the simulator (required)")
	simIP := flag.String("simulator-ip", "127.0.0.1", "IP address for the simulator")
	simPort := flag.Int("simulator-port", 921, "Port for the simulator")
	bridgeIP := flag.String("bridge-ip", "127.0.0.1", "IP address for the Fairway Bridge")
	bridgePort := flag.Int("bridge-port", 2483, "Port for the Fairway Bridge")
	httpIP := flag.String("http-ip", "127.0.0.1", "IP address for the Fairway Bridge HTTP Server")
	httpPort := flag.Int("http-port", 2484, "Port for the Fairway Bridge HTTP Server")
	logFile := flag.String("bridge-log-file", "fairway-bridge.log", "Log file path")
	logType := flag.String("bridge-log-type", "CONSOLE", "Log type (console, json)")
	shotFile := flag.String("bridge-shot-file", "./shots.csv", "File to save shot data")
	camera := flag.String("camera", "", "Name of the camera")
	videoDir := flag.String("camera-video-dir", "./recordings/", "Directory to save video files")
	autoStopSeconds := flag.Int("camera-auto-stop-seconds", 5, "Recording duration (in seconds) before auto-stop")
	overrideVideo := flag.Bool("camera-override-video", false, "If true, always save to the same video file instead of creating new ones")
	networkIP := flag.String("camera-network-ip", "10.5.5.100", "Local IP address to use for outbound connections")

	// Parse all flags.
	flag.Parse()

	// Check required flags.
	if *launchMonitor == "" || *simulator == "" || *camera == "" {
		flag.Usage()
		fmt.Println("\nError: -launch-monitor, -simulator & -camera are required.")
		os.Exit(1)
	}

	// Build the configuration from the provided flags.
	return Config{
		LaunchMonitor: LaunchMonitor{
			Name: strings.ToUpper(*launchMonitor),
		},
		Simulator: Simulator{
			Name:      strings.ToUpper(*simulator),
			IPAddress: *simIP,
			Port:      *simPort,
		},
		Bridge: Bridge{
			IPAddress: *bridgeIP,
			Port:      *bridgePort,
			LogFile:   *logFile,
			LogType:   strings.ToUpper(*logType),
			ShotFile:  *shotFile,
			Version:   Version,
		},
		Camera: Camera{
			Name:            strings.ToUpper(*camera),
			VideoDir:        *videoDir,
			AutoStopSeconds: *autoStopSeconds,
			OverrideVideo:   *overrideVideo,
			NetworkIP:       *networkIP,
		},
		HTTP: HTTP{
			IPAddress: *httpIP,
			Port:      *httpPort,
		},
	}
}

// PrintConfig prints the configuration to the provided logger.
func (config *Config) PrintConfig(logger *zap.Logger) {
	if logger == nil {
		return
	}
	log := logger.With(zap.String("component", "CONFIG")).Sugar()
	log.Infof("Launch Monitor:\n  - Name: %s\n", config.LaunchMonitor.Name)
	log.Infof("Simulator:\n  - Name: %s\n  - IP Address: %s\n  - Port: %d\n",
		config.Simulator.Name, config.Simulator.IPAddress, config.Simulator.Port)
	log.Infof("Fairway Bridge:\n  - IP Address: %s\n  - Port: %d\n  - Log File: %s\n  - Shot File: %s\n  - Version: %s\n",
		config.Bridge.IPAddress, config.Bridge.Port, config.Bridge.LogFile, config.Bridge.ShotFile, config.Bridge.Version)
	log.Infof("HTTP Server:\n  - IP Address: %s\n  - Port: %d\n",
		config.HTTP.IPAddress, config.HTTP.Port)
	log.Infof("Camera:\n  - Name: %s\n  - Video Directory: %s\n  - Auto Stop (s): %d\n  - Override Video: %t\n  - Network IP: %s\n",
		config.Camera.Name, config.Camera.VideoDir, config.Camera.AutoStopSeconds, config.Camera.OverrideVideo, config.Camera.NetworkIP)
}
