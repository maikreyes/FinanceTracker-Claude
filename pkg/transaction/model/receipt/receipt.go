package receipt

// Interpreted es lo que un intérprete de IA extrae de la foto de un
// recibo. Amount es puntero para distinguir "no se encontró un monto"
// (nil) de "el monto es 0". CategoryName/PaymentMethod pueden venir vacíos
// si la IA no pudo determinarlos con certeza: el movimiento se registra
// igual sin esos campos (ver
// specs/features/007-interpretar-recibo-con-ia/spec.md).
type Interpreted struct {
	Amount        *float64
	Concept       string
	CategoryName  string
	PaymentMethod string
}
