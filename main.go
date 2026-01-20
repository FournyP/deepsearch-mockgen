package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/FournyP/deepsearch-mockgen/generator"
	"github.com/FournyP/deepsearch-mockgen/tui"
)

func main() {
	// Define CLI flags
	searchDir := flag.String("search", "", "Directory to search for interfaces")
	outputDir := flag.String("output", "", "Directory to save generated mocks")

	var acceptAll bool
	flag.BoolVar(&acceptAll, "A", false, "Generate mocks for all interfaces without prompting")
	flag.BoolVar(&acceptAll, "all", false, "Generate mocks for all interfaces without prompting")

	// Parse flags
	flag.Parse()

	// Prompt for missing values
	if *searchDir == "" {
		*searchDir = tui.PromptInput("Enter the search directory:")
	}

	if *outputDir == "" {
		*outputDir = tui.PromptInput("Enter the output directory:")
	}

	// Find interfaces
	interfaces, err := generator.FindInterfaces(*searchDir)
	if err != nil {
		log.Fatal(err)
	}

	if len(interfaces) == 0 {
		log.Println("No interfaces found")
		return
	}

	// Prompt user for each interface (or accept all when `-A/--all` provided)
	finalPaths := make(map[string]string)
	for iface, ifacePath := range interfaces {
		if !acceptAll {
			generate := tui.PromptYesNoWithDefaultValue(fmt.Sprintf("Generate mock for %s?:", iface), true)

			if !generate {
				continue
			}
		}

		// Compute default mock path
		defaultMockPath := generator.ComputeMockPath(*searchDir, *outputDir, ifacePath, iface)

		// Ask the user to confirm or modify the path
		mockPath := tui.PromptInputWithDefault(
			fmt.Sprintf("Mock path for %s:", iface),
			defaultMockPath,
		)

		finalPaths[iface] = mockPath
	}

	// Generate mocks
	for iface, mockPath := range finalPaths {
		err := generator.GenerateMock(iface, interfaces[iface], mockPath)
		if err != nil {
			log.Printf("Failed to generate mock for %s: %v\n", iface, err)
		} else {
			log.Printf("Mock generated for %s at %s\n", iface, mockPath)
		}
	}
}
