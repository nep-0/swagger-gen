package main

import (
	_ "embed"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed swagger-ui-bundle.js
var swaggerJs string

//go:embed swagger-ui.css
var swaggerCss string

//go:embed swagger-ui-bundle.js.LICENSE.txt
var swaggerLicense string

// HTML template
const htmlTemplate = `<!--
{{ .License }}

MIT License

Copyright (c) 2026 Jeff

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

-->
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>API Documentation</title>
  <style>
    {{ .CSS }}
  </style>
</head>
<body>
  <div id="your-app-docs"></div>
  <script>
    {{ .JS }}
  </script>
  <script>
    window.onload = function () {
      const ui = SwaggerUIBundle({
        url: "{{ .DataURL }}",
        dom_id: "#your-app-docs",
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        layout: "BaseLayout"
      })
    }
  </script>
</body>
</html>`

type PageData struct {
	CSS     string
	JS      string
	DataURL string
	License string
}

func main() {
	inputPath := flag.String("i", "", "Path to the OpenAPI JSON/YAML file")
	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Usage: swagger-gen -i <input_file> [output_file]")
		os.Exit(1)
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

	// Determine mime type
	ext := strings.ToLower(filepath.Ext(*inputPath))
	var mimeType string
	switch ext {
	case ".json":
		mimeType = "application/json"
	case ".yaml", ".yml":
		mimeType = "application/yaml"
	default:
		// Fallback or guess based on content? Defaulting to yaml as it is common
		mimeType = "application/yaml"
	}

	// Create Data URI
	encodedContent := base64.StdEncoding.EncodeToString(content)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, encodedContent)

	data := PageData{
		CSS:     swaggerCss,
		JS:      swaggerJs,
		DataURL: dataURL,
		License: swaggerLicense,
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
