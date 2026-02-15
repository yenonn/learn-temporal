package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"learn-temporal/sample-app/model"
	"learn-temporal/sample-app/workflow"
)

const taskQueue = "order-processing-queue"

func main() {
	orderID := flag.String("order-id", "order-001", "Order ID")
	customerID := flag.String("customer-id", "cust-123", "Customer ID")
	email := flag.String("email", "customer@example.com", "Customer email")
	amount := flag.Float64("amount", 99.99, "Total order amount")
	flag.Parse()

	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalf("Failed to connect to Temporal: %v", err)
	}
	defer c.Close()

	order := model.Order{
		OrderID:    *orderID,
		CustomerID: *customerID,
		Email:      *email,
		Items: []model.Item{
			{SKU: "ITEM-001", Name: "Wireless Mouse", Quantity: 1, Price: 29.99},
			{SKU: "ITEM-002", Name: "USB-C Cable", Quantity: 2, Price: 14.99},
			{SKU: "ITEM-003", Name: "Laptop Stand", Quantity: 1, Price: 40.02},
		},
		Address: model.Address{
			Street:  "123 Main St",
			City:    "San Francisco",
			State:   "CA",
			ZipCode: "94105",
			Country: "US",
		},
		TotalAmount: *amount,
	}

	workflowID := fmt.Sprintf("order-%s", order.OrderID)

	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQueue,
	}, workflow.OrderWorkflow, order)
	if err != nil {
		log.Fatalf("Failed to start workflow: %v", err)
	}

	log.Printf("Started workflow: WorkflowID=%s, RunID=%s", we.GetID(), we.GetRunID())
	log.Printf("View in Temporal UI: http://localhost:8233/namespaces/default/workflows/%s/%s", we.GetID(), we.GetRunID())

	var result model.OrderResult
	if err := we.Get(context.Background(), &result); err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	log.Printf("Workflow completed!")
	log.Printf("  Order ID:       %s", result.OrderID)
	log.Printf("  Status:         %s", result.Status)
	log.Printf("  Transaction ID: %s", result.Payment.TransactionID)
	log.Printf("  Tracking:       %s", result.Shipment.TrackingNumber)
	log.Printf("  Carrier:        %s", result.Shipment.Carrier)
}
