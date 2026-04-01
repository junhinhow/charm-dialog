package dialog

import (
	"strings"
	"testing"
)

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
		t.Errorf("tipo esperado TypeAlert, obteve %d", d.Type)
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

func TestRenderOverlay(t *testing.T) {
	d := NewAlert("Info", "Mensagem de teste")
	output := d.RenderOverlay(80, 24)
	if output == "" {
		t.Error("overlay nao deveria retornar string vazia")
	}
}
