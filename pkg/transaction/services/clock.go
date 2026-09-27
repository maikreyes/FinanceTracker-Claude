package services

import "time"

// Colombia no tiene horario de verano, así que un offset fijo evita
// depender de que el host tenga la base de zonas horarias instalada.
var bogota = time.FixedZone("America/Bogota", -5*60*60)

// nowBogota devuelve la hora actual en Bogotá: la fecha de un movimiento
// es la del día del usuario, no la del servidor donde corre el bot (en un
// host en UTC, un gasto de las 8pm caería en el día siguiente).
func nowBogota() time.Time { return time.Now().In(bogota) }
