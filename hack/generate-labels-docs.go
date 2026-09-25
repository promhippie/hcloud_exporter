//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/promhippie/hcloud_exporter/pkg/config"
)

type section struct {
	Name   string
	Labels []string
}

func main() {
	f, err := os.Create("docs/partials/labels.md")

	if err != nil {
		fmt.Printf("failed to create file")
		os.Exit(1)
	}

	defer f.Close()

	sections := []section{
		{Name: "Floating IP", Labels: config.FloatingIPLabels()},
		{Name: "Image", Labels: config.ImageLabels()},
		{Name: "Server", Labels: config.ServerLabels()},
		{Name: "Server Metrics", Labels: config.ServerMetricsLabels()},
		{Name: "Load Balancer", Labels: config.LoadBalancerLabels()},
		{Name: "SSH Key", Labels: config.SSHKeyLabels()},
		{Name: "Volume", Labels: config.VolumeLabels()},
		{Name: "StorageBox", Labels: config.StorageBoxLabels()},
	}

	last := sections[len(sections)-1]
	for _, row := range sections {
		f.WriteString(fmt.Sprintf(
			"### %s Labels\n\n",
			row.Name,
		))

		for _, label := range row.Labels {
			f.WriteString(fmt.Sprintf(
				"* %s\n",
				label,
			))
		}

		if row.Name != last.Name {
			f.WriteString("\n")
		}
	}
}
