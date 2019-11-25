package pixelgl

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Joystick is a joystick or controller.
type Joystick int

// List all of the joysticks.
const (
	Joystick1  = Joystick(glfw.Joystick1)
	Joystick2  = Joystick(glfw.Joystick2)
	Joystick3  = Joystick(glfw.Joystick3)
	Joystick4  = Joystick(glfw.Joystick4)
	Joystick5  = Joystick(glfw.Joystick5)
	Joystick6  = Joystick(glfw.Joystick6)
	Joystick7  = Joystick(glfw.Joystick7)
	Joystick8  = Joystick(glfw.Joystick8)
	Joystick9  = Joystick(glfw.Joystick9)
	Joystick10 = Joystick(glfw.Joystick10)
	Joystick11 = Joystick(glfw.Joystick11)
	Joystick12 = Joystick(glfw.Joystick12)
	Joystick13 = Joystick(glfw.Joystick13)
	Joystick14 = Joystick(glfw.Joystick14)
	Joystick15 = Joystick(glfw.Joystick15)
	Joystick16 = Joystick(glfw.Joystick16)

	JoystickLast = Joystick(glfw.JoystickLast)
)

// JoystickPresent returns if the joystick is currently connected.
//
// This API is experimental.
func (w *Window) JoystickPresent(js Joystick) bool {
	return w.currJoy[js].connected
}

// JoystickName returns the name of the joystick. A disconnected joystick will return an
// empty string.
//
// This API is experimental.
func (w *Window) JoystickName(js Joystick) string {
	return w.currJoy[js].name
}

// JoystickButtonCount returns the number of buttons a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickButtonCount(js Joystick) int {
	return len(w.currJoy[js].buttons)
}

// JoystickAxisCount returns the number of axes a connected joystick has.
//
// This API is experimental.
func (w *Window) JoystickAxisCount(js Joystick) int {
	return len(w.currJoy[js].axes)
}

// JoystickPressed returns whether the joystick Button is currently pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickPressed(js Joystick, button int) bool {
	// Check that the joystick and button is valid, return false by default
	if w.currJoy[js].buttons == nil || button >= len(w.currJoy[js].buttons) || button < 0 {
		return false
	}

	return w.currJoy[js].buttons[button] == glfw.Repeat
}

// JoystickJustPressed returns whether the joystick Button has just been pressed down.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustPressed(js Joystick, button int) bool {
	// Check that the joystick and button is valid, return false by default
	if w.currJoy[js].buttons == nil || button >= len(w.currJoy[js].buttons) || button < 0 {
		return false
	}
	return w.currJoy[js].buttons[button] == glfw.Press
}

// JoystickJustReleased returns whether the joystick Button has just been released up.
// If the button index is out of range, this will return false.
//
// This API is experimental.
func (w *Window) JoystickJustReleased(js Joystick, button int) bool {
	return w.currJoy[js].buttons[button] == glfw.Release
}

// JoystickAxis returns the value of a joystick axis at the last call to Window.Update.
// If the axis index is out of range, this will return 0.
//
// This API is experimental.
func (w *Window) JoystickAxes(js Joystick, axes int) float64 {
	// Check that the joystick and axis is valid, return 0 by default.
	if w.currJoy[js].axes == nil || axes >= len(w.currJoy[js].axes) || axes < 0 {
		return 0
	}
	return float64(w.currJoy[js].axes[axes])
}

// Used internally during Window.UpdateInput to update the state of the joysticks.
func (w *Window) updateJoystickInput() {
	for js := Joystick1; js <= JoystickLast; js++ {
		if glfw.Joystick(js).IsGamepad() {
			if !w.currJoy[js].connected {
				w.currJoy[js].connected = true
				w.currJoy[js].name = glfw.Joystick(js).GetName()
			}

			w.currJoy[js].buttons = glfw.Joystick(js).GetButtons()
			w.currJoy[js].axes = glfw.Joystick(js).GetAxes()
		} else {
			if w.currJoy[js].connected {
				w.currJoy[js].buttons = nil
				w.currJoy[js].axes = nil

				w.currJoy[js].name = ""

				w.currJoy[js].connected = false
			}
		}
	}
}

// var joycon[JoystickLast + 1]joystickState

type joystickState struct {
	connected bool
	name      string
	buttons   []glfw.Action
	axes      []float32
}
