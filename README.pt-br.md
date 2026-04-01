# charm-dialog

Caixas de dialogo modais para aplicacoes [Bubble Tea](https://github.com/charmbracelet/bubbletea) — confirmacao, alerta, entrada de texto, customizado.

> **[Read in English](README.md)**

## Funcionalidades

- Dialogos de **confirmacao** com botoes Sim/Nao
- Dialogos de **alerta** com botao OK
- Dialogos de **entrada** com campo de texto
- Dialogos **customizados** com botoes arbitrarios
- Totalmente estilizado com [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Navegacao por teclado entre botoes
- Renderizacao em overlay para posicionamento modal

## Instalacao

```bash
go get github.com/junhinhow/charm-dialog@latest
```

## Uso

```go
package main

import (
    "fmt"
    "github.com/junhinhow/charm-dialog"
)

func main() {
    // Dialogo de confirmacao
    d := dialog.NewConfirm("Deletar Arquivo", "Tem certeza que deseja deletar?")
    fmt.Println(d.Render(80))

    // Dialogo de alerta
    a := dialog.NewAlert("Sucesso", "Operacao concluida com sucesso.")
    fmt.Println(a.Render(80))

    // Dialogo de entrada
    i := dialog.NewInput("Renomear", "Digite o novo nome:")
    fmt.Println(i.Render(80))

    // Dialogo customizado
    c := dialog.NewCustom("Salvar", "Voce tem alteracoes nao salvas.", []string{"Salvar", "Descartar", "Cancelar"})
    fmt.Println(c.Render(80))
}
```

## Tipos de Dialogo

| Tipo | Construtor | Botoes |
|------|-----------|--------|
| Confirmacao | `NewConfirm(titulo, mensagem)` | Yes, No |
| Alerta | `NewAlert(titulo, mensagem)` | OK |
| Entrada | `NewInput(titulo, prompt)` | Submit, Cancel |
| Customizado | `NewCustom(titulo, mensagem, botoes)` | Definido pelo usuario |

## Estilizacao

Todos os estilos sao customizaveis via struct `Styles`:

```go
d := dialog.NewConfirm("Titulo", "Mensagem")
s := dialog.DefaultStyles()
s.TitleStyle = s.TitleStyle.Background(lipgloss.Color("#FF0000"))
d = d.WithStyles(s)
```

## Licenca

[MIT](LICENSE) - junhinhow
