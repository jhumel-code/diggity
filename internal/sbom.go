package sbom

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/carbonetes/diggity/internal/docker"
	"github.com/carbonetes/diggity/internal/file"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/output"
	"github.com/carbonetes/diggity/internal/ui"
	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser"
	"github.com/carbonetes/diggity/pkg/parser/stream"
	tm "github.com/vmware/transport-go/model"

	"os"

	"github.com/schollz/progressbar/v3"
)

const (
	tarFile      string = "tar"
	image        string = "image"
	dir          string = "dir"
	unknown      string = "Unknown"
	defaultTag   string = "latest"
	tagSeparator string = ":"
)

var (
	log = logger.GetLogger()
)

// Start SBOM extraction
func Start(arguments *model.Arguments) {
	start := time.Now()
	if *arguments.Quiet {
		log = logger.SetQuietMode(log)
	}
	//check image if DIR
	source, spinnerMsg := file.CheckUserInput(arguments)
	if source == image && !strings.Contains(*arguments.Image, tagSeparator) {
		log.Printf("Using default tag:" + defaultTag)
	}

	extractSpinner := ui.InitSpinner(spinnerMsg)

	if *arguments.Stream {
		handler, err := stream.Subscribe(stream.PackageChannel)
		if err != nil {
			log.Error(err)
		}

		handler.Handle(func(m *tm.Message) {
			pkg, _ := json.MarshalIndent(m.Payload, "", " ")
			log.Printf("%s\n", pkg)
		}, func(err error) {
			log.Error(err)
		})

		extractImage(source, arguments, extractSpinner)
		parser.Start(arguments)
		log.Printf("\nScanning finished in %.2fs", time.Since(start).Seconds())
		os.Exit(0)
	}
	//Extract Image
	if !*arguments.Quiet {
		// Pull (if needed) and Extract Image
		spinnerMsg = extractImage(source, arguments, extractSpinner)

		// Run Parsers
		parseSpinner := ui.InitSpinner(spinnerMsg)
		go ui.RunSpinner(parseSpinner)
		parser.Start(arguments)
		ui.DoneSpinner(parseSpinner)
	} else {
		extractImage(source, arguments, extractSpinner)
		parser.Start(arguments)
	}

	//Print Results and Cleanup
	output.PrintResults()
	log.Printf("\nScanning finished in %.2fs", time.Since(start).Seconds())
}

// Extract image
func extractImage(source string, arguments *model.Arguments, spinner *progressbar.ProgressBar) string {
	switch source {
	case tarFile:
		dir := *arguments.Tar
		if file.Exists(dir) {
			docker.ExtractFromDir(arguments.Tar)
			return "Parsing Tar file..."
		}
		log.Printf("%s not found\n", *arguments.Tar)
		os.Exit(0)
	case image:
		if !strings.Contains(*arguments.Image, tagSeparator) {
			modifiedTag := *arguments.Image + tagSeparator + defaultTag
			arguments.Image = &modifiedTag
		}
		docker.ExtractImage(arguments, spinner)
		return "Parsing Image..."
	case dir:
		dir := *arguments.Dir
		if file.Exists(dir) {
			err := file.GetFilesFromDir(dir)
			if err != nil {
				panic(err)
			}

			docker.CreateTempDir()
			return "Parsing Directory..."
		}
		log.Printf("%s not found\n", *arguments.Dir)
		os.Exit(0)
	}
	return ""
}
