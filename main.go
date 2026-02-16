package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"GoGame/renderer"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	glm "github.com/go-gl/mathgl/mgl32"

	"neilpa.me/go-stbi"
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

	// Texture
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// We set the texture wrapping/filtering options
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Load and generate the texture
	var width, height int32
	data, err := stbi.Load("textures/container.jpg")
	if err != nil {
		panic(fmt.Sprintf("failed to load the texture%v", err))
	}

	width = int32(data.Rect.Dx())
	height = int32(data.Rect.Dy())

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&data.Pix[0]))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// We create the shader program from our shader struct that we have created externally
	shader0, err := renderer.NewShader("Shaders/vShader.glsl", "Shaders/fShader.glsl")
	if err != nil {
		panic(fmt.Sprintf("Shader creation failed%v", err))
	}
	defer shader0.Delete()

	// We create the vertex array object that stores all the information from the vertices
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)

	// We create the vertex buffer objects that stores the amount of vertices
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	// Stores data into the Vertex buffer object, with the created vertices
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	// We link the vertex attributes and enable them
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 20, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		deltaTime = float32(currentFrame) - lastFrame
		lastFrame = float32(currentFrame)

		processInput(window)
		width, height := window.GetFramebufferSize()

		gl.ClearColor(0.2, 0.3, 0.3, 1.)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.BindTexture(gl.TEXTURE_2D, texture)

		shader0.Use()

		// We pass the projection matrix to the shader
		projection := glm.Perspective(glm.DegToRad(camera.Zoom), float32(width)/float32(height), 0.1, 100.0)
		shader0.SetMat4("projection", projection)

		// Camera/view trasformation
		view := camera.GetViewMatrix()
		shader0.SetMat4("view", view)

		gl.BindVertexArray(VAO)

		for i, pos := range cubePos {
			// We create an identity matrix
			model := mgl32.Ident4()
			// We chain the translation to it
			model = model.Mul4(mgl32.Translate3D(pos[0], pos[1], pos[2]))

			// We chaing the rotation to the matrix
			angle := float32(20.0 * float32(i))
			model = model.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{1.0, 0.3, 0.5}))

			// Send the matrix to the shader
			shader0.SetMat4("model", model)

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)

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
