package Garmin_R10

import (
	"Fairway_Bridge/Shared"
	"bufio"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"
)

const (
	heartbeatInterval   = 10 * time.Second
	pingTimeoutDuration = 3 * time.Second
)

// LaunchMonitor encapsulates the R10 server and connection state.
type LaunchMonitor struct {
	config          Shared.Config
	port            int
	localIP         string
	listener        net.Listener
	client          net.Conn
	log             *zap.SugaredLogger
	didReceivePong  bool
	pingTimer       *time.Timer
	heartbeatTicker *time.Ticker
	ballData        BallData
	clubData        ClubData
	clubType        ClubType
	onShotCallback  Shared.ShotHandlerFunc
}

// NewLaunchMonitor creates a new instance of LaunchMonitor.
func NewLaunchMonitor(localIP string, port int, logger *zap.Logger, config Shared.Config) *LaunchMonitor {
	log := logger.With(zap.String("component", "LAUNCH_MONITOR"), zap.String("type", "R10")).Sugar()
	log.Info("initializing R10 server")
	return &LaunchMonitor{
		config:         config,
		localIP:        localIP,
		port:           port,
		log:            log,
		clubType:       "7Iron",
		didReceivePong: true,
	}
}

// Connect starts the TCP server and listens for connections.
func (r *LaunchMonitor) Connect() error {
	addr := fmt.Sprintf("%s:%d", r.localIP, r.port)
	r.log.Infof("binding to address %s", addr)
	ln, err := net.Listen("tcp", addr)
	r.listener = ln
	if err != nil {
		r.log.Errorf("starting server: %v", err)
		return fmt.Errorf("starting server: %w", err)
	}
	r.log.Infof("server listening at %s", addr)

	go func() {
		for {
			r.log.Infof("waiting for a new connection...")
			conn, err := r.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
					r.log.Warnf("listener closed, stopping accept loop")
					break
				}
				r.log.Errorf("accepting connection: %v", err)
				continue
			}
			r.log.Infof("connection accepted")
			r.handleConnection(conn)
		}
	}()

	return nil
}

// Close stops the server and closes the connection.
func (r *LaunchMonitor) Close() error {
	err := r.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

// handleConnection handles a new client connection using a bufio.Scanner with a custom JSON split.
func (r *LaunchMonitor) handleConnection(conn net.Conn) {
	r.log.Infof("established a new connection")
	r.client = conn

	// Start heartbeat: send a ping every heartbeatInterval.
	r.heartbeatTicker = time.NewTicker(heartbeatInterval)
	go func() {
		for range r.heartbeatTicker.C {
			r.log.Debugf("heatbeat ticker; sending ping...")
			r.sendPing()
		}
	}()

	// Use a scanner with a custom split function to detect complete JSON objects.
	scanner := bufio.NewScanner(conn)
	scanner.Split(jsonSplit)

	for scanner.Scan() {
		token := scanner.Bytes()
		r.log.Debugf("received JSON: %s", token)

		var msg Message
		if err := json.Unmarshal(token, &msg); err != nil {
			r.log.Errorf("parsing JSON: %v", err)
			continue
		}
		r.log.Infof("parsed message: %+v", msg)
		r.handleIncomingData(msg)
	}
	if err := scanner.Err(); err != nil {
		r.log.Errorf("scanner error: %v", err)
	}
	r.log.Infof("scanner loop ended, disconnecting")
	r.handleDisconnect()
}

// handleIncomingData processes messages by type.
func (r *LaunchMonitor) handleIncomingData(msg Message) {
	r.log.Infof("processing message type: %s", msg.Type)
	r.log.Debugf("processing message: %v", msg)

	switch msg.Type {
	case "Handshake":
		r.log.Infof("received Handshake")
		r.sendMessage(getHandshakeMessage(1))
	case "Challenge":
		r.log.Infof("received Challenge")
		r.sendMessage(getHandshakeMessage(2))
	case "Close":
		r.log.Infof("received Close")
		r.handleDisconnect()
	case "Pong":
		r.log.Infof("received Pong")
		r.handlePong()
	case "SetClubType":
		r.log.Infof("received SetClubType: %s", msg.ClubType)
		r.updateClubType(msg.ClubType)
	case "SetBallData":
		r.log.Infof("received SetBallData: %+v", msg.BallData)
		r.setBallData(msg.BallData)
	case "SetClubData":
		r.log.Infof("received SetClubData: %+v", msg.ClubData)
		r.setClubData(msg.ClubData)
	case "SendShot":
		r.log.Infof("received SendShot")
		go r.sendShot()
	default:
		r.log.Warnf("no match for message type: %s", msg.Type)
	}
}

// updateClubType updates the club type and replies with success.
func (r *LaunchMonitor) updateClubType(clubType ClubType) {
	r.log.Infof("changing club type to %s", clubType)
	r.clubType = clubType
	r.sendMessage(getSuccessMessage("SetClubType"))
}

// setBallData converts and saves ball data then replies with success.
func (r *LaunchMonitor) setBallData(bd *BallData) {
	if bd == nil {
		r.log.Errorf("received nil ball data")
		return
	}
	r.log.Infof("received ball data: %+v", bd)

	spinAxis := bd.SpinAxis

	// Modify the spinAxis based on the simulator type.
	if r.config.LaunchMonitor.Name == "GSPro" {
		// Convert from 0° to 360° to -90° to 90° range
		if spinAxis > 90 {
			spinAxis -= 360
		}
		spinAxis = spinAxis * -1
	}

	r.ballData = BallData{
		BallSpeed:       bd.BallSpeed,
		SpinAxis:        spinAxis,
		TotalSpin:       bd.TotalSpin,
		LaunchDirection: bd.LaunchDirection,
		LaunchAngle:     bd.LaunchAngle,
	}
	r.log.Infof("updated ball data: %+v", r.ballData)
	r.sendMessage(getSuccessMessage("SetBallData"))
}

// setClubData saves club data and replies with success.
func (r *LaunchMonitor) setClubData(cd *ClubData) {
	if cd == nil {
		r.log.Infof("received nil club data")
		return
	}
	r.clubData = *cd
	r.log.Infof("updated club data: %+v", r.clubData)
	r.sendMessage(getSuccessMessage("SetClubData"))
}

// sendShot sends the shot sequence.
func (r *LaunchMonitor) sendShot() {
	r.log.Infof("processing shot with ballData: %+v and clubData: %+v", r.ballData, r.clubData)
	r.sendMessage(getSuccessMessage("SendShot"))
	time.AfterFunc(300*time.Millisecond, func() {
		r.log.Infof("sending ShotComplete message")
		r.sendMessage(getShotCompleteMessage())
	})

	// Invoke the callback with the converted data.
	r.log.Infof("invoking shot callback")
	standardBall, standardClub := ConvertToStandard(r.ballData, r.clubData)
	r.onShotCallback(standardBall, standardClub, Shared.ShotDataOptions{
		ContainsBallData:          true,
		ContainsClubData:          true,
		LaunchMonitorBallDetected: true,
		ClubType:                  string(r.clubType),
	})

	// Respond to the launch monitor to ready for the next shot.
	time.AfterFunc(700*time.Millisecond, func() {
		r.log.Infof("sending Disarm command")
		r.sendMessage(getSimCommand("Disarm"))
	})
	time.AfterFunc(1000*time.Millisecond, func() {
		r.log.Infof("sending Arm command")
		r.sendMessage(getSimCommand("Arm"))
	})
}

// SetOnShotCallback sets the callback function for shot events.
func (r *LaunchMonitor) SetOnShotCallback(callback Shared.ShotHandlerFunc) {
	r.onShotCallback = callback
}

// handlePong marks that a pong was received.
func (r *LaunchMonitor) handlePong() {
	r.log.Infof("pong received, resetting ping timer")
	r.didReceivePong = true
	if r.pingTimer != nil {
		r.pingTimer.Stop()
	}
}

// sendPing sends a ping and starts a timeout.
func (r *LaunchMonitor) sendPing() {
	r.log.Infof("initiating ping")
	r.didReceivePong = false
	if r.client != nil {
		r.log.Infof("sending ping to client at %s", r.client.RemoteAddr().String())
		r.sendMessage(getSimCommand("Ping"))
		r.pingTimer = time.AfterFunc(pingTimeoutDuration, func() {
			if !r.didReceivePong {
				r.log.Warnf("ping timeout - the R10 has stopped responding")
				r.handleDisconnect()
			} else {
				r.log.Infof("pong received in time")
			}
		})
	} else {
		r.log.Warnf("no client to send ping")
	}
}

// sendMessage writes a JSON message (with newline) to the client.
func (r *LaunchMonitor) sendMessage(msg string) {
	r.log.Infof("preparing to send message: %s", msg)
	if r.client != nil {
		_, err := r.client.Write([]byte(msg + "\n"))
		if err != nil {
			r.log.Errorf("sending message: %v", err)
		} else {
			r.log.Infof("message sent successfully")
		}
	} else {
		r.log.Warnf("no client connected to send message")
	}
}

// handleDisconnect cleans up the client connection.
func (r *LaunchMonitor) handleDisconnect() {
	r.log.Infof("disconnecting client...")
	if r.client != nil {
		_ = r.client.Close()
		r.client = nil
	}
	if r.heartbeatTicker != nil {
		r.heartbeatTicker.Stop()
	}
	r.log.Infof("client disconnected")
}

// getHandshakeMessage returns a handshake message based on step.
func getHandshakeMessage(step int) string {
	if step == 1 {
		msg := map[string]interface{}{
			"Challenge":               "gQW3om37uK4OOU4FXQH9GWgljxOrNcL5MvubVHAtQC0x6Z1AwJTgAIKyamJJMzm9",
			"E6Version":               "2, 0, 0, 0",
			"ProtocolVersion":         "1.0.0.5",
			"RequiredProtocolVersion": "1.0.0.0",
			"Type":                    "Handshake",
		}
		data, _ := json.Marshal(msg)
		return string(data)
	}
	// For step 2 (Challenge response)
	msg := map[string]interface{}{
		"Success": "true",
		"Type":    "Authentication",
	}
	data, _ := json.Marshal(msg)
	return string(data)
}

// getSuccessMessage returns an ACK message.
func getSuccessMessage(subType string) string {
	msg := map[string]interface{}{
		"Details": "Success.",
		"SubType": subType,
		"Type":    "ACK",
	}
	data, _ := json.Marshal(msg)
	return string(data)
}

// getSimCommand returns a SimCommand message.
func getSimCommand(cmd string) string {
	msg := map[string]interface{}{
		"SubType": cmd,
		"Type":    "SimCommand",
	}
	data, _ := json.Marshal(msg)
	return string(data)
}

// getShotCompleteMessage returns a fixed ShotComplete message.
func getShotCompleteMessage() string {
	msg := map[string]interface{}{
		"Details": map[string]interface{}{
			"Apex": 62.2087860107422,
			"BallData": map[string]interface{}{
				"BackSpin":        4690.28662109375,
				"BallSpeed":       151.587356567383,
				"LaunchAngle":     17.7735958099365,
				"LaunchDirection": -5.00650501251221,
				"SideSpin":        -542.832092285156,
				"SpinAxis":        353.398223876953,
				"TotalSpin":       4721.59423828125,
			},
			"BallInHole":          false,
			"BallLocation":        "Fringe",
			"CarryDeviationAngle": 357.429321289063,
			"CarryDeviationFeet":  -19.5566101074219,
			"CarryDistance":       436.027191162109,
			"ClubData": map[string]interface{}{
				"ClubAngleFace":    -2.42121529579163,
				"ClubAnglePath":    -10.2835702896118,
				"ClubHeadSpeed":    110.317367553711,
				"ClubHeadSpeedMPH": 75.2163848876953,
				"ClubType":         "7Iron",
				"SmashFactor":      1.37410235404968,
			},
			"DistanceToPin":       122.404106140137,
			"TotalDeviationAngle": 356.053466796875,
			"TotalDeviationFeet":  -32.0723648071289,
			"TotalDistance":       465.995697021484,
		},
		"SubType": "ShotComplete",
		"Type":    "SimCommand",
	}
	data, _ := json.Marshal(msg)
	return string(data)
}

// LaunchShot simulates launching a shot.
func (r *LaunchMonitor) LaunchShot() {
	r.log.Infof("simulating shot launch")
	r.sendShot()
}
