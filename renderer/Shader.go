package renderer

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ID uint32
}

func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	// Read the source files
	vertexCode, err := os.ReadFile(vertexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vertex shader file: %w", err)
	}

	fragmentCode, err := os.ReadFile(fragmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read fragment shader file: %w", err)
	}

	// Compile shaders
	vertexShader, err := compileShader(string(vertexCode), gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fragmentShader, err := compileShader(string(fragmentCode), gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	// Create the program and link it
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		if logLength <= 0 {
			return nil, fmt.Errorf("shader compilation failed, but no info log available")
		}
		// Create a buffer filled with \0 bytes
		infoLog := strings.Repeat("\x00", int(logLength))

		gl.GetShaderInfoLog(shaderProgram, logLength, nil, gl.Str(infoLog))
		// We cut at first \0
		msg := strings.TrimRight(infoLog, "\x00")

		return nil, fmt.Errorf("failed to compile the shader \n%s", msg)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &Shader{ID: shaderProgram}, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	cSource, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, cSource, nil)
	gl.CompileShader(shader)
	free()

	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		if logLength <= 0 {
			return 0, fmt.Errorf("shader compilation failed, but no info log available")
		}
		// Create a buffer filled with \0 bytes
		infoLog := strings.Repeat("\x00", int(logLength))

		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(infoLog))
		// We cut at first \0
		msg := strings.TrimRight(infoLog, "\x00")

		return 0, fmt.Errorf("failed to compile the shader \n%s", msg)
	}

	return shader, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) SetBool(name string, value bool) {
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), int32(boolToInt(value)))
}

func (s *Shader) SetInt(name string, value int) {
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), int32(value))
}

func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetMat4(name string, mat mgl32.Mat4) {
	location := gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(location, 1, false, &mat[0])
}

func (s *Shader) Delete() {
	gl.DeleteProgram(s.ID)
}

func boolToInt(value bool) int32 {
	if value {
		return 1
	}
	return 0
}
