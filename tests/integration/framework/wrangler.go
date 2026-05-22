package framework

import (
	"context"

	corecontrollers "github.com/rancher/wrangler/pkg/generated/controllers/core"
	"k8s.io/client-go/kubernetes"
)

type WranglerContext struct {
	Context     context.Context
	Cancel      context.CancelFunc
	CoreFactory *corecontrollers.Factory
	Clientset   *kubernetes.Clientset
}

func StartWrangler(setupControllers func(*corecontrollers.Factory)) (*WranglerContext, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	coreFactory, err := corecontrollers.NewFactoryFromConfig(cfg)
	if err != nil {
		cancel()
		return nil, err
	}

	if setupControllers != nil {
		setupControllers(coreFactory)
	}

	if err := coreFactory.Start(ctx, 5); err != nil {
		cancel()
		return nil, err
	}

	if err := coreFactory.Sync(ctx); err != nil {
		cancel()
		return nil, err
	}

	return &WranglerContext{
		Context:     ctx,
		Cancel:      cancel,
		CoreFactory: coreFactory,
		Clientset:   clientset,
	}, nil
}
