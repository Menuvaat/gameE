package renderer

/*

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	glm "github.com/go-gl/mathgl/mgl32"
)

type Camera_Movement int

const (
	FORWARD Camera_Movement = iota
	BACKWARD
	LEFT,
	RIGHT
)

const (
	YAW = float32(-90.)
	PITCH = float32(0.)
	SPEED = float32(2.5)
	SENSITIVITY = float32(0.1)
	ZOOM = float32(45.)
)

struct Camera {
	var Position glm.Vec3
	var Front glm.Vec3
	var Up glm.Vec3
	var Right glm.Vec3
	var WorldUp glm.Vec3

	var Yaw float32
	var Pitch float32
	var MovementSpeed float32
	var MouseSensitivity float32
	var Zoom float32
}


func NewCam(position := glm.Vec3{0., 0., 0.}, up := glm.Vec3{0., 1., 0.}, yaw := YAW, pitch := PITCH) (*Camera, error) {
	return &Camera{
		Position: position,
		WorldUp: up,
		Yaw: yaw,
		Pitch: pitch,
		update
	}
}

func (c *Camera) updateCameraVectors() {
	var front glm.Vec3

	yawRad := glm.DegToRad(c.Yaw)
	pitchRad := glm.DegToRad(c.Pitch)

	front[0] = float32(math.Cos(float64(yawRad)) * math.Cos(float64(pitchRad)))
	front[1] = float32(math.Sin(float64(pitchRad)))
	front[2] = float32(math.Sin(float64(yawRad)) * math.Cos(float64(pitchRad)))

	c.Front = front.No
}
*/
