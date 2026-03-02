package renderer

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // Register decoders
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/bloeys/assimp-go/asig"
	"github.com/go-gl/gl/v3.3-core/gl"
	glm "github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	meshes          []Mesh
	directory       string
	textures_loaded []Texture
}

func NewModel(path string) *Model {
	m := &Model{}
	m.LoadModel(path)
	return m
}

func (m *Model) Draw(shader Shader) {
	for i := 0; i < len(m.meshes); i++ {
		m.meshes[i].Draw(shader)
	}
}

func (m *Model) LoadModel(path string) {
	// We load the model
	scene, release, err := asig.ImportFile(path, asig.PostProcessTriangulate|asig.PostProcessFlipUVs)

	if err != nil {
		fmt.Printf("ERROR::ASSIMP::%s\n", err.Error())
		return
	}
	// We check if the scene and the root node of the scene are not null adn check one of its flags to see if the returned data is incomplete
	if scene.Flags&asig.SceneFlagIncomplete != 0 || scene.RootNode == nil {
		fmt.Println("ERROR::ASSIMP::Scene is incomplete or has no root node")
		release()
		return
	}
	defer release()

	m.directory = filepath.Dir(path)
	// If all is good, we process all of the scene's nodes
	m.ProcessNode(scene.RootNode, scene) // We pass the root node, to process this node, and then process its children nodes
}

func (m *Model) ProcessNode(node *asig.Node, scene *asig.Scene) {
	// Process all the node's meshes if any
	for i := 0; i < len(node.MeshIndicies); i++ {
		mesh := scene.Meshes[node.MeshIndicies[i]]
		m.meshes = append(m.meshes, *m.ProcessMesh(mesh, scene))
	}
	// Then do the same for each of its children
	for i := 0; i < len(node.Children); i++ {
		m.ProcessNode(node.Children[i], scene)
	}
}

func (m *Model) ProcessMesh(mesh *asig.Mesh, scene *asig.Scene) *Mesh {
	var vertices []Vertex
	var indices []uint32
	var textures []Texture

	for i := 0; i < len(mesh.Vertices); i++ {
		var vertex Vertex
		var vector glm.Vec3

		// We set the vertex positions of the mesh
		vector = glm.Vec3{mesh.Vertices[i].X(), mesh.Vertices[i].Y(), mesh.Vertices[i].Z()}
		vertex.Position = vector
		// We set the normals of the mesh
		vector = glm.Vec3{mesh.Normals[i].X(), mesh.Normals[i].Y(), mesh.Normals[i].Z()}
		vertex.Normal = vector
		// Setting the texture coordinates of the mesh
		if mesh.TexCoords[0] != nil { // Does the mesh contain texture coordinates?
			var vec glm.Vec2
			vec = glm.Vec2{mesh.TexCoords[0][i].X(), mesh.TexCoords[0][i].Y()}
			vertex.TexCoords = vec
		} else {
			vertex.TexCoords = glm.Vec2{0.0, 0.0}
		}
		// We add the vertex to the vector
		vertices = append(vertices, vertex)
	}
	// We iterate through all the mesh and get the indices of the vertices to know in which order they have to be drawn
	for i := 0; i < len(mesh.Faces); i++ {
		face := mesh.Faces[i]
		for j := 0; j < len(face.Indices); j++ {
			indices = append(indices, uint32(face.Indices[j]))
		}
	}
	// Here we get all the materials and textures from the model, the diffuse and specular maps, and we add all them to the textures vector
	if mesh.MaterialIndex >= 0 {
		var material *asig.Material = scene.Materials[mesh.MaterialIndex]
		var diffuseMaps []Texture = m.LoadMaterialTextures(material, asig.TextureTypeDiffuse, "texture_diffuse")
		textures = append(textures, diffuseMaps...)

		var specularMaps []Texture = m.LoadMaterialTextures(material, asig.TextureTypeSpecular, "texture_specular")
		textures = append(textures, specularMaps...)
	}
	// Finally we create a mesh with all the data saved early
	return NewMesh(vertices, indices, textures)
}

// Function to load the textures from the model
func (m *Model) LoadMaterialTextures(mat *asig.Material, mType asig.TextureType, typeName string) []Texture {
	var textures []Texture
	count := asig.GetMaterialTextureCount(mat, mType)
	// We iterate through all the textures
	for i := 0; i < count; i++ {
		// We get the path of the textures
		path, err := asig.GetMaterialTexture(mat, mType, uint(i))
		// This var is to skip any texture if it is repeated
		var skip bool = false
		if err != nil {
			fmt.Printf("%s", "Error loading the texture from path: "+path.Path)
			continue
		}
		// It checks if the texture that we have saved is the same that we have saved in our textures_loaded variable, If it is repeated, we load that saved texture
		for j := 0; j < len(m.textures_loaded); j++ {
			if m.textures_loaded[j].path == path.Path {
				textures = append(textures, m.textures_loaded[j])
				skip = true
				break
			}
		}
		// If the texture doesn't have to be skipped we get access to the texture features and save them in our texture var, then we add this texture to the vector of textures and to the vector of loaded textures to not load it again
		if !skip {
			var texture Texture
			id, err := TextureFromFile(path.Path, m.directory)
			if err != nil {
				fmt.Printf("Failed to load texture from the file: %s", path.Path)
			}
			texture.id = uint(id)
			texture.textureType = typeName
			texture.path = path.Path
			textures = append(textures, texture)
			m.textures_loaded = append(m.textures_loaded, texture)
		}
	}
	// We return all the textures
	return textures
}

// Function that reads the textures and processes them
func TextureFromFile(path string, directory string) (uint32, error) {
	// We get the file name of the texture
	fileName := filepath.Join(directory, path)

	// We open and decode the image
	imgFile, err := os.Open(fileName)
	if err != nil {
		return 0, fmt.Errorf("Failed to open texture file %v", err)
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, fmt.Errorf("Failed to decode image %v", err)
	}

	// Convert the image saved into RGBA pixels
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Flip the pixels vertically
	flipVertical(rgba)

	// We generate the texture
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	// Tells OpenGL that rows are not necessarily padded to 4 bytes
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	// We save the width, height and pixels of the texture
	width := int32(rgba.Bounds().Dx())
	height := int32(rgba.Bounds().Dy())
	pix := gl.Ptr(rgba.Pix)

	// We set the values for OpenGL to know how the texture is
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, pix)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// Set parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// We return the texture
	return textureID, nil
}

func flipVertical(rgba *image.RGBA) {
	height := rgba.Rect.Dy()
	// Calculates the stride of each row
	rowStride := rgba.Stride
	tmpRow := make([]uint8, rowStride)

	for i := 0; i < height/2; i++ {
		topRow := i * rowStride
		bottomRow := (height - 1 - i) * rowStride
		// We swap the entire row of bytes at once
		copy(tmpRow, rgba.Pix[topRow:topRow+rowStride])                                  // It takes the top row of the image and copies it into a temporary buffer
		copy(rgba.Pix[topRow:topRow+rowStride], rgba.Pix[bottomRow:bottomRow+rowStride]) // It takes the bottom row of the image and copies it into the top row's memory space
		copy(rgba.Pix[bottomRow:bottomRow+rowStride], tmpRow)                            // It takes the data that we saved in the temporary buffer and copies it into the bottom row's memory space
	}
}
