# charm-dialog

Caixas de diálogo modais para aplicações [Bubble Tea](https://github.com/charmbracelet/bubbletea) — confirmação, alerta, entrada de texto, customizado.

> **[Read in English](README.md)**

## Funcionalidades

- Diálogos de **confirmação** com botões Sim/Não
- Diálogos de **alerta** com botão OK
- Diálogos de **entrada** com digitação de verdade (texto, backspace, ctrl+u)
- Diálogos **customizados** com botões arbitrários
- Totalmente estilizado com [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Navegação por teclado entre botões (tab, shift+tab, setas)
- Renderização em overlay para posicionamento modal, com exatamente a altura pedida
- A renderização nunca ultrapassa a largura pedida

## Instalação

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
    // Diálogo de confirmação
    d := dialog.NewConfirm("Deletar Arquivo", "Tem certeza que deseja deletar?")
    fmt.Println(d.Render(80))

    // Diálogo de alerta
    a := dialog.NewAlert("Sucesso", "Operacao concluida com sucesso.")
    fmt.Println(a.Render(80))

    // Diálogo de entrada
    i := dialog.NewInput("Renomear", "Digite o novo nome:")
    fmt.Println(i.Render(80))

    // Diálogo customizado
    c := dialog.NewCustom("Salvar", "Voce tem alteracoes nao salvas.", []string{"Salvar", "Descartar", "Cancelar"})
    fmt.Println(c.Render(80))
}
```

## Tipos de Diálogo

| Tipo | Construtor | Botões |
|------|-----------|--------|
| Confirmação | `NewConfirm(titulo, mensagem)` | Yes, No |
| Alerta | `NewAlert(titulo, mensagem)` | OK |
| Entrada | `NewInput(titulo, prompt)` | Submit, Cancel |
| Customizado | `NewCustom(titulo, mensagem, botoes)` | Definido pelo usuário |

## Tratamento de teclas

`Dialog.SetKeyMsg` processa uma tecla do Bubble Tea v2 (`tea.KeyPressMsg`) e
`Dialog.HandleMsg` despacha qualquer `tea.Msg` para ele. Os dois devolvem a
ação disparada e informam se o diálogo consumiu a tecla — quando a tecla não
é consumida, a aplicação continua o processamento normal dela.

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

| Tecla | Efeito |
|-------|--------|
| `tab` | seleciona o próximo botão |
| `shift+tab` | seleciona o botão anterior |
| `left` / `right` | seleciona o botão anterior / próximo |
| `enter` | confirma o diálogo (`ActionSubmit`) |
| `esc` | cancela o diálogo (`ActionCancel`) |
| caracteres imprimíveis | inseridos no campo de entrada (apenas em `TypeInput`) |
| `backspace` | apaga o último caractere do campo de entrada |
| `ctrl+u` | limpa o campo de entrada |

O cursor do texto não pode ser movido: o texto é sempre digitado no fim do
campo.

## Renderização

`Render(width)` centraliza o diálogo dentro de `width`. Quando o conteúdo não
cabe, a mensagem é quebrada em vez de ser cortada, então a saída nunca
ultrapassa a largura pedida.

`RenderOverlay(width, height)` centraliza o diálogo vertical e horizontalmente
dentro da área informada e devolve **exatamente `height` linhas** (sem `\n` no
final). Quando o diálogo é mais alto que a área, o conteúdo é truncado.

## Estilização

Todos os estilos são customizáveis via struct `Styles`:

```go
d := dialog.NewConfirm("Titulo", "Mensagem")
s := dialog.DefaultStyles()
s.TitleStyle = s.TitleStyle.Background(lipgloss.Color("#FF0000"))
d.WithStyles(s)
```

`WithStyles` tem receptor por ponteiro: os estilos são aplicados no diálogo
original, então não é preciso reatribuir o resultado. Uma cópia do diálogo
também é devolvida, o que permite o uso encadeado (`d = d.WithStyles(s)`).

## Licença

[MIT](LICENSE) - junhinhow
