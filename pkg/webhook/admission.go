package webhook

import (
	"context"
	"encoding/json"
	"fmt"
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
	admissionReview := admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := a.Handle(context.Background(), admission.Request{
		AdmissionRequest: *admissionReview.Request,
	})

	admissionReview.Response = &response.AdmissionResponse
	err = json.NewEncoder(w).Encode(admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// InjectDecoder implements the inject.Decoder interface
func (a *AdmissionController) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
