# Order Processing Sample App

A Go application demonstrating Temporal workflow orchestration with a realistic order processing pipeline.

[image]{./image/temporal-example.png}

## Workflow Steps

1. **ValidateOrder** — checks items, customer info, and amounts
2. **ProcessPayment** — simulates charging (10% random transient failure to demo retries)
3. **ShipOrder** — generates tracking number, selects carrier
4. **SendNotification** — sends confirmation email (best-effort, won't fail the workflow)

## Prerequisites

- Go 1.21+
- Temporal server running on `localhost:7233` (use `docker compose up -d` from the project root)

## Usage

**Start the worker:**

```bash
cd sample-app
go run main.go
```

**Start a workflow (in another terminal):**

```bash
cd sample-app/starter
go run main.go
```

With custom flags:

```bash
go run main.go --order-id=order-042 --customer-id=cust-456 --email=jane@example.com --amount=149.99
```

**View in Temporal UI:** <http://localhost:8233>

## Project Structure

```
sample-app/
├── main.go              # Worker: connects to Temporal, registers workflow + activities
├── starter/main.go      # CLI to start a workflow with sample order data
├── workflow/order.go     # Workflow definition: orchestrates 4 sequential steps
├── activity/order.go     # Activity implementations (simulated business logic)
└── model/order.go        # Shared structs: Order, Item, Address, result types
```
