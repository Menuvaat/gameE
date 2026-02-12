package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"GEngineGo/renderer"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	glm "github.com/go-gl/mathgl/mgl32"
	"neilpa.me/go-stbi"
)

const (
	wWidth  = 800
	wHeight = 600
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

	offsetX   float32 = 0.
	offsetY   float32 = 0.
	offsetZ   float32 = -3.
	moveSpeed float32 = 0.015
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize GLFW: %v", err))
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(wWidth, wHeight, wTitle, nil, nil)
	if err != nil {
		glfw.Terminate()
		panic(fmt.Sprintf("Failed to create GLFW window: %v", err))
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	if err := gl.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize OpenGL: %v", err))
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version: ", version)

	glfw.SwapInterval(1)
	gl.Enable(gl.DEPTH_TEST)

	//---------------------
	//Texture
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	var width, height int32
	data, err := stbi.Load("textures/container.jpg")
	if err != nil {
		panic(fmt.Sprintf("Failed to load texture: %v", err))
	}

	width = int32(data.Rect.Dx())
	height = int32(data.Rect.Dy())

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&data.Pix[0]))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	//-----------------------------

	shader0, err := renderer.NewShader("shaders/shader.vert", "shaders/shader.frag")
	if err != nil {
		panic(fmt.Sprintf("Shader creation failed: %v", err))
	}
	defer shader0.Delete()

	//-------------------------

	modelLoc := gl.GetUniformLocation(shader0.ID, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(shader0.ID, gl.Str("view\x00"))
	projectionLoc := gl.GetUniformLocation(shader0.ID, gl.Str("projection\x00"))

	//-------------------------

	var VBO, VAO uint32
	gl.GenBuffers(1, &VBO)
	//gl.GenBuffers(1, &EBO)
	gl.GenVertexArrays(1, &VAO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, unsafe.Pointer(&vertices[0]), gl.STATIC_DRAW)

	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	//gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, unsafe.Pointer(&indices[0]), gl.STATIC_DRAW)

	//Position atribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 20, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	//Texture atribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	//Main loop
	for !window.ShouldClose() {
		processInput(window)
		width, height := window.GetFramebufferSize()

		gl.ClearColor(0.2, 0.3, 0.3, 1.)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//Matrices
		model := glm.HomogRotate3D(float32(glfw.GetTime())*glm.DegToRad(50.), glm.Vec3{0.5, 1., 0.}.Normalize())
		view := glm.Translate3D(offsetX, offsetY, offsetZ)
		projection := glm.Perspective(glm.DegToRad(45.), float32(width)/float32(height), 0.1, 100.)

		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
		gl.UniformMatrix4fv(projectionLoc, 1, false, &projection[0])

		//-------------------------------------------

		gl.BindTexture(gl.TEXTURE_2D, texture)
		shader0.Use()
		gl.BindVertexArray(VAO)

		for i, pos := range cubePos {
			model = glm.Translate3D(pos[0], pos[1], pos[2])
			angle := 20. * i
			model = glm.HomogRotate3D(glm.DegToRad(float32(angle)), glm.Vec3{1., 0.3, 0.5}.Normalize()).Mul4(model)
			gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}

	fmt.Println("Window closed cleanly.")
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
		return
	}

	if window.GetKey(glfw.KeyW) == glfw.Press {
		offsetZ += moveSpeed
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		offsetZ -= moveSpeed
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		offsetX += moveSpeed
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		offsetX -= moveSpeed
	}
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		offsetY -= moveSpeed
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		offsetY += moveSpeed
	}

	/*
		if offsetX < -1.2 {
			offsetX = -1.2
		}
		if offsetX > 1.2 {
			offsetX = 1.2
		}
		if offsetY < -1.2 {
			offsetY = -1.2
		}
		if offsetY > 1.2 {
			offsetY = 1.2
		}
	*/
}

func framebufferSizeCallback(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}
