package GSPro

import (
	"Fairway_Bridge/Shared"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"sync/atomic"
	"time"
)

const HeartbeatInterval = 5 * time.Second

// Simulator handles connecting and sending shot data to GSPro.
type Simulator struct {
	IPAddress    string
	Port         int
	DeviceID     string
	Units        string
	ShotNumber   *int32
	APIVersion   string
	conn         net.Conn
	log          *zap.SugaredLogger
	shutdownChan chan struct{}
}

// NewSimulator creates a new Simulator instance.
func NewSimulator(ip string, port int, logger *zap.Logger) *Simulator {
	log := logger.With(zap.String("component", "SIMULATOR"), zap.String("type", "GSPRO")).Sugar()

	log.Infof("initializing GSPro client")
	startingShotNumber := int32(0)
	return &Simulator{
		IPAddress:    ip,
		Port:         port,
		DeviceID:     "GSPro LM 1.1",
		Units:        "Yards",
		ShotNumber:   &startingShotNumber,
		APIVersion:   "1",
		log:          log,
		shutdownChan: make(chan struct{}),
	}
}

// Connect establishes the TCP connection to GSPro.
func (g *Simulator) Connect() error {
	address := fmt.Sprintf("%s:%d", g.IPAddress, g.Port)
	g.log.Infof("connecting to GSPro at %s", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		g.log.Errorf("connecting to GSPro: %v", err)
		return err
	}
	g.conn = conn
	g.log.Infof("connected to GSPro Connect at %s", address)

	go g.readResponses()
	go g.startHeartbeat()

	return nil
}

// startHeartbeat sends a heartbeat message every HeartbeatInterval.
func (g *Simulator) startHeartbeat() {
	go func() {
		for {
			time.Sleep(HeartbeatInterval)
			_ = g.LaunchShot(Shared.StandardizedBallData{}, Shared.StandardizedClubData{}, Shared.ShotDataOptions{
				IsHeartBeat:               true,
				LaunchMonitorIsReady:      true,
				LaunchMonitorBallDetected: true,
				ContainsBallData:          false,
				ContainsClubData:          false,
			})
		}
	}()
}

// readResponses continuously reads responses from GSPro.
func (g *Simulator) readResponses() {
	reader := bufio.NewReader(g.conn) // Buffered Reader for efficiency
	var buffer bytes.Buffer           // Buffer to accumulate fragmented messages

	defer func() { // Recover from unexpected errors
		if r := recover(); r != nil {
			g.log.Errorf("panic reading response: %v", r)
		}
	}()

	for {
		select {
		case <-g.shutdownChan:
			g.log.Info("shutdown signal received; exiting run loop.")
			return
		default:
		}

		// Read into a temporary buffer
		temp := make([]byte, 1024) // Read chunks of data
		n, err := reader.Read(temp)
		if err != nil {
			if err == io.EOF { // Handle disconnection gracefully
				g.log.Warnf("connection closed by GSPro.")
				return
			}
			g.log.Errorf("reading response: %v", err)
			continue
		}

		buffer.Write(temp[:n])

		// Try to decode JSON messages
		for {
			var resp GSProResponse
			decoder := json.NewDecoder(&buffer)
			if err := decoder.Decode(&resp); err != nil {
				if err == io.EOF {
					// Incomplete JSON, wait for more data
					break
				}
				g.log.Errorf("parsing response: %v", err)
				buffer.Reset() // Clear buffer to avoid corrupted state
				break
			}

			g.log.Infof("✅ GSPro Response: Code=%d, Message=%s, Player=%+v", resp.Code, resp.Message, resp.Player)

			buffer.Reset()
		}
	}
}

// LaunchShot sends a shot message to GSPro.
func (g *Simulator) LaunchShot(ballData Shared.StandardizedBallData, clubData Shared.StandardizedClubData, shotDataOptions Shared.ShotDataOptions) error {

	ballOut, clubOut := ConvertToSimulator(ballData, clubData)

	shot := ShotMessage{
		DeviceID:        g.DeviceID,
		Units:           g.Units,
		ShotNumber:      atomic.AddInt32(g.ShotNumber, 1),
		APIversion:      g.APIVersion,
		BallData:        ballOut,
		ClubData:        clubOut,
		ShotDataOptions: shotDataOptions,
	}
	data, err := json.Marshal(shot)
	if err != nil {
		g.log.Errorf("marshaling shot message: %v", err)
		return err
	}
	g.log.Infof("sending shot: %s", data)
	if g.conn == nil {
		g.log.Errorf("no connection to GSPro")
		return fmt.Errorf("no connection")
	}
	// Write JSON followed by newline.
	_, err = g.conn.Write(append(data, '\n'))
	if err != nil {
		g.log.Errorf("sending shot message: %v", err)
		return err
	}
	g.log.Infof("shot sent successfully")
	return nil
}

// Close closes the connection.
func (g *Simulator) Close() error {
	close(g.shutdownChan)
	if g.conn != nil {
		g.log.Infof("disconnecting from GSPro")
		err := g.conn.Close()
		if err != nil {
			return err
		}
		g.conn = nil
	}
	return nil
}
