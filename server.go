package main

import (
	"context"
	"inventory-system/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	telebot "gopkg.in/tucnak/telebot.v2"
)

var medicamentoCollection *mongo.Collection

func init() {
	medicamentoCollection = db.Collection("medicamentos")
}

func startServer() {
	r := gin.Default()

	r.GET("/medicamentos", func(c *gin.Context) {
		var medicamentos []models.Medicamento
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cursor, err := medicamentoCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var medicamento models.Medicamento
			cursor.Decode(&medicamento)
			medicamentos = append(medicamentos, medicamento)
		}

		c.JSON(http.StatusOK, medicamentos)
	})

	r.POST("/medicamentos", func(c *gin.Context) {
		var newMedicamento models.Medicamento
		if err := c.ShouldBindJSON(&newMedicamento); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newMedicamento.ID = uuid.New().String()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := medicamentoCollection.InsertOne(ctx, newMedicamento)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newMedicamento)
	})

	r.POST("/webhook", func(c *gin.Context) {
		bot, err := telebot.NewBot(telebot.Settings{
			Token:  os.Getenv("TELEGRAM_TOKEN"),
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			log.Fatal(err)
		}

		update := telebot.Update{}
		if err := c.ShouldBindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bot.ProcessUpdate(update)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}
