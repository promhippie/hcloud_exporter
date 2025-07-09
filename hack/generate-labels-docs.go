//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/promhippie/hcloud_exporter/pkg/config"
)

func main() {
	f, err := os.Create("docs/partials/labels.md")

	if err != nil {
		fmt.Printf("failed to create file")
		os.Exit(1)
	}

	defer f.Close()

	f.WriteString("### StorageBox Labels\n\n")
	for _, row := range config.StorageBoxLabels() {
		f.WriteString(fmt.Sprintf(
			"* %s\n",
			row,
		))
	}
}
