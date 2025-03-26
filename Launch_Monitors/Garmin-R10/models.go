package Garmin_R10

import "Fairway_Bridge/Shared"

// BallData holds ball flight measurements.
type BallData struct {
	BallSpeed       float64 `json:"BallSpeed"`
	SpinAxis        float64 `json:"SpinAxis"`
	TotalSpin       float64 `json:"TotalSpin"`
	LaunchDirection float64 `json:"LaunchDirection"`
	LaunchAngle     float64 `json:"LaunchAngle"`
}

// ClubData holds club-related measurements.
type ClubData struct {
	ClubHeadSpeed float64 `json:"ClubHeadSpeed"`
	ClubAngleFace float64 `json:"ClubAngleFace"`
	ClubAnglePath float64 `json:"ClubAnglePath"`
}

// ClubType represents the type of golf club.
type ClubType string

const (
	Driver        ClubType = "Driver"
	ThreeWood     ClubType = "3Wood"
	FiveWood      ClubType = "5Wood"
	SevenWood     ClubType = "7Wood"
	TwoHybrid     ClubType = "2Hybrid"
	ThreeHybrid   ClubType = "3Hybrid"
	FourHybrid    ClubType = "4Hybrid"
	FiveHybrid    ClubType = "5Hybrid"
	SixHybrid     ClubType = "6Hybrid"
	OneIron       ClubType = "1Iron"
	TwoIron       ClubType = "2Iron"
	ThreeIron     ClubType = "3Iron"
	FourIron      ClubType = "4Iron"
	FiveIron      ClubType = "5Iron"
	SixIron       ClubType = "6Iron"
	SevenIron     ClubType = "7Iron"
	EightIron     ClubType = "8Iron"
	NineIron      ClubType = "9Iron"
	PitchingWedge ClubType = "PitchingWedge"
	GapWedge      ClubType = "GapWedge"
	SandWedge     ClubType = "SandWedge"
	LobWedge      ClubType = "LobWedge"
	Putter        ClubType = "Putter"
)

// ValidClubTypes is a set of known valid club types
var ValidClubTypes = map[ClubType]bool{
	Driver:        true,
	ThreeWood:     true,
	FiveWood:      true,
	SevenWood:     true,
	TwoHybrid:     true,
	ThreeHybrid:   true,
	FourHybrid:    true,
	FiveHybrid:    true,
	SixHybrid:     true,
	OneIron:       true,
	TwoIron:       true,
	ThreeIron:     true,
	FourIron:      true,
	FiveIron:      true,
	SixIron:       true,
	SevenIron:     true,
	EightIron:     true,
	NineIron:      true,
	PitchingWedge: true,
	GapWedge:      true,
	SandWedge:     true,
	LobWedge:      true,
	Putter:        true,
}

// Message represents an incoming message from the launch monitor.
type Message struct {
	Type     string    `json:"Type"`
	BallData *BallData `json:"BallData,omitempty"`
	ClubData *ClubData `json:"ClubData,omitempty"`
	ClubType ClubType  `json:"ClubType,omitempty"`
}

// ConvertToStandard converts the launch monitor's BallData and ClubData into their standardized equivalents.
func ConvertToStandard(ball BallData, club ClubData) (Shared.StandardizedBallData, Shared.StandardizedClubData) {
	standardBall := Shared.StandardizedBallData{
		Speed:     ball.BallSpeed,
		SpinAxis:  ball.SpinAxis,
		TotalSpin: ball.TotalSpin,
		HLA:       ball.LaunchDirection,
		VLA:       ball.LaunchAngle,
	}

	standardClub := Shared.StandardizedClubData{
		Speed:        club.ClubHeadSpeed,
		FaceToTarget: club.ClubAngleFace,
		Path:         club.ClubAnglePath,
	}

	return standardBall, standardClub
}

// ConvertToSimulator converts the StandardizedBallData and StandardizedClubData to the launch monitor's types.
func ConvertToSimulator(ball Shared.StandardizedBallData, club Shared.StandardizedClubData) (BallData, ClubData) {
	simBall := BallData{
		BallSpeed:       ball.Speed,
		SpinAxis:        ball.SpinAxis,
		TotalSpin:       ball.TotalSpin,
		LaunchDirection: ball.HLA,
		LaunchAngle:     ball.VLA,
	}

	simClub := ClubData{
		ClubHeadSpeed: club.Speed,
		ClubAngleFace: club.FaceToTarget,
		ClubAnglePath: club.Path,
	}

	return simBall, simClub
}
