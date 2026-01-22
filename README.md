# swagger-gen

Generate single-file HTML documentation from OpenAPI v3 specifications.

## Installation

Download the latest binary for your platform from [Releases](../../releases).

## Usage

```bash
swagger-gen -i openapi.yaml [output.html]
```

If no output file is specified, defaults to `docs.html`.

## CDN Options

Use CDN for smaller output files (requires internet access):

```bash
swagger-gen -i openapi.yaml -jsdelivr docs.html    # jsdelivr
swagger-gen -i openapi.yaml -unpkg docs.html       # unpkg
swagger-gen -i openapi.yaml -staticfile docs.html  # staticfile (China)
swagger-gen -i openapi.yaml -cdn <URL> docs.html   # custom CDN
```

Defaults to embedded mode (self-contained, works offline).

## License

MIT © 2026 Jeff
