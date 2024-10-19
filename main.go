package main

import (
	"os"

	"github.com/sankie005/my-admission-controller/pkg/webhook"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Port: 8080, // Using port 8080 for non-TLS
	})
	if err != nil {
		os.Exit(1)
	}

	wh := &webhook.AdmissionController{}
	mgr.GetWebhookServer().Register("/validate-jobs", wh)

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		os.Exit(1)
	}
}
