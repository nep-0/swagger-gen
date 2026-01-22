package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

//go:embed swagger-ui-bundle.js
var swaggerJs string

//go:embed swagger-ui.css
var swaggerCss string

//go:embed swagger-ui-standalone-preset.js
var swaggerPresetJs string

//go:embed swagger-ui-bundle.js.LICENSE.txt
var swaggerLicense string

func main() {
	inputPath := flag.String("i", "", "Path to the OpenAPI JSON/YAML file")
	cdnURL := flag.String("cdn", "", "CDN base URL for Swagger UI assets (e.g., https://cdn.jsdelivr.net/npm/swagger-ui-dist@5)")
	useJSDelivr := flag.Bool("jsdelivr", false, "Use jsdelivr CDN (shortcut for -cdn https://cdn.jsdelivr.net/npm/swagger-ui-dist@5)")
	useUnpkg := flag.Bool("unpkg", false, "Use unpkg CDN (shortcut for -cdn https://unpkg.com/swagger-ui-dist@5)")
	useStaticfile := flag.Bool("staticfile", false, "Use staticfile CDN (shortcut for -cdn https://cdn.staticfile.net/swagger-ui-dist/5)")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Usage: swagger-gen -i <input_file> [output_file]")
		os.Exit(1)
	}

	// Determine CDN URL from flags
	if *useJSDelivr {
		*cdnURL = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5"
	} else if *useUnpkg {
		*cdnURL = "https://unpkg.com/swagger-ui-dist@5"
	} else if *useStaticfile {
		*cdnURL = "https://cdn.staticfile.net/swagger-ui/5.18.2"
	}

	outputPath := "docs.html"
	if args := flag.Args(); len(args) > 0 {
		outputPath = args[0]
	}

	// Read input file
	content, err := os.ReadFile(*inputPath)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Parse the input file and convert to JSON for JavaScript
	ext := strings.ToLower(filepath.Ext(*inputPath))
	var specJSON []byte
	var err2 error

	switch ext {
	case ".json":
		// For JSON files, just validate and use as-is
		var jsonData any
		err2 = json.Unmarshal(content, &jsonData)
		if err2 != nil {
			log.Fatalf("Error parsing JSON file: %v", err2)
		}
		specJSON = content
	case ".yaml", ".yml":
		// For YAML files, convert to JSON while preserving order
		var yamlData yaml.MapSlice
		err2 = yaml.Unmarshal(content, &yamlData)
		if err2 != nil {
			log.Fatalf("Error parsing YAML file: %v", err2)
		}
		// Convert YAML data to JSON-compatible format
		jsonData := convertToJSONType(yamlData)
		specJSON, err2 = orderedJSONMarshal(jsonData)
		if err2 != nil {
			log.Fatalf("Error converting YAML to JSON: %v", err2)
		}
	default:
		// Try YAML first, then JSON
		var yamlData yaml.MapSlice
		err2 = yaml.Unmarshal(content, &yamlData)
		if err2 != nil {
			// Try JSON
			var jsonData any
			err2 = json.Unmarshal(content, &jsonData)
			if err2 != nil {
				log.Fatalf("Error parsing input file as either YAML or JSON: %v", err2)
			}
			specJSON = content
		} else {
			// Successfully parsed as YAML
			jsonData := convertToJSONType(yamlData)
			specJSON, err2 = orderedJSONMarshal(jsonData)
			if err2 != nil {
				log.Fatalf("Error converting YAML to JSON: %v", err2)
			}
		}
	}

	// Extract title from spec
	title := extractTitle(specJSON)

	data := PageData{
		CSS:      swaggerCss,
		JS:       swaggerJs,
		PresetJS: swaggerPresetJs,
		Spec:     string(specJSON),
		License:  swaggerLicense,
		Title:    title,
		CDNURL:   *cdnURL,
	}

	t, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outFile.Close()

	err = t.Execute(outFile, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	fmt.Printf("Documentation generated at %s\n", outputPath)
}
