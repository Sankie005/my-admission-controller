package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	maxMemory = "512Mi" // Example memory limit
	maxCPU    = "500m"  // Example CPU limit
)

// AdmissionController implements the admission.Handler interface
type AdmissionController struct {
	decoder *admission.Decoder
}

// Handle implements the admission.Handler interface
func (a *AdmissionController) Handle(ctx context.Context, req admission.Request) admission.Response {
	// Decode the admission request
	cronJob := &batchv1.CronJob{}
	if err := a.decoder.Decode(req, cronJob); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// Check the resource limits
	for _, container := range cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers {
		if container.Resources.Requests.Memory().Cmp(resource.MustParse(maxMemory)) > 0 ||
			container.Resources.Requests.Cpu().Cmp(resource.MustParse(maxCPU)) > 0 {
			return admission.Denied("CronJob exceeds resource limits")
		}
	}

	fmt.Printf("Handling admission request for: %s\n", req.Name)
	return admission.Allowed("Resource limits are within acceptable range")
}

// ServeHTTP implements the http.Handler interface
func (a *AdmissionController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// Verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	admissionReview := admissionv1.AdmissionReview{}
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		http.Error(w, fmt.Sprintf("could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	response := a.Handle(context.Background(), admission.Request{
		AdmissionRequest: *admissionReview.Request,
	})

	admissionReview.Response = &response.AdmissionResponse
	responseBody, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// InjectDecoder implements the inject.Decoder interface
func (a *AdmissionController) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}

// StartWebhookServer starts the webhook server without TLS
func StartWebhookServer(addr string) error {
	ac := &AdmissionController{}
	http.HandleFunc("/validate", ac.ServeHTTP)
	fmt.Printf("Starting webhook server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
