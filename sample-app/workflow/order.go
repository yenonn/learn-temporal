package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"learn-temporal/sample-app/activity"
	"learn-temporal/sample-app/model"
)

func OrderWorkflow(ctx workflow.Context, order model.Order) (model.OrderResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Order workflow started", "orderID", order.OrderID)

	var activities *activity.Activities
	result := model.OrderResult{
		OrderID: order.OrderID,
		Status:  "processing",
	}

	// Step 1: Validate Order
	validateCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    3,
		},
	})
	if err := workflow.ExecuteActivity(validateCtx, activities.ValidateOrder, order).Get(ctx, nil); err != nil {
		result.Status = "validation_failed"
		return result, err
	}
	logger.Info("Order validated", "orderID", order.OrderID)

	// Step 2: Process Payment
	paymentCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    5,
		},
	})
	if err := workflow.ExecuteActivity(paymentCtx, activities.ProcessPayment, order).Get(ctx, &result.Payment); err != nil {
		result.Status = "payment_failed"
		return result, err
	}
	logger.Info("Payment processed", "orderID", order.OrderID, "transactionID", result.Payment.TransactionID)

	// Step 3: Ship Order
	shipCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    3,
		},
	})
	if err := workflow.ExecuteActivity(shipCtx, activities.ShipOrder, order).Get(ctx, &result.Shipment); err != nil {
		result.Status = "shipping_failed"
		return result, err
	}
	logger.Info("Order shipped", "orderID", order.OrderID, "tracking", result.Shipment.TrackingNumber)

	// Step 4: Send Notification (best-effort — failure does not fail the workflow)
	notifyCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 20 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    2,
		},
	})
	if err := workflow.ExecuteActivity(notifyCtx, activities.SendNotification, order, result).Get(ctx, nil); err != nil {
		logger.Warn("Failed to send notification, continuing anyway", "orderID", order.OrderID, "error", err)
	}

	result.Status = "completed"
	logger.Info("Order workflow completed", "orderID", order.OrderID)
	return result, nil
}
