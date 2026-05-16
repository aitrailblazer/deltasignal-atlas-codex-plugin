package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type registryResponse struct {
	Servers []struct {
		Server struct {
			Name        string `json:"name"`
			Version     string `json:"version"`
			Description string `json:"description"`
			WebsiteURL  string `json:"websiteUrl"`
			Remotes     []struct {
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"remotes"`
		} `json:"server"`
		Meta struct {
			Official struct {
				IsLatest  bool   `json:"isLatest"`
				UpdatedAt string `json:"updatedAt"`
			} `json:"io.modelcontextprotocol.registry/official"`
		} `json:"_meta"`
	} `json:"servers"`
}

func main() {
	expectedVersion := flag.String("version", "1.1.4", "expected latest official registry version")
	expectedRemote := flag.String("remote", "https://api.aitrailblazer.net/mcp", "expected Streamable HTTP MCP endpoint")
	expectedWebsite := flag.String("website", "https://aitrailblazer.github.io/deltasignal-atlas-codex-plugin", "expected website URL")
	registryURL := flag.String("registry", "https://registry.modelcontextprotocol.io/v0.1/servers?search=deltasignal&limit=10", "official registry search URL")
	storkURL := flag.String("stork", "https://www.stork.ai/mcp/io-github-aitrailblazer-deltasignal-atlas-7", "Stork listing URL")
	timeout := flag.Duration("timeout", 20*time.Second, "HTTP timeout")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	registryBody := fetch(ctx, *registryURL)
	storkBody := fetch(ctx, *storkURL)

	latest := latestRegistryRecord(registryBody)
	registryOK := latest.Server.Version == *expectedVersion &&
		latest.Server.WebsiteURL == *expectedWebsite &&
		hasRemote(latest.Server.Remotes, *expectedRemote)

	storkText := string(storkBody)
	storkHasVersion := strings.Contains(storkText, *expectedVersion)
	storkHasRemote := strings.Contains(storkText, *expectedRemote)
	storkHasWebsite := strings.Contains(storkText, *expectedWebsite)
	storkHasFakeNPM := strings.Contains(storkText, "@deltasignal-atlas-7/mcp-server")
	storkOK := storkHasVersion && storkHasRemote && storkHasWebsite && !storkHasFakeNPM

	result := map[string]any{
		"ok":                  registryOK && storkOK,
		"registry_ok":         registryOK,
		"stork_ok":            storkOK,
		"registry_version":    latest.Server.Version,
		"registry_updated_at": latest.Meta.Official.UpdatedAt,
		"registry_website":    latest.Server.WebsiteURL,
		"registry_remotes":    latest.Server.Remotes,
		"stork_has_version":   storkHasVersion,
		"stork_has_remote":    storkHasRemote,
		"stork_has_website":   storkHasWebsite,
		"stork_has_fake_npm":  storkHasFakeNPM,
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	must("marshal result", err)
	fmt.Println(string(encoded))

	if !registryOK || !storkOK {
		os.Exit(1)
	}
}

func latestRegistryRecord(body []byte) struct {
	Server struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		WebsiteURL  string `json:"websiteUrl"`
		Remotes     []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"remotes"`
	} `json:"server"`
	Meta struct {
		Official struct {
			IsLatest  bool   `json:"isLatest"`
			UpdatedAt string `json:"updatedAt"`
		} `json:"io.modelcontextprotocol.registry/official"`
	} `json:"_meta"`
} {
	var decoded registryResponse
	must("decode registry response", json.Unmarshal(body, &decoded))
	for _, item := range decoded.Servers {
		if item.Server.Name == "io.github.aitrailblazer/deltasignal-atlas-7" && item.Meta.Official.IsLatest {
			return item
		}
	}
	fatalf("no latest registry record found for io.github.aitrailblazer/deltasignal-atlas-7")
	return decoded.Servers[0]
}

func hasRemote(remotes []struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}, expected string) bool {
	for _, remote := range remotes {
		if remote.Type == "streamable-http" && remote.URL == expected {
			return true
		}
	}
	return false
}

func fetch(ctx context.Context, rawURL string) []byte {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	must("build request", err)
	resp, err := http.DefaultClient.Do(req)
	must("fetch "+rawURL, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	must("read response", err)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fatalf("%s returned HTTP %d: %s", rawURL, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return body
}

func must(label string, err error) {
	if err != nil {
		fatalf("%s: %v", label, err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
