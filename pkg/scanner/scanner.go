// Package scan provides functionality for scanning Docker images.
package scanner

import (
	"log"

	"github.com/carbonetes/diggity/internal/slsa"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser"
	"github.com/carbonetes/diggity/pkg/parser/bom"
	"github.com/carbonetes/diggity/pkg/parser/util"
)

// Diggity scans the Docker images, Tar Files, and Codebases(directories) specified in the given model.Arguments struct and returns a sbom(model.Result) struct.
func Scan(arguments *model.Arguments) (*model.SBOM, *[]error) {
	requirements, err := bom.InitParsers(arguments)
	if err != nil {
		log.Fatal(err)
	}
	requirements.WG.Add(len(parser.FindFunctions))
	for _, parser := range parser.FindFunctions {
		go parser(requirements)
	}
	requirements.WG.Wait()
	defer util.CleanUp(requirements)

	if *arguments.Provenance != "" {
		requirements.SBOM.SLSA = slsa.Provenance(requirements)
	}

	return requirements.SBOM, requirements.Errors
}
