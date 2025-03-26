package Storage

import (
	"Fairway_Bridge/Shared"
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

// FileStorage implements the Storage interface for file-based storage.
type FileStorage struct {
	file   *os.File
	writer *csv.Writer
	mutex  sync.Mutex
	log    *zap.SugaredLogger
}

// NewFileStorage initializes a new FileStorage instance.
func NewFileStorage(logger *zap.Logger, config Shared.Config) (*FileStorage, error) {
	log := logger.With(zap.String("component", "STORAGE")).Sugar()
	log.Infof("creating file storage at %s", config.Bridge.ShotFile)

	// Open file for appending; create if it doesn't exist.
	file, err := os.OpenFile(config.Bridge.ShotFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(file)

	log.Infof("file storage at %s created.", config.Bridge.ShotFile)

	// If file is new, write header.
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		header := []string{
			"Timestamp", "ShotUUID", "ClubType",
			// Raw Ball Data
			"BallSpeed", "BallSpinAxis", "BallTotalSpin", "BallBackSpin",
			"BallSideSpin", "BallHLA", "BallVLA", "BallCarryDistance",
			// Raw Club Data
			"ClubSpeed", "ClubSpeedAtImpact", "ClubPath", "ClubAngleOfAttack",
			"ClubClosureRate", "ClubLie", "ClubLoft", "ClubFaceToTarget",
			"ClubVerticalFaceImpact", "ClubHorizontalFaceImpact",
			// Adjusted Ball Data (duplicate columns for modifications)
			"AdjBallSpeed", "AdjBallSpinAxis", "AdjBallTotalSpin", "AdjBallBackSpin",
			"AdjBallSideSpin", "AdjBallHLA", "AdjBallVLA", "AdjBallCarryDistance",
			// Adjusted Club Data (duplicate columns for modifications)
			"AdjClubSpeed", "AdjClubSpeedAtImpact", "AdjClubPath", "AdjClubAngleOfAttack",
			"AdjClubClosureRate", "AdjClubLie", "AdjClubLoft", "AdjClubFaceToTarget",
			"AdjClubVerticalFaceImpact", "AdjClubHorizontalFaceImpact",
		}
		if err := writer.Write(header); err != nil {
			return nil, err
		}
		writer.Flush()
	}

	return &FileStorage{file: file, writer: writer, log: log}, nil
}

// generateUUID generates a random UUID (version 4).
func generateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant (2) bits per RFC 4122
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:]), nil
}

// SaveShot saves the shot data to the file.
func (s *FileStorage) SaveShot(
	timestamp time.Time,
	ball Shared.StandardizedBallData,
	club Shared.StandardizedClubData,
	adjustedBall Shared.StandardizedBallData,
	adjustedClub Shared.StandardizedClubData,
	clubType string,
) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.log.Infof("saving shot data to file storage at %s", s.file.Name())

	// Generate unique Shot UUID
	shotUUID, err := generateUUID()
	if err != nil {
		return err
	}

	// Raw Ball Data
	ballSpeed := fmt.Sprintf("%v", ball.Speed)
	ballSpinAxis := fmt.Sprintf("%v", ball.SpinAxis)
	ballTotalSpin := fmt.Sprintf("%v", ball.TotalSpin)
	ballBackSpin := fmt.Sprintf("%v", ball.BackSpin)
	ballSideSpin := fmt.Sprintf("%v", ball.SideSpin)
	ballHLA := fmt.Sprintf("%v", ball.HLA) // Horizontal Launch Angle
	ballVLA := fmt.Sprintf("%v", ball.VLA) // Vertical Launch Angle
	ballCarryDistance := fmt.Sprintf("%v", ball.CarryDistance)

	// Raw Club Data
	clubSpeed := fmt.Sprintf("%v", club.Speed)
	clubSpeedAtImpact := fmt.Sprintf("%v", club.SpeedAtImpact)
	clubPath := fmt.Sprintf("%v", club.Path)
	clubAngleOfAttack := fmt.Sprintf("%v", club.AngleOfAttack)
	clubClosureRate := fmt.Sprintf("%v", club.ClosureRate)
	clubLie := fmt.Sprintf("%v", club.Lie)
	clubLoft := fmt.Sprintf("%v", club.Loft)
	clubFaceToTarget := fmt.Sprintf("%v", club.FaceToTarget)
	clubVerticalFaceImpact := fmt.Sprintf("%v", club.VerticalFaceImpact)
	clubHorizontalFaceImpact := fmt.Sprintf("%v", club.HorizontalFaceImpact)

	// Adjusted Ball Data
	adjBallSpeed := fmt.Sprintf("%v", adjustedBall.Speed)
	adjBallSpinAxis := fmt.Sprintf("%v", adjustedBall.SpinAxis)
	adjBallTotalSpin := fmt.Sprintf("%v", adjustedBall.TotalSpin)
	adjBallBackSpin := fmt.Sprintf("%v", adjustedBall.BackSpin)
	adjBallSideSpin := fmt.Sprintf("%v", adjustedBall.SideSpin)
	adjBallHLA := fmt.Sprintf("%v", adjustedBall.HLA)
	adjBallVLA := fmt.Sprintf("%v", adjustedBall.VLA)
	adjBallCarryDistance := fmt.Sprintf("%v", adjustedBall.CarryDistance)

	// Adjusted Club Data
	adjClubSpeed := fmt.Sprintf("%v", adjustedClub.Speed)
	adjClubSpeedAtImpact := fmt.Sprintf("%v", adjustedClub.SpeedAtImpact)
	adjClubPath := fmt.Sprintf("%v", adjustedClub.Path)
	adjClubAngleOfAttack := fmt.Sprintf("%v", adjustedClub.AngleOfAttack)
	adjClubClosureRate := fmt.Sprintf("%v", adjustedClub.ClosureRate)
	adjClubLie := fmt.Sprintf("%v", adjustedClub.Lie)
	adjClubLoft := fmt.Sprintf("%v", adjustedClub.Loft)
	adjClubFaceToTarget := fmt.Sprintf("%v", adjustedClub.FaceToTarget)
	adjClubVerticalFaceImpact := fmt.Sprintf("%v", adjustedClub.VerticalFaceImpact)
	adjClubHorizontalFaceImpact := fmt.Sprintf("%v", adjustedClub.HorizontalFaceImpact)

	// Construct row for CSV
	row := []string{
		timestamp.Format(time.RFC3339),
		shotUUID,
		clubType,
		// Raw Ball Data
		ballSpeed, ballSpinAxis, ballTotalSpin, ballBackSpin,
		ballSideSpin, ballHLA, ballVLA, ballCarryDistance,
		// Raw Club Data
		clubSpeed, clubSpeedAtImpact, clubPath, clubAngleOfAttack,
		clubClosureRate, clubLie, clubLoft, clubFaceToTarget,
		clubVerticalFaceImpact, clubHorizontalFaceImpact,
		// Adjusted Ball Data
		adjBallSpeed, adjBallSpinAxis, adjBallTotalSpin, adjBallBackSpin,
		adjBallSideSpin, adjBallHLA, adjBallVLA, adjBallCarryDistance,
		// Adjusted Club Data
		adjClubSpeed, adjClubSpeedAtImpact, adjClubPath, adjClubAngleOfAttack,
		adjClubClosureRate, adjClubLie, adjClubLoft, adjClubFaceToTarget,
		adjClubVerticalFaceImpact, adjClubHorizontalFaceImpact,
	}

	// Write data to CSV
	if err := s.writer.Write(row); err != nil {
		return err
	}
	s.writer.Flush()

	s.log.Infof("shot data saved successfully to %s", s.file.Name())
	return nil
}
