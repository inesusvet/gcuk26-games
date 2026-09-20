package input

import (
	"testing"
)

func TestInputManagerInitialization(t *testing.T) {
	mgr := NewManager()
	if mgr == nil {
		t.Fatal("Expected NewManager() to return a non-nil Manager")
	}

	// Without active game window/gamepad attached in test env, these should return default false values gracefully without panicking.
	_ = mgr.IsGamepadConnected()
	_ = mgr.GetGamepadName()
	_ = mgr.Is8BitDoM30()
}

func TestInputManagerActionMethods(t *testing.T) {
	mgr := NewManager()

	// Verify all action queries return non-panicking boolean values in non-headless test environment
	if mgr.IsMoveRight() {
		t.Log("IsMoveRight returned true")
	}
	if mgr.IsMoveLeft() {
		t.Log("IsMoveLeft returned true")
	}
	if mgr.IsJumpJustPressed() {
		t.Log("IsJumpJustPressed returned true")
	}
	if mgr.IsInteractJustPressed() {
		t.Log("IsInteractJustPressed returned true")
	}
	if mgr.IsTurboJustPressed() {
		t.Log("IsTurboJustPressed returned true")
	}
	if mgr.IsConfirmJustPressed() {
		t.Log("IsConfirmJustPressed returned true")
	}
	if mgr.IsColor1JustPressed() {
		t.Log("IsColor1JustPressed returned true")
	}
	if mgr.IsColor2JustPressed() {
		t.Log("IsColor2JustPressed returned true")
	}
	if mgr.IsColor3JustPressed() {
		t.Log("IsColor3JustPressed returned true")
	}
}
