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
		translators.NewUserTranslator(),
		translators.NewFileTranslator(),
		translators.NewRunCmdTranslator(),
	)

	butaneCfg, err := engine.Run(cloudCfg)
	if err != nil {
		log.Fatalf("Transpilation failed: %v", err)
	}

	outData, err := butane.Generate(butaneCfg)
	if err != nil {
		log.Fatalf("Error generating Butane YAML: %v", err)
	}

	fmt.Println(string(outData))
}
