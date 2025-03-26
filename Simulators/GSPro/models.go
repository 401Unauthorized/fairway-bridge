package GSPro

import "Fairway_Bridge/Shared"

// BallData corresponds to GSPro Ball Data.
type BallData struct {
	Speed         float64 `json:"Speed"`
	SpinAxis      float64 `json:"SpinAxis"`
	TotalSpin     float64 `json:"TotalSpin"`
	BackSpin      float64 `json:"BackSpin,omitempty"`
	SideSpin      float64 `json:"SideSpin,omitempty"`
	HLA           float64 `json:"HLA"`
	VLA           float64 `json:"VLA"`
	CarryDistance float64 `json:"CarryDistance,omitempty"`
}

// ClubData corresponds to GSPro Club Data.
type ClubData struct {
	Speed                float64 `json:"Speed"`
	AngleOfAttack        float64 `json:"AngleOfAttack"`
	FaceToTarget         float64 `json:"FaceToTarget"`
	Lie                  float64 `json:"Lie"`
	Loft                 float64 `json:"Loft"`
	Path                 float64 `json:"Path"`
	SpeedAtImpact        float64 `json:"SpeedAtImpact"`
	VerticalFaceImpact   float64 `json:"VerticalFaceImpact"`
	HorizontalFaceImpact float64 `json:"HorizontalFaceImpact"`
	ClosureRate          float64 `json:"ClosureRate"`
}

// ShotMessage is the complete JSON structure to send a shot.
type ShotMessage struct {
	DeviceID        string                 `json:"DeviceID"`
	Units           string                 `json:"Units"`
	ShotNumber      int32                  `json:"ShotNumber"`
	APIversion      string                 `json:"APIversion"`
	BallData        BallData               `json:"BallData"`
	ClubData        ClubData               `json:"ClubData"`
	ShotDataOptions Shared.ShotDataOptions `json:"ShotDataOptions"`
}

// GSProResponse represents a typical response from GSPro Connect.
type GSProResponse struct {
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Player  *PlayerInfo `json:"Player,omitempty"`
}

// PlayerInfo contains player details.
type PlayerInfo struct {
	Handed           string  `json:"Handed"`
	Club             string  `json:"Club"`
	DistanceToTarget float64 `json:"DistanceToTarget"`
	Surface          string  `json:"Surface"`
}

// ConvertToStandard converts the simulator's BallData and ClubData to the standardized versions.
// It returns a StandardizedBallData and StandardizedClubData from the Shared package.
func ConvertToStandard(ball BallData, club ClubData) (Shared.StandardizedBallData, Shared.StandardizedClubData) {
	standardBall := Shared.StandardizedBallData{
		Speed:         ball.Speed,
		SpinAxis:      ball.SpinAxis,
		TotalSpin:     ball.TotalSpin,
		BackSpin:      ball.BackSpin,
		SideSpin:      ball.SideSpin,
		HLA:           ball.HLA,
		VLA:           ball.VLA,
		CarryDistance: ball.CarryDistance,
	}

	standardClub := Shared.StandardizedClubData{
		Speed:                club.Speed,
		AngleOfAttack:        club.AngleOfAttack,
		FaceToTarget:         club.FaceToTarget,
		Lie:                  club.Lie,
		Loft:                 club.Loft,
		Path:                 club.Path,
		SpeedAtImpact:        club.SpeedAtImpact,
		VerticalFaceImpact:   club.VerticalFaceImpact,
		HorizontalFaceImpact: club.HorizontalFaceImpact,
		ClosureRate:          club.ClosureRate,
	}

	return standardBall, standardClub
}

// ConvertToSimulator converts the standardized BallData and ClubData back to the simulator's types.
func ConvertToSimulator(ball Shared.StandardizedBallData, club Shared.StandardizedClubData) (BallData, ClubData) {
	simBall := BallData{
		Speed:         ball.Speed,
		SpinAxis:      ball.SpinAxis,
		TotalSpin:     ball.TotalSpin,
		BackSpin:      ball.BackSpin,
		SideSpin:      ball.SideSpin,
		HLA:           ball.HLA,
		VLA:           ball.VLA,
		CarryDistance: ball.CarryDistance,
	}

	simClub := ClubData{
		Speed:                club.Speed,
		AngleOfAttack:        club.AngleOfAttack,
		FaceToTarget:         club.FaceToTarget,
		Lie:                  club.Lie,
		Loft:                 club.Loft,
		Path:                 club.Path,
		SpeedAtImpact:        club.SpeedAtImpact,
		VerticalFaceImpact:   club.VerticalFaceImpact,
		HorizontalFaceImpact: club.HorizontalFaceImpact,
		ClosureRate:          club.ClosureRate,
	}

	return simBall, simClub
}
