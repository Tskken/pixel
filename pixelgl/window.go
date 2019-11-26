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

// Option is the type for option functions given to NewWindow() for configuring the window settings on startup.
type Option func(*options)

type options struct {
	title       string
	icon        []pixel.Picture
	monitor     *Monitor
	resizable   int
	decorated   int
	noIconify   int
	alwaysOnTop bool
	vsync       bool

	transparentFramebuffer int
}

var defaultOptions = options{
	resizable:              glfw.False,
	decorated:              glfw.True,
	noIconify:              glfw.False,
	transparentFramebuffer: glfw.False,
}

// SetTitle is an option function for NewWindow which sets the windows Title.
//
// Default is an empty string.
func SetTitle(title string) Option {
	return func(o *options) {
		o.title = title
	}
}

// SetIcon is the option function for NewWindow to set the icon for a window.
//
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
//
// By default no icon is active and the value for icon will just be a nil slice.
func SetIcon(icon []pixel.Picture) Option {
	return func(o *options) {
		o.icon = icon
	}
}

// SetMonitor is the option function for NewWindow to set the monitor of the progarm.
//
// By default monitor is nil which sets the window to windowed mode.
// Using this function with any valid value for monitor other then nill will set the window
// to fullscreen on that monitor
func SetMonitor(monitor *Monitor) Option {
	return func(o *options) {
		o.monitor = monitor
	}
}

// Resizable is the option function for NewWindow to sets the window to be resizable.
//
// By default windows are not resizable and this needs to be given to NewWindow to set a window
// to be resizable.
func Resizable() Option {
	return func(o *options) {
		o.resizable = glfw.True
	}
}

// Undecorated is the option function for NewWindow to set the window to be Undecorated.
//
// This means the window will not have things like window borders, close and minimize buttons, etc...
// This can be used for things like windowed borderless.
//
// By default windows are decorated, meaning they have there borders and buttons.
func Undecorated() Option {
	return func(o *options) {
		o.decorated = glfw.False
	}
}

// NoIconify is the option function for NewWindow to set the window to be noIconify.
//
// NoIconify specifies whether fullscreen windows should not automatically
// iconify (and restore the previous video mode) on focus loss.
//
// By default this value is false.
func NoIconify() Option {
	return func(o *options) {
		o.noIconify = glfw.True
	}
}

// AlwaysOnTop is the option function for NewWindow to set the window to always be on top.
//
// By default this value is false.
func AlwaysOnTop() Option {
	return func(o *options) {
		o.alwaysOnTop = true
	}
}

// VSyncEnabled is the option function for NewWindow to enable VSync.
//
// By default VSync is off so this needs to be used to enable it.
func VSyncEnabled() Option {
	return func(o *options) {
		o.vsync = true
	}
}

// TransparentWindowEnabled is the option function for NewWindow to enable transparent windows.
//
// This specificity sets the glfw WindowHint for TransparentFramebuffer to true which enables windows to handle Alpha for background colors.
// This means you can have fully (or parshaly) transparent glfw windows if this is called.
//
// By default this value is false.
func TransparentWindowEnabled() Option {
	return func(o *options) {
		o.transparentFramebuffer = glfw.True
	}
}

// Window is a window handler. Use this type to manipulate a window (input, drawing, etc.).
type Window struct {
	window *glfw.Window

	bounds             pixel.Rect
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
func NewWindow(width, height int, options ...Option) (*Window, error) {
	w := &Window{cursorVisible: true}

	o := defaultOptions
	for _, fnc := range options {
		fnc(&o)
	}

	err := mainthread.CallErr(func() error {
		var err error

		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		glfw.WindowHint(glfw.Resizable, o.resizable)
		glfw.WindowHint(glfw.Decorated, o.decorated)
		glfw.WindowHint(glfw.Floating, o.decorated)
		glfw.WindowHint(glfw.AutoIconify, o.noIconify)
		glfw.WindowHint(glfw.TransparentFramebuffer, o.transparentFramebuffer)

		var share *glfw.Window
		if currWin != nil {
			share = currWin.window
		}
		w.window, err = glfw.CreateWindow(
			width,
			height,
			o.title,
			nil,
			share,
		)
		if err != nil {
			return err
		}

		// enter the OpenGL context
		w.begin()
		glhf.Init()
		w.end()

		return nil
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
			w.window.SetIcon(imgs)
		})
	}

	w.SetVSync(o.vsync)

	w.initInput()
	w.SetMonitor(o.monitor)

	w.canvas = NewCanvas(pixel.R(0, 0, float64(width), float64(height)))
	w.Update()

	runtime.SetFinalizer(w, (*Window).Destroy)

	return w, nil
}

// Destroy destroys the Window. The Window can't be used any further.
func (w *Window) Destroy() {
	mainthread.Call(func() {
		w.window.Destroy()
	})
}

// Update swaps buffers and polls events. Call this method at the end of each frame.
func (w *Window) Update() {
	mainthread.Call(func() {
		_, _, oldW, oldH := intBounds(w.bounds)
		newW, newH := w.window.GetSize()
		w.bounds = w.bounds.ResizedMin(w.bounds.Size().Add(pixel.V(
			float64(newW-oldW),
			float64(newH-oldH),
		)))
	})

	w.canvas.SetBounds(w.bounds)

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
	mainthread.Call(func() {
		w.window.SetShouldClose(closed)
	})
}

// Closed returns the closed flag of the Window, which reports whether the Window should be closed.
//
// The closed flag is automatically set when a user attempts to close the Window.
func (w *Window) Closed() bool {
	var closed bool
	mainthread.Call(func() {
		closed = w.window.ShouldClose()
	})
	return closed
}

// SetTitle changes the title of the Window.
func (w *Window) SetTitle(title string) {
	mainthread.Call(func() {
		w.window.SetTitle(title)
	})
}

// SetBounds sets the bounds of the Window in pixels. Bounds can be fractional, but the actual size
// of the window will be rounded to integers.
func (w *Window) SetBounds(bounds pixel.Rect) {
	w.bounds = bounds
	mainthread.Call(func() {
		_, _, width, height := intBounds(bounds)
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
		left, top := int(pos.X), int(pos.Y)
		w.window.SetPos(left, top)
	})
}

// GetPos gets the position, in screen coordinates, of the upper-left corner
// of the client area of the window. The position is rounded to integers.
func (w *Window) GetPos() pixel.Vec {
	var v pixel.Vec
	mainthread.Call(func() {
		x, y := w.window.GetPos()
		v = pixel.V(float64(x), float64(y))
	})
	return v
}

// Bounds returns the current bounds of the Window.
func (w *Window) Bounds() pixel.Rect {
	return w.bounds
}

func (w *Window) setFullscreen(monitor *Monitor) {
	mainthread.Call(func() {
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
	})
}

func (w *Window) setWindowed() {
	mainthread.Call(func() {
		w.window.SetMonitor(
			nil,
			w.restore.xpos,
			w.restore.ypos,
			w.restore.width,
			w.restore.height,
			0,
		)
	})
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
	var monitor *glfw.Monitor
	mainthread.Call(func() {
		monitor = w.window.GetMonitor()
	})
	if monitor == nil {
		return nil
	}
	return &Monitor{
		monitor: monitor,
	}
}

// Focused returns true if the Window has input focus.
func (w *Window) Focused() bool {
	var focused bool
	mainthread.Call(func() {
		focused = w.window.GetAttrib(glfw.Focused) == glfw.True
	})
	return focused
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
	mainthread.Call(func() {
		if visible {
			w.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		} else {
			w.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		}
	})
}

// CursorVisible returns the visibility status of the mouse cursor.
func (w *Window) CursorVisible() bool {
	return w.cursorVisible
}

// Note: must be called inside the main thread.
func (w *Window) begin() {
	if currWin != w {
		w.window.MakeContextCurrent()
		currWin = w
	}
}

// Note: must be called inside the main thread.
func (w *Window) end() {
	// nothing, really
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
