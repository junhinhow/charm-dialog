package dialog

import (
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// semANSI remove os códigos de escape ANSI do texto, para que as comparações
// de conteúdo não sejam afetadas pela estilização.
func semANSI(s string) string {
	return ansi.Strip(s)
}

// larguraMaxima devolve a maior largura de exibição entre as linhas de um
// texto (usa a largura de exibição, e não a contagem de bytes).
func larguraMaxima(s string) int {
	maior := 0
	for _, linha := range strings.Split(s, "\n") {
		if w := lipgloss.Width(linha); w > maior {
			maior = w
		}
	}
	return maior
}

// indiceDaBorda devolve o índice da linha que contém a borda superior.
func indiceDaBorda(linhas []string) int {
	for i, linha := range linhas {
		if strings.Contains(linha, "╭") {
			return i
		}
	}
	return -1
}

func TestNewConfirm(t *testing.T) {
	d := NewConfirm("Confirmar", "Tem certeza?")
	if d.Title != "Confirmar" {
		t.Errorf("titulo esperado 'Confirmar', obteve '%s'", d.Title)
	}
	if d.Message != "Tem certeza?" {
		t.Errorf("mensagem esperada 'Tem certeza?', obteve '%s'", d.Message)
	}
	if len(d.Buttons) != 2 {
		t.Fatalf("esperados 2 botoes, obteve %d", len(d.Buttons))
	}
	if d.Buttons[0] != "Yes" || d.Buttons[1] != "No" {
		t.Errorf("botoes esperados [Yes, No], obteve %v", d.Buttons)
	}
	if d.Type != TypeConfirm {
		t.Errorf("tipo esperado TypeConfirm, obteve %d", d.Type)
	}
}

func TestNewAlert(t *testing.T) {
	d := NewAlert("Aviso", "Operacao concluida")
	if len(d.Buttons) != 1 || d.Buttons[0] != "OK" {
		t.Errorf("botao esperado [OK], obteve %v", d.Buttons)
	}
	if d.Type != TypeAlert {
		t.Errorf("tipo esperado TypeAlert,obtained %d", d.Type)
	}
}

func TestNewInput(t *testing.T) {
	d := NewInput("Entrada", "Digite seu nome:")
	if d.InputPrompt != "Digite seu nome:" {
		t.Errorf("prompt esperado 'Digite seu nome:', obteve '%s'", d.InputPrompt)
	}
	if d.Type != TypeInput {
		t.Errorf("tipo esperado TypeInput, obteve %d", d.Type)
	}
	if len(d.Buttons) != 2 {
		t.Fatalf("esperados 2 botoes, obteve %d", len(d.Buttons))
	}
}

func TestNewCustom(t *testing.T) {
	btns := []string{"Salvar", "Descartar", "Cancelar"}
	d := NewCustom("Opcoes", "Escolha uma acao", btns)
	if len(d.Buttons) != 3 {
		t.Fatalf("esperados 3 botoes, obteve %d", len(d.Buttons))
	}
	if d.Type != TypeCustom {
		t.Errorf("tipo esperado TypeCustom, obteve %d", d.Type)
	}
}

func TestNavigation(t *testing.T) {
	d := NewConfirm("Teste", "nav")

	if d.Selected != 0 {
		t.Errorf("selecao inicial esperada 0, obteve %d", d.Selected)
	}

	d.NextButton()
	if d.Selected != 1 {
		t.Errorf("apos NextButton esperado 1, obteve %d", d.Selected)
	}

	// Nao deve ultrapassar o limite
	d.NextButton()
	if d.Selected != 1 {
		t.Errorf("nao deveria ultrapassar limite, esperado 1, obteve %d", d.Selected)
	}

	d.PrevButton()
	if d.Selected != 0 {
		t.Errorf("apos PrevButton esperado 0, obteve %d", d.Selected)
	}

	// Nao deve ir abaixo de zero
	d.PrevButton()
	if d.Selected != 0 {
		t.Errorf("nao deveria ir abaixo de 0, esperado 0, obteve %d", d.Selected)
	}
}

func TestSelectedButton(t *testing.T) {
	d := NewConfirm("Teste", "sel")
	if d.SelectedButton() != "Yes" {
		t.Errorf("botao selecionado esperado 'Yes', obteve '%s'", d.SelectedButton())
	}
	d.NextButton()
	if d.SelectedButton() != "No" {
		t.Errorf("botao selecionado esperado 'No', obteve '%s'", d.SelectedButton())
	}
}

func TestRender(t *testing.T) {
	d := NewConfirm("Confirmar", "Deseja continuar?")
	output := d.Render(80)
	if output == "" {
		t.Error("render nao deveria retornar string vazia")
	}
	if !strings.Contains(output, "Confirmar") {
		t.Error("output deveria conter o titulo")
	}
	if !strings.Contains(output, "Deseja continuar?") {
		t.Error("output deveria conter a mensagem")
	}
}

func TestRenderTemBorda(t *testing.T) {
	d := NewAlert("Info", "Mensagem")
	output := d.Render(80)

	for _, canto := range []string{"╭", "╮", "╰", "╯"} {
		if !strings.Contains(output, canto) {
			t.Errorf("output deveria conter o canto %q, obteve:\n%s", canto, output)
		}
	}
}

func TestRenderRespeitaLarguraMaxima(t *testing.T) {
	d := NewConfirm("Confirmar", strings.Repeat("palavra muito longa ", 40))
	output := d.Render(40)

	if largura := larguraMaxima(output); largura > 40 {
		t.Errorf("saida tem %d colunas, o maximo e 40:\n%s", largura, output)
	}
}

func TestRenderLarguraNegativaNaoPanica(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Render entrou em panico: %v", r)
		}
	}()

	d := NewAlert("Info", "Mensagem")
	if output := d.Render(-10); output == "" {
		t.Error("render nao deveria retornar string vazia")
	}
}

func TestRenderDialogoCurtoMantemLarguraNatural(t *testing.T) {
	d := NewAlert("Info", "Mensagem")
	natural := larguraMaxima(d.Render(0))
	if natural >= 80 {
		t.Fatalf("dialogo curto deveria ter largura natural menor que 80, obteve %d", natural)
	}

	output := d.Render(80)
	if largura := larguraMaxima(output); largura != 80 {
		t.Errorf("saida centralizada deveria ter 80 colunas, obteve %d", largura)
	}
}

func TestRenderTypeInput(t *testing.T) {
	d := NewInput("Entrada", "Digite seu nome:")
	d.InputValue = "abc"
	output := semANSI(d.Render(80))

	if !strings.Contains(output, "Digite seu nome:") {
		t.Error("output deveria conter o prompt")
	}
	if !strings.Contains(output, "abc_") {
		t.Errorf("output deveria conter o valor do campo de entrada, obteve:\n%s", output)
	}
}

func TestRenderTypeInputSemLinhaEmBrancoEntrePromptECampo(t *testing.T) {
	d := NewInput("Entrada", "Digite seu nome:")
	d.InputValue = "abc"
	linhas := strings.Split(semANSI(d.Render(80)), "\n")

	promptIdx, inputIdx := -1, -1
	for i, linha := range linhas {
		if promptIdx == -1 && strings.Contains(linha, "Digite seu nome:") {
			promptIdx = i
		}
		if inputIdx == -1 && strings.Contains(linha, "abc_") {
			inputIdx = i
		}
	}

	if promptIdx == -1 {
		t.Fatal("prompt nao encontrado na saida")
	}
	if inputIdx == -1 {
		t.Fatal("campo de entrada nao encontrado na saida")
	}
	if inputIdx != promptIdx+1 {
		t.Errorf("campo deveria ficar na linha seguinte ao prompt (prompt=%d, campo=%d)", promptIdx, inputIdx)
	}
}

func TestRenderOverlay(t *testing.T) {
	d := NewAlert("Info", "Mensagem de teste")
	output := d.RenderOverlay(80, 24)
	if output == "" {
		t.Error("overlay nao deveria retornar string vazia")
	}
}

func TestRenderOverlayCentralizaVerticalmente(t *testing.T) {
	d := NewAlert("Info", "Mensagem de teste")

	const altura = 24
	linhas := strings.Split(d.RenderOverlay(80, altura), "\n")

	if len(linhas) != altura {
		t.Fatalf("overlay deveria ter exatamente %d linhas, obteve %d", altura, len(linhas))
	}

	alturaDialogo := len(strings.Split(d.Render(80), "\n"))
	topo := (altura - alturaDialogo) / 2

	if idx := indiceDaBorda(linhas); idx != topo {
		t.Errorf("borda superior do dialogo deveria estar na linha %d, obteve %d", topo, idx)
	}

	// As linhas acima do diálogo só contêm espaços (centralização vertical).
	for i := range topo {
		if strings.TrimSpace(linhas[i]) != "" {
			t.Errorf("linha %d acima do dialogo deveria estar em branco, obteve %q", i, linhas[i])
		}
	}
}

func TestRenderOverlayCentralizaHorizontalmente(t *testing.T) {
	d := NewAlert("Info", "Mensagem")
	linhas := strings.Split(d.RenderOverlay(80, 24), "\n")

	idx := indiceDaBorda(linhas)
	if idx == -1 {
		t.Fatal("borda superior do dialogo nao encontrada")
	}

	coluna := len(linhas[idx]) - len(strings.TrimLeft(linhas[idx], " "))
	if coluna <= 0 {
		t.Fatalf("dialogo deveria estar deslocado da margem esquerda, coluna %d", coluna)
	}

	// Todas as linhas do diálogo começam na mesma coluna.
	for i := idx; i < len(linhas); i++ {
		conteudo := strings.TrimLeft(linhas[i], " ")
		if conteudo == "" {
			continue
		}
		if got := len(linhas[i]) - len(conteudo); got != coluna {
			t.Errorf("linha %d comeca na coluna %d, esperado %d", i, got, coluna)
		}
	}
}

func TestRenderOverlayNaoExcedeALargura(t *testing.T) {
	d := NewConfirm("Confirmar", strings.Repeat("palavra muito longa ", 40))
	output := d.RenderOverlay(40, 24)

	if largura := larguraMaxima(output); largura > 40 {
		t.Errorf("overlay tem %d colunas, o maximo e 40", largura)
	}
}

func TestRenderOverlayLarguraNegativaNaoPanica(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RenderOverlay entrou em panico: %v", r)
		}
	}()

	d := NewAlert("Info", "Mensagem")
	output := d.RenderOverlay(-10, 10)
	if output == "" {
		t.Error("overlay nao deveria retornar string vazia")
	}
}

func TestRenderOverlayAlturaNegativaNaoPanica(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RenderOverlay entrou em panico: %v", r)
		}
	}()

	d := NewAlert("Info", "Mensagem")
	if output := d.RenderOverlay(80, -10); output != "" {
		t.Errorf("altura negativa deveria devolver string vazia, obteve %q", output)
	}
}

func TestRenderOverlayNaoTerminaComNewline(t *testing.T) {
	d := NewAlert("Info", "Mensagem")
	output := d.RenderOverlay(80, 24)

	if strings.HasSuffix(output, "\n") {
		t.Error("overlay nao deveria terminar com \\n")
	}

	// Exatamente 24 linhas: 23 separadores e nenhuma linha extra no fim.
	if n := strings.Count(output, "\n"); n != 23 {
		t.Errorf("overlay de 24 linhas deveria ter 23 quebras de linha, obteve %d", n)
	}
}

func TestRenderOverlayAlturaMenorQueODialogo(t *testing.T) {
	d := NewAlert("Info", "Mensagem")

	const altura = 2
	linhas := strings.Split(d.RenderOverlay(80, altura), "\n")
	if len(linhas) != altura {
		t.Errorf("overlay deveria ter exatamente %d linhas, obteve %d", altura, len(linhas))
	}
}

func TestWithStylesAplicaNoDialogoOriginal(t *testing.T) {
	d := NewConfirm("Titulo", "Mensagem")

	estilos := DefaultStyles()
	estilos.BorderStyle = estilos.BorderStyle.Border(lipgloss.ASCIIBorder())

	// Sem reatribuir: os estilos precisam valer no diálogo original.
	d.WithStyles(estilos)

	saida := d.Render(80)
	if !strings.Contains(saida, "+") {
		t.Errorf("esperava a borda ASCII depois de WithStyles, obteve:\n%s", saida)
	}
	if strings.Contains(saida, "╭") {
		t.Errorf("a borda padrao nao deveria mais aparecer, obteve:\n%s", saida)
	}
}

func TestWithStylesDevolveCopia(t *testing.T) {
	d := NewConfirm("Titulo", "Mensagem")
	estiloOriginal := d.Styles.TitleStyle

	estilos := DefaultStyles()
	estilos.TitleStyle = estilos.TitleStyle.Bold(false)

	copia := d.WithStyles(estilos)

	if copia.Styles.TitleStyle.GetBold() {
		t.Error("a copia deveria ter os novos estilos")
	}
	if d.Styles.TitleStyle.GetBold() {
		t.Error("o dialogo original deveria ter recebido os novos estilos")
	}
	if !estiloOriginal.GetBold() {
		t.Error("o estilo original nao deveria ter sido alterado")
	}
}

func TestSetKeyMsgInsereTexto(t *testing.T) {
	d := NewInput("Entrada", "Prompt:")

	acao, consumida := d.SetKeyMsg(tea.KeyPressMsg{Text: "a", Code: 'a'})
	if !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if acao != ActionNone {
		t.Errorf("esperava ActionNone, obteve %v", acao)
	}
	if d.InputValue != "a" {
		t.Errorf("esperava 'a', obteve %q", d.InputValue)
	}

	d.SetKeyMsg(tea.KeyPressMsg{Text: "b", Code: 'b'})
	if d.InputValue != "ab" {
		t.Errorf("esperava 'ab', obteve %q", d.InputValue)
	}
}

func TestSetKeyMsgBackspace(t *testing.T) {
	d := NewInput("Entrada", "Prompt:")
	d.InputValue = "abc"

	_, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if d.InputValue != "ab" {
		t.Errorf("esperava 'ab', obteve %q", d.InputValue)
	}
}

func TestSetKeyMsgBackspaceComCampoVazio(t *testing.T) {
	d := NewInput("Entrada", "Prompt:")

	if _, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyBackspace}); !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if d.InputValue != "" {
		t.Errorf("esperava campo vazio, obteve %q", d.InputValue)
	}
}

func TestSetKeyMsgCtrlULimpaCampo(t *testing.T) {
	d := NewInput("Entrada", "Prompt:")
	d.InputValue = "abc"

	_, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl})
	if !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if d.InputValue != "" {
		t.Errorf("esperava campo vazio, obteve %q", d.InputValue)
	}
}

func TestSetKeyMsgEnterConfirma(t *testing.T) {
	d := NewConfirm("Titulo", "Mensagem")

	acao, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if acao != ActionSubmit {
		t.Errorf("esperava ActionSubmit, obteve %v", acao)
	}
}

func TestSetKeyMsgEscCancela(t *testing.T) {
	d := NewConfirm("Titulo", "Mensagem")

	acao, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if acao != ActionCancel {
		t.Errorf("esperava ActionCancel, obteve %v", acao)
	}
}

func TestSetKeyMsgNavegaEntreBotoes(t *testing.T) {
	d := NewConfirm("Titulo", "Mensagem")

	if _, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyTab}); !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if d.Selected != 1 {
		t.Errorf("apos tab esperado botao 1, obteve %d", d.Selected)
	}

	if _, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyRight}); !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if d.Selected != 1 {
		t.Errorf("apos right o botao nao deveria passar do limite, obteve %d", d.Selected)
	}

	if _, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyLeft}); !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if d.Selected != 0 {
		t.Errorf("apos left esperado botao 0, obteve %d", d.Selected)
	}
}

func TestSetKeyMsgShiftTabVoltaParaOBotaoAnterior(t *testing.T) {
	d := NewConfirm("Titulo", "Mensagem")
	d.NextButton()

	if _, consumida := d.SetKeyMsg(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}); !consumida {
		t.Error("a tecla deveria ter sido consumida")
	}
	if d.Selected != 0 {
		t.Errorf("apos shift+tab esperado botao 0, obteve %d", d.Selected)
	}
}

func TestSetKeyMsgNaoConsumidaEmDialogoSemEntrada(t *testing.T) {
	d := NewConfirm("Titulo", "Mensagem")

	_, consumida := d.SetKeyMsg(tea.KeyPressMsg{Text: "a", Code: 'a'})
	if consumida {
		t.Error("a tecla nao deveria ter sido consumida por um dialogo sem entrada")
	}
	if d.InputValue != "" {
		t.Errorf("o campo de entrada nao deveria ser alterado, obteve %q", d.InputValue)
	}
}

func TestHandleMsg(t *testing.T) {
	d := NewInput("Entrada", "Prompt:")

	acao, consumida := d.HandleMsg(tea.KeyPressMsg{Text: "x", Code: 'x'})
	if !consumida {
		t.Error("a mensagem de tecla deveria ter sido consumida")
	}
	if acao != ActionNone {
		t.Errorf("esperava ActionNone, obteve %v", acao)
	}
	if d.InputValue != "x" {
		t.Errorf("esperava 'x', obteve %q", d.InputValue)
	}

	// Mensagens que não são de tecla são ignoradas.
	if _, consumida := d.HandleMsg(nil); consumida {
		t.Error("uma mensagem desconhecida nao deveria ter sido consumida")
	}
}

func TestActionString(t *testing.T) {
	if ActionNone.String() != "none" {
		t.Errorf("esperava 'none', obteve %q", ActionNone.String())
	}
	if ActionSubmit.String() != "submit" {
		t.Errorf("esperava 'submit', obteve %q", ActionSubmit.String())
	}
	if ActionCancel.String() != "cancel" {
		t.Errorf("esperava 'cancel', obteve %q", ActionCancel.String())
	}
}
