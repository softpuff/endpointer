package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"softpuff/endpointer/config"
	"softpuff/endpointer/utils"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	kubeconfig    string
	namespace     string
	ctx           = context.Background()
	appName       string
	verbose       bool
	allNamespaces bool
)

func main() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig path")
	flag.StringVar(&namespace, "namespace", "", "namespace")
	flag.StringVar(&appName, "app", "", "app name")
	flag.BoolVar(&allNamespaces, "all-namespaces", false, "all namespaces")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.Parse()

	var apps []utils.App

	var ok bool
	if kubeconfig == "" {
		kubeconfig, ok = os.LookupEnv("KUBECONFIG")
		if !ok {
			fmt.Fprint(os.Stderr, "Kubeconfig not given")
			os.Exit(1)
		}
	}

	cs := config.NewConfig(config.WithKubeconfig(kubeconfig))

	if namespace == "" {
		namespace = cs.Namespace
	}

	if allNamespaces {
		namespace = ""
	}

	var lo v1.ListOptions

	if appName != "" {
		lo.LabelSelector = fmt.Sprintf("app=%s", appName)
	}

	endpoints, err := cs.CoreV1().Endpoints(namespace).List(ctx, lo)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, e := range endpoints.Items {

		if verbose {
			fmt.Printf("%s/%s\n", e.Namespace, e.Name)
		}

		var addr []string
		var badAddr []string

		for _, a := range e.Subsets {
			for _, ip := range a.Addresses {
				addr = append(addr, ip.IP)
			}
			for _, bIp := range a.NotReadyAddresses {
				badAddr = append(badAddr, bIp.IP)
			}
		}

		app := utils.NewApp(addr, badAddr, e.Name, e.Namespace, len(addr)+len(badAddr))
		apps = append(apps, app)
	}

	utils.PrintApps(apps)
}
