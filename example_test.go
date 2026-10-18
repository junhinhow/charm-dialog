package dialog

import (
	"fmt"

	"charm.land/bubbletea/v2"
)

func ExampleNewConfirm() {
	d := NewConfirm("Delete File", "Are you sure you want to delete this file?")
	fmt.Println(d.Render(80))
}

func ExampleNewInput() {
	d := NewInput("Rename", "Enter new name:")
	fmt.Println(d.Render(80))
}

func ExampleDialog_RenderOverlay() {
	d := NewAlert("Success", "Operation completed successfully.")

	// The overlay has exactly 24 lines and is centered in the area.
	fmt.Println(d.RenderOverlay(80, 24))
}

func ExampleDialog_SetKeyMsg() {
	d := NewInput("Rename", "Enter new name:")

	// Simulates typing "docs".
	for _, r := range "docs" {
		d.SetKeyMsg(tea.KeyPressMsg{Text: string(r), Code: r})
	}

	fmt.Println(d.InputValue)
	// Output: docs
}

func ExampleDialog_HandleMsg() {
	d := NewConfirm("Delete File", "Are you sure?")

	// Tab moves to the next button, enter confirms the dialog.
	d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyTab})
	action, _ := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyEnter})

	fmt.Println(d.SelectedButton())
	fmt.Println(action)
	// Output:
	// No
	// submit
}
