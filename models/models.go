package models

import "time"

type Medicamento struct {
	ID               string    `json:"id" bson:"_id,omitempty"`
	Nombre           string    `json:"nombre" bson:"nombre"`
	PrincipioActivo  string    `json:"principio_activo" bson:"principio_activo"`
	Marca            string    `json:"marca" bson:"marca"`
	Cantidad         int       `json:"cantidad" bson:"cantidad"`
	FechaVencimiento time.Time `json:"fecha_vencimiento" bson:"fecha_vencimiento"`
}
