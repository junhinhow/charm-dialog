// Package dialog fornece caixas de diálogo modais para aplicações Bubble Tea.
//
// Suporta diálogos de confirmação, alerta, entrada de texto e customizados,
// com estilização via lipgloss.
package dialog

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// DialogType representa o tipo do diálogo.
type DialogType int

const (
	// TypeConfirm diálogo de confirmação com botões Sim/Não.
	TypeConfirm DialogType = iota
	// TypeAlert diálogo de alerta com botão OK.
	TypeAlert
	// TypeInput diálogo com campo de entrada de texto.
	TypeInput
	// TypeCustom diálogo com botões personalizados.
	TypeCustom
)

// Styles agrupa os estilos usados para renderizar o diálogo.
type Styles struct {
	// TitleStyle estilo do título do diálogo.
	TitleStyle lipgloss.Style
	// MessageStyle estilo da mensagem do diálogo.
	MessageStyle lipgloss.Style
	// ButtonStyle estilo dos botões inativos.
	ButtonStyle lipgloss.Style
	// ActiveButtonStyle estilo do botão selecionado.
	ActiveButtonStyle lipgloss.Style
	// BorderStyle estilo da borda do diálogo.
	BorderStyle lipgloss.Style
}

// DefaultStyles retorna os estilos padrão para o diálogo.
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

// Dialog representa uma caixa de diálogo modal.
type Dialog struct {
	// Title título exibido no topo do diálogo.
	Title string
	// Message mensagem principal do diálogo.
	Message string
	// Buttons lista de rótulos dos botões.
	Buttons []string
	// Selected índice do botão atualmente selecionado.
	Selected int
	// Type tipo do diálogo.
	Type DialogType
	// InputValue valor do campo de entrada (apenas para TypeInput).
	InputValue string
	// InputPrompt prompt exibido antes do campo de entrada.
	InputPrompt string
	// Styles estilos de renderização.
	Styles Styles
}

// NewConfirm cria um diálogo de confirmação com botões "Yes" e "No".
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

// NewAlert cria um diálogo de alerta com botão "OK".
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

// NewInput cria um diálogo com campo de entrada de texto.
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

// NewCustom cria um diálogo com botões personalizados.
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

// WithStyles substitui os estilos do diálogo.
//
// O método tem receptor por ponteiro: os estilos são aplicados ao diálogo
// original, portanto não é preciso reatribuir o resultado. Uma cópia do
// diálogo também é devolvida, o que permite o uso encadeado:
//
//	d := dialog.NewConfirm("Título", "Mensagem")
//	d = d.WithStyles(estilos)
func (d *Dialog) WithStyles(s Styles) Dialog {
	d.Styles = s
	return *d
}

// NextButton avança a seleção para o próximo botão.
func (d *Dialog) NextButton() {
	if d.Selected < len(d.Buttons)-1 {
		d.Selected++
	}
}

// PrevButton retrocede a seleção para o botão anterior.
func (d *Dialog) PrevButton() {
	if d.Selected > 0 {
		d.Selected--
	}
}

// SelectedButton retorna o rótulo do botão selecionado.
func (d Dialog) SelectedButton() string {
	if d.Selected >= 0 && d.Selected < len(d.Buttons) {
		return d.Buttons[d.Selected]
	}
	return ""
}

// Render renderiza o diálogo como string centralizada na largura fornecida.
//
// Quando width é maior que zero, o conteúdo é quebrado para caber nessa
// largura: a saída nunca ultrapassa a largura pedida. Valores menores ou
// iguais a zero renderizam o diálogo com a largura natural.
func (d Dialog) Render(width int) string {
	if width < 0 {
		width = 0
	}

	boxed := d.renderBox(width)

	// Centraliza horizontalmente
	if width > 0 {
		boxed = lipgloss.PlaceHorizontal(width, lipgloss.Center, boxed)
	}

	return boxed
}

// RenderOverlay renderiza o diálogo centralizado vertical e horizontalmente
// dentro das dimensões fornecidas, útil para sobrepor em telas existentes.
//
// A saída tem exatamente height linhas (sem "\n" ao final) e nenhuma linha
// ultrapassa a largura informada. Quando o diálogo é mais alto que a área
// disponível, o conteúdo é truncado.
func (d Dialog) RenderOverlay(width, height int) string {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	if height == 0 {
		return ""
	}

	lines := strings.Split(d.renderBox(width), "\n")

	// O diálogo não pode ser mais alto que a área disponível.
	if len(lines) > height {
		lines = lines[:height]
	}

	// Centraliza verticalmente: o espaço restante é dividido entre o topo e
	// a base, de forma que a soma das duas margens seja exata.
	topPad := 0
	if height > len(lines) {
		topPad = (height - len(lines)) / 2
	}
	bottomPad := height - len(lines) - topPad

	// Centraliza horizontalmente o bloco inteiro do diálogo (e não cada
	// linha isoladamente), para que a borda e o conteúdo fiquem alinhados.
	dialogWidth := 0
	for _, line := range lines {
		if w := lipgloss.Width(line); w > dialogWidth {
			dialogWidth = w
		}
	}
	leftPad := ""
	if width > dialogWidth {
		leftPad = strings.Repeat(" ", (width-dialogWidth)/2)
	}

	// As linhas de margem são preenchidas com espaços para que a saída não
	// termine com uma quebra de linha extra.
	blankLine := ""
	if width > 0 {
		blankLine = strings.Repeat(" ", width)
	}

	out := make([]string, 0, height)
	for range topPad {
		out = append(out, blankLine)
	}
	for _, line := range lines {
		out = append(out, leftPad+line)
	}
	for range bottomPad {
		out = append(out, blankLine)
	}

	return strings.Join(out, "\n")
}

// renderBox renderiza o diálogo com a borda, limitando a largura total a
// maxWidth quando ela é maior que zero.
//
// O limite é aplicado aos blocos internos antes da borda, para que a
// mensagem longa seja quebrada em várias linhas em vez de estourar a tela.
func (d Dialog) renderBox(maxWidth int) string {
	s := d.Styles
	frame := s.BorderStyle.GetHorizontalFrameSize()

	// Conteúdo natural, sem limite de largura.
	content := d.renderContent(0)

	// Só restringe quando o conteúdo não cabe: assim um diálogo curto continua
	// sendo exibido com a largura natural.
	if maxWidth > 0 && lipgloss.Width(content)+frame > maxWidth {
		inner := maxWidth - frame
		if inner < 1 {
			inner = 1
		}
		content = d.renderContent(inner)
	}

	return s.BorderStyle.Render(content)
}

// renderContent monta o conteúdo interno do diálogo (sem a borda).
//
// Quando inner é maior que zero, cada bloco é limitado a essa largura.
func (d Dialog) renderContent(inner int) string {
	s := d.Styles

	titleStyle := s.TitleStyle
	messageStyle := s.MessageStyle
	if inner > 0 {
		titleStyle = titleStyle.Width(inner)
		messageStyle = messageStyle.Width(inner)
	}

	// Renderiza título
	titleStr := titleStyle.Render(d.Title)

	// Renderiza mensagem ou prompt de entrada
	var bodyStr string
	if d.Type == TypeInput {
		inputLine := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Underline(true).
			Render(d.InputValue + "_")

		// Prompt e campo são renderizados juntos: assim o padding vertical do
		// MessageStyle fica ao redor do bloco, sem criar uma linha em branco
		// entre o prompt e o campo.
		bodyStr = messageStyle.Render(d.InputPrompt + "\n" + inputLine)
	} else {
		bodyStr = messageStyle.Render(d.Message)
	}

	// Renderiza botões
	var buttons []string
	for i, btn := range d.Buttons {
		style := s.ButtonStyle
		if i == d.Selected {
			style = s.ActiveButtonStyle
		}
		buttons = append(buttons, style.Render(btn))
	}
	buttonsStr := lipgloss.JoinHorizontal(lipgloss.Center, buttons...)
	if inner > 0 {
		buttonsStr = lipgloss.NewStyle().Width(inner).Render(buttonsStr)
	}

	// Monta conteúdo interno
	return lipgloss.JoinVertical(lipgloss.Center,
		titleStr,
		bodyStr,
		"",
		buttonsStr,
	)
}
