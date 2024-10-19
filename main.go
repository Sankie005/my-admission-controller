package main

import (
	"fmt"
	"os"

	"github.com/sankie005/my-admission-controller/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	webhookruntime "sigs.k8s.io/controller-runtime/pkg/webhook"
)

func main() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Error getting config: %v\n", err)
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		fmt.Printf("Error creating manager: %v\n", err)
		os.Exit(1)
	}

	// Setup webhooks
	fmt.Println("Setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	fmt.Println("Registering webhooks to the webhook server")
	hookServer.Register("/validate", &webhookruntime.Admission{Handler: &webhook.AdmissionController{}})

	fmt.Println("Starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		fmt.Printf("Error starting manager: %v\n", err)
		os.Exit(1)
	}
}
