package service

import (
	"fmt"
	tx "github.com/tittuvarghese/core/storage"
	"github.com/tittuvarghese/order-management-service/core/database"
	"github.com/tittuvarghese/order-management-service/models"
	"gorm.io/gorm"
)

func CreateOrder(order models.Order, storage *database.RelationalDatabase) error {

	// err := storage.Instance.Insert(&order)
	var transaction tx.AtomicTransaction

	// Building transaction
	// 1. Order creation
	orderCreation := tx.Operation{
		Model:   &order,
		Command: tx.CreateCommand,
	}

	transaction.Operations = append(transaction.Operations, orderCreation)

	// Update Quantity
	for _, item := range order.Items {
		condition := map[string]interface{}{"id": item.ProductID}
		qtyUpdate := tx.Operation{
			Model:     &models.Product{},
			Command:   tx.UpdateCommand,
			Condition: condition,
			Expr: tx.Expr{
				Column: "quantity",
				Value:  gorm.Expr("quantity - ?", item.Quantity),
			},
		}
		transaction.Operations = append(transaction.Operations, qtyUpdate)
	}

	err := storage.Instance.Transaction(transaction)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil

}
