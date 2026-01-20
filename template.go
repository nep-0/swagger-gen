package main

import "encoding/json"

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
  <title>{{ .Title }}</title>
  <style>
    {{ .CSS }}
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script>
    {{ .JS }}
  </script>
  <script>
    {{ .PresetJS }}
  </script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        spec: {{ .Spec }},
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`

// PageData holds the template data
type PageData struct {
	CSS      string
	JS       string
	PresetJS string
	Spec     string
	License  string
	Title    string
}

// extractTitle extracts the title from the spec JSON
func extractTitle(specJSON []byte) string {
	var spec map[string]any
	if err := json.Unmarshal(specJSON, &spec); err != nil {
		return "API Documentation"
	}

	if info, ok := spec["info"].(map[string]any); ok {
		if title, ok := info["title"].(string); ok && title != "" {
			return title
		}
	}

	return "API Documentation"
}
