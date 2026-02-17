package renderer

import (
	"math"

	glm "github.com/go-gl/mathgl/mgl32"
)

type CameraMovement int

const (
	Forward CameraMovement = iota
	Backward
	Left
	Right
	Up
	Down
)

const (
	Yaw         = float32(-90.)
	Pitch       = float32(0.)
	Speed       = float32(2.5)
	Sensitivity = float32(0.1)
	Zoom        = float32(45.)
)

type Camera struct {
	Position glm.Vec3
	Front    glm.Vec3
	Up       glm.Vec3
	Right    glm.Vec3
	WorldUp  glm.Vec3

	Yaw   float32
	Pitch float32

	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

func NewCam(position glm.Vec3, up glm.Vec3, yaw float32, pitch float32) *Camera {
	if position == (glm.Vec3{}) {
		position = glm.Vec3{0, 0, 0}
	}
	if up == (glm.Vec3{}) {
		position = glm.Vec3{0, 1, 0}
	}
	if yaw == 0 {
		yaw = Yaw
	}
	if pitch == 0 {
		pitch = Pitch
	}

	c := &Camera{
		Position:         position,
		WorldUp:          up,
		Yaw:              yaw,
		Pitch:            pitch,
		MovementSpeed:    Speed,
		MouseSensitivity: Sensitivity,
		Zoom:             Zoom,
	}
	c.updateCameraVectors()

	return c
}

func (c *Camera) updateCameraVectors() {
	yawRad := glm.DegToRad(c.Yaw)
	pitchRad := glm.DegToRad(c.Pitch)

	// We convert spherical coordinates into cartesian
	front := glm.Vec3{
		float32(math.Cos(float64(yawRad)) * math.Cos(float64(pitchRad))),
		float32(math.Sin(float64(pitchRad))),
		float32(math.Sin(float64(yawRad)) * math.Cos(float64(pitchRad))),
	}
	// We calculate the different vectors from the camera
	c.Front = front.Normalize()
	// Right vector = Front * WorldUp
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	// Up vector = Right * Front
	c.Up = c.Right.Cross(c.Front).Normalize()
}

// Process Keyboard from the user
func (c *Camera) ProcessKeyBoard(direction CameraMovement, deltaTime float32) {
	velocity := c.MovementSpeed * deltaTime

	switch direction {
	case Forward:
		c.Position = c.Position.Add(glm.Vec3{c.Front[0] * velocity, 0.0, c.Front[2] * velocity})
	case Backward:
		c.Position = c.Position.Sub(glm.Vec3{c.Front[0] * velocity, 0.0, c.Front[2] * velocity})
	case Right:
		c.Position = c.Position.Add(glm.Vec3{c.Right[0] * velocity, 0.0, c.Right[2] * velocity})
	case Left:
		c.Position = c.Position.Sub(glm.Vec3{c.Right[0] * velocity, 0.0, c.Right[2] * velocity})
	case Up:
		c.Position = c.Position.Add(glm.Vec3{0.0, c.Up[1] * velocity, 0.0})
	case Down:
		c.Position = c.Position.Sub(glm.Vec3{0.0, c.Up[1] * velocity, 0.0})
	}
}

// Process the mouse movement to rotate the camera
func (c *Camera) ProcessMouseMovement(xoffset, yoffset float32, constrainPitch bool) {
	xoffset *= c.MouseSensitivity
	yoffset *= c.MouseSensitivity

	c.Yaw += xoffset
	c.Pitch += yoffset

	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}
	c.updateCameraVectors()
}

// Process input from the mouse scroll to zoom
func (c *Camera) ProcessMouseScroll(yoffset float32) {
	c.Zoom -= yoffset
	if c.Zoom < 1.0 {
		c.Zoom = 1.0
	}
	if c.Zoom > 45.0 {
		c.Zoom = 45.0
	}
}

// To get the view matrix
func (c *Camera) GetViewMatrix() glm.Mat4 {
	return glm.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}
