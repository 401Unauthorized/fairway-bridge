package Shared

// StandardizedBallData provides a standardized structure for ball data.
type StandardizedBallData struct {
	Speed         float64 `json:"Speed"`
	SpinAxis      float64 `json:"SpinAxis"`
	TotalSpin     float64 `json:"TotalSpin"`
	BackSpin      float64 `json:"BackSpin,omitempty"`
	SideSpin      float64 `json:"SideSpin,omitempty"`
	HLA           float64 `json:"HLA"`
	VLA           float64 `json:"VLA"`
	CarryDistance float64 `json:"CarryDistance,omitempty"`
}

// ApplyAdjustment applies multipliers to the ball data.
// It returns a new StandardizedBallData with the adjusted values.
func (original StandardizedBallData) ApplyAdjustment(multipliers StandardizedBallData) StandardizedBallData {
	return StandardizedBallData{
		Speed:         original.Speed * multipliers.Speed,
		SpinAxis:      original.SpinAxis * multipliers.SpinAxis,
		TotalSpin:     original.TotalSpin * multipliers.TotalSpin,
		BackSpin:      original.BackSpin * multipliers.BackSpin,
		SideSpin:      original.SideSpin * multipliers.SideSpin,
		HLA:           original.HLA * multipliers.HLA,
		VLA:           original.VLA * multipliers.VLA,
		CarryDistance: original.CarryDistance * multipliers.CarryDistance,
	}
}

// StandardizedClubData provides a standardized structure for club data.
type StandardizedClubData struct {
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

// ApplyAdjustment applies multipliers to the club data.
// It returns a new StandardizedClubData with the adjusted values.
func (original StandardizedClubData) ApplyAdjustment(multipliers StandardizedClubData) StandardizedClubData {
	return StandardizedClubData{
		Speed:                original.Speed * multipliers.Speed,
		AngleOfAttack:        original.AngleOfAttack * multipliers.AngleOfAttack,
		FaceToTarget:         original.FaceToTarget * multipliers.FaceToTarget,
		Lie:                  original.Lie * multipliers.Lie,
		Loft:                 original.Loft * multipliers.Loft,
		Path:                 original.Path * multipliers.Path,
		SpeedAtImpact:        original.SpeedAtImpact * multipliers.SpeedAtImpact,
		VerticalFaceImpact:   original.VerticalFaceImpact * multipliers.VerticalFaceImpact,
		HorizontalFaceImpact: original.HorizontalFaceImpact * multipliers.HorizontalFaceImpact,
		ClosureRate:          original.ClosureRate * multipliers.ClosureRate,
	}
}

type ModifierData struct {
	BallData StandardizedBallData `json:"ball_data"`
	ClubData StandardizedClubData `json:"club_data"`
}

var DefaultModifiers = ModifierData{
	BallData: StandardizedBallData{
		Speed:         1.0,
		SpinAxis:      1.0,
		TotalSpin:     1.0,
		BackSpin:      1.0,
		SideSpin:      1.0,
		HLA:           1.0,
		VLA:           1.0,
		CarryDistance: 1.0,
	},
	ClubData: StandardizedClubData{
		Speed:                1.0,
		AngleOfAttack:        1.0,
		FaceToTarget:         1.0,
		Lie:                  1.0,
		Loft:                 1.0,
		Path:                 1.0,
		SpeedAtImpact:        1.0,
		VerticalFaceImpact:   1.0,
		HorizontalFaceImpact: 1.0,
		ClosureRate:          1.0,
	},
}

var InteractiveModifiers = ModifierData{
	BallData: DefaultModifiers.BallData,
	ClubData: DefaultModifiers.ClubData,
}

func GetModifiers() ModifierData {
	return InteractiveModifiers
}

func SetModifiers(md ModifierData) {
	InteractiveModifiers = md
}

// ShotDataOptions provides options for shot data.
type ShotDataOptions struct {
	ContainsBallData          bool   `json:"ContainsBallData"`
	ContainsClubData          bool   `json:"ContainsClubData"`
	LaunchMonitorIsReady      bool   `json:"LaunchMonitorIsReady,omitempty"`
	LaunchMonitorBallDetected bool   `json:"LaunchMonitorBallDetected,omitempty"`
	IsHeartBeat               bool   `json:"IsHeartBeat,omitempty"`
	ClubType                  string `json:"ClubType,omitempty"`
}

// ShotHandlerFunc defines a callback function type that handles shot data.
type ShotHandlerFunc func(standardBall StandardizedBallData, standardClub StandardizedClubData, shotDataOptions ShotDataOptions)
