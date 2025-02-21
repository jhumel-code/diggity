package attestation

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/carbonetes/diggity/internal/cli"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/ui"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/google/uuid"
)

var (
	log        = logger.GetLogger()
	cosign     = "cosign"
	sbomPrefix = "diggity-sbom-"

	// Arguments sbom args
	Arguments model.Arguments
)

// Attest runs SBOM Attestation
func Attest(image string, attestationOptions *model.AttestationOptions) {
	var predicate string

	// Verify if cosign is installed
	err := checkCosign()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			log.Error("Unable to run cosign. Make sure it is installed first.")
		} else {
			log.Error(err)
		}
		return
	}

	// Generate SBOM as needed
	if *attestationOptions.Predicate == "" {
		predicate = generateBom(image, attestationOptions.BomArgs, *attestationOptions.OutputType, *attestationOptions.Provenance)
	} else {
		predicate = *attestationOptions.Predicate
	}

	// Attest specified BOM file
	ui.OnSbomAttestation()
	err = attestBom(image, predicate, attestationOptions)
	if err != nil {
		log.Error("Error occurred when running SBOM attestation. Please make sure that the paths or fields specified are correct.")
		ui.DoneSpinner()
		return
	}
	ui.DoneSpinner()

	// Get Attestation
	ui.OnVerifyingAttestation()
	err = getAttestation(image, attestationOptions)
	if err != nil {
		log.Error("Error occurred when verifying attestation. Please make sure that the paths or fields specified are correct.")
		ui.DoneSpinner()
		return
	}
	ui.DoneSpinner()
}

// Check if cosign is installed on machine
func checkCosign() error {
	cmd := exec.Command("cosign")
	return cmd.Run()
}

// Attest SBOM
func attestBom(image string, predicate string, attestationOptions *model.AttestationOptions) error {
	args := fmt.Sprintf("attest --yes --key %+v --type %+v --predicate %+v %+v",
		*attestationOptions.Key, *attestationOptions.AttestType, predicate, image)
	attest := strings.Split(args, " ")

	cmd := exec.Command(cosign, attest...)
	cmd.Stdin = strings.NewReader(*attestationOptions.Password)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Get Attestation by Verifying
func getAttestation(image string, attestationOptions *model.AttestationOptions) error {
	args := fmt.Sprintf("verify-attestation --key %+v --type %+v %+v", *attestationOptions.Pub, *attestationOptions.AttestType, image)
	verify := strings.Split(args, " ")

	if *attestationOptions.OutputFile != "" {
		fileArg := fmt.Sprintf("--output-file %+v", *attestationOptions.OutputFile)
		verify = append(verify, strings.Split(fileArg, " ")...)
	}

	cmd := exec.Command(cosign, verify...)
	if *attestationOptions.OutputFile == "" {
		cmd.Stdout = os.Stdout
	}
	return cmd.Run()
}

// Generate SBOM
func generateBom(image string, arguments *model.Arguments, outputType string, provenance string) string {
	// Generate Temp Bom Filename
	var bomFileName string
	switch outputType {
	case model.JSON.ToOutput(), model.CycloneDXJSON, model.SPDXJSON, model.GithubJSON:
		bomFileName = sbomPrefix + uuid.NewString() + ".json"
	case model.CycloneDXXML:
		bomFileName = sbomPrefix + uuid.NewString() + ".cdx"
	case model.SPDXTagValue:
		bomFileName = sbomPrefix + uuid.NewString() + ".spdx"
	default:
		bomFileName = sbomPrefix + uuid.NewString() + ".json"
	}

	// Init Args
	bomPath := filepath.Join(".", bomFileName)
	Arguments = *arguments
	Arguments.Image = &image
	Arguments.OutputFile = &bomPath
	Arguments.Provenance = &provenance
	Arguments.Output = (*model.Output)(&outputType)

	// Start SBOM
	cli.Start(&Arguments)

	return bomPath
}
