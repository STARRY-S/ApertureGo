package shader

import (
	"fmt"
	"os"
	"strings"

	"github.com/STARRY-S/aperture/utils"
	"github.com/engoengine/glm"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShaderObj struct {
	vertex   string
	fragment string
	geometry string
	id       uint32
}

// Load function loads shader program from file, the geometry shader can be empty
func (s *ShaderObj) Load(vs, fs, gs string) error {
	if s == nil {
		return utils.ErrInvalidPointer
	}
	if vs == "" || fs == "" {
		return utils.ErrInvalidFilePath
	}
	s.vertex, s.fragment, s.geometry = vs, fs, gs

	vData, err := os.ReadFile(vs)
	if err != nil {
		return err
	}
	if len(vData) <= 0 {
		return fmt.Errorf("Empty file [%v]", vs)
	}
	fData, err := os.ReadFile(fs)
	if err != nil {
		return err
	}
	if len(fData) <= 0 {
		return fmt.Errorf("Empty file [%v]", fs)
	}
	// geometry shader can be empty
	var gData []byte
	if gs != "" {
		gData, err = os.ReadFile(gs)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	shaderID, err := newProgram(string(vData), string(fData), string(gData))
	if err != nil {
		return err
	}
	if shaderID == 0 {
		return fmt.Errorf("failed to create shader program")
	}
	s.id = shaderID

	return nil
}

// LoadString function loads shader program from memory, geometry shader can be empty
func (s *ShaderObj) LoadMemory(vs, fs, gs string) error {
	shaderID, err := newProgram(vs, fs, gs)
	if err != nil {
		return err
	}
	if shaderID == 0 {
		return fmt.Errorf("failed to create shader program")
	}
	s.id = shaderID

	return nil
}

// Set function set values to shader program
func (s ShaderObj) Set(name string, value interface{}) error {
	if s.id == 0 {
		// shader not initialized, return
		return nil
	}
	var location int32
	location = gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	if location < 0 {
		return fmt.Errorf("Failed to set [%T] for [%s] at location [%v]", value, name, location)
	}
	switch value.(type) {
	case int32:
		gl.Uniform1i(location, value.(int32))
	case uint32:
		gl.Uniform1ui(location, value.(uint32))
	case float32:
		gl.Uniform1f(location, value.(float32))
	case glm.Vec2:
		v := value.(glm.Vec2)
		gl.Uniform2fv(location, 1, &v[0])
	case glm.Vec3:
		v := value.(glm.Vec3)
		gl.Uniform3fv(location, 1, &v[0])
	case glm.Vec4:
		v := value.(glm.Vec4)
		gl.Uniform4fv(location, 1, &v[0])
	case glm.Mat2:
		m := value.(glm.Mat2)
		gl.UniformMatrix2fv(location, 1, false, &m[0])
	case glm.Mat3:
		m := value.(glm.Mat3)
		gl.UniformMatrix3fv(location, 1, false, &m[0])
	case glm.Mat4:
		m := value.(glm.Mat4)
		gl.UniformMatrix4fv(location, 1, false, &m[0])
	case glm.Mat2x3:
		m := value.(glm.Mat2x3)
		gl.UniformMatrix2x3fv(location, 1, false, &m[0])
	case glm.Mat3x4:
		m := value.(glm.Mat3x4)
		gl.UniformMatrix3x4fv(location, 1, false, &m[0])
	default:
		return fmt.Errorf("Unsupported value type [%T]", value)
	}
	return nil
}

func newProgram(vertexShaderSource, fragmentShaderSource, geometryShaderSource string) (uint32, error) {
	var vertexShader, fragmentShader, geometryShader uint32
	var err error

	vertexShader, err = compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err = compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	if geometryShaderSource != "" {
		geometryShader, err = compileShader(geometryShaderSource, gl.GEOMETRY_SHADER)
		if err != nil {
			return 0, err
		}
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	if geometryShader != 0 {
		gl.AttachShader(program, geometryShader)
	}
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v:\n%v", source, log)
	}

	return shader, nil
}

func (s ShaderObj) GetID() uint32 {
	return s.id
}
