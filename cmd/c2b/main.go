package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/transpiler"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/transpiler/translators"
)

func main() {
	inputFile := flag.String("i", "", "Path to the input cloud-config.yaml file")
	outputFile := flag.String("o", "", "Path to output Butane YAML config")

	flag.Parse()

	if *inputFile == "" {
		log.Fatalf("path to input file not provided")
	}

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	cloudCfg, err := cloudinit.Parse(data)
	if err != nil {
		log.Fatalf("Error parsing cloud-init YAML: %v", err)
	}

	engine := transpiler.NewEngine(
		translators.NewBootCmdTranslator(),
		translators.NewUserTranslator(),
		translators.NewFileTranslator(),
		translators.NewRunCmdTranslator(),
		translators.NewNTPTranslator(),
		translators.NewCACertsTranslator(),
	)

	butaneCfg, err := engine.Run(cloudCfg)
	if err != nil {
		log.Fatalf("Transpilation failed: %v", err)
	}

	outData, err := butane.Generate(butaneCfg)
	if err != nil {
		log.Fatalf("Error generating Butane YAML: %v", err)
	}

	if *outputFile != "" {
		err := os.WriteFile(*outputFile, outData, 0644)
		if err != nil {
			log.Fatalf("failed to write YAML: %v", err)
		}
		fmt.Printf("Success: Wrote Butane configuration to %s\n", *outputFile)
	} else {
		if _, err := os.Stdout.Write(outData); err != nil {
			log.Fatalf("Failed to write to stdout: %v", err)
		}
	}
}
