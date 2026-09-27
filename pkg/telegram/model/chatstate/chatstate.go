// Package chatstate define lo que el bot recuerda de un chat entre un
// update y el siguiente.
package chatstate

import (
	"time"

	"finance-tracker/pkg/transaction/model/transaction"
)

// PendingTTL es cuánto tiempo puede esperar el bot una respuesta antes de
// olvidarla. Sin esto, un botón tocado por error dejaría al chat
// interpretando su próximo mensaje normal como el nombre de una categoría.
const PendingTTL = 30 * time.Minute

// WizardStep es el paso actual del flujo guiado de /egreso o /ingreso.
type WizardStep int

const (
	StepAmount WizardStep = iota
	StepConcept
	StepPaymentMethod
	StepCategory
)

// Wizard es el flujo guiado de /egreso o /ingreso en curso.
type Wizard struct {
	TxType  transaction.Type
	Step    WizardStep
	Concept string
	Amount  float64
	Method  transaction.PaymentMethod
}

// PendingKind identifica qué está esperando el bot de un chat. Un chat
// espera como máximo una cosa a la vez, y ChatState solo se construye con
// los constructores de abajo, así que activar una espera descarta la
// anterior por construcción.
type PendingKind string

const (
	PendingNone         PendingKind = ""
	PendingPhotoType    PendingKind = "photo_type"    // foto sin caption, esperando Ingreso/Egreso
	PendingWizard       PendingKind = "wizard"        // flujo guiado en curso
	PendingCategoryName PendingKind = "category_name" // esperando el nombre de una categoría nueva
)

// ChatState es todo lo que el bot recuerda de un chat entre un update y el
// siguiente. El valor cero significa "nada pendiente".
type ChatState struct {
	Kind PendingKind `json:"kind,omitempty"`

	// PendingPhotoType: file_id de Telegram de la foto (no los bytes, para
	// no meter imágenes en un store externo) y message_id de la pregunta
	// Ingreso/Egreso, para que solo esos botones puedan resolverla y no los
	// de una foto anterior.
	PhotoFileID     string `json:"photo_file_id,omitempty"`
	PromptMessageID int    `json:"prompt_message_id,omitempty"`

	// PendingWizard.
	Wizard *Wizard `json:"wizard,omitempty"`

	ExpiresAt time.Time `json:"expires_at"`
}

func NewPhotoPending(fileID string, promptMessageID int, now time.Time) ChatState {
	return ChatState{Kind: PendingPhotoType, PhotoFileID: fileID, PromptMessageID: promptMessageID, ExpiresAt: now.Add(PendingTTL)}
}

func NewWizardPending(w Wizard, now time.Time) ChatState {
	return ChatState{Kind: PendingWizard, Wizard: &w, ExpiresAt: now.Add(PendingTTL)}
}

func NewCategoryPending(now time.Time) ChatState {
	return ChatState{Kind: PendingCategoryName, ExpiresAt: now.Add(PendingTTL)}
}

func (s ChatState) Expired(now time.Time) bool {
	return s.Kind != PendingNone && now.After(s.ExpiresAt)
}
