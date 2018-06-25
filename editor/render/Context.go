package render

import (
	"github.com/inkyblackness/hacked/ui/opengl"

	mgl "github.com/go-gl/mathgl/mgl32"
)

// Context provides current render properties.
type Context struct {
	OpenGL opengl.OpenGL

	ViewMatrix       *mgl.Mat4
	ProjectionMatrix mgl.Mat4
}
