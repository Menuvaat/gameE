package renderer

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	glm "github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position  glm.Vec3
	Normal    glm.Vec3
	TexCoords glm.Vec2
}

type Texture struct {
	id          uint
	textureType string
	path        string
}

type Mesh struct {
	Vertices      []Vertex
	Indices       []uint32
	Textures      []Texture
	vao, vbo, ebo uint32
}

// Constructor function for the mesh to assign the different values on the mesh vectors
func NewMesh(vertices []Vertex, indices []uint32, textures []Texture) *Mesh {
	m := Mesh{
		Vertices: vertices,
		Indices:  indices,
		Textures: textures,
	}
	m.SetupMesh()
	return &m
}

func (m *Mesh) SetupMesh() {
	// We set the vertex array object that stores all the information from the vertices
	gl.GenVertexArrays(1, &m.vao)
	// We create the vertex buffer objects that stores the amount of vertices
	gl.GenBuffers(1, &m.vbo)
	// We create the element buffer object it stores the indices of the vertices to define how they should be connected
	gl.GenBuffers(1, &m.ebo)

	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	if len(m.Vertices) > 0 {
		// Stores data into the Vertex buffer object, with the created vertices
		gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*int(unsafe.Sizeof(Vertex{})), unsafe.Pointer(&m.Vertices[0]), gl.STATIC_DRAW)
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	// Stores data into the element buffer object, with the indices of the vertices
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*int(unsafe.Sizeof(uint32(0))), unsafe.Pointer(&m.Indices[0]), gl.STATIC_DRAW)

	stride := int32(unsafe.Sizeof(Vertex{}))
	// Set vertex positions
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, 0)

	// Vertex normals
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, stride, unsafe.Offsetof(Vertex{}.Normal))

	// Texture coords
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, stride, unsafe.Offsetof(Vertex{}.TexCoords))

	gl.BindVertexArray(0)
}

func (m *Mesh) Draw(shader Shader) {
	// It calculates the n-component per texture type and concatenates them to the texture's type string to get the appropiate uniform name
	var diffuseNr uint = 1
	var specularNr uint = 1
	for i := 0; i < len(m.Textures); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i)) // Activates the proper texture unit befor binding it
		var number string
		name := m.Textures[i].textureType
		// Here we check which type of texture is, if it is diffuse or specular and increase the number of them to have all them saved to use
		if name == "texture_diffuse" {
			number = fmt.Sprintf("%v", diffuseNr)
			diffuseNr++
		} else if name == "texture_specular" {
			number = fmt.Sprintf("%v", specularNr)
			specularNr++
		}
		// We locate the appropiate sampler and bind the texture
		shader.SetInt(("material." + name + number), i)
		gl.BindTexture(gl.TEXTURE_2D, uint32(m.Textures[i].id))
	}
	gl.ActiveTexture(gl.TEXTURE0)

	// draw mesh
	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}
