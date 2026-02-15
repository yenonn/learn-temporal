package activity

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"

	"learn-temporal/sample-app/model"
)

type Activities struct{}

func (a *Activities) ValidateOrder(ctx context.Context, order model.Order) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Validating order", "orderID", order.OrderID)

	if order.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}
	if order.CustomerID == "" {
		return fmt.Errorf("customer ID is required")
	}
	if len(order.Items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}
	if order.TotalAmount <= 0 {
		return fmt.Errorf("total amount must be positive")
	}
	if order.Email == "" {
		return fmt.Errorf("email is required")
	}

	// Simulate validation latency
	time.Sleep(200 * time.Millisecond)

	logger.Info("Order validated successfully", "orderID", order.OrderID)
	return nil
}

func (a *Activities) ProcessPayment(ctx context.Context, order model.Order) (model.PaymentResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Processing payment", "orderID", order.OrderID, "amount", order.TotalAmount)

	// 10% simulated transient failure to demonstrate retries
	if rand.Float64() < 0.1 {
		return model.PaymentResult{}, fmt.Errorf("payment gateway temporarily unavailable")
	}

	// Simulate payment processing latency
	time.Sleep(500 * time.Millisecond)

	result := model.PaymentResult{
		TransactionID: fmt.Sprintf("txn-%s-%d", order.OrderID, time.Now().UnixMilli()),
		Status:        "charged",
	}

	logger.Info("Payment processed", "orderID", order.OrderID, "transactionID", result.TransactionID)
	return result, nil
}

func (a *Activities) ShipOrder(ctx context.Context, order model.Order) (model.ShipmentResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Shipping order", "orderID", order.OrderID)

	// Simulate shipping latency
	time.Sleep(300 * time.Millisecond)

	carriers := []string{"FedEx", "UPS", "USPS", "DHL"}
	carrier := carriers[rand.Intn(len(carriers))]

	result := model.ShipmentResult{
		TrackingNumber: fmt.Sprintf("TRK-%s-%d", order.OrderID, time.Now().UnixMilli()),
		Carrier:        carrier,
	}

	logger.Info("Order shipped", "orderID", order.OrderID, "tracking", result.TrackingNumber, "carrier", result.Carrier)
	return result, nil
}

func (a *Activities) SendNotification(ctx context.Context, order model.Order, result model.OrderResult) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending notification",
		"orderID", order.OrderID,
		"email", order.Email,
		"tracking", result.Shipment.TrackingNumber,
	)

	// Simulate sending email
	time.Sleep(100 * time.Millisecond)

	logger.Info("Notification sent", "orderID", order.OrderID, "email", order.Email)
	return nil
}
