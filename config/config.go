package config

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Namespace string
	*kubernetes.Clientset
}

type ConfigOpts func(*Config)

func WithKubeconfig(kubeconfig string) ConfigOpts {
	return func(c *Config) {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			os.Exit(1)
		}
		loadedConfig, err := clientcmd.LoadFromFile(kubeconfig)
		if err != nil {
			os.Exit(1)
		}

		c.Namespace = loadedConfig.Contexts[loadedConfig.CurrentContext].Namespace

		cl, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			os.Exit(1)
		}
		c.Clientset = cl
	}
}

func NewConfig(opts ...ConfigOpts) *Config {
	c := &Config{}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
