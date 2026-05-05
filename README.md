# DeltaSignal ATLAS-7 - Native Codex MCP Server

Real-time issuer intelligence as a native MCP server for Codex CLI and other MCP agents.

DeltaSignal ATLAS-7 gives your agent structured tools for crypto public-company stress, risk analysis, peer ranking, alpha signals, and fundamentals, with a free starter tier and x402 payment fallback. No subscription required.

## Quick Start

### 1. Install

```bash
npx codex-marketplace add aitrailblazer/deltasignal-atlas-codex-plugin --plugin --project
```

For local development from this checkout:

```bash
codex plugin install .
```

### 2. Payment Setup

```bash
npx @coinbase/payments-mcp install
```

- Select `Codex CLI`
- Sign in with email
- Fund your wallet with USDC

Restart Codex CLI after payment setup.

### 3. Use

```bash
@DeltaSignal check readiness
@DeltaSignal show top stressed issuers
@DeltaSignal alpha signals for MARA
@DeltaSignal peer ranking for RIOT
@DeltaSignal company fundamentals for COIN
```

## Access Model

- First 5 live calls are free per user.
- After the free tier, protected live routes return x402 payment requirements for USDC payment-capable clients.
- Internal and local development profiles can test without payment.
- All MCP tools are read-only, strictly validated, and closed-world.

## Available Tools

- `deltasignal_readiness`
- `deltasignal_top_stressed`
- `deltasignal_covenant_stress`
- `deltasignal_peer_ranking`
- `deltasignal_alpha_signals`
- `deltasignal_company_fundamentals`
- `deltasignal_risk_distribution`
- `deltasignal_daily_changes`

## Development Modes

Production marketplace profile:

- `.mcp.json`
- `DELTASIGNAL_PAYMENT_MODE=live`
- `DELTASIGNAL_API_BASE_URL=https://api.aitrailblazer.net`
- No embedded test key

Local no-payment profile:

- `.mcp.local.json`
- `DELTASIGNAL_PAYMENT_MODE=local`
- `DELTASIGNAL_API_BASE_URL=http://127.0.0.1:8080`
- Uses `local-dev-key` for the local gateway

Run the local gateway with:

```bash
GOWORK=off GATEWAY_PORT=8080 CRAWLER_ENABLED=false MCP_API_KEY=local-dev-key go run ./cmd/search-gateway
```

Deployed no-payment test profile:

- `.mcp.deployed-test.json`
- `DELTASIGNAL_PAYMENT_MODE=internal`
- Requires `DELTASIGNAL_API_KEY` in the MCP process environment

Use an internal pre-authorized key from Key Vault or local secret storage for deployed smoke tests. Do not commit that key into the plugin.

## Technical

- Bundled STDIO MCP server
- Strict input schemas with bounded arguments
- Ticker, date, filter, and response-size validation
- Structured MCP tool results
- Compatible with Codex and other MCP clients

## Optional Environment

- `DELTASIGNAL_API_BASE_URL` defaults to `https://api.aitrailblazer.net`
- `DELTASIGNAL_PAYMENT_MODE=live` uses the deployed x402/free-tier path
- `DELTASIGNAL_PAYMENT_MODE=internal` requires `DELTASIGNAL_API_KEY` and bypasses x402 through the internal API-key path
- `DELTASIGNAL_PAYMENT_MODE=local` marks tool results as local mode and sends `X-Test-Mode: free`
- `DELTASIGNAL_CODEX_USER` sets the free-tier user identity header
- `DELTASIGNAL_API_KEY` uses the internal/pre-authorized API key path
