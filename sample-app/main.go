package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"learn-temporal/sample-app/activity"
	"learn-temporal/sample-app/workflow"
)

const taskQueue = "order-processing-queue"

func main() {
	numWorkers := 1
	if v := os.Getenv("NUM_WORKERS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			log.Fatalf("Invalid NUM_WORKERS value: %s", v)
		}
		numWorkers = n
	}

	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalf("Failed to connect to Temporal: %v", err)
	}
	defer c.Close()

	workers := make([]worker.Worker, numWorkers)
	for i := range workers {
		w := worker.New(c, taskQueue, worker.Options{})
		w.RegisterWorkflow(workflow.OrderWorkflow)
		w.RegisterActivity(&activity.Activities{})
		workers[i] = w
	}

	// Graceful shutdown on SIGINT/SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down workers...")
		for _, w := range workers {
			w.Stop()
		}
	}()

	log.Printf("Starting %d worker(s) on task queue: %s", numWorkers, taskQueue)

	var wg sync.WaitGroup
	for i, w := range workers {
		wg.Add(1)
		go func(id int, w worker.Worker) {
			defer wg.Done()
			if err := w.Run(worker.InterruptCh()); err != nil {
				log.Printf("Worker %d failed: %v", id, err)
			}
		}(i, w)
	}
	wg.Wait()
}