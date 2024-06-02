package main

import (
	"context"
	"fmt"
	"inventory-system/models"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	telebot "gopkg.in/tucnak/telebot.v2"
)

func startBot() {
	token := "7016470919:AAEeljC0NV3VKsPmZVdmxWBBVjsnjz4HJNw"
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(m *telebot.Message) {
		bot.Send(m.Sender, "Bienvenido al sistema de inventario de medicamentos. Usa /add para agregar un medicamento y /list para ver todos los medicamentos.")
	})

	bot.Handle("/add", func(m *telebot.Message) {
		bot.Send(m.Sender, "Envía los detalles del medicamento en el siguiente formato:\nNombre, Principio Activo, Marca, Cantidad, Fecha de Vencimiento (YYYY-MM-DD)")
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		if m.Text == "/list" {
			var medicamentos []models.Medicamento
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			cursor, err := medicamentoCollection.Find(ctx, bson.M{})
			if err != nil {
				bot.Send(m.Sender, "Error al obtener los medicamentos.")
				return
			}
			defer cursor.Close(ctx)

			for cursor.Next(ctx) {
				var medicamento models.Medicamento
				cursor.Decode(&medicamento)
				medicamentos = append(medicamentos, medicamento)
			}

			var response string
			for _, med := range medicamentos {
				response += fmt.Sprintf("Nombre: %s, Principio Activo: %s, Marca: %s, Cantidad: %d, Fecha de Vencimiento: %s\n",
					med.Nombre, med.PrincipioActivo, med.Marca, med.Cantidad, med.FechaVencimiento.Format("2006-01-02"))
			}
			bot.Send(m.Sender, response)
		} else {
			var med models.Medicamento
			_, err := fmt.Sscanf(m.Text, "%s, %s, %s, %d, %s",
				&med.Nombre, &med.PrincipioActivo, &med.Marca, &med.Cantidad, &med.FechaVencimiento)
			if err != nil {
				bot.Send(m.Sender, "Error en el formato. Usa: Nombre, Principio Activo, Marca, Cantidad, Fecha de Vencimiento (YYYY-MM-DD)")
				return
			}
			med.ID = uuid.New().String()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err = medicamentoCollection.InsertOne(ctx, med)
			if err != nil {
				bot.Send(m.Sender, "Error al agregar el medicamento.")
				return
			}

			bot.Send(m.Sender, "Medicamento agregado con éxito.")
		}
	})

	bot.Start()
}
