package handler

import (
	"context"
	"strings"
	"testing"
	"time"

	"finance-tracker/pkg/telegram/model/chatstate"
	"finance-tracker/pkg/telegram/model/update"
	"finance-tracker/pkg/transaction/model/receipt"
	"finance-tracker/pkg/transaction/model/summary"
	"finance-tracker/pkg/transaction/model/transaction"
)

func amountOf(v float64) *float64 { return &v }

// --- autorización ---

func TestDispatch_ChatNoAutorizado_NoHaceNada(t *testing.T) {
	h := newHarness(t)

	h.handler.Dispatch(context.Background(), update.Update{Message: &update.Message{Chat: update.Chat{ID: 999}, Text: "egreso mecato 13000 efectivo"}})
	h.handler.Dispatch(context.Background(), update.Update{CallbackQuery: &update.CallbackQuery{
		ID: "cq", Data: catCallbackAdd, Message: &update.Message{Chat: update.Chat{ID: 999}, MessageID: 5},
	}})

	if h.writer.created != 0 {
		t.Error("un chat no autorizado no debería registrar nada")
	}
	if len(h.messenger.texts)+len(h.messenger.edits)+len(h.messenger.answered) != 0 {
		t.Error("un chat no autorizado no debería recibir ninguna respuesta")
	}
	if h.store.sets != 0 {
		t.Error("un chat no autorizado no debería dejar estado")
	}
}

func TestDispatch_CallbackSinMensaje_SeIgnora(t *testing.T) {
	h := newHarness(t)

	h.handler.Dispatch(context.Background(), update.Update{CallbackQuery: &update.CallbackQuery{ID: "cq", Data: catCallbackAdd}})

	if len(h.messenger.answered) != 0 || h.store.sets != 0 {
		t.Error("un callback sin mensaje asociado no se puede autorizar y debe ignorarse")
	}
}

// --- mensajes que no traen nada que interpretar ---

func TestDispatch_MensajeSinTextoNiFoto_NoAvanzaElFlujo(t *testing.T) {
	h := newHarness(t)
	h.say("/egreso")
	h.say("13000")
	h.say("mecato")
	h.say("efectivo")
	before := h.state(t)
	sent := len(h.messenger.texts)

	// Sticker, audio o documento: llegan sin texto ni foto.
	h.handler.Dispatch(context.Background(), update.Update{Message: &update.Message{Chat: update.Chat{ID: testChat}}})

	after := h.state(t)
	if after.Wizard == nil || after.Wizard.Step != before.Wizard.Step {
		t.Errorf("el flujo debería seguir en %v, quedó en %+v", before.Wizard, after.Wizard)
	}
	if h.writer.created != 0 {
		t.Error("un mensaje vacío no debería registrar nada")
	}
	if len(h.messenger.texts) != sent {
		t.Error("un mensaje vacío no debería recibir respuesta")
	}
}

// --- estado pendiente: exclusión mutua ---

func TestDispatch_FotoCancelaFlujoGuiado(t *testing.T) {
	h := newHarness(t)
	h.say("/egreso")

	h.sendPhoto("file-abc", "")

	state := h.state(t)
	if state.Kind != chatstate.PendingPhotoType || state.Wizard != nil {
		t.Fatalf("estado = %+v, want solo chatstate.PendingPhotoType", state)
	}
	if state.PhotoFileID != "file-abc" {
		t.Errorf("PhotoFileID = %q, want file-abc (la de mayor calidad)", state.PhotoFileID)
	}
	if state.PromptMessageID != 1 {
		t.Errorf("PromptMessageID = %d, want 1", state.PromptMessageID)
	}
}

func TestDispatch_FotoConCaptionDescartaFotoPendienteVieja(t *testing.T) {
	h := newHarness(t)
	h.sendPhoto("vieja", "")

	h.sendPhoto("nueva", "egreso mecato 13000 efectivo")

	if state := h.state(t); state.Kind != chatstate.PendingNone {
		t.Errorf("estado = %+v, la foto vieja no debería seguir pendiente", state)
	}
	if h.writer.created != 1 {
		t.Errorf("registros = %d, want 1 (el de la foto con caption)", h.writer.created)
	}
}

func TestDispatch_ComienzoDeFlujoDescartaFotoPendiente(t *testing.T) {
	h := newHarness(t)
	h.sendPhoto("file-1", "")

	h.say("/ingreso")

	state := h.state(t)
	if state.Kind != chatstate.PendingWizard || state.PhotoFileID != "" {
		t.Errorf("estado = %+v, want solo chatstate.PendingWizard sin foto", state)
	}
}

func TestDispatch_AgregarCategoriaCancelaFlujoGuiado(t *testing.T) {
	h := newHarness(t)
	h.say("/egreso")

	h.tap(10, catCallbackAdd)

	state := h.state(t)
	if state.Kind != chatstate.PendingCategoryName || state.Wizard != nil {
		t.Errorf("estado = %+v, want solo chatstate.PendingCategoryName", state)
	}
}

func TestDispatch_ComandosYResumenCancelanFlujoGuiado(t *testing.T) {
	for _, text := range []string{"/ayuda", "/ultimo", "/categorias", "/resumen"} {
		t.Run(text, func(t *testing.T) {
			h := newHarness(t)
			h.say("/egreso")
			h.say("13000")

			h.say(text)

			if state := h.state(t); state.Kind != chatstate.PendingNone {
				t.Errorf("estado = %+v, %q debería cancelar el flujo", state, text)
			}
			if h.writer.created != 0 {
				t.Errorf("%q no debería registrar nada", text)
			}
		})
	}
}

func TestDispatch_TextoNormalNoTocaFotoPendiente(t *testing.T) {
	h := newHarness(t)
	h.sendPhoto("file-1", "")
	before := h.state(t)
	setsBefore := h.store.sets

	h.say("egreso mecato 13000 efectivo")

	if h.writer.created != 1 {
		t.Fatalf("registros = %d, want 1", h.writer.created)
	}
	if after := h.state(t); after != before {
		t.Errorf("estado = %+v, want %+v (la foto sigue esperando sus botones)", after, before)
	}
	if h.store.sets != setsBefore {
		t.Error("un mensaje que no cambia lo pendiente no debería escribir en el store")
	}
}

// --- vencimiento ---

func TestDispatch_EsperaDeCategoriaVencidaNoSeTragaElSiguienteMensaje(t *testing.T) {
	h := newHarness(t)
	h.tap(10, catCallbackAdd)
	h.clock = h.clock.Add(chatstate.PendingTTL + time.Minute)

	h.say("egreso mecato 13000 efectivo")

	if len(h.categories.created) != 0 {
		t.Errorf("categorías creadas = %v, la espera vencida no debería crear nada", h.categories.created)
	}
	if h.writer.created != 1 {
		t.Error("el mensaje debería registrarse como movimiento normal")
	}
}

func TestDispatch_BotonesDeFotoVencidosYaNoAplican(t *testing.T) {
	h := newHarness(t)
	h.sendPhoto("file-1", "")
	h.clock = h.clock.Add(chatstate.PendingTTL + time.Minute)

	h.tap(1, tipoCallbackData(transaction.TypeExpense))

	if h.writer.created != 0 {
		t.Error("una pregunta vencida no debería registrar")
	}
	if h.messenger.lastEdit() != staleQuestionText {
		t.Errorf("edit = %q, want %q", h.messenger.lastEdit(), staleQuestionText)
	}
}

// --- fallo al leer el estado ---

func TestDispatch_FalloAlLeerEstado_NoPisaElEstado(t *testing.T) {
	h := newHarness(t)
	h.say("/egreso")
	h.say("13000")
	before := h.state(t)
	h.store.getErr = errBoom
	setsBefore := h.store.sets

	h.say("mecato")

	if h.store.sets != setsBefore {
		t.Error("con el store caído no se debe escribir encima del estado real")
	}
	if !strings.Contains(h.messenger.lastText(), "problema técnico") {
		t.Errorf("respuesta = %q, debería avisar del problema", h.messenger.lastText())
	}
	h.store.getErr = nil
	if after := h.state(t); after.Wizard == nil || after.Wizard.Step != before.Wizard.Step {
		t.Errorf("estado = %+v, want el flujo intacto en %+v", after.Wizard, before.Wizard)
	}
}

// --- panics ---

func TestDispatch_PanicEnHandlerNoTumbaElProceso(t *testing.T) {
	h := newHarness(t)
	h.handler.Summarizer = fakeSummarizer{panics: true}

	h.say("/resumen") // no debe propagar el panic
}

// --- flujo guiado ---

func TestDispatch_FlujoGuiadoCompleto(t *testing.T) {
	h := newHarness(t)

	for _, text := range []string{"/egreso", "13000", "mecato", "efectivo", "comida"} {
		h.say(text)
	}

	if state := h.state(t); state.Kind != chatstate.PendingNone {
		t.Errorf("estado = %+v, el flujo debería haber terminado", state)
	}
	got := h.writer.got
	if got.Name != "mecato" || got.Amount != 13000 || got.CategoryID != "cat-123" || got.PaymentMethod != transaction.PaymentMethodCash {
		t.Errorf("transacción registrada = %+v, no coincide con el flujo contestado", got)
	}
}

func TestDispatch_FlujoGuiado_MontoInvalidoNoAvanza(t *testing.T) {
	h := newHarness(t)
	h.say("/egreso")

	for _, bad := range []string{"NaN", "Inf", "-5", "0", "abc"} {
		h.say(bad)
		if state := h.state(t); state.Wizard == nil || state.Wizard.Step != chatstate.StepAmount {
			t.Fatalf("después de %q el flujo debería seguir en chatstate.StepAmount, estado = %+v", bad, state)
		}
	}
}

func TestDispatch_FlujoGuiado_CategoriaInvalidaPermiteReintentar(t *testing.T) {
	h := newHarness(t)
	for _, text := range []string{"/egreso", "13000", "-", "efectivo", "inventada"} {
		h.say(text)
	}

	if !strings.Contains(h.messenger.lastText(), "inventada") || !strings.Contains(h.messenger.lastText(), "comida") {
		t.Errorf("respuesta = %q, debería nombrar la categoría inválida y las válidas", h.messenger.lastText())
	}
	if state := h.state(t); state.Wizard == nil || state.Wizard.Step != chatstate.StepCategory {
		t.Fatalf("estado = %+v, want el flujo esperando otra categoría", state)
	}

	h.say("comida")
	if h.writer.created != 1 {
		t.Error("con una categoría válida debería registrar")
	}
}

func TestAdvanceWizard(t *testing.T) {
	w := chatstate.Wizard{TxType: transaction.TypeExpense, Step: chatstate.StepAmount}

	w, _, ok := advanceWizard(w, "13000")
	if !ok || w.Amount != 13000 || w.Step != chatstate.StepConcept {
		t.Fatalf("monto: ok=%v wizard=%+v", ok, w)
	}
	w, _, ok = advanceWizard(w, "-")
	if !ok || w.Concept != "" || w.Step != chatstate.StepPaymentMethod {
		t.Fatalf("concepto omitido: ok=%v wizard=%+v", ok, w)
	}
	if _, _, ok = advanceWizard(w, "bitcoin"); ok {
		t.Fatal("un medio de pago desconocido no debería avanzar")
	}
	w, _, ok = advanceWizard(w, "banco")
	if !ok || w.Method != transaction.PaymentMethodBank || w.Step != chatstate.StepCategory {
		t.Fatalf("medio de pago: ok=%v wizard=%+v", ok, w)
	}
}

// --- categorías ---

func TestDispatch_CategoriasAgregar_CreaCategoria(t *testing.T) {
	h := newHarness(t)

	h.tap(10, catCallbackAdd)
	if state := h.state(t); state.Kind != chatstate.PendingCategoryName {
		t.Fatalf("estado = %+v, want chatstate.PendingCategoryName", state)
	}

	h.say("Mascotas")

	if len(h.categories.created) != 1 || h.categories.created[0] != "Mascotas" {
		t.Errorf("categories.created = %v, want [Mascotas]", h.categories.created)
	}
	if state := h.state(t); state.Kind != chatstate.PendingNone {
		t.Errorf("estado = %+v, la espera del nombre debería haber terminado", state)
	}
}

func TestDispatch_CategoriasListar(t *testing.T) {
	h := newHarness(t)
	h.categories.names = []string{"comida", "ahorro"}

	h.tap(10, catCallbackList)

	if got := h.messenger.lastEdit(); !strings.Contains(got, "comida") || !strings.Contains(got, "ahorro") {
		t.Errorf("edit = %q, debería listar las categorías", got)
	}
	if len(h.messenger.answered) != 1 {
		t.Errorf("answered = %d, want 1 (el toque se reconoce una sola vez)", len(h.messenger.answered))
	}
}

// --- foto sin caption: botones Ingreso/Egreso ---

func TestDispatch_FotoSinCaption_PreguntaTipo(t *testing.T) {
	h := newHarness(t)

	h.sendPhoto("file-1", "")

	if len(h.messenger.keyboards) != 1 || len(h.messenger.keyboards[0]) != 2 {
		t.Fatalf("teclados = %+v, want uno con dos botones", h.messenger.keyboards)
	}
	if len(h.messenger.downloaded) != 0 {
		t.Error("la foto no debería descargarse hasta que el usuario elija el tipo")
	}
}

func TestDispatch_FotoSinCaption_TecladoFallidoNoDejaEstado(t *testing.T) {
	h := newHarness(t)
	h.messenger.keyboardErr = errBoom

	h.sendPhoto("file-1", "")

	if state := h.state(t); state.Kind != chatstate.PendingNone {
		t.Errorf("estado = %+v, sin botones en pantalla no debería quedar nada pendiente", state)
	}
}

func TestDispatch_Tipo_RegistraConIAYConsumeElEstado(t *testing.T) {
	h := newHarness(t)
	h.interpreter.result = receipt.Interpreted{Amount: amountOf(45000), Concept: "supermercado"}
	h.sendPhoto("file-1", "")

	h.tap(1, tipoCallbackData(transaction.TypeExpense))

	if h.writer.created != 1 || h.writer.got.Amount != 45000 || h.writer.got.Type != transaction.TypeExpense {
		t.Fatalf("registro = %+v (creados: %d)", h.writer.got, h.writer.created)
	}
	if got := h.messenger.downloaded; len(got) != 1 || got[0] != "file-1" {
		t.Errorf("descargas = %v, want [file-1]", got)
	}
	if !strings.Contains(h.messenger.lastText(), "Interpretado por IA") {
		t.Errorf("respuesta = %q, debería avisar que lo interpretó la IA", h.messenger.lastText())
	}
	if state := h.state(t); state.Kind != chatstate.PendingNone {
		t.Errorf("estado = %+v, debería estar consumido", state)
	}
}

func TestDispatch_Tipo_DobleToqueRegistraUnaSolaVez(t *testing.T) {
	h := newHarness(t)
	h.interpreter.result = receipt.Interpreted{Amount: amountOf(45000)}
	h.sendPhoto("file-1", "")

	h.tap(1, tipoCallbackData(transaction.TypeExpense))
	h.tap(1, tipoCallbackData(transaction.TypeExpense))

	if h.writer.created != 1 {
		t.Errorf("registros = %d, want 1", h.writer.created)
	}
}

func TestDispatch_Tipo_BotonesDeFotoAnteriorNoApliquenALaNueva(t *testing.T) {
	h := newHarness(t)
	h.interpreter.result = receipt.Interpreted{Amount: amountOf(45000)}
	h.sendPhoto("vieja", "") // pregunta = mensaje 1
	h.sendPhoto("nueva", "") // pregunta = mensaje 2

	h.tap(1, tipoCallbackData(transaction.TypeExpense))

	if h.writer.created != 0 {
		t.Error("los botones de la foto vieja no deberían registrar la nueva")
	}
	if h.messenger.lastEdit() != staleQuestionText {
		t.Errorf("edit = %q, want %q", h.messenger.lastEdit(), staleQuestionText)
	}
	if state := h.state(t); state.Kind != chatstate.PendingPhotoType || state.PhotoFileID != "nueva" {
		t.Errorf("estado = %+v, la foto nueva debería seguir pendiente", state)
	}

	h.tap(2, tipoCallbackData(transaction.TypeExpense))
	if h.writer.created != 1 {
		t.Error("los botones de la foto nueva sí deberían registrar")
	}
}

func TestDispatch_Tipo_DatoDesconocidoNoConsumeElEstado(t *testing.T) {
	h := newHarness(t)
	h.sendPhoto("file-1", "")

	h.tap(1, tipoCallbackPrefix+"algo_raro")

	if state := h.state(t); state.Kind != chatstate.PendingPhotoType {
		t.Errorf("estado = %+v, un callback inválido no debería perder la foto", state)
	}
}

func TestDispatch_Tipo_FalloDeDescargaPermiteReintentar(t *testing.T) {
	h := newHarness(t)
	h.interpreter.result = receipt.Interpreted{Amount: amountOf(45000)}
	h.sendPhoto("file-1", "")
	h.messenger.downloadErr = errBoom

	h.tap(1, tipoCallbackData(transaction.TypeExpense))

	if state := h.state(t); state.Kind != chatstate.PendingPhotoType {
		t.Fatalf("estado = %+v, la foto debería seguir pendiente para reintentar", state)
	}

	h.messenger.downloadErr = nil
	h.tap(1, tipoCallbackData(transaction.TypeExpense))
	if h.writer.created != 1 {
		t.Errorf("registros = %d, want 1 en el reintento", h.writer.created)
	}
}

func TestDispatch_Tipo_FalloDeIAPermiteReintentar(t *testing.T) {
	h := newHarness(t)
	h.sendPhoto("file-1", "")
	h.interpreter.err = errBoom

	h.tap(1, tipoCallbackData(transaction.TypeExpense))

	if state := h.state(t); state.Kind != chatstate.PendingPhotoType {
		t.Fatalf("estado = %+v, la foto debería seguir pendiente tras un fallo de la IA", state)
	}
	if !strings.Contains(h.messenger.lastText(), "Tocá el botón de nuevo") {
		t.Errorf("respuesta = %q, debería indicar cómo reintentar", h.messenger.lastText())
	}
}

func TestDispatch_Tipo_SinMontoNoReintenta(t *testing.T) {
	h := newHarness(t)
	h.sendPhoto("file-1", "") // el intérprete devuelve sin monto

	h.tap(1, tipoCallbackData(transaction.TypeExpense))

	if state := h.state(t); state.Kind != chatstate.PendingNone {
		t.Errorf("estado = %+v, reintentar la misma foto daría lo mismo", state)
	}
	if h.writer.created != 0 {
		t.Error("sin monto no debería registrar nada")
	}
}

func TestDispatch_Tipo_ConOtraEsperaPendienteLaDevuelveIntacta(t *testing.T) {
	h := newHarness(t)
	h.sendPhoto("file-1", "")
	h.say("/egreso") // el flujo guiado reemplaza a la foto
	before := h.state(t)

	h.tap(1, tipoCallbackData(transaction.TypeExpense))

	if after := h.state(t); after.Kind != chatstate.PendingWizard || after.Wizard == nil || *after.Wizard != *before.Wizard {
		t.Errorf("estado = %+v, want el flujo guiado intacto", after)
	}
	if h.messenger.lastEdit() != staleQuestionText {
		t.Errorf("edit = %q, want %q", h.messenger.lastEdit(), staleQuestionText)
	}
}

// --- comandos de solo lectura ---

func TestDispatch_ResumenYComandos(t *testing.T) {
	h := newHarness(t)
	h.handler.Summarizer = fakeSummarizer{summary: summary.Monthly{Egreso: 13000, Ingreso: 3000000, CountEgreso: 1, CountIngreso: 1}}

	h.say("/resumen")
	if !strings.Contains(h.messenger.lastText(), "2987000") {
		t.Errorf("resumen = %q, debería traer el balance", h.messenger.lastText())
	}

	h.say("/ayuda")
	if h.messenger.lastText() != ayudaText {
		t.Errorf("/ayuda = %q", h.messenger.lastText())
	}

	h.say("/ultimo")
	if !strings.Contains(h.messenger.lastText(), "Todavía no hay movimientos") {
		t.Errorf("/ultimo = %q", h.messenger.lastText())
	}

	h.say("/categorias")
	if len(h.messenger.keyboards) != 1 {
		t.Errorf("teclados = %d, /categorias debería mandar los botones", len(h.messenger.keyboards))
	}
}
