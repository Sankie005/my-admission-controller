package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sankie005/my-admission-controller/pkg/webhook"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Create a new AdmissionController
	wh := &webhook.AdmissionController{}

	// Set up the HTTP server
	http.HandleFunc("/validate-jobs", wh.ServeHTTP)

	// Start the webhook server
	fmt.Println("Starting webhook server on :443")
	if err := http.ListenAndServe(":443", nil); err != nil {
		log.Fatalf("Error starting webhook server: %v", err)
		os.Exit(1)
	}
}
