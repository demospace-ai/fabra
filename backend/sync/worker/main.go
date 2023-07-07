package main

import (
	"log"

	"go.fabra.io/server/common/database"
	"go.fabra.io/sync/temporal"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
)

const WORKER_PEM_KEY = "projects/932264813910/secrets/temporal-worker-pem/versions/latest"
const WORKER_KEY_KEY = "projects/932264813910/secrets/temporal-worker-key/versions/latest"

func main() {
	c, err := temporal.CreateClient(WORKER_PEM_KEY, WORKER_KEY_KEY)
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
		return
	}

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, temporal.SyncTaskQueue, worker.Options{
		// Create interceptor that will unwrap CustomerVisibleError and set it at top level
		Interceptors: []interceptor.WorkerInterceptor{
			temporal.NewErrorInterceptor(),
		},
	})

	activities := &temporal.Activities{
		Db: db,
	}

	w.RegisterActivity(activities)
	w.RegisterWorkflow(temporal.SyncWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
