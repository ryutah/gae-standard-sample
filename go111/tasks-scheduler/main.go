package main

import (
	"fmt"
	"net/http"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

func main() {
	http.HandleFunc("/cron", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("=======================CRON=======================")
		for k, v := range r.Header {
			fmt.Printf("%v: %v\n", k, v)
		}
		fmt.Println("==================================================")
	})
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("=======================TASKS=======================")
		for k, v := range r.Header {
			fmt.Printf("%v: %v\n", k, v)
		}
		fmt.Println("===================================================")
	})
	http.HandleFunc("/tasks/insert", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := cloudtasks.NewClient(ctx)
		if err != nil {
			panic(err)
		}
		defer client.Close()

		project := os.Getenv("GOOGLE_CLOUD_PROJECT")
		task, err := client.CreateTask(ctx, &taskspb.CreateTaskRequest{
			Parent: fmt.Sprintf("projects/%s/locations/asia-northeast1/queues/example-queue", project),
			Task: &taskspb.Task{
				PayloadType: &taskspb.Task_AppEngineHttpRequest{
					AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
						AppEngineRouting: &taskspb.AppEngineRouting{Service: "tasks"},
						RelativeUri:      "/tasks",
						HttpMethod:       taskspb.HttpMethod_GET,
					},
				},
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Task: ", task.GetName())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
