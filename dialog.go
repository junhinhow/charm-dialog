// Package dialog fornece caixas de dialogo modais para aplicacoes Bubble Tea.
//
// Suporta dialogos de confirmacao, alerta, entrada de texto e customizados,
// com estilizacao via lipgloss.
package dialog

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// DialogType representa o tipo do dialogo.
type DialogType int

const (
	// TypeConfirm dialogo de confirmacao com botoes Sim/Nao.
	TypeConfirm DialogType = iota
	// TypeAlert dialogo de alerta com botao OK.
	TypeAlert
	// TypeInput dialogo com campo de entrada de texto.
	TypeInput
	// TypeCustom dialogo com botoes personalizados.
	TypeCustom
)

// Styles agrupa os estilos usados para renderizar o dialogo.
type Styles struct {
	// TitleStyle estilo do titulo do dialogo.
	TitleStyle lipgloss.Style
	// MessageStyle estilo da mensagem do dialogo.
	MessageStyle lipgloss.Style
	// ButtonStyle estilo dos botoes inativos.
	ButtonStyle lipgloss.Style
	// ActiveButtonStyle estilo do botao selecionado.
	ActiveButtonStyle lipgloss.Style
	// BorderStyle estilo da borda do dialogo.
	BorderStyle lipgloss.Style
}

// DefaultStyles retorna os estilos padrao para o dialogo.
func DefaultStyles() Styles {
	return Styles{
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1),
		MessageStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD")).
			Padding(1, 0),
		ButtonStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 2).
			Margin(0, 1),
		ActiveButtonStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2).
			Margin(0, 1),
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2),
	}
}

// Dialog representa uma caixa de dialogo modal.
type Dialog struct {
	// Title titulo exibido no topo do dialogo.
	Title string
	// Message mensagem principal do dialogo.
	Message string
	// Buttons lista de rotulos dos botoes.
	Buttons []string
	// Selected indice do botao atualmente selecionado.
	Selected int
	// Type tipo do dialogo.
	Type DialogType
	// InputValue valor do campo de entrada (apenas para TypeInput).
	InputValue string
	// InputPrompt prompt exibido antes do campo de entrada.
	InputPrompt string
	// Styles estilos de renderizacao.
	Styles Styles
}

// NewConfirm cria um dialogo de confirmacao com botoes "Yes" e "No".
func NewConfirm(title, message string) Dialog {
	return Dialog{
		Title:    title,
		Message:  message,
		Buttons:  []string{"Yes", "No"},
		Selected: 0,
		Type:     TypeConfirm,
		Styles:   DefaultStyles(),
	}
}

// NewAlert cria um dialogo de alerta com botao "OK".
func NewAlert(title, message string) Dialog {
	return Dialog{
		Title:    title,
		Message:  message,
		Buttons:  []string{"OK"},
		Selected: 0,
		Type:     TypeAlert,
		Styles:   DefaultStyles(),
	}
}

// NewInput cria um dialogo com campo de entrada de texto.
func NewInput(title, prompt string) Dialog {
	return Dialog{
		Title:       title,
		Message:     "",
		Buttons:     []string{"Submit", "Cancel"},
		Selected:    0,
		Type:        TypeInput,
		InputPrompt: prompt,
		Styles:      DefaultStyles(),
	}
}

// NewCustom cria um dialogo com botoes personalizados.
func NewCustom(title, message string, buttons []string) Dialog {
	return Dialog{
		Title:    title,
		Message:  message,
		Buttons:  buttons,
		Selected: 0,
		Type:     TypeCustom,
		Styles:   DefaultStyles(),
	}
}

// WithStyles substitui os estilos do dialogo.
func (d Dialog) WithStyles(s Styles) Dialog {
	d.Styles = s
	return d
}

// NextButton avanca a selecao para o proximo botao.
func (d *Dialog) NextButton() {
	if d.Selected < len(d.Buttons)-1 {
		d.Selected++
	}
}

// PrevButton retrocede a selecao para o botao anterior.
func (d *Dialog) PrevButton() {
	if d.Selected > 0 {
		d.Selected--
	}
}

// SelectedButton retorna o rotulo do botao selecionado.
func (d Dialog) SelectedButton() string {
	if d.Selected >= 0 && d.Selected < len(d.Buttons) {
		return d.Buttons[d.Selected]
	}
	return ""
}

// Render renderiza o dialogo como string centralizada na largura fornecida.
func (d Dialog) Render(width int) string {
	s := d.Styles

	// Renderiza titulo
	titleStr := s.TitleStyle.Render(d.Title)

	// Renderiza mensagem ou prompt de entrada
	var bodyStr string
	if d.Type == TypeInput {
		promptLine := s.MessageStyle.Render(d.InputPrompt)
		inputLine := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Underline(true).
			Render(d.InputValue + "_")
		bodyStr = promptLine + "\n" + inputLine
	} else {
		bodyStr = s.MessageStyle.Render(d.Message)
	}

	// Renderiza botoes
	var buttons []string
	for i, btn := range d.Buttons {
		if i == d.Selected {
			buttons = append(buttons, s.ActiveButtonStyle.Render(btn))
		} else {
			buttons = append(buttons, s.ButtonStyle.Render(btn))
		}
	}
	buttonsStr := lipgloss.JoinHorizontal(lipgloss.Center, buttons...)

	// Monta conteudo interno
	content := lipgloss.JoinVertical(lipgloss.Center,
		titleStr,
		bodyStr,
		"",
		buttonsStr,
	)

	// Aplica borda
	boxed := s.BorderStyle.Render(content)

	// Centraliza horizontalmente
	if width > 0 {
		boxed = lipgloss.PlaceHorizontal(width, lipgloss.Center, boxed)
	}

	return boxed
}

// RenderOverlay renderiza o dialogo centralizado vertical e horizontalmente
// dentro das dimensoes fornecidas, util para sobrepor em telas existentes.
func (d Dialog) RenderOverlay(width, height int) string {
	rendered := d.Render(0)

	// Calcula linhas para centralizar verticalmente
	lines := strings.Count(rendered, "\n") + 1
	topPad := 0
	if height > lines {
		topPad = (height - lines) / 2
	}

	var b strings.Builder
	for range topPad {
		b.WriteString(strings.Repeat(" ", width) + "\n")
	}

	// Centraliza cada linha horizontalmente
	for _, line := range strings.Split(rendered, "\n") {
		lineWidth := lipgloss.Width(line)
		pad := 0
		if width > lineWidth {
			pad = (width - lineWidth) / 2
		}
		b.WriteString(strings.Repeat(" ", pad) + line + "\n")
	}

	return b.String()
}
