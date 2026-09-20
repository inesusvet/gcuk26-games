package input

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const Deadzone = 0.3

// Manager provides a unified interface for Keyboard and Gamepad (8bitdo M30) controls.
type Manager struct{}

// NewManager creates a new Input Manager.
func NewManager() *Manager {
	return &Manager{}
}

// IsGamepadConnected returns true if any gamepad (such as 8bitdo M30) is currently connected.
func (m *Manager) IsGamepadConnected() bool {
	return len(ebiten.AppendGamepadIDs(nil)) > 0
}

// GetGamepadName returns the name of the first connected gamepad, or empty string.
func (m *Manager) GetGamepadName() string {
	gamepads := ebiten.AppendGamepadIDs(nil)
	if len(gamepads) == 0 {
		return ""
	}
	return ebiten.GamepadName(gamepads[0])
}

// Is8BitDoM30 returns true if a connected controller is recognized as an 8bitdo M30.
func (m *Manager) Is8BitDoM30() bool {
	name := strings.ToLower(m.GetGamepadName())
	return strings.Contains(name, "8bitdo") || strings.Contains(name, "m30")
}

// IsMoveRight checks if pedal right / move right action is active (Keyboard or Gamepad).
func (m *Manager) IsMoveRight() bool {
	// Keyboard
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		return true
	}

	// Gamepad
	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftRight) {
				return true
			}
			if ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) > Deadzone {
				return true
			}
		} else {
			// Raw fallback
			axisCount := ebiten.GamepadAxisCount(id)
			if axisCount > 0 && ebiten.GamepadAxisValue(id, 0) > Deadzone {
				return true
			}
		}
	}

	return false
}

// IsMoveLeft checks if pedal left / move left action is active (Keyboard or Gamepad).
func (m *Manager) IsMoveLeft() bool {
	// Keyboard
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		return true
	}

	// Gamepad
	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftLeft) {
				return true
			}
			if ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) < -Deadzone {
				return true
			}
		} else {
			// Raw fallback
			axisCount := ebiten.GamepadAxisCount(id)
			if axisCount > 0 && ebiten.GamepadAxisValue(id, 0) < -Deadzone {
				return true
			}
		}
	}

	return false
}

// IsJumpJustPressed checks if jump was just pressed on Keyboard or Gamepad.
func (m *Manager) IsJumpJustPressed() bool {
	// Keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return true
	}

	// Gamepad (A, B, D-Pad Up, Stick Up)
	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) ||
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightRight) ||
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftTop) {
				return true
			}
		} else {
			// Raw buttons fallback
			btnCount := ebiten.GamepadButtonCount(id)
			for b := 0; b < btnCount && b < 4; b++ {
				if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton(b)) {
					return true
				}
			}
		}
	}

	return false
}

// IsInteractJustPressed checks if interact action was triggered (Jump / E / Enter / Start).
func (m *Manager) IsInteractJustPressed() bool {
	if m.IsJumpJustPressed() {
		return true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return true
	}

	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonCenterRight) || // Start
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonCenterLeft) { // Select
				return true
			}
		}
	}

	return false
}

// IsTurboJustPressed checks if turbo boost action was triggered (Shift, J, X, Y, Shoulders).
func (m *Manager) IsTurboJustPressed() bool {
	// Keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) || inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		return true
	}

	// Gamepad
	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightLeft) || // X button
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightTop) || // Y button
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonFrontTopRight) || // R / C button
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonFrontBottomRight) { // R2 / Z button
				return true
			}
		} else {
			btnCount := ebiten.GamepadButtonCount(id)
			if btnCount >= 6 {
				if inpututil.IsGamepadButtonJustPressed(id, 4) || inpututil.IsGamepadButtonJustPressed(id, 5) {
					return true
				}
			}
		}
	}

	return false
}

// IsConfirmJustPressed checks if confirm / start action was pressed (Space, Enter, A, B, Start).
func (m *Manager) IsConfirmJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return true
	}

	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) ||
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightRight) ||
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonCenterRight) {
				return true
			}
		} else {
			btnCount := ebiten.GamepadButtonCount(id)
			for b := 0; b < btnCount && b < 4; b++ {
				if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton(b)) {
					return true
				}
			}
		}
	}

	return false
}

// IsColor1JustPressed checks if Bike Color 1 (Red) was selected.
func (m *Manager) IsColor1JustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
		return true
	}

	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) ||
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftLeft) {
				return true
			}
		}
	}

	return false
}

// IsColor2JustPressed checks if Bike Color 2 (Blue) was selected.
func (m *Manager) IsColor2JustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
		return true
	}

	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightRight) ||
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftTop) {
				return true
			}
		}
	}

	return false
}

// IsColor3JustPressed checks if Bike Color 3 (Gold) was selected.
func (m *Manager) IsColor3JustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
		return true
	}

	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightLeft) ||
				inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonLeftRight) {
				return true
			}
		}
	}

	return false
}
