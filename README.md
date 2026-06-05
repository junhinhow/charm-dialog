# charm-dialog

Modal dialog boxes for [Bubble Tea](https://github.com/charmbracelet/bubbletea) applications — confirm, alert, input, custom.

> **[Leia em Portugues (PT-BR)](README.pt-br.md)**

## Features

- **Confirm** dialogs with Yes/No buttons
- **Alert** dialogs with OK button
- **Input** dialogs with text entry
- **Custom** dialogs with arbitrary buttons
- Fully styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Keyboard navigation between buttons
- Overlay rendering for modal placement

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

## Styling

All styles are customizable via the `Styles` struct:

```go
d := dialog.NewConfirm("Title", "Message")
s := dialog.DefaultStyles()
s.TitleStyle = s.TitleStyle.Background(lipgloss.Color("#FF0000"))
d = d.WithStyles(s)
```

## License

[MIT](LICENSE) - junhinhow
