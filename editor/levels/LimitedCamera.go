package levels

import (
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

// LimitedCamera is a camera implementation with ranged controls for zooming and moving.
type LimitedCamera struct {
	viewportWidth, viewportHeight float32
	minZoom, maxZoom              float32
	minPos, maxPos                float32

	requestedZoomLevel       float32
	viewOffsetX, viewOffsetY float32

	viewMatrix mgl.Mat4
}

// NewLimitedCamera returns a new instance of a LimitedCamera.
func NewLimitedCamera(minZoom, maxZoom float32, minPos, maxPos float32) *LimitedCamera {
	cam := &LimitedCamera{
		viewportWidth:  1.0,
		viewportHeight: 1.0,
		minZoom:        minZoom,
		maxZoom:        maxZoom,
		minPos:         minPos,
		maxPos:         maxPos,
		viewMatrix:     mgl.Ident4()}

	return cam
}

// SetViewportSize notifies the camera how big the view is.
func (cam *LimitedCamera) SetViewportSize(width, height float32) {
	if (cam.viewportWidth != width) || (cam.viewportHeight != height) {
		cam.viewportWidth, cam.viewportHeight = width, height
		cam.updateViewMatrix()
	}
}

// ViewMatrix implements the Viewer interface.
func (cam *LimitedCamera) ViewMatrix() *mgl.Mat4 {
	return &cam.viewMatrix
}

// MoveBy adjusts the requested view offset by given delta values in world coordinates.
func (cam *LimitedCamera) MoveBy(dx, dy float32) {
	cam.MoveTo(cam.viewOffsetX+dx, cam.viewOffsetY+dy)
}

// MoveTo sets the requested view offset to the given world coordinates.
func (cam *LimitedCamera) MoveTo(worldX, worldY float32) {
	cam.viewOffsetX = cam.limitValue(worldX, -cam.maxPos, cam.minPos)
	cam.viewOffsetY = cam.limitValue(worldY, -cam.maxPos, cam.minPos)
	cam.updateViewMatrix()
}

// ZoomAt adjusts the requested zoom level by given delta, centered around given world position.
// Positive values zoom in.
func (cam *LimitedCamera) ZoomAt(levelDelta float32, x, y float32) {
	cam.requestedZoomLevel = cam.limitValue(cam.requestedZoomLevel+levelDelta, cam.minZoom, cam.maxZoom)

	focusPoint := mgl.Vec4{x, y, 0.0, 1.0}
	oldPixel := cam.viewMatrix.Mul4x1(focusPoint)

	cam.updateViewMatrix()

	newPixel := cam.viewMatrix.Mul4x1(focusPoint)
	scaleFactor := cam.scaleFactor()
	cam.MoveBy(-(newPixel[0]-oldPixel[0])/scaleFactor, +(newPixel[1]-oldPixel[1])/scaleFactor)
}

func (cam *LimitedCamera) limitValue(value float32, min, max float32) float32 {
	result := value

	if result < min {
		result = min
	}
	if result > max {
		result = max
	}

	return result
}

func (cam *LimitedCamera) scaleFactor() float32 {
	return float32(math.Pow(2.0, float64(cam.requestedZoomLevel)))
}

func (cam *LimitedCamera) updateViewMatrix() {
	scaleFactor := cam.scaleFactor()
	cam.viewMatrix = mgl.Ident4().
		Mul4(mgl.Translate3D(cam.viewportWidth/2.0, cam.viewportHeight/2.0, 0)).
		Mul4(mgl.Scale3D(scaleFactor, -scaleFactor, 1.0)).
		Mul4(mgl.Translate3D(cam.viewOffsetX, cam.viewOffsetY, 0))
}
