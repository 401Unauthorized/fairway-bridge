package Cameras

// CameraController is an interface that defines the methods for controlling a camera.
type CameraController interface {
	Connect() error
	SaveLastRecording() error
	DeleteLastRecording() error
	StopCapture() error
	StartCapture() error
	LoopCapture() error
	Shutdown() error
}
