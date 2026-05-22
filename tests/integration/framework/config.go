package framework

import "flag"

var (
	Kubeconfig string
	Namespace  string
)

func init() {
	flag.StringVar(
		&Kubeconfig,
		"kubeconfig",
		"",
		"Path to kubeconfig",
	)

	flag.StringVar(
		&Namespace,
		"namespace",
		"wrangler-tests",
		"Test namespace",
	)
}