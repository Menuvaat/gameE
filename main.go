package main

import (
	"fmt"
	"runtime"

	"gayEngine/renderer"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	glm "github.com/go-gl/mathgl/mgl32"
)

const (
	wWidth  = 1280
	wHeight = 720
	wTitle  = "GEngineG"
)

var (
	vertices = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	cubePos = []glm.Vec3{
		{0., 0., -1.},
		{0.2, -0.9, 0.},
		{-5., 3., 4.},
	}

	// Timing
	deltaTime float32 = 0.
	lastFrame float32 = 0.

	// Camera
	camera     *renderer.Camera
	lastX      float64
	lastY      float64
	firstMouse bool = true
)

func init() {
	// Glfw and OpenGL must run on the main thread
	runtime.LockOSThread()
	// Camera
	camera = renderer.NewCam(
		glm.Vec3{0, 0, 3},
		glm.Vec3{0, 1, 0},
		renderer.Yaw,
		renderer.Pitch,
	)
}

func main() {
	// Initialize glfw and ensure if there is any error
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// Configure glfw
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// We create a window object
	window, err := glfw.CreateWindow(wWidth, wHeight, "Hello Go", nil, nil)
	if err != nil {
		glfw.Terminate()
		panic(err)
	}
	defer window.Destroy()

	// We make the context of the specified window curreent on the calling thread
	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	window.SetCursorPosCallback(mouseCallBack)
	window.SetScrollCallback(scrollCallBack)

	// Tell GLFW to capture our mouse
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Initialize glad
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version: ", version)

	glfw.SwapInterval(1)
	gl.Enable(gl.DEPTH_TEST)

	// Model
	model0 := renderer.NewModel("resources/objects/backpack/backpack.obj")

	// We create the shader program from our shader struct that we have created externally
	shader0, err := renderer.NewShader("shaders/vShader.glsl", "shaders/fShader.glsl")
	if err != nil {
		panic(fmt.Sprintf("Shader creation failed%v", err))
	}
	defer shader0.Delete()

	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		deltaTime = float32(currentFrame) - lastFrame
		lastFrame = float32(currentFrame)

		processInput(window)
		width, height := window.GetFramebufferSize()

		gl.ClearColor(0.2, 0.3, 0.3, 1.)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shader0.Use()

		// We pass the projection matrix to the shader
		projection := glm.Perspective(glm.DegToRad(camera.Zoom), float32(width)/float32(height), 0.1, 100.0)
		shader0.SetMat4("projection", projection)

		// Camera/view trasformation
		view := camera.GetViewMatrix()
		shader0.SetMat4("view", view)

		// Render the loaded model
		model := mgl32.Ident4()
		// We set in the origin of coordinates, and scale it into 1 dimension
		translation := mgl32.Translate3D(0.0, 0.0, 0.0)
		model = model.Mul4(translation)
		scale := mgl32.Scale3D(1.0, 1.0, 1.0)
		model = model.Mul4(scale)
		// We set the model matrix and draw the model
		shader0.SetMat4("model", model)
		model0.Draw(*shader0)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	glfw.Terminate()

	fmt.Println("Window closed cleanly.")
}

// Function to process input from the user
func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera.ProcessKeyBoard(0, deltaTime)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera.ProcessKeyBoard(1, deltaTime)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera.ProcessKeyBoard(2, deltaTime)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera.ProcessKeyBoard(3, deltaTime)
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		camera.ProcessKeyBoard(4, deltaTime)
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		camera.ProcessKeyBoard(5, deltaTime)
	}
}

func framebufferSizeCallback(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

// Function to process the mouse movement
func mouseCallBack(window *glfw.Window, xposIn, yposIn float64) {
	xpos := xposIn
	ypos := yposIn

	if firstMouse {
		lastX = xpos
		lastY = ypos
		firstMouse = false
	}

	xoffset := xpos - lastX
	yoffset := lastY - ypos

	lastX = xpos
	lastY = ypos

	camera.ProcessMouseMovement(float32(xoffset), float32(yoffset), true)
}

// Function to process the mouse scroll movement
func scrollCallBack(window *glfw.Window, xoffset, yoffset float64) {
	camera.ProcessMouseScroll(float32(yoffset))
}
