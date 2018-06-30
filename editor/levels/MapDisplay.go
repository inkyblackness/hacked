package levels

import (
	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ui/input"
	"github.com/inkyblackness/hacked/ui/opengl"
)

// MapDisplay renders a level map.
type MapDisplay struct {
	context render.Context
	camera  *LimitedCamera

	background *BackgroundGrid
	mapGrid    *MapGrid

	moveCapture func(pixelX, pixelY float32)
}

// NewMapDisplay returns a new instance.
func NewMapDisplay(gl opengl.OpenGL, guiScale float32) *MapDisplay {
	tilesPerMapSide := float32(64)

	tileBaseLength := fineCoordinatesPerTileSide
	tileBaseHalf := tileBaseLength / 2.0
	camLimit := tilesPerMapSide*tileBaseLength - tileBaseHalf
	zoomShift := guiScale - 1.0
	zoomLevelMin := float32(-5) + zoomShift
	zoomLevelMax := float32(1) + zoomShift

	display := &MapDisplay{
		context: render.Context{
			OpenGL:           gl,
			ProjectionMatrix: mgl.Ident4(),
		},
		camera:      NewLimitedCamera(zoomLevelMin, zoomLevelMax, -tileBaseHalf, camLimit),
		moveCapture: func(float32, float32) {},
	}
	display.context.ViewMatrix = display.camera.ViewMatrix()
	display.background = NewBackgroundGrid(&display.context)
	display.mapGrid = NewMapGrid(&display.context)

	centerX, centerY := (tilesPerMapSide*tileBaseLength)/-2.0, (tilesPerMapSide*tileBaseLength)/-2.0
	display.camera.ZoomAt(-3+zoomShift, centerX, centerY)
	display.camera.MoveTo(centerX, centerY)

	return display
}

// Render renders the whole map display.
func (display *MapDisplay) Render(level *level.Level) {
	display.background.Render()
	display.mapGrid.Render(level)
}

// WindowResized must be called to notify of a change in window geometry.
func (display *MapDisplay) WindowResized(width int, height int) {
	display.context.ProjectionMatrix = mgl.Ortho2D(0.0, float32(width), float32(height), 0.0)
	display.camera.SetViewportSize(float32(width), float32(height))
}

func (display *MapDisplay) unprojectPixel(pixelX, pixelY float32) (x, y float32) {
	pixelVec := mgl.Vec4{pixelX, pixelY, 0.0, 1.0}
	invertedView := display.camera.ViewMatrix().Inv()
	result := invertedView.Mul4x1(pixelVec)

	return result[0], result[1]
}

// MouseButtonDown must be called when a button was pressed.
func (display *MapDisplay) MouseButtonDown(mouseX, mouseY float32, button uint32) {
	if button == input.MousePrimary {
		lastPixelX, lastPixelY := mouseX, mouseY

		display.moveCapture = func(pixelX, pixelY float32) {
			lastWorldX, lastWorldY := display.unprojectPixel(lastPixelX, lastPixelY)
			worldX, worldY := display.unprojectPixel(pixelX, pixelY)

			display.camera.MoveBy(worldX-lastWorldX, worldY-lastWorldY)
			lastPixelX, lastPixelY = pixelX, pixelY
		}
	}
}

// MouseButtonUp must be called when a button was released.
func (display *MapDisplay) MouseButtonUp(mouseX, mouseY float32, button uint32) {
	if button == input.MousePrimary {
		display.moveCapture = func(float32, float32) {}
	}
}

// MouseMoved must be called for a mouse move.
func (display *MapDisplay) MouseMoved(mouseX, mouseY float32) {
	display.moveCapture(mouseX, mouseY)
}

// MouseScrolled must be called for a mouse scroll
func (display *MapDisplay) MouseScrolled(mouseX, mouseY float32, deltaX, deltaY float32) {
	worldX, worldY := display.unprojectPixel(mouseX, mouseY)

	if deltaY < 0 {
		display.camera.ZoomAt(-0.5, worldX, worldY)
	}
	if deltaY > 0 {
		display.camera.ZoomAt(0.5, worldX, worldY)
	}
}
