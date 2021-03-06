package pixelgl

import (
	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// // Joystick is a joystick or controller.
// type Joystick int

// List all of the joysticks.
const (
	Joystick1  = glfw.Joystick1
	Joystick2  = glfw.Joystick2
	Joystick3  = glfw.Joystick3
	Joystick4  = glfw.Joystick4
	Joystick5  = glfw.Joystick5
	Joystick6  = glfw.Joystick6
	Joystick7  = glfw.Joystick7
	Joystick8  = glfw.Joystick8
	Joystick9  = glfw.Joystick9
	Joystick10 = glfw.Joystick10
	Joystick11 = glfw.Joystick11
	Joystick12 = glfw.Joystick12
	Joystick13 = glfw.Joystick13
	Joystick14 = glfw.Joystick14
	Joystick15 = glfw.Joystick15
	Joystick16 = glfw.Joystick16

	JoystickLast = glfw.JoystickLast
)

// JoystickPresent returns if the joystick is currently connected.
//
// This API is experimental.
func (w *Window) JoystickPresent(js glfw.Joystick) bool {
	return w.currJoy[js].connected
}

// JoystickName returns the name of the joystick. A disconnected joystick will return an
// empty string.
//
// This API is experimental.
func (w *Window) JoystickName(js glfw.Joystick) string {
	return w.currJoy[js].name
}

// JoystickButtonCount returns the number of buttons a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickButtonCount(js glfw.Joystick) int {
	return len(w.currJoy[js].buttons)
}

// JoystickAxisCount returns the number of axes a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickAxisCount(js glfw.Joystick) int {
	return len(w.currJoy[js].axes)
}

// JoystickPressed returns whether the joystick Button is currently pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickPressed(js glfw.Joystick, button int) bool {
	return w.currJoy[js].getButton(button)
}

// JoystickJustPressed returns whether the joystick Button has just been pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustPressed(js glfw.Joystick, button int) bool {
	return w.currJoy[js].getButton(button) && !w.prevJoy[js].getButton(button)
}

// JoystickJustReleased returns whether the joystick Button has just been released up.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustReleased(js glfw.Joystick, button int) bool {
	return !w.currJoy[js].getButton(button) && w.prevJoy[js].getButton(button)
}

// JoystickAxis returns the value of a joystick axis at the last call to Window.Update.
// If the axis index is out of range, this will return 0.
//
// This API is experimental.
func (w *Window) JoystickAxis(js glfw.Joystick, axis int) float64 {
	return w.currJoy[js].getAxis(axis)
}

// Used internally during Window.UpdateInput to update the state of the joysticks.
func (w *Window) updateJoystickInput() {
	for js := Joystick1; js <= JoystickLast; js++ {
		w.prevJoy = w.currJoy
		mainthread.Call(func() {
			if glfw.Joystick(js).IsGamepad() {
				if !w.currJoy[js].connected {
					w.currJoy[js].connected = true
					w.currJoy[js].name = js.GetName()
				}

				w.currJoy[js].buttons = js.GetButtons()
				w.currJoy[js].axes = js.GetAxes()
			} else {
				if w.currJoy[js].connected {
					w.currJoy[js] = joystickState{}
				}
			}
		})
	}
}

type joystickState struct {
	connected bool
	name      string
	buttons   []glfw.Action
	axes      []float32
}

// Returns if a button on a joystick is down, returning false if the button or joystick is invalid.
func (js *joystickState) getButton(button int) bool {
	// Check that the joystick and button is valid, return false by default
	if js.buttons == nil || button >= len(js.buttons) || button < 0 {
		return false
	}
	return js.buttons[button] == glfw.Press
}

// Returns the value of a joystick axis, returning 0 if the button or joystick is invalid.
func (js *joystickState) getAxis(axis int) float64 {
	// Check that the joystick and axis is valid, return 0 by default.
	if js.axes == nil || axis >= len(js.axes) || axis < 0 {
		return 0
	}
	return float64(js.axes[axis])
}
