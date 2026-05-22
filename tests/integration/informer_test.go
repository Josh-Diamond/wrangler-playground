package integration

import (
	"fmt"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/josh-diamond/wrangler-playground/tests/integration/framework"
)

func TestConfigMapInformerReceivesCreate(t *testing.T) {
	name := fmt.Sprintf("integration-test-%d", time.Now().UnixNano())

	framework.WaitForConfigMapEvent(framework.ConfigMapEventInput{
		T:     t,
		Name:  name,
		Event: framework.EventCreate,
		MatchFunc: func(cm *corev1.ConfigMap) bool {
			return cm != nil && cm.Name == name
		},
		Timeout: 10 * time.Second,
	})
}

func TestConfigMapInformerReceivesUpdate(t *testing.T) {
	name := fmt.Sprintf("integration-test-%d", time.Now().UnixNano())

	framework.WaitForConfigMapEvent(framework.ConfigMapEventInput{
		T:     t,
		Name:  name,
		Event: framework.EventUpdate,
		MatchFunc: func(cm *corev1.ConfigMap) bool {
			return cm != nil &&
				cm.Name == name &&
				cm.Annotations["updated"] == "true"
		},
		Mutate: func(cm *corev1.ConfigMap) {
			if cm.Annotations == nil {
				cm.Annotations = map[string]string{}
			}
			cm.Annotations["updated"] = "true"
		},
		Timeout: 10 * time.Second,
	})
}

func TestConfigMapInformerReceivesDelete(t *testing.T) {
	name := fmt.Sprintf("integration-test-%d", time.Now().UnixNano())

	framework.WaitForConfigMapEvent(framework.ConfigMapEventInput{
		T:     t,
		Name:  name,
		Event: framework.EventDelete,
		MatchFunc: func(cm *corev1.ConfigMap) bool {
			return cm == nil
		},
		Timeout: 10 * time.Second,
	})
}
