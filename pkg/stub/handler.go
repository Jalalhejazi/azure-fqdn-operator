package stub

import (
	"context"
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/api/core/v1"
)

// NewHandler function
func NewHandler() sdk.Handler {
	return &Handler{}
}

// Handler struct
type Handler struct {
}

// Handle function - reacts to listener
func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1.Service:

		if event.Deleted {
			fmt.Printf("Service Deleted - do nothing yet")
		} else {
			s := o

			// Get annotation
			annotations := s.Annotations

			// If annotation is present, watch service for PIP
			if _, ok := annotations["azure-fqdn-kill"]; !ok {
				if _, ok := annotations["azure-fqdn-value"]; ok {
					readyService(s)
				}
			}
		}
	}
	return nil
}
