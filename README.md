# DeltaSignal ATLAS-7 - Native Codex MCP Server

**Real-time intelligence signals** as a **native MCP server** for Codex CLI and other MCP agents.

DeltaSignal ATLAS-7 gives your agent structured, validated tools for market signals, risk analysis, peer ranking, and fundamentals. Common workflows are exposed as composite MCP presets so agents can call one reliable tool instead of hand-orchestrating several low-level calls. No subscription required.

Also available on [Smithery](https://smithery.ai/servers/aitrailblazer/deltasignal-atlas-7) and as a [Glama Connector](https://glama.ai/mcp/connectors/net.aitrailblazer.api/delta-signal-atlas-7).

## Quick Start

### 1. Install

```bash
codex plugin install aitrailblazer/deltasignal-atlas-codex-plugin
```

Or from an npm environment:

```bash
npx codex-marketplace add aitrailblazer/deltasignal-atlas-codex-plugin --plugin --project
```

### 2. Payment Setup

```bash
npx @coinbase/payments-mcp install
```

- Select `Codex CLI`
- Sign in with email
- Fund your wallet with USDC

Restart Codex CLI after this step.

### 3. Use

```bash
@DeltaSignal morning_brief
@DeltaSignal company_report ticker:RIOT
@DeltaSignal pressure_board
@DeltaSignal alpha_sweep
@DeltaSignal quick_ticker_check ticker:MARA
@DeltaSignal alpha_signals ticker:COIN
@DeltaSignal top_stressed tickers:MARA,RIOT,HUT,CLSK
@DeltaSignal company_fundamentals ticker:MSTR
@DeltaSignal covenant_stress ticker:MARA
@DeltaSignal risk_distribution tickers:COIN,MSTR,MARA,RIOT
```

## Access Model

- First 5 calls are free per user.
- After the free tier, supported clients receive automatic x402 payment requirements in USDC.
- Internal testing can use a pre-authorized API key without payment.
- All tools are read-only, schema-validated, and closed-world.

## Available MCP Tools

Composite presets:

- `deltasignal_morning_brief` - daily scan: readiness + daily changes + risk distribution + top stressed(limit=10) + alpha opportunities(limit=10). Public Builder price: 18 credits / $0.18.
- `deltasignal_company_report` - ticker report: readiness + fundamentals + alpha signals + peer ranking + covenant stress. Public Builder price: 30 credits / $0.30.
- `deltasignal_pressure_board` - risk view: readiness + top stressed(limit=15) + risk distribution. Public Builder price: 14 credits / $0.14.
- `deltasignal_alpha_sweep` - opportunity screen: readiness + alpha opportunities(limit=15) + daily changes. Public Builder price: 14 credits / $0.14.
- `deltasignal_quick_ticker_check` - fast ticker check: readiness + covenant stress + alpha signals. Public Builder price: 18 credits / $0.18.

Granular tools:

- `deltasignal_readiness`
- `deltasignal_top_stressed`
- `deltasignal_covenant_stress`
- `deltasignal_peer_ranking`
- `deltasignal_alpha_signals`
- `deltasignal_company_fundamentals`
- `deltasignal_risk_distribution`
- `deltasignal_daily_changes`

Credit packs are not implemented yet. Pricing metadata uses `1 credit = $0.01` of DeltaSignal usage value.

## Development Modes

- Local: `DELTASIGNAL_PAYMENT_MODE=local`
- Internal testing: `DELTASIGNAL_PAYMENT_MODE=internal` plus `DELTASIGNAL_API_KEY`
- Live production: `DELTASIGNAL_PAYMENT_MODE=live`

## Technical

- Bundled STDIO MCP server in `mcp-stdio/`
- Remote MCP endpoint: `https://api.aitrailblazer.net/mcp`
- OpenAPI surface: `https://api.aitrailblazer.net/openapi.json`
- Strict input validation and bounded responses
- Compatible with Codex, Claude, OpenClaw, and other MCP clients
