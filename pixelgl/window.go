package pixelgl

import (
	"image"
	"image/color"
	"runtime"

	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/faiface/pixel"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/pkg/errors"
)

// // WindowConfig is a structure for specifying all possible properties of a Window. Properties are
// // chosen in such a way, that you usually only need to set a few of them - defaults (zeros) should
// // usually be sensible.
// //
// // Note that you always need to set the Bounds of a Window.
// type WindowConfig struct {
// 	// Title at the top of the Window.
// 	Title string

// 	// Icon specifies the icon images available to be used by the window. This is usually
// 	// displayed in the top bar of the window or in the task bar of the desktop environment.
// 	//
// 	// If passed one image, it will use that image, if passed an array of images those of or
// 	// closest to the sizes desired by the system are selected. The desired image sizes varies
// 	// depending on platform and system settings. The selected images will be rescaled as
// 	// needed. Good sizes include 16x16, 32x32 and 48x48.
// 	//
// 	// Note: Setting this value doesn't have an effect on OSX. You'll need to set the icon when
// 	// bundling your application for release.
// 	Icon []pixel.Picture

// 	// Bounds specify the bounds of the Window in pixels.
// 	Bounds pixel.Rect

// 	// If set to nil, the Window will be windowed. Otherwise it will be fullscreen on the
// 	// specified Monitor.
// 	Monitor *Monitor

// 	// Whether the Window is resizable.
// 	Resizable bool

// 	// Undecorated Window ommits the borders and decorations (close button, etc.).
// 	Undecorated bool

// 	// NoIconify specifies whether fullscreen windows should not automatically
// 	// iconify (and restore the previous video mode) on focus loss.
// 	NoIconify bool

// 	// AlwaysOnTop specifies whether the windowed mode window will be floating
// 	// above other regular windows, also called topmost or always-on-top.
// 	// This is intended primarily for debugging purposes and cannot be used to
// 	// implement proper full screen windows.
// 	AlwaysOnTop bool

// 	// VSync (vertical synchronization) synchronizes Window's framerate with the framerate of
// 	// the monitor.
// 	VSync bool
// }

type Option func(*options)

type options struct {
	// Title at the top of the Window.
	title string

	// Icon specifies the icon images available to be used by the window. This is usually
	// displayed in the top bar of the window or in the task bar of the desktop environment.
	//
	// If passed one image, it will use that image, if passed an array of images those of or
	// closest to the sizes desired by the system are selected. The desired image sizes varies
	// depending on platform and system settings. The selected images will be rescaled as
	// needed. Good sizes include 16x16, 32x32 and 48x48.
	//
	// Note: Setting this value doesn't have an effect on OSX. You'll need to set the icon when
	// bundling your application for release.
	icon []pixel.Picture

	// If set to nil, the Window will be windowed. Otherwise it will be fullscreen on the
	// specified Monitor.
	monitor *Monitor

	// Whether the Window is resizable.
	resizable int

	decorated int

	// NoIconify specifies whether fullscreen windows should not automatically
	// iconify (and restore the previous video mode) on focus loss.
	noIconify int

	// AlwaysOnTop specifies whether the windowed mode window will be floating
	// above other regular windows, also called topmost or always-on-top.
	// This is intended primarily for debugging purposes and cannot be used to
	// implement proper full screen windows.
	alwaysOnTop int

	// VSync (vertical synchronization) synchronizes Window's framerate with the framerate of
	// the monitor.
	vsync bool
}

var defaultOptions = options{
	title:       "",
	icon:        nil,
	monitor:     nil,
	resizable:   glfw.False,
	decorated:   glfw.True,
	noIconify:   glfw.False,
	alwaysOnTop: glfw.False,
	vsync:       false,
}

// Window is a window handler. Use this type to manipulate a window (input, drawing, etc.).
type Window struct {
	window             *glfw.Window
	canvas             *Canvas
	vsync              bool
	cursorVisible      bool
	cursorInsideWindow bool

	// need to save these to correctly restore a fullscreen window
	restore struct {
		xpos, ypos, width, height int
	}

	prevInp, currInp inputState

	prevJoy, currJoy [JoystickLast + 1]joystickState
}

var currWin *Window

// NewWindow creates a new Window with it's properties specified in the provided config.
//
// If Window creation fails, an error is returned (e.g. due to unavailable graphics device).
func NewWindow(width, height int, ops ...Option) (win *Window, err error) {
	o := defaultOptions

	for _, fnc := range ops {
		fnc(&o)
	}

	mainthread.Call(func() {
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		glfw.WindowHint(glfw.Resizable, o.resizable)
		glfw.WindowHint(glfw.Decorated, o.decorated)
		glfw.WindowHint(glfw.Floating, o.alwaysOnTop)
		glfw.WindowHint(glfw.AutoIconify, o.noIconify)

		var share *glfw.Window
		if currWin != nil {
			share = currWin.window
		}

		win = &Window{}

		win.window, err = glfw.CreateWindow(
			width,
			height,
			o.title,
			nil,
			share,
		)
		if err != nil {
			return
		}

		// enter the OpenGL context
		win.begin()
		glhf.Init()
		win.end()
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating window failed")
	}

	if len(o.icon) > 0 {
		imgs := make([]image.Image, len(o.icon))
		for i, icon := range o.icon {
			pic := pixel.PictureDataFromPicture(icon)
			imgs[i] = pic.Image()
		}
		mainthread.Call(func() {
			win.window.SetIcon(imgs)
		})
	}

	win.SetVSync(o.vsync)

	win.initInput()
	win.SetMonitor(o.monitor)

	win.canvas = NewCanvas(pixel.R(0, 0, float64(width), float64(height)))
	win.Update()

	runtime.SetFinalizer(win, (*Window).Destroy)

	return win, nil
}

// Destroy destroys the Window. The Window can't be used any further.
func (w *Window) Destroy() {
	mainthread.Call(func() {
		w.window.Destroy()
	})
}

// Update swaps buffers and polls events. Call this method at the end of each frame.
func (w *Window) Update() {
	width, height := w.window.GetSize()
	w.canvas.SetBounds(pixel.R(0, 0, float64(width), float64(height)))

	mainthread.Call(func() {
		w.begin()

		framebufferWidth, framebufferHeight := w.window.GetFramebufferSize()
		glhf.Bounds(0, 0, framebufferWidth, framebufferHeight)

		glhf.Clear(0, 0, 0, 0)
		w.canvas.gf.Frame().Begin()
		w.canvas.gf.Frame().Blit(
			nil,
			0, 0, w.canvas.Texture().Width(), w.canvas.Texture().Height(),
			0, 0, framebufferWidth, framebufferHeight,
		)
		w.canvas.gf.Frame().End()

		if w.vsync {
			glfw.SwapInterval(1)
		} else {
			glfw.SwapInterval(0)
		}
		w.window.SwapBuffers()
		w.end()
	})

	w.UpdateInput()
}

// SetClosed sets the closed flag of the Window.
//
// This is useful when overriding the user's attempt to close the Window, or just to close the
// Window from within the program.
func (w *Window) SetClosed(closed bool) {
	w.window.SetShouldClose(closed)
}

// Closed returns the closed flag of the Window, which reports whether the Window should be closed.
//
// The closed flag is automatically set when a user attempts to close the Window.
func (w *Window) Closed() bool {
	return w.window.ShouldClose()
}

// SetTitle changes the title of the Window.
func (w *Window) SetTitle(title string) {
	mainthread.Call(func() {
		w.window.SetTitle(title)
	})
}

// SetSize sets the width and height of the window
func (w *Window) SetSize(width, height int) {
	mainthread.Call(func() {
		w.window.SetSize(width, height)
	})
}

// SetPos sets the position, in screen coordinates, of the upper-left corner
// of the client area of the window. Position can be fractional, but the actual position
// of the window will be rounded to integers.
//
// If it is a full screen window, this function does nothing.
func (w *Window) SetPos(pos pixel.Vec) {
	mainthread.Call(func() {
		w.window.SetPos(int(pos.X), int(pos.Y))
	})
}

// GetPos gets the position, in screen coordinates, of the upper-left corner
// of the client area of the window. The position is rounded to integers.
func (w *Window) GetPos() pixel.Vec {
	x, y := w.window.GetPos()
	return pixel.V(float64(x), float64(y))
}

// Bounds returns the current bounds of the Window.
func (w *Window) Bounds() pixel.Rect {
	width, height := w.window.GetSize()
	return pixel.R(0, 0, float64(width), float64(height))
}

func (w *Window) setFullscreen(monitor *Monitor) {
	w.restore.xpos, w.restore.ypos = w.window.GetPos()
	w.restore.width, w.restore.height = w.window.GetSize()

	mode := monitor.monitor.GetVideoMode()

	w.window.SetMonitor(
		monitor.monitor,
		0,
		0,
		mode.Width,
		mode.Height,
		mode.RefreshRate,
	)
}

func (w *Window) setWindowed() {
	w.window.SetMonitor(
		nil,
		w.restore.xpos,
		w.restore.ypos,
		w.restore.width,
		w.restore.height,
		0,
	)
}

// SetMonitor sets the Window fullscreen on the given Monitor. If the Monitor is nil, the Window
// will be restored to windowed state instead.
//
// The Window will be automatically set to the Monitor's resolution. If you want a different
// resolution, you will need to set it manually with SetBounds method.
func (w *Window) SetMonitor(monitor *Monitor) {
	if w.Monitor() != monitor {
		if monitor != nil {
			w.setFullscreen(monitor)
		} else {
			w.setWindowed()
		}
	}
}

// Monitor returns a monitor the Window is fullscreen on. If the Window is not fullscreen, this
// function returns nil.
func (w *Window) Monitor() *Monitor {
	if monitor := w.window.GetMonitor(); monitor != nil {
		return &Monitor{
			monitor: monitor,
		}
	}
	return nil
}

// Focused returns true if the Window has input focus.
func (w *Window) Focused() bool {
	return w.window.GetAttrib(glfw.Focused) == glfw.True
}

// SetVSync sets whether the Window's Update should synchronize with the monitor refresh rate.
func (w *Window) SetVSync(vsync bool) {
	w.vsync = vsync
}

// VSync returns whether the Window is set to synchronize with the monitor refresh rate.
func (w *Window) VSync() bool {
	return w.vsync
}

// SetCursorVisible sets the visibility of the mouse cursor inside the Window client area.
func (w *Window) SetCursorVisible(visible bool) {
	w.cursorVisible = visible
	if visible {
		w.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		w.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
}

// CursorVisible returns the visibility status of the mouse cursor.
func (w *Window) CursorVisible() bool {
	return w.cursorVisible
}

func (w *Window) begin() {
	if currWin != w {
		w.window.MakeContextCurrent()
		currWin = w
	}
}

func (w *Window) end() {
	glfw.DetachCurrentContext()
}

// MakeTriangles generates a specialized copy of the supplied Triangles that will draw onto this
// Window.
//
// Window supports TrianglesPosition, TrianglesColor and TrianglesPicture.
func (w *Window) MakeTriangles(t pixel.Triangles) pixel.TargetTriangles {
	return w.canvas.MakeTriangles(t)
}

// MakePicture generates a specialized copy of the supplied Picture that will draw onto this Window.
//
// Window supports PictureColor.
func (w *Window) MakePicture(p pixel.Picture) pixel.TargetPicture {
	return w.canvas.MakePicture(p)
}

// SetMatrix sets a Matrix that every point will be projected by.
func (w *Window) SetMatrix(m pixel.Matrix) {
	w.canvas.SetMatrix(m)
}

// SetColorMask sets a global color mask for the Window.
func (w *Window) SetColorMask(c color.Color) {
	w.canvas.SetColorMask(c)
}

// SetComposeMethod sets a Porter-Duff composition method to be used in the following draws onto
// this Window.
func (w *Window) SetComposeMethod(cmp pixel.ComposeMethod) {
	w.canvas.SetComposeMethod(cmp)
}

// SetSmooth sets whether the stretched Pictures drawn onto this Window should be drawn smooth or
// pixely.
func (w *Window) SetSmooth(smooth bool) {
	w.canvas.SetSmooth(smooth)
}

// Smooth returns whether the stretched Pictures drawn onto this Window are set to be drawn smooth
// or pixely.
func (w *Window) Smooth() bool {
	return w.canvas.Smooth()
}

// Clear clears the Window with a single color.
func (w *Window) Clear(c color.Color) {
	w.canvas.Clear(c)
}

// Color returns the color of the pixel over the given position inside the Window.
func (w *Window) Color(at pixel.Vec) pixel.RGBA {
	return w.canvas.Color(at)
}

// Canvas returns the window's underlying Canvas
func (w *Window) Canvas() *Canvas {
	return w.canvas
}
