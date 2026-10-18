# charm-dialog

Modal dialog boxes for [Bubble Tea](https://github.com/charmbracelet/bubbletea) applications — confirm, alert, input, custom.

> **[Leia em Português (PT-BR)](README.pt-br.md)**

## Features

- **Confirm** dialogs with Yes/No buttons
- **Alert** dialogs with OK button
- **Input** dialogs with real text entry (typing, backspace, ctrl+u)
- **Custom** dialogs with arbitrary buttons
- Fully styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Keyboard navigation between buttons (tab, shift+tab, arrows)
- Overlay rendering for modal placement, with exactly the requested height
- Rendering never exceeds the requested width

## Install

```bash
go get github.com/junhinhow/charm-dialog@latest
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/junhinhow/charm-dialog"
)

func main() {
    // Confirm dialog
    d := dialog.NewConfirm("Delete File", "Are you sure you want to delete this file?")
    fmt.Println(d.Render(80))

    // Alert dialog
    a := dialog.NewAlert("Success", "Operation completed successfully.")
    fmt.Println(a.Render(80))

    // Input dialog
    i := dialog.NewInput("Rename", "Enter new name:")
    fmt.Println(i.Render(80))

    // Custom dialog
    c := dialog.NewCustom("Save Changes", "You have unsaved changes.", []string{"Save", "Discard", "Cancel"})
    fmt.Println(c.Render(80))
}
```

## Dialog Types

| Type | Constructor | Buttons |
|------|-------------|---------|
| Confirm | `NewConfirm(title, message)` | Yes, No |
| Alert | `NewAlert(title, message)` | OK |
| Input | `NewInput(title, prompt)` | Submit, Cancel |
| Custom | `NewCustom(title, message, buttons)` | User-defined |

## Keyboard handling

`Dialog.SetKeyMsg` processes a Bubble Tea v2 key press (`tea.KeyPressMsg`) and
`Dialog.HandleMsg` dispatches any `tea.Msg` to it. Both return the action that
was triggered and whether the dialog consumed the key — when the key is not
consumed the application keeps processing it normally.

```go
case tea.KeyPressMsg:
    if action, ok := d.HandleMsg(msg); ok {
        switch action {
        case dialog.ActionSubmit:
            m.confirmed = true
        case dialog.ActionCancel:
            m.cancelled = true
        }
    }
```

| Key | Effect |
|-----|--------|
| `tab` | selects the next button |
| `shift+tab` | selects the previous button |
| `left` / `right` | selects the previous / next button |
| `enter` | confirms the dialog (`ActionSubmit`) |
| `esc` | cancels the dialog (`ActionCancel`) |
| printable characters | inserted into the input field (`TypeInput` only) |
| `backspace` | deletes the last character of the input field |
| `ctrl+u` | clears the input field |

The text cursor cannot be moved: the text is always typed at the end of the
field.

## Rendering

`Render(width)` centers the dialog within `width`. When the content does not
fit, the message is wrapped instead of being cut, so the output never exceeds
the requested width.

`RenderOverlay(width, height)` centers the dialog both vertically and
horizontally inside the given area and returns **exactly `height` lines**
(without a trailing newline). When the dialog is taller than the area, the
content is truncated.

## Styling

All styles are customizable via the `Styles` struct:

```go
d := dialog.NewConfirm("Title", "Message")
s := dialog.DefaultStyles()
s.TitleStyle = s.TitleStyle.Background(lipgloss.Color("#FF0000"))
d.WithStyles(s)
```

`WithStyles` has a pointer receiver: the styles are applied to the original
dialog, so the result does not need to be reassigned. A copy of the dialog is
also returned, which allows chained use (`d = d.WithStyles(s)`).

## License

[MIT](LICENSE) - junhinhow
