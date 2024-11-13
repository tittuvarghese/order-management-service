package service

import (
	tx "github.com/tittuvarghese/core/storage"
	"github.com/tittuvarghese/order-management-service/core/database"
	"github.com/tittuvarghese/order-management-service/models"
)

func CreateOrder(order models.Order, storage *database.RelationalDatabase) error {

	var transaction = database.DbTxn

	// Building transaction
	// 1. Order creation
	orderCreation := tx.Operation{
		Model:   &order,
		Command: database.CreateCommand,
	}

	transaction.Operations = append(transaction.Operations, orderCreation)

	// Update Quantity
	for _, item := range order.Items {
		condition := map[string]interface{}{"id": item.ProductID}
		var queryUpdate = database.DbOps
		queryUpdate.Model = &models.Product{}
		queryUpdate.Condition = condition
		queryUpdate.Command = database.UpdateCommand

		var queryExpr = database.DbExpr
		queryExpr.Column = "quantity"
		queryExpr.Value = storage.Instance.BuildExpr("quantity - ?", item.Quantity) // gorm.Expr("quantity - ?", item.Quantity)
		queryUpdate.Expr = queryExpr

		//queryUpdate.Expr = tx.Expr{
		//	Column: "quantity",
		//	Value:  gorm.Expr("quantity - ?", item.Quantity),
		//}
		//var queryUpdate = database.DbOps
		//queryUpdate.M

		transaction.Operations = append(transaction.Operations, queryUpdate)
	}

	err := storage.Instance.Transaction(transaction)

	if err != nil {
		return err
	}
	return nil

}
