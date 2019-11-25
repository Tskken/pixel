package pixelgl

import (
	"github.com/faiface/mainthread"
	"github.com/faiface/pixel"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Pressed returns whether the given button is currently pressed down.
func (w *Window) Pressed(button int) bool {
	return w.currInp.buttons[button]
}

// JustPressed returns whether the given button has just been pressed down.
func (w *Window) JustPressed(button int) bool {
	return w.currInp.buttons[button] && !w.prevInp.buttons[button]
}

// JustReleased returns whether the given button has just been released up.
func (w *Window) JustReleased(button int) bool {
	return !w.currInp.buttons[button] && w.prevInp.buttons[button]
}

// Repeated returns whether a repeat event has been triggered on button.
//
// Repeat event occurs repeatedly when a button is held down for some time.
func (w *Window) Repeated(button int) bool {
	return w.currInp.repeat[button]
}

// MousePosition returns the current mouse position in the Window's Bounds.
func (w *Window) MousePosition() pixel.Vec {
	return w.currInp.mouse
}

// MousePreviousPosition returns the previous mouse position in the Window's Bounds.
func (w *Window) MousePreviousPosition() pixel.Vec {
	return w.prevInp.mouse
}

// SetMousePosition positions the mouse cursor anywhere within the Window's Bounds.
func (w *Window) SetMousePosition(v pixel.Vec) {
	mainthread.Call(func() {
		if (v.X >= 0 && v.X <= w.bounds.W()) &&
			(v.Y >= 0 && v.Y <= w.bounds.H()) {
			w.window.SetCursorPos(
				v.X+w.bounds.Min.X,
				(w.bounds.H()-v.Y)+w.bounds.Min.Y,
			)
			w.prevInp.mouse = v
			w.currInp.mouse = v
			w.tempInp.mouse = v
		}
	})
}

// MouseInsideWindow returns true if the mouse position is within the Window's Bounds.
func (w *Window) MouseInsideWindow() bool {
	return w.cursorInsideWindow
}

// MouseScroll returns the mouse scroll amount (in both axes) since the last call to Window.Update.
func (w *Window) MouseScroll() pixel.Vec {
	return w.currInp.scroll
}

// Typed returns the text typed on the keyboard since the last call to Window.Update.
func (w *Window) Typed() string {
	return w.currInp.typed
}

// List of all mouse buttons.
const (
	MouseButton1      = glfw.MouseButton1
	MouseButton2      = glfw.MouseButton2
	MouseButton3      = glfw.MouseButton3
	MouseButton4      = glfw.MouseButton4
	MouseButton5      = glfw.MouseButton5
	MouseButton6      = glfw.MouseButton6
	MouseButton7      = glfw.MouseButton7
	MouseButton8      = glfw.MouseButton8
	MouseButtonLast   = glfw.MouseButtonLast
	MouseButtonLeft   = glfw.MouseButtonLeft
	MouseButtonRight  = glfw.MouseButtonRight
	MouseButtonMiddle = glfw.MouseButtonMiddle
)

// List of all keyboard buttons.
const (
	KeyUnknown      = glfw.KeyUnknown
	KeySpace        = glfw.KeySpace
	KeyApostrophe   = glfw.KeyApostrophe
	KeyComma        = glfw.KeyComma
	KeyMinus        = glfw.KeyMinus
	KeyPeriod       = glfw.KeyPeriod
	KeySlash        = glfw.KeySlash
	Key0            = glfw.Key0
	Key1            = glfw.Key1
	Key2            = glfw.Key2
	Key3            = glfw.Key3
	Key4            = glfw.Key4
	Key5            = glfw.Key5
	Key6            = glfw.Key6
	Key7            = glfw.Key7
	Key8            = glfw.Key8
	Key9            = glfw.Key9
	KeySemicolon    = glfw.KeySemicolon
	KeyEqual        = glfw.KeyEqual
	KeyA            = glfw.KeyA
	KeyB            = glfw.KeyB
	KeyC            = glfw.KeyC
	KeyD            = glfw.KeyD
	KeyE            = glfw.KeyE
	KeyF            = glfw.KeyF
	KeyG            = glfw.KeyG
	KeyH            = glfw.KeyH
	KeyI            = glfw.KeyI
	KeyJ            = glfw.KeyJ
	KeyK            = glfw.KeyK
	KeyL            = glfw.KeyL
	KeyM            = glfw.KeyM
	KeyN            = glfw.KeyN
	KeyO            = glfw.KeyO
	KeyP            = glfw.KeyP
	KeyQ            = glfw.KeyQ
	KeyR            = glfw.KeyR
	KeyS            = glfw.KeyS
	KeyT            = glfw.KeyT
	KeyU            = glfw.KeyU
	KeyV            = glfw.KeyV
	KeyW            = glfw.KeyW
	KeyX            = glfw.KeyX
	KeyY            = glfw.KeyY
	KeyZ            = glfw.KeyZ
	KeyLeftBracket  = glfw.KeyLeftBracket
	KeyBackslash    = glfw.KeyBackslash
	KeyRightBracket = glfw.KeyRightBracket
	KeyGraveAccent  = glfw.KeyGraveAccent
	KeyWorld1       = glfw.KeyWorld1
	KeyWorld2       = glfw.KeyWorld2
	KeyEscape       = glfw.KeyEscape
	KeyEnter        = glfw.KeyEnter
	KeyTab          = glfw.KeyTab
	KeyBackspace    = glfw.KeyBackspace
	KeyInsert       = glfw.KeyInsert
	KeyDelete       = glfw.KeyDelete
	KeyRight        = glfw.KeyRight
	KeyLeft         = glfw.KeyLeft
	KeyDown         = glfw.KeyDown
	KeyUp           = glfw.KeyUp
	KeyPageUp       = glfw.KeyPageUp
	KeyPageDown     = glfw.KeyPageDown
	KeyHome         = glfw.KeyHome
	KeyEnd          = glfw.KeyEnd
	KeyCapsLock     = glfw.KeyCapsLock
	KeyScrollLock   = glfw.KeyScrollLock
	KeyNumLock      = glfw.KeyNumLock
	KeyPrintScreen  = glfw.KeyPrintScreen
	KeyPause        = glfw.KeyPause
	KeyF1           = glfw.KeyF1
	KeyF2           = glfw.KeyF2
	KeyF3           = glfw.KeyF3
	KeyF4           = glfw.KeyF4
	KeyF5           = glfw.KeyF5
	KeyF6           = glfw.KeyF6
	KeyF7           = glfw.KeyF7
	KeyF8           = glfw.KeyF8
	KeyF9           = glfw.KeyF9
	KeyF10          = glfw.KeyF10
	KeyF11          = glfw.KeyF11
	KeyF12          = glfw.KeyF12
	KeyF13          = glfw.KeyF13
	KeyF14          = glfw.KeyF14
	KeyF15          = glfw.KeyF15
	KeyF16          = glfw.KeyF16
	KeyF17          = glfw.KeyF17
	KeyF18          = glfw.KeyF18
	KeyF19          = glfw.KeyF19
	KeyF20          = glfw.KeyF20
	KeyF21          = glfw.KeyF21
	KeyF22          = glfw.KeyF22
	KeyF23          = glfw.KeyF23
	KeyF24          = glfw.KeyF24
	KeyF25          = glfw.KeyF25
	KeyKP0          = glfw.KeyKP0
	KeyKP1          = glfw.KeyKP1
	KeyKP2          = glfw.KeyKP2
	KeyKP3          = glfw.KeyKP3
	KeyKP4          = glfw.KeyKP4
	KeyKP5          = glfw.KeyKP5
	KeyKP6          = glfw.KeyKP6
	KeyKP7          = glfw.KeyKP7
	KeyKP8          = glfw.KeyKP8
	KeyKP9          = glfw.KeyKP9
	KeyKPDecimal    = glfw.KeyKPDecimal
	KeyKPDivide     = glfw.KeyKPDivide
	KeyKPMultiply   = glfw.KeyKPMultiply
	KeyKPSubtract   = glfw.KeyKPSubtract
	KeyKPAdd        = glfw.KeyKPAdd
	KeyKPEnter      = glfw.KeyKPEnter
	KeyKPEqual      = glfw.KeyKPEqual
	KeyLeftShift    = glfw.KeyLeftShift
	KeyLeftControl  = glfw.KeyLeftControl
	KeyLeftAlt      = glfw.KeyLeftAlt
	KeyLeftSuper    = glfw.KeyLeftSuper
	KeyRightShift   = glfw.KeyRightShift
	KeyRightControl = glfw.KeyRightControl
	KeyRightAlt     = glfw.KeyRightAlt
	KeyRightSuper   = glfw.KeyRightSuper
	KeyMenu         = glfw.KeyMenu
	KeyLast         = glfw.KeyLast
)

func (w *Window) initInput() {
	mainthread.Call(func() {
		w.window.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
			switch action {
			case glfw.Press:
				w.tempInp.buttons[button] = true
			case glfw.Release:
				w.tempInp.buttons[button] = false
			}
		})

		w.window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if key == glfw.KeyUnknown {
				return
			}
			switch action {
			case glfw.Press:
				w.tempInp.buttons[key] = true
			case glfw.Release:
				w.tempInp.buttons[key] = false
			case glfw.Repeat:
				w.tempInp.repeat[key] = true
			}
		})

		w.window.SetCursorEnterCallback(func(_ *glfw.Window, entered bool) {
			w.cursorInsideWindow = entered
		})

		w.window.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
			w.tempInp.mouse = pixel.V(
				x+w.bounds.Min.X,
				(w.bounds.H()-y)+w.bounds.Min.Y,
			)
		})

		w.window.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
			w.tempInp.scroll.X += xoff
			w.tempInp.scroll.Y += yoff
		})

		w.window.SetCharCallback(func(_ *glfw.Window, r rune) {
			w.tempInp.typed += string(r)
		})
	})
}

// UpdateInput polls window events. Call this function to poll window events
// without swapping buffers. Note that the Update method invokes UpdateInput.
func (w *Window) UpdateInput() {
	mainthread.Call(func() {
		glfw.PollEvents()
	})

	w.prevInp = w.currInp
	w.currInp = w.tempInp

	w.tempInp.repeat = [KeyLast + 1]bool{}
	w.tempInp.scroll = pixel.ZV
	w.tempInp.typed = ""

	w.updateJoystickInput()
}

type inputState struct {
	mouse   pixel.Vec
	buttons [KeyLast + 1]bool
	repeat  [KeyLast + 1]bool
	scroll  pixel.Vec
	typed   string
}
