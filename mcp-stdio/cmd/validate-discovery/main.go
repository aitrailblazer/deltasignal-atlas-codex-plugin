package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type arazzoDoc struct {
	Workflows []workflow `json:"workflows"`
}

type workflow struct {
	WorkflowID string `json:"workflowId"`
	Steps      []step `json:"steps"`
}

type step struct {
	StepID      string `json:"stepId"`
	OperationID string `json:"operationId"`
	MCPTool     string `json:"x-mcp-tool"`
}

type openAPIDoc struct {
	Paths map[string]map[string]struct {
		OperationID string `json:"operationId"`
	} `json:"paths"`
}

type mcpToolsResponse struct {
	Result struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	} `json:"result"`
	Error any `json:"error"`
}

func main() {
	root := flag.String("root", "..", "plugin repository root")
	openAPIURL := flag.String("openapi", "https://api.aitrailblazer.net/openapi.json", "live OpenAPI URL")
	mcpURL := flag.String("mcp", "https://api.aitrailblazer.net/mcp", "live MCP JSON-RPC URL")
	apiKey := flag.String("api-key", strings.TrimSpace(os.Getenv("DELTASIGNAL_API_KEY")), "optional DeltaSignal API key for live MCP tools/list")
	timeout := flag.Duration("timeout", 20*time.Second, "HTTP timeout")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	repoRoot, err := filepath.Abs(*root)
	must("resolve root", err)

	checkMirror(repoRoot, "deltasignal-arazzo.yaml")
	checkMirror(repoRoot, "deltasignal-arazzo.json")

	arazzo := loadArazzo(repoRoot)
	openAPI := fetchOpenAPI(ctx, *openAPIURL)
	operationIDs := collectOperationIDs(openAPI)
	missingOps := missingArazzoOperations(arazzo, operationIDs)

	mcpRefs := collectArazzoMCPTools(arazzo)
	missingTools := []string{}
	mcpToolCount := 0
	if strings.TrimSpace(*apiKey) != "" {
		liveTools := fetchMCPTools(ctx, *mcpURL, *apiKey)
		mcpToolCount = len(liveTools)
		missingTools = missingMCPTools(mcpRefs, liveTools)
	}

	summary := map[string]any{
		"ok":                      len(missingOps) == 0 && len(missingTools) == 0,
		"workflow_count":          len(arazzo.Workflows),
		"openapi_operation_count": len(operationIDs),
		"referenced_mcp_tools":    len(mcpRefs),
		"live_mcp_tool_count":     mcpToolCount,
		"mcp_tools_checked":       strings.TrimSpace(*apiKey) != "",
		"missing_operation_ids":   missingOps,
		"missing_mcp_tools":       missingTools,
	}
	encoded, err := json.MarshalIndent(summary, "", "  ")
	must("marshal summary", err)
	fmt.Println(string(encoded))
	if len(missingOps) > 0 || len(missingTools) > 0 {
		os.Exit(1)
	}
}

func checkMirror(root, name string) {
	canonical, err := os.ReadFile(filepath.Join(root, "arazzo", name))
	must("read canonical "+name, err)
	mirror, err := os.ReadFile(filepath.Join(root, name))
	must("read mirror "+name, err)
	if !bytes.Equal(canonical, mirror) {
		fatalf("%s root mirror does not match arazzo/%s", name, name)
	}
}

func loadArazzo(root string) arazzoDoc {
	raw, err := os.ReadFile(filepath.Join(root, "arazzo", "deltasignal-arazzo.json"))
	must("read Arazzo JSON", err)
	var doc arazzoDoc
	must("decode Arazzo JSON", json.Unmarshal(raw, &doc))
	if len(doc.Workflows) == 0 {
		fatalf("Arazzo JSON has no workflows")
	}
	return doc
}

func fetchOpenAPI(ctx context.Context, rawURL string) openAPIDoc {
	body := fetch(ctx, rawURL, "")
	var doc openAPIDoc
	must("decode OpenAPI", json.Unmarshal(body, &doc))
	if len(doc.Paths) == 0 {
		fatalf("OpenAPI has no paths")
	}
	return doc
}

func fetchMCPTools(ctx context.Context, rawURL, apiKey string) map[string]bool {
	payload := []byte(`{"jsonrpc":"2.0","id":"tools","method":"tools/list","params":{}}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(payload))
	must("build MCP request", err)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	resp, err := http.DefaultClient.Do(req)
	must("fetch MCP tools", err)
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	must("read MCP tools", err)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fatalf("MCP tools/list returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var decoded mcpToolsResponse
	must("decode MCP tools", json.Unmarshal(body, &decoded))
	if decoded.Error != nil {
		fatalf("MCP tools/list returned error: %v", decoded.Error)
	}
	tools := map[string]bool{}
	for _, tool := range decoded.Result.Tools {
		if tool.Name != "" {
			tools[tool.Name] = true
		}
	}
	return tools
}

func fetch(ctx context.Context, rawURL, apiKey string) []byte {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	must("build GET request", err)
	if strings.TrimSpace(apiKey) != "" {
		req.Header.Set("x-api-key", apiKey)
	}
	resp, err := http.DefaultClient.Do(req)
	must("fetch "+rawURL, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	must("read "+rawURL, err)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fatalf("%s returned HTTP %d: %s", rawURL, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return body
}

func collectOperationIDs(doc openAPIDoc) map[string]bool {
	ids := map[string]bool{}
	for _, pathItem := range doc.Paths {
		for _, op := range pathItem {
			if op.OperationID != "" {
				ids[op.OperationID] = true
			}
		}
	}
	return ids
}

func missingArazzoOperations(doc arazzoDoc, live map[string]bool) []string {
	missing := []string{}
	for _, wf := range doc.Workflows {
		for _, step := range wf.Steps {
			if step.OperationID != "" && !live[step.OperationID] {
				missing = append(missing, wf.WorkflowID+"."+step.StepID+":"+step.OperationID)
			}
		}
	}
	sort.Strings(missing)
	return missing
}

func collectArazzoMCPTools(doc arazzoDoc) []string {
	seen := map[string]bool{}
	for _, wf := range doc.Workflows {
		for _, step := range wf.Steps {
			if step.MCPTool != "" {
				seen[step.MCPTool] = true
			}
		}
	}
	tools := []string{}
	for tool := range seen {
		tools = append(tools, tool)
	}
	sort.Strings(tools)
	return tools
}

func missingMCPTools(refs []string, live map[string]bool) []string {
	missing := []string{}
	for _, ref := range refs {
		if !live[ref] {
			missing = append(missing, ref)
		}
	}
	sort.Strings(missing)
	return missing
}

func must(label string, err error) {
	if err != nil {
		fatalf("%s: %v", label, err)
	}
}

func fatalf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
