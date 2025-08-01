package parser

import (
	"LocalDex/types"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

func ParseMetadata(path string, d *[]byte) (string, error) {
	var err error

	// Always load static config to get default (*) and #not_found routes
	staticPath := filepath.Join(os.Getenv("ETC_DIR"), "config", "static.route.json")
	staticData, err := os.ReadFile(staticPath)
	if err != nil {
		return "", fmt.Errorf("failed to read static metadata file: %w", err)
	}

	var staticRoutes []types.RouteMetadata
	if err := json.Unmarshal(staticData, &staticRoutes); err != nil {
		return "", fmt.Errorf("failed to parse static metadata JSON: %w", err)
	}

	var defaultMeta, notFoundMeta *types.RouteMetadata

	for i := range staticRoutes {
		switch staticRoutes[i].Path {
		case "*":
			defaultMeta = &staticRoutes[i]
		case "#not_found":
			notFoundMeta = &staticRoutes[i]
		}
	}

	// Load dynamic metadata if provided
	var routes []types.RouteMetadata

	if d != nil && *d != nil && len(*d) != 0 {
		if err := json.Unmarshal(*d, &routes); err != nil {
			return "", fmt.Errorf("failed to parse dynamic JSON: %w", err)
		}
	} else {
		// If no dynamic metadata passed, use the static file again
		routes = staticRoutes
	}

	// Find the route metadata for this path
	var currentMeta *types.RouteMetadata
	for i := range routes {
		if routes[i].Path == path {
			currentMeta = &routes[i]
			break
		}
	}

	if currentMeta == nil {
		currentMeta = notFoundMeta
	}

	if currentMeta == nil {
		return "", nil // No metadata to render
	}

	// Merge default + current
	var mergedMeta = types.RouteMetadata{
		Path:  currentMeta.Path,
		Title: currentMeta.Title,
	}

	if defaultMeta != nil {
		mergedMeta.Meta = append(defaultMeta.Meta, currentMeta.Meta...)
		mergedMeta.Link = append(defaultMeta.Link, currentMeta.Link...)
	} else {
		mergedMeta.Meta = currentMeta.Meta
		mergedMeta.Link = currentMeta.Link
	}

	// Render HTML
	var builder strings.Builder

	if mergedMeta.Title != "" {
		builder.WriteString(fmt.Sprintf(`<title id="__SERVER_PROPS__">%s</title>`+"\n", html.EscapeString(mergedMeta.Title)))
	}

	for _, meta := range mergedMeta.Meta {
		switch {
		case meta.Charset != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" charset="%s">`+"\n", html.EscapeString(meta.Charset)))
		case meta.Name != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" name="%s" content="%s">`+"\n", html.EscapeString(meta.Name), html.EscapeString(meta.Content)))
		case meta.Property != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" property="%s" content="%s">`+"\n", html.EscapeString(meta.Property), html.EscapeString(meta.Content)))
		case meta.HTTPEquiv != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" http-equiv="%s" content="%s">`+"\n", html.EscapeString(meta.HTTPEquiv), html.EscapeString(meta.Content)))
		}
	}

	for _, link := range mergedMeta.Link {
		builder.WriteString(`<link id="__SERVER_PROPS__"`)
		builder.WriteString(fmt.Sprintf(` rel="%s"`, html.EscapeString(link.Rel)))
		builder.WriteString(fmt.Sprintf(` href="%s"`, html.EscapeString(link.Href)))

		if link.Type != "" {
			builder.WriteString(fmt.Sprintf(` type="%s"`, html.EscapeString(link.Type)))
		}
		if link.Crossorigin != "" {
			builder.WriteString(fmt.Sprintf(` crossorigin="%s"`, html.EscapeString(link.Crossorigin)))
		}
		if link.Media != "" {
			builder.WriteString(fmt.Sprintf(` media="%s"`, html.EscapeString(link.Media)))
		}
		if link.Sizes != "" {
			builder.WriteString(fmt.Sprintf(` sizes="%s"`, html.EscapeString(link.Sizes)))
		}
		if link.As != "" {
			builder.WriteString(fmt.Sprintf(` as="%s"`, html.EscapeString(link.As)))
		}
		if link.Referrer != "" {
			builder.WriteString(fmt.Sprintf(` referrerpolicy="%s"`, html.EscapeString(link.Referrer)))
		}
		if link.Title != "" {
			builder.WriteString(fmt.Sprintf(` title="%s"`, html.EscapeString(link.Title)))
		}

		builder.WriteString(`>` + "\n")
	}

	return builder.String(), nil
}

func ParseStaticMetadataForPaths(paths []string) (string, error) {
	// Load static metadata
	staticPath := filepath.Join(os.Getenv("ETC_DIR"), "config", "static.route.json")
	staticData, err := os.ReadFile(staticPath)
	if err != nil {
		return "", fmt.Errorf("failed to read static metadata file: %w", err)
	}

	var staticRoutes []types.RouteMetadata
	if err := json.Unmarshal(staticData, &staticRoutes); err != nil {
		return "", fmt.Errorf("failed to parse static metadata JSON: %w", err)
	}

	// Filter routes for only explicitly provided paths
	var mergedMeta types.RouteMetadata
	for _, path := range paths {
		for i := range staticRoutes {
			if staticRoutes[i].Path == path {
				// First non-empty title wins
				if mergedMeta.Title == "" && staticRoutes[i].Title != "" {
					mergedMeta.Title = staticRoutes[i].Title
				}
				mergedMeta.Meta = append(mergedMeta.Meta, staticRoutes[i].Meta...)
				mergedMeta.Link = append(mergedMeta.Link, staticRoutes[i].Link...)
				break
			}
		}
	}

	// Render HTML
	var builder strings.Builder

	if mergedMeta.Title != "" {
		builder.WriteString(fmt.Sprintf(`<title id="__SERVER_PROPS__">%s</title>`+"\n", html.EscapeString(mergedMeta.Title)))
	}

	for _, meta := range mergedMeta.Meta {
		switch {
		case meta.Charset != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" charset="%s">`+"\n", html.EscapeString(meta.Charset)))
		case meta.Name != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" name="%s" content="%s">`+"\n", html.EscapeString(meta.Name), html.EscapeString(meta.Content)))
		case meta.Property != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" property="%s" content="%s">`+"\n", html.EscapeString(meta.Property), html.EscapeString(meta.Content)))
		case meta.HTTPEquiv != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" http-equiv="%s" content="%s">`+"\n", html.EscapeString(meta.HTTPEquiv), html.EscapeString(meta.Content)))
		}
	}

	for _, link := range mergedMeta.Link {
		builder.WriteString(`<link id="__SERVER_PROPS__"`)
		builder.WriteString(fmt.Sprintf(` rel="%s"`, html.EscapeString(link.Rel)))
		builder.WriteString(fmt.Sprintf(` href="%s"`, html.EscapeString(link.Href)))

		if link.Type != "" {
			builder.WriteString(fmt.Sprintf(` type="%s"`, html.EscapeString(link.Type)))
		}
		if link.Crossorigin != "" {
			builder.WriteString(fmt.Sprintf(` crossorigin="%s"`, html.EscapeString(link.Crossorigin)))
		}
		if link.Media != "" {
			builder.WriteString(fmt.Sprintf(` media="%s"`, html.EscapeString(link.Media)))
		}
		if link.Sizes != "" {
			builder.WriteString(fmt.Sprintf(` sizes="%s"`, html.EscapeString(link.Sizes)))
		}
		if link.As != "" {
			builder.WriteString(fmt.Sprintf(` as="%s"`, html.EscapeString(link.As)))
		}
		if link.Referrer != "" {
			builder.WriteString(fmt.Sprintf(` referrerpolicy="%s"`, html.EscapeString(link.Referrer)))
		}
		if link.Title != "" {
			builder.WriteString(fmt.Sprintf(` title="%s"`, html.EscapeString(link.Title)))
		}

		builder.WriteString(`>` + "\n")
	}

	return builder.String(), nil
}
