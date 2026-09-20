# Unified Input Action Mapping & Rebinding Reference

This module covers input device abstraction, action rebind mapping, analog stick deadzone filtering, 8bitdo M30 / 6-button gamepad integration, raw layout fallbacks, and multi-input state handling for Ebitengine (v2).

---

## 1. Action Mapping Abstraction Layer

Decouple raw physical keys/gamepad buttons from logical game actions (`ActionJump`, `ActionAttack`, `ActionPause`):

```go
type Action int

const (
	ActionMoveLeft Action = iota
	ActionMoveRight
	ActionJump
	ActionAttack
	ActionPause
)

type Binding struct {
	Keys     []ebiten.Key
	Buttons  []ebiten.StandardGamepadButton
	Mouses   []ebiten.MouseButton
}

type InputManager struct {
	Bindings map[Action]Binding
}

func NewInputManager() *InputManager {
	im := &InputManager{
		Bindings: make(map[Action]Binding),
	}
	// Default Bindings
	im.Bindings[ActionMoveLeft] = Binding{Keys: []ebiten.Key{ebiten.KeyA, ebiten.KeyLeft}}
	im.Bindings[ActionMoveRight] = Binding{Keys: []ebiten.Key{ebiten.KeyD, ebiten.KeyRight}}
	im.Bindings[ActionJump] = Binding{
		Keys:    []ebiten.Key{ebiten.KeySpace, ebiten.KeyW, ebiten.KeyUp},
		Buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonRightBottom}, // South (A/B)
	}
	im.Bindings[ActionAttack] = Binding{
		Keys:    []ebiten.Key{ebiten.KeyJ, ebiten.KeyZ},
		Mouses:  []ebiten.MouseButton{ebiten.MouseButtonLeft},
		Buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonRightLeft}, // West (X/Y)
	}
	return im
}
```

---

## 2. Ebitengine v2 Gamepad API & 8bitdo M30 Mapping Reference

Ebitengine v2 uses Web Standard Gamepad Button & Axis layouts.

### Standard Gamepad Button Constants
- **Face Buttons**:
  - `ebiten.StandardGamepadButtonRightBottom` (South: A on Xbox / B on Switch)
  - `ebiten.StandardGamepadButtonRightRight` (East: B on Xbox / A on Switch)
  - `ebiten.StandardGamepadButtonRightLeft` (West: X on Xbox / Y on Switch)
  - `ebiten.StandardGamepadButtonRightTop` (North: Y on Xbox / X on Switch)
- **D-Pad**:
  - `ebiten.StandardGamepadButtonLeftLeft` (D-Pad Left)
  - `ebiten.StandardGamepadButtonLeftRight` (D-Pad Right)
  - `ebiten.StandardGamepadButtonLeftTop` (D-Pad Up)
  - `ebiten.StandardGamepadButtonLeftBottom` (D-Pad Down)
- **Shoulders / Triggers (M30 C & Z Extra Face Buttons)**:
  - `ebiten.StandardGamepadButtonFrontTopRight` (R / C button)
  - `ebiten.StandardGamepadButtonFrontBottomRight` (R2 / Z button)
  - `ebiten.StandardGamepadButtonFrontTopLeft` (L button)
  - `ebiten.StandardGamepadButtonFrontBottomLeft` (L2 button)
- **Menu / Center**:
  - `ebiten.StandardGamepadButtonCenterRight` (Start / Plus)
  - `ebiten.StandardGamepadButtonCenterLeft` (Select / Minus)

### Standard vs Raw Gamepad Inspection
When standard SDL layout mapping is missing (e.g. DInput or unmapped Bluetooth mode), fallback to raw button and axis indices:

```go
gamepads := ebiten.AppendGamepadIDs(nil)
for _, id := range gamepads {
	if ebiten.IsStandardGamepadLayoutAvailable(id) {
		if ebiten.IsStandardGamepadButtonPressed(id, ebiten.StandardGamepadButtonLeftRight) {
			return true
		}
		if ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) > 0.3 {
			return true
		}
	} else {
		// Raw Fallback for Unmapped Gamepads
		if ebiten.GamepadAxisCount(id) > 0 && ebiten.GamepadAxisValue(id, 0) > 0.3 {
			return true
		}
		for b := 0; b < ebiten.GamepadButtonCount(id) && b < 4; b++ {
			if ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(b)) {
				return true
			}
		}
	}
}
```

---

## 3. Action State Evaluation

Query whether an action is currently held, just pressed, or just released across all assigned physical inputs:

```go
func (im *InputManager) IsActionPressed(action Action) bool {
	binding, ok := im.Bindings[action]
	if !ok { return false }

	for _, k := range binding.Keys {
		if ebiten.IsKeyPressed(k) { return true }
	}
	for _, m := range binding.Mouses {
		if ebiten.IsMouseButtonPressed(m) { return true }
	}
	
	// Gamepad support
	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		for _, b := range binding.Buttons {
			if ebiten.IsStandardGamepadButtonPressed(id, b) { return true }
		}
	}
	return false
}

func (im *InputManager) IsActionJustPressed(action Action) bool {
	binding, ok := im.Bindings[action]
	if !ok { return false }

	for _, k := range binding.Keys {
		if inpututil.IsKeyJustPressed(k) { return true }
	}
	for _, m := range binding.Mouses {
		if inpututil.IsMouseButtonJustPressed(m) { return true }
	}

	gamepads := ebiten.AppendGamepadIDs(nil)
	for _, id := range gamepads {
		for _, b := range binding.Buttons {
			if inpututil.IsStandardGamepadButtonJustPressed(id, b) { return true }
		}
	}
	return false
}
```

---

## 4. Gamepad Analog Stick Radial Deadzone Filtering

Analog thumbsticks on gamepads exhibit hardware drift near center. Apply radial deadzone filtering:

```go
func ApplyRadialDeadzone(axisX, axisY, deadzone float64) (float64, float64) {
	magnitude := math.Sqrt(axisX*axisX + axisY*axisY)
	if magnitude < deadzone {
		return 0, 0
	}
	// Rescale remaining range [deadzone..1.0] to [0.0..1.0]
	normalizedMag := (magnitude - deadzone) / (1.0 - deadzone)
	if normalizedMag > 1.0 { normalizedMag = 1.0 }

	nx := (axisX / magnitude) * normalizedMag
	ny := (axisY / magnitude) * normalizedMag
	return nx, ny
}
```

---

## 5. Dual Input UI / HUD Prompt Hints

Always format prompt text to display both Gamepad and Keyboard input choices:
- `PRESS A / SPACE / ENTER TO START`
- `Controls: Move (D-Pad / Arrows / A-D) | Jump (A / Space) | Turbo (X / Shift)`
- `SELECT COLOR: Press [A/1] RED | [B/2] BLUE | [X/3] GOLD`

---

## 6. Input Rebinding & Persistence

Serialize custom keybindings to disk as JSON:

```go
type SerializableBindings map[string][]string

func (im *InputManager) SaveBindings(filepath string) error {
	data, err := json.MarshalIndent(im.Bindings, "", "  ")
	if err != nil { return err }
	return os.WriteFile(filepath, data, 0644)
}
```
