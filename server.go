package main

import (
	"context"
	"inventory-system/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	r.GET("/medicamentos/:nombre", func(c *gin.Context) {
		nombre := c.Param("nombre")
		var medicamento models.Medicamento

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := medicamentoCollection.FindOne(ctx, bson.M{"nombre": nombre}).Decode(&medicamento)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"message": "Medicamento no encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, medicamento)
	})

	r.Run(":8080")
}
