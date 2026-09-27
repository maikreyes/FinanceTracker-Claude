package register

import "finance-tracker/pkg/transaction/model/transaction"

// Nombre y tipo con los que se sube una foto de Telegram: Telegram siempre
// entrega las fotos comprimidas como JPEG, sin importar el formato original.
const (
	ReceiptFilename    = "recibo.jpg"
	ReceiptContentType = "image/jpeg"
)

// ParsedMessage es la salida cruda de parsear un mensaje: Transaction trae
// todo lo que ya mapea 1:1 a Notion; CategoryName es el texto de categoría
// tal como vino en el mensaje (vacío si no hubo), sin resolver.
type ParsedMessage struct {
	Transaction  transaction.Transaction
	CategoryName string
}

// Result trae todo lo necesario para confirmarle al usuario qué se
// registró, sin que el caller tenga que volver a resolver nada contra
// Notion.
type Result struct {
	PageID       string
	Transaction  transaction.Transaction
	CategoryName string // tal como vino en el mensaje; vacío si no hubo
	Interpreted  bool   // true si los datos vinieron de IA, no del usuario
}

// ConversationInput trae los campos ya separados y validados que
// recolectó el flujo guiado de /egreso o /ingreso: a diferencia de un
// mensaje de texto, no hay una línea que parsear.
type ConversationInput struct {
	Type          transaction.Type
	Concept       string // vacío → se usa string(Type) como Name
	Amount        float64
	PaymentMethod transaction.PaymentMethod
	CategoryName  string // vacío = sin categoría
}
