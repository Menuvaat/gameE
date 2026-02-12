package renderer

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader struct {
	ID uint32
}

func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	// Read source files
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

	// Link program
	programID := gl.CreateProgram()
	gl.AttachShader(programID, vertexShader)
	gl.AttachShader(programID, fragmentShader)
	gl.LinkProgram(programID)

	//Check linking errors
	var success int32
	gl.GetProgramiv(programID, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programID, logLength, nil, gl.Str(log))
		gl.DeleteProgram(programID)
		gl.DeleteShader(vertexShader)
		gl.DeleteShader(fragmentShader)
		return nil, fmt.Errorf("program linking error:\n%s", log)
	}

	// Cleanup
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &Shader{ID: programID}, nil 
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

func (s *Shader) Delete() {
	gl.DeleteProgram(s.ID)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csource, nil)
	free()

	gl.CompileShader(shader)

	// Check for errors 
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		gl.DeleteShader(shader)
		typ := "VERTEX"
		if shaderType == gl.FRAGMENT_SHADER {
			typ = "FRAGMENT"
		}
		return 0, fmt.Errorf("shader compilation error (%s):\n%s", typ, log)
	}

	return shader, nil 
}

func boolToInt(b bool) int32 {
	if b {
		return 1 
	}
	return 0 
}


