package utils

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

type App struct {
	Endpoints      []string `json:"endpoints"`
	EndpointsCount int      `json:"endpointsCount"`
	Name           string   `json:"name"`
	Replicas       int      `json:"replicas"`
	Unready        []string `json:"unready"`
	UnreadyCount   int      `json:"unreadyCount"`
	Namespace      string   `json:"namespace"`
}

func NewApp(ep, un []string, name, namespace string, r int) App {
	return App{
		Endpoints:      ep,
		Name:           name,
		Replicas:       r,
		EndpointsCount: len(ep),
		Unready:        un,
		UnreadyCount:   len(un),
		Namespace:      namespace,
	}
}

func PrintApps(apps []App) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	format := "%s\t%s\t%s\t%s\t%s\t%s\n"
	fmt.Fprintf(tw, format, "Namespace", "Name", "EndpointsCount", "Replicas", "Unready", "UnreadyCount")

	for _, a := range apps {
		fmt.Fprintf(tw, format, a.Namespace, a.Name, strconv.Itoa(a.EndpointsCount), strconv.Itoa(a.Replicas), a.Unready, strconv.Itoa(a.UnreadyCount))
	}
	tw.Flush()
}
