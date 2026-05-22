package framework

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corecontrollers "github.com/rancher/wrangler/pkg/generated/controllers/core"
	"github.com/stretchr/testify/require"
)

// EventType defines what action to perform after create
type EventType string

const (
	EventCreate EventType = "create"
	EventUpdate EventType = "update"
	EventDelete EventType = "delete"
)

type ConfigMapEventInput struct {
	T         *testing.T
	Name      string
	Event     EventType
	MatchFunc func(*corev1.ConfigMap) bool
	Mutate    func(*corev1.ConfigMap)
	Timeout   time.Duration
}

func WaitForConfigMapEvent(in ConfigMapEventInput) {
	var observed atomic.Bool

	ctx := context.Background()

	wranglerCtx, err := StartWrangler(func(factory *corecontrollers.Factory) {
		configmaps := factory.Core().V1().ConfigMap()

		configmaps.OnChange(
			ctx,
			"configmap-event-test",
			func(key string, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
				if in.MatchFunc(cm) {
					observed.Store(true)
				}
				return cm, nil
			},
		)
	})

	require.NoError(in.T, err)
	defer wranglerCtx.Cancel()

	client := wranglerCtx.Clientset.CoreV1().ConfigMaps(Namespace)

	// CREATE
	cm, err := client.Create(wranglerCtx.Context, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: in.Name,
		},
	}, metav1.CreateOptions{})
	require.NoError(in.T, err)

	// OPTIONAL MUTATION (update/delete)
	if in.Event == EventUpdate && in.Mutate != nil {
		cm = cm.DeepCopy()
		in.Mutate(cm)

		_, err = client.Update(wranglerCtx.Context, cm, metav1.UpdateOptions{})
		require.NoError(in.T, err)
	}

	if in.Event == EventDelete {
		err = client.Delete(wranglerCtx.Context, cm.Name, metav1.DeleteOptions{})
		require.NoError(in.T, err)
	}

	require.Eventually(in.T,
		func() bool { return observed.Load() },
		in.Timeout,
		500*time.Millisecond,
		"expected ConfigMap event",
	)
}
