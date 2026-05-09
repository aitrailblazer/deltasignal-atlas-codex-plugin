package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultBaseURL = "https://api.aitrailblazer.net"

const (
	maxLimit         = 100
	maxOffset        = 1000
	maxTickerLength  = 16
	maxResponseBytes = 2 << 20
)

type deltaClient struct {
	baseURL    string
	httpClient *http.Client
	codexUser  string
	apiKey     string
	mode       string
}

type readinessArgs struct{}

type topStressedArgs struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type covenantStressArgs struct {
	Ticker       string `json:"ticker"`
	Period       string `json:"period"`
	SourceDate   string `json:"source_date"`
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	RiskTier     string `json:"risk_tier"`
	QualityFlag  string `json:"quality_flag"`
	DebtCoverage string `json:"debt_coverage_status"`
	MinStress    int    `json:"min_stress"`
	LinkbaseOnly bool   `json:"linkbase_only"`
}

type tickerArgs struct {
	Ticker     string `json:"ticker"`
	Period     string `json:"period"`
	SourceDate string `json:"source_date"`
}

type outputModeArgs struct {
	OutputMode string `json:"output_mode"`
}

type companyReportArgs struct {
	Ticker              string `json:"ticker"`
	Period              string `json:"period"`
	IncludeSegments     bool   `json:"include_segments"`
	IncludeRelatedParty bool   `json:"include_related_party"`
	OutputMode          string `json:"output_mode"`
}

type quickTickerCheckArgs struct {
	Ticker     string `json:"ticker"`
	OutputMode string `json:"output_mode"`
}

type compositeLeg struct {
	Name        string
	Path        string
	Query       url.Values
	Args        map[string]any
	Criticality string
	PriceUSD    float64
}

func main() {
	client := newDeltaClient()
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "deltasignal-atlas-7",
		Version: "0.2.0",
	}, &mcp.ServerOptions{
		Instructions: "DeltaSignal ATLAS-7 delivers SEC-grounded financial intelligence for crypto public companies. Prefer composite MCP tools for common workflows: deltasignal_morning_brief, deltasignal_company_report, deltasignal_pressure_board, deltasignal_alpha_sweep, and deltasignal_quick_ticker_check. First 5 granular calls are free where supported; paid live usage is charged through x402 USDC micropayments.",
	})

	mcp.AddTool(server, tool("deltasignal_morning_brief", "Server-enforced daily DeltaSignal scan. Internally calls readiness, daily_changes, risk_distribution, top_stressed(limit=10), and alpha_opportunities(limit=10). Use this for morning briefs and daily scans instead of manually chaining granular tools.", noArgCompositeSchema()), client.morningBrief)
	mcp.AddTool(server, tool("deltasignal_company_report", "Server-enforced single-issuer report. Internally calls readiness, company_fundamentals, alpha_signals, peer_ranking, and covenant_stress for one normalized ticker. SPECTRA is not included by default.", companyReportSchema()), client.companyReport)
	mcp.AddTool(server, tool("deltasignal_pressure_board", "Server-enforced risk view. Internally calls readiness, top_stressed(limit=15), and risk_distribution with no offset or issuer drilldowns.", noArgCompositeSchema()), client.pressureBoard)
	mcp.AddTool(server, tool("deltasignal_alpha_sweep", "Server-enforced opportunity view. Internally calls readiness, alpha_opportunities(limit=15), and daily_changes with no offset or issuer drilldowns.", noArgCompositeSchema()), client.alphaSweep)
	mcp.AddTool(server, tool("deltasignal_quick_ticker_check", "Server-enforced fast ticker check. Internally calls readiness, covenant_stress(ticker), and alpha_signals(ticker).", quickTickerCheckSchema()), client.quickTickerCheck)

	mcp.AddTool(server, tool("deltasignal_readiness", "Checks the live ATLAS-7 data plane before analysis, including service readiness, active data freshness, issuer coverage, and whether the current signal slice is safe to query.", map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"additionalProperties": false,
	}), client.readiness)

	mcp.AddTool(server, tool("deltasignal_top_stressed", "Ranks the most stressed crypto public companies in the active DeltaSignal slice, using covenant stress, debt coverage, quality flags, and current issuer risk signals for triage workflows.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"limit":  boundedIntSchema("Maximum issuers to return.", 1, maxLimit),
			"offset": boundedIntSchema("Offset for pagination.", 0, maxOffset),
		},
		"additionalProperties": false,
	}), client.topStressed)

	mcp.AddTool(server, tool("deltasignal_covenant_stress", "Fetches ATLAS-7 covenant stress for a single ticker or screens the active issuer universe with filters for risk tier, quality flags, debt coverage status, linkbase-backed rows, and minimum stress score.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"ticker":               boundedStringSchema("Optional ticker. When set, returns the detail route for that issuer.", 1, maxTickerLength),
			"period":               dateSchema("Optional YYYY-MM-DD period for ticker detail."),
			"source_date":          dateSchema("Optional YYYY-MM-DD source date for list mode."),
			"limit":                boundedIntSchema("Maximum rows for list mode.", 1, maxLimit),
			"offset":               boundedIntSchema("Offset for list mode.", 0, maxOffset),
			"risk_tier":            boundedStringSchema("Optional risk tier filter for list mode.", 1, 64),
			"quality_flag":         boundedStringSchema("Optional quality flag filter for list mode.", 1, 64),
			"debt_coverage_status": boundedStringSchema("Optional debt coverage status filter for list mode.", 1, 64),
			"min_stress":           boundedIntSchema("Optional minimum stress score for list mode.", 0, 100),
			"linkbase_only":        boolSchema("When true, restrict list mode to linkbase-backed issuers."),
		},
		"additionalProperties": false,
	}), client.covenantStress)

	mcp.AddTool(server, tickerTool("deltasignal_peer_ranking", "Compares one crypto public company against its peer group and returns relative covenant stress ranking, percentile context, and peer-comparison signals for issuer research."), client.peerRanking)
	mcp.AddTool(server, tickerTool("deltasignal_alpha_signals", "Returns high-conviction alpha and edge signals for one crypto public company ticker, including resilience, treasury, regime-fit, and phase-one opportunity indicators."), client.alphaSignals)
	mcp.AddTool(server, tickerTool("deltasignal_company_fundamentals", "Returns SEC XBRL-backed company fundamentals for one ticker, with optional segment facts and related-party facts for deeper financial diligence."), client.companyFundamentals)

	mcp.AddTool(server, tool("deltasignal_risk_distribution", "Summarizes the active issuer universe by ATLAS-7 risk tier so agents can understand market-wide stress concentration before drilling into individual companies.", map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"additionalProperties": false,
	}), client.riskDistribution)

	mcp.AddTool(server, tool("deltasignal_daily_changes", "Returns the latest daily signal-change snapshot, highlighting new SEC/XBRL-driven changes, issuer movements, and fresh DeltaSignal updates for monitoring workflows.", map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"additionalProperties": false,
	}), client.dailyChanges)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("deltasignal-atlas-7 MCP failed: %v", err)
		os.Exit(1)
	}
}

func newDeltaClient() *deltaClient {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("DELTASIGNAL_API_BASE_URL")), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	codexUser := strings.TrimSpace(os.Getenv("DELTASIGNAL_CODEX_USER"))
	if codexUser == "" {
		if current, err := user.Current(); err == nil && current != nil {
			codexUser = current.Username
		}
	}
	if codexUser == "" {
		codexUser = "deltasignal-mcp-local"
	}
	return &deltaClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		codexUser: codexUser,
		apiKey:    strings.TrimSpace(os.Getenv("DELTASIGNAL_API_KEY")),
		mode:      normalizeMode(os.Getenv("DELTASIGNAL_PAYMENT_MODE")),
	}
}

func (c *deltaClient) readiness(ctx context.Context, _ *mcp.CallToolRequest, _ readinessArgs) (*mcp.CallToolResult, any, error) {
	return c.get(ctx, "/v1/readiness", nil)
}

func (c *deltaClient) topStressed(ctx context.Context, _ *mcp.CallToolRequest, args topStressedArgs) (*mcp.CallToolResult, any, error) {
	if err := validateBounds(args.Limit, args.Offset, 0); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	q := url.Values{}
	addInt(q, "limit", args.Limit)
	addInt(q, "offset", args.Offset)
	return c.get(ctx, "/v1/top-stressed", q)
}

func (c *deltaClient) covenantStress(ctx context.Context, _ *mcp.CallToolRequest, args covenantStressArgs) (*mcp.CallToolResult, any, error) {
	if err := validateBounds(args.Limit, args.Offset, args.MinStress); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	if err := validateDate("period", args.Period); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	if err := validateDate("source_date", args.SourceDate); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	if err := validateShortFilters(args.RiskTier, args.QualityFlag, args.DebtCoverage); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	ticker := strings.ToUpper(strings.TrimSpace(args.Ticker))
	if ticker != "" {
		if err := validateTicker(ticker); err != nil {
			return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
		}
		q := url.Values{}
		addString(q, "period", args.Period)
		return c.get(ctx, "/v1/covenant-stress/"+url.PathEscape(ticker), q)
	}
	q := url.Values{}
	addString(q, "period", args.Period)
	addString(q, "source_date", args.SourceDate)
	addInt(q, "limit", args.Limit)
	addInt(q, "offset", args.Offset)
	addString(q, "risk_tier", args.RiskTier)
	addString(q, "quality_flag", args.QualityFlag)
	addString(q, "debt_coverage_status", args.DebtCoverage)
	addInt(q, "min_stress", args.MinStress)
	if args.LinkbaseOnly {
		q.Set("linkbase_only", "true")
	}
	return c.get(ctx, "/v1/covenant-stress", q)
}

func (c *deltaClient) peerRanking(ctx context.Context, _ *mcp.CallToolRequest, args tickerArgs) (*mcp.CallToolResult, any, error) {
	return c.getTicker(ctx, "/v1/peer-ranking/", args.Ticker, "period", args.Period)
}

func (c *deltaClient) alphaSignals(ctx context.Context, _ *mcp.CallToolRequest, args tickerArgs) (*mcp.CallToolResult, any, error) {
	return c.getTicker(ctx, "/v1/alpha-signals/", args.Ticker, "source_date", args.SourceDate)
}

func (c *deltaClient) companyFundamentals(ctx context.Context, _ *mcp.CallToolRequest, args tickerArgs) (*mcp.CallToolResult, any, error) {
	return c.getTicker(ctx, "/v1/company-fundamentals/", args.Ticker, "period", args.Period)
}

func (c *deltaClient) riskDistribution(ctx context.Context, _ *mcp.CallToolRequest, _ readinessArgs) (*mcp.CallToolResult, any, error) {
	return c.get(ctx, "/v1/risk-distribution", nil)
}

func (c *deltaClient) dailyChanges(ctx context.Context, _ *mcp.CallToolRequest, _ readinessArgs) (*mcp.CallToolResult, any, error) {
	return c.get(ctx, "/v1/daily-changes/latest", nil)
}

func (c *deltaClient) morningBrief(ctx context.Context, _ *mcp.CallToolRequest, args outputModeArgs) (*mcp.CallToolResult, any, error) {
	if err := validateOutputMode(args.OutputMode); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	return c.runComposite(ctx, "deltasignal_morning_brief", map[string]any{"output_mode": "compact"}, []compositeLeg{
		{Name: "deltasignal_readiness", Path: "/v1/readiness", Criticality: "advisory", PriceUSD: 0.05},
		{Name: "deltasignal_daily_changes", Path: "/v1/daily-changes/latest", Criticality: "required", PriceUSD: 0.04},
		{Name: "deltasignal_risk_distribution", Path: "/v1/risk-distribution", Criticality: "required", PriceUSD: 0.05},
		{Name: "deltasignal_top_stressed", Path: "/v1/top-stressed", Query: values(map[string]string{"limit": "10"}), Args: map[string]any{"limit": 10}, Criticality: "required", PriceUSD: 0.06},
		{Name: "deltasignal_alpha_opportunities", Path: "/v1/alpha-opportunities", Query: values(map[string]string{"limit": "10"}), Args: map[string]any{"limit": 10}, Criticality: "required", PriceUSD: 0.06},
	}, 0.18)
}

func (c *deltaClient) companyReport(ctx context.Context, _ *mcp.CallToolRequest, args companyReportArgs) (*mcp.CallToolResult, any, error) {
	if err := validateOutputMode(args.OutputMode); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	ticker := strings.ToUpper(strings.TrimSpace(args.Ticker))
	if err := validateTicker(ticker); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	if err := validateDate("period", args.Period); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	q := url.Values{}
	addString(q, "period", args.Period)
	request := map[string]any{"ticker": ticker, "output_mode": "compact"}
	if args.Period != "" {
		request["period"] = args.Period
	}
	return c.runComposite(ctx, "deltasignal_company_report", request, []compositeLeg{
		{Name: "deltasignal_readiness", Path: "/v1/readiness", Criticality: "advisory", PriceUSD: 0.05},
		{Name: "deltasignal_company_fundamentals", Path: "/v1/company-fundamentals/" + url.PathEscape(ticker), Query: q, Args: map[string]any{"ticker": ticker}, Criticality: "required", PriceUSD: 0.07},
		{Name: "deltasignal_alpha_signals", Path: "/v1/alpha-signals/" + url.PathEscape(ticker), Args: map[string]any{"ticker": ticker}, Criticality: "required", PriceUSD: 0.09},
		{Name: "deltasignal_peer_ranking", Path: "/v1/peer-ranking/" + url.PathEscape(ticker), Query: q, Args: map[string]any{"ticker": ticker}, Criticality: "required", PriceUSD: 0.08},
		{Name: "deltasignal_covenant_stress", Path: "/v1/covenant-stress/" + url.PathEscape(ticker), Query: q, Args: map[string]any{"ticker": ticker}, Criticality: "required", PriceUSD: 0.15},
	}, 0.30)
}

func (c *deltaClient) pressureBoard(ctx context.Context, _ *mcp.CallToolRequest, args outputModeArgs) (*mcp.CallToolResult, any, error) {
	if err := validateOutputMode(args.OutputMode); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	return c.runComposite(ctx, "deltasignal_pressure_board", map[string]any{"output_mode": "compact"}, []compositeLeg{
		{Name: "deltasignal_readiness", Path: "/v1/readiness", Criticality: "advisory", PriceUSD: 0.05},
		{Name: "deltasignal_top_stressed", Path: "/v1/top-stressed", Query: values(map[string]string{"limit": "15"}), Args: map[string]any{"limit": 15}, Criticality: "required", PriceUSD: 0.06},
		{Name: "deltasignal_risk_distribution", Path: "/v1/risk-distribution", Criticality: "required", PriceUSD: 0.05},
	}, 0.14)
}

func (c *deltaClient) alphaSweep(ctx context.Context, _ *mcp.CallToolRequest, args outputModeArgs) (*mcp.CallToolResult, any, error) {
	if err := validateOutputMode(args.OutputMode); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	return c.runComposite(ctx, "deltasignal_alpha_sweep", map[string]any{"output_mode": "compact"}, []compositeLeg{
		{Name: "deltasignal_readiness", Path: "/v1/readiness", Criticality: "advisory", PriceUSD: 0.05},
		{Name: "deltasignal_alpha_opportunities", Path: "/v1/alpha-opportunities", Query: values(map[string]string{"limit": "15"}), Args: map[string]any{"limit": 15}, Criticality: "required", PriceUSD: 0.06},
		{Name: "deltasignal_daily_changes", Path: "/v1/daily-changes/latest", Criticality: "required", PriceUSD: 0.04},
	}, 0.14)
}

func (c *deltaClient) quickTickerCheck(ctx context.Context, _ *mcp.CallToolRequest, args quickTickerCheckArgs) (*mcp.CallToolResult, any, error) {
	if err := validateOutputMode(args.OutputMode); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	ticker := strings.ToUpper(strings.TrimSpace(args.Ticker))
	if err := validateTicker(ticker); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	return c.runComposite(ctx, "deltasignal_quick_ticker_check", map[string]any{"ticker": ticker, "output_mode": "compact"}, []compositeLeg{
		{Name: "deltasignal_readiness", Path: "/v1/readiness", Criticality: "advisory", PriceUSD: 0.05},
		{Name: "deltasignal_covenant_stress", Path: "/v1/covenant-stress/" + url.PathEscape(ticker), Args: map[string]any{"ticker": ticker}, Criticality: "required", PriceUSD: 0.15},
		{Name: "deltasignal_alpha_signals", Path: "/v1/alpha-signals/" + url.PathEscape(ticker), Args: map[string]any{"ticker": ticker}, Criticality: "required", PriceUSD: 0.09},
	}, 0.18)
}

func (c *deltaClient) runComposite(ctx context.Context, workflow string, request map[string]any, legs []compositeLeg, bundlePrice float64) (*mcp.CallToolResult, any, error) {
	start := time.Now()
	data := map[string]any{
		"schema_version": "2026-05-09",
		"workflow":       workflow,
		"request":        request,
		"payment_mode":   c.mode,
	}
	ledger := make([]map[string]any, 0, len(legs))
	errors := make([]map[string]any, 0)
	requiredFailed := 0
	requiredTotal := 0
	for _, leg := range legs {
		if leg.Criticality == "required" {
			requiredTotal++
		}
		legStart := time.Now()
		_, payloadAny, err := c.get(ctx, leg.Path, leg.Query)
		elapsed := time.Since(legStart).Milliseconds()
		payload, _ := payloadAny.(map[string]any)
		status := 0
		if payload != nil {
			if raw, ok := payload["status"].(float64); ok {
				status = int(raw)
			} else if raw, ok := payload["status"].(int); ok {
				status = raw
			}
		}
		ok := err == nil && status >= 200 && status < 300
		row := map[string]any{
			"tool":             leg.Name,
			"route":            routeName(leg.Path),
			"args":             legArgs(leg.Args),
			"criticality":      leg.Criticality,
			"equivalent_price": fmt.Sprintf("$%.2f", leg.PriceUSD),
			"price_usd":        leg.PriceUSD,
			"elapsed_ms":       elapsed,
			"status":           "ok",
		}
		if !ok {
			row["status"] = "error"
			row["error"] = errorText(err, payload)
			errors = append(errors, row)
			if leg.Criticality == "required" {
				requiredFailed++
			}
		} else {
			data[resultKey(leg.Name)] = payload["data"]
		}
		ledger = append(ledger, row)
	}
	partial := requiredFailed > 0
	status := "ok"
	if partial {
		status = "partial"
	}
	if requiredTotal > 0 && requiredFailed == requiredTotal {
		status = "error"
	}
	data["status"] = status
	data["partial"] = partial
	data["internal_calls"] = ledger
	data["errors"] = errors
	data["mcp_billing_estimate"] = map[string]any{
		"payment_flow_used":             false,
		"payment_mode":                  c.mode,
		"bundle_x402_equivalent_price":  fmt.Sprintf("$%.2f", bundlePrice),
		"bundle_x402_price_usd":         bundlePrice,
		"charged_once":                  false,
		"internal_call_count":           len(legs),
		"note":                          "Local plugin composite executes deterministic granular reads. The hosted /mcp composite is the canonical server-enforced workflow.",
		"public_credit_packs_available": false,
	}
	data["elapsed_ms"] = time.Since(start).Milliseconds()
	return result(data, status == "error")
}

func (c *deltaClient) getTicker(ctx context.Context, pathPrefix, ticker, dateField, dateValue string) (*mcp.CallToolResult, any, error) {
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	if err := validateTicker(ticker); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	if err := validateDate(dateField, dateValue); err != nil {
		return result(map[string]any{"error": err.Error(), "payment_mode": c.mode}, true)
	}
	q := url.Values{}
	addString(q, dateField, dateValue)
	return c.get(ctx, pathPrefix+url.PathEscape(ticker), q)
}

func (c *deltaClient) get(ctx context.Context, path string, q url.Values) (*mcp.CallToolResult, any, error) {
	endpoint, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, nil, err
	}
	if len(q) > 0 {
		endpoint.RawQuery = q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Codex-User", c.codexUser)
	if c.mode == "local" {
		req.Header.Set("X-Test-Mode", "free")
	}
	if c.mode == "internal" && c.apiKey == "" {
		return result(map[string]any{
			"error":        "DELTASIGNAL_API_KEY is required when DELTASIGNAL_PAYMENT_MODE=internal",
			"payment_mode": c.mode,
			"url":          endpoint.String(),
		}, true)
	}
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return result(map[string]any{"error": err.Error(), "url": endpoint.String()}, true)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if readErr != nil {
		return result(map[string]any{"error": readErr.Error(), "status": resp.StatusCode, "url": endpoint.String()}, true)
	}
	if len(body) > maxResponseBytes {
		return result(map[string]any{"error": "response body too large", "status": resp.StatusCode, "url": endpoint.String()}, true)
	}
	payload := map[string]any{
		"status":       resp.StatusCode,
		"url":          endpoint.String(),
		"payment_mode": c.mode,
		"data":         parseJSON(body),
	}
	if freeTier := resp.Header.Get("X-DeltaSignal-Free-Tier"); freeTier != "" {
		payload["free_tier"] = map[string]any{
			"mode":      freeTier,
			"remaining": resp.Header.Get("X-DeltaSignal-Free-Tier-Remaining"),
			"limit":     resp.Header.Get("X-DeltaSignal-Free-Tier-Limit"),
		}
	}
	if resp.StatusCode == http.StatusPaymentRequired {
		payload["payment_required"] = map[string]any{
			"payment_required": resp.Header.Get("PAYMENT-REQUIRED"),
			"www_authenticate": resp.Header.Get("WWW-Authenticate"),
			"note":             "Free tier is exhausted or unavailable. Use an x402-capable client against the same /v1 route, or use the /mpp/v1 route with a Tempo MPP client.",
		}
		return result(payload, true)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result(payload, true)
	}
	return result(payload, false)
}

func normalizeMode(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "local", "dev", "test", "free":
		return "local"
	case "internal", "api-key", "apikey", "deployed-test":
		return "internal"
	default:
		return "live"
	}
}

func result(payload map[string]any, isError bool) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		IsError: isError,
	}, payload, nil
}

func parseJSON(body []byte) any {
	if len(strings.TrimSpace(string(body))) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		return string(body)
	}
	return v
}

func tool(name, description string, schema map[string]any) *mcp.Tool {
	return &mcp.Tool{
		Name:        name,
		Title:       name,
		Description: description,
		Annotations: &mcp.ToolAnnotations{
			Title:           name,
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(false),
		},
		InputSchema: schema,
	}
}

func tickerTool(name, description string) *mcp.Tool {
	properties := map[string]any{
		"ticker": boundedStringSchema("Required issuer ticker, normalized to uppercase.", 1, maxTickerLength),
	}
	switch name {
	case "deltasignal_peer_ranking", "deltasignal_company_fundamentals":
		properties["period"] = dateSchema("Optional YYYY-MM-DD period.")
	case "deltasignal_alpha_signals":
		properties["source_date"] = dateSchema("Optional YYYY-MM-DD source date.")
	}
	return tool(name, description, map[string]any{
		"type":                 "object",
		"properties":           properties,
		"required":             []string{"ticker"},
		"additionalProperties": false,
	})
}

func noArgCompositeSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"output_mode": boundedStringSchema("Optional response mode. Only compact is accepted.", 1, 16),
		},
		"additionalProperties": false,
	}
}

func companyReportSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"ticker":                boundedStringSchema("Required issuer ticker, normalized to uppercase.", 1, maxTickerLength),
			"period":                dateSchema("Optional YYYY-MM-DD period."),
			"include_segments":      boolSchema("Reserved for hosted MCP compatibility."),
			"include_related_party": boolSchema("Reserved for hosted MCP compatibility."),
			"output_mode":           boundedStringSchema("Optional response mode. Only compact is accepted.", 1, 16),
		},
		"required":             []string{"ticker"},
		"additionalProperties": false,
	}
}

func quickTickerCheckSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"ticker":      boundedStringSchema("Required issuer ticker, normalized to uppercase.", 1, maxTickerLength),
			"output_mode": boundedStringSchema("Optional response mode. Only compact is accepted.", 1, 16),
		},
		"required":             []string{"ticker"},
		"additionalProperties": false,
	}
}

func stringSchema(description string) map[string]any {
	return map[string]any{"type": "string", "description": description}
}

func intSchema(description string) map[string]any {
	return map[string]any{"type": "integer", "description": description}
}

func boundedIntSchema(description string, min, max int) map[string]any {
	schema := intSchema(description)
	schema["minimum"] = min
	schema["maximum"] = max
	return schema
}

func boundedStringSchema(description string, min, max int) map[string]any {
	schema := stringSchema(description)
	schema["minLength"] = min
	schema["maxLength"] = max
	return schema
}

func dateSchema(description string) map[string]any {
	schema := boundedStringSchema(description, 10, 10)
	schema["pattern"] = `^\d{4}-\d{2}-\d{2}$`
	return schema
}

func boolSchema(description string) map[string]any {
	return map[string]any{"type": "boolean", "description": description}
}

func boolPtr(v bool) *bool {
	return &v
}

func addString(q url.Values, key, value string) {
	value = strings.TrimSpace(value)
	if value != "" {
		q.Set(key, value)
	}
}

func addInt(q url.Values, key string, value int) {
	if value > 0 {
		q.Set(key, strconv.Itoa(value))
	}
}

func validateBounds(limit, offset, minStress int) error {
	if limit < 0 || limit > maxLimit {
		return fmt.Errorf("limit must be between 1 and %d when set", maxLimit)
	}
	if offset < 0 || offset > maxOffset {
		return fmt.Errorf("offset must be between 0 and %d", maxOffset)
	}
	if minStress < 0 || minStress > 100 {
		return fmt.Errorf("min_stress must be between 0 and 100")
	}
	return nil
}

func validateTicker(ticker string) error {
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	if ticker == "" {
		return fmt.Errorf("ticker is required")
	}
	if len(ticker) > maxTickerLength {
		return fmt.Errorf("ticker is too long")
	}
	for _, r := range ticker {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' {
			continue
		}
		return fmt.Errorf("ticker contains unsupported characters")
	}
	return nil
}

func validateDate(field, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if len(value) != 10 || value[4] != '-' || value[7] != '-' {
		return fmt.Errorf("%s must be YYYY-MM-DD", field)
	}
	for i, r := range value {
		if i == 4 || i == 7 {
			continue
		}
		if r < '0' || r > '9' {
			return fmt.Errorf("%s must be YYYY-MM-DD", field)
		}
	}
	return nil
}

func validateShortFilters(values ...string) error {
	for _, value := range values {
		if len(strings.TrimSpace(value)) > 64 {
			return fmt.Errorf("filter value is too long")
		}
	}
	return nil
}

func validateOutputMode(value string) error {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || value == "compact" {
		return nil
	}
	return fmt.Errorf("output_mode must be compact when set")
}

func values(raw map[string]string) url.Values {
	q := url.Values{}
	for key, value := range raw {
		q.Set(key, value)
	}
	return q
}

func legArgs(args map[string]any) map[string]any {
	if args == nil {
		return map[string]any{}
	}
	return args
}

func routeName(path string) string {
	for _, prefix := range []string{"/v1/company-fundamentals/", "/v1/alpha-signals/", "/v1/peer-ranking/", "/v1/covenant-stress/"} {
		if strings.HasPrefix(path, prefix) {
			return "GET " + strings.TrimRight(prefix, "/") + "/:ticker"
		}
	}
	return "GET " + path
}

func resultKey(toolName string) string {
	return strings.TrimPrefix(toolName, "deltasignal_")
}

func errorText(err error, payload map[string]any) string {
	if err != nil {
		return err.Error()
	}
	if payload == nil {
		return "request failed"
	}
	if data, ok := payload["data"].(map[string]any); ok {
		if msg, ok := data["error"].(string); ok && strings.TrimSpace(msg) != "" {
			return msg
		}
	}
	if msg, ok := payload["error"].(string); ok && strings.TrimSpace(msg) != "" {
		return msg
	}
	if status, ok := payload["status"]; ok {
		return fmt.Sprintf("request failed with status %v", status)
	}
	return "request failed"
}

func _printf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}
