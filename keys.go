package dialog

import (
	"charm.land/bubbletea/v2"
)

// Action representa uma ação disparada pelo usuário dentro do diálogo.
type Action int

const (
	// ActionNone indica que nenhuma ação foi disparada.
	ActionNone Action = iota
	// ActionSubmit indica que o usuário confirmou o diálogo (Enter).
	ActionSubmit
	// ActionCancel indica que o usuário cancelou o diálogo (Esc).
	ActionCancel
)

// String devolve o nome da ação, útil em logs e mensagens de erro.
func (a Action) String() string {
	switch a {
	case ActionSubmit:
		return "submit"
	case ActionCancel:
		return "cancel"
	default:
		return "none"
	}
}

// SetKeyMsg processa uma tecla no diálogo.
//
// O primeiro valor devolvido é a ação disparada pela tecla e o segundo indica
// se a tecla foi consumida pelo diálogo. Quando a tecla não é consumida, a
// aplicação deve seguir com o processamento normal do evento.
//
// Teclas tratadas:
//
//	"left"       seleciona o botão anterior
//	"right"      seleciona o próximo botão
//	"tab"        seleciona o próximo botão
//	"shift+tab"  seleciona o botão anterior
//	"enter"      confirma o diálogo (ActionSubmit)
//	"esc"        cancela o diálogo (ActionCancel)
//
// Em diálogos do tipo TypeInput os caracteres imprimíveis são inseridos no
// campo de entrada, "backspace" apaga o último caractere e "ctrl+u" limpa o
// campo. O cursor do texto não é móvel: o texto é sempre digitado no fim.
func (d *Dialog) SetKeyMsg(msg tea.KeyPressMsg) (Action, bool) {
	key := msg.Key()

	switch key.Code {
	case tea.KeyEscape:
		return ActionCancel, true
	case tea.KeyEnter:
		return ActionSubmit, true
	case tea.KeyTab:
		if key.Mod.Contains(tea.ModShift) {
			d.PrevButton()
		} else {
			d.NextButton()
		}
		return ActionNone, true
	case tea.KeyLeft:
		d.PrevButton()
		return ActionNone, true
	case tea.KeyRight:
		d.NextButton()
		return ActionNone, true
	case tea.KeyBackspace:
		if d.Type == TypeInput {
			d.deleteLastInputChar()
			return ActionNone, true
		}
		return ActionNone, false
	}

	// "ctrl+u" limpa o campo de entrada.
	if d.Type == TypeInput && key.Mod.Contains(tea.ModCtrl) &&
		(key.Code == 'u' || key.Code == 'U') {
		d.InputValue = ""
		return ActionNone, true
	}

	// Caracteres imprimíveis vão para o campo de entrada.
	if d.Type == TypeInput && key.Text != "" &&
		!key.Mod.Contains(tea.ModCtrl) && !key.Mod.Contains(tea.ModAlt) {
		d.InputValue += key.Text
		return ActionNone, true
	}

	return ActionNone, false
}

// HandleMsg processa uma mensagem do Bubble Tea.
//
// As mensagens de tecla são delegadas a [Dialog.SetKeyMsg]; as demais são
// ignoradas. É a forma prática de delegar o teclado ao diálogo dentro do
// Update de um modelo:
//
//	case tea.KeyPressMsg:
//	    if action, ok := d.HandleMsg(msg); ok {
//	        switch action {
//	        case dialog.ActionSubmit:
//	            // confirmado
//	        case dialog.ActionCancel:
//	            // cancelado
//	        }
//	    }
func (d *Dialog) HandleMsg(msg tea.Msg) (Action, bool) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		return d.SetKeyMsg(key)
	}
	return ActionNone, false
}

// deleteLastInputChar apaga o último caractere do campo de entrada.
func (d *Dialog) deleteLastInputChar() {
	runes := []rune(d.InputValue)
	if len(runes) == 0 {
		return
	}
	d.InputValue = string(runes[:len(runes)-1])
}
