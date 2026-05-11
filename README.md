# DeltaSignal ATLAS-7 - Native Codex MCP Server

**Real-time intelligence signals** as a **native MCP server** for Codex CLI and other MCP agents.

DeltaSignal ATLAS-7 gives your agent structured, validated tools for market signals, risk analysis, peer ranking, and fundamentals. Common workflows are exposed as composite MCP presets so agents can call one reliable tool instead of hand-orchestrating several low-level calls. No subscription required.

Also available on [Smithery](https://smithery.ai/servers/aitrailblazer/deltasignal-atlas-7) and as a [Glama Connector](https://glama.ai/mcp/connectors/net.aitrailblazer.api/delta-signal-atlas-7).

## Quick Start

### 1. Payment Setup

```bash
npx @coinbase/payments-mcp install
```

- Select `Codex CLI`
- Sign in with email
- Fund your wallet with Base USDC

Restart Codex CLI after this step.

### 2. Install

```bash
codex plugin marketplace add aitrailblazer/deltasignal-atlas-codex-plugin
```

### 3. Coinbase x402 / Bazaar Route

Use Coinbase x402 and Agentic Market on the Base route family:

- Base URL: `https://api.aitrailblazer.net`
- x402 routes: `/v1/*`
- payment: Base USDC
- seller `payTo`: `0x6D91ADF2c545047cbbC5b37a5f457cce081B48d3`
- discovery: Bazaar/CDP merchant resources for the seller `payTo`

No-payment validation before spending:

```text
GET https://api.aitrailblazer.net/v1/readiness
```

Expected payment contract:

- `402 Payment Required`
- `x402Version=2`
- `network=eip155:8453`
- `amount=40000` atomic USDC (`$0.04`)
- `extensions.bazaar.routeTemplate=/v1/readiness`

Do not claim Agentic Market first-party verification until the Agentic Market UI or Coinbase/Agentic Market team confirms the branded DeltaSignal ATLAS-7 listing.

### 4. Use

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
@DeltaSignal daily_changes
@DeltaSignal daily_change_evidence ticker:ARKB source_date:2026-05-08 limit:100
```

## Access Model

- First 5 calls are free per user.
- After the free tier, supported clients receive automatic x402 payment requirements in USDC.
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
- `deltasignal_daily_changes` - compact Daily Monitoring; no raw tag arrays; typical public route price $0.03.
- `deltasignal_daily_change_evidence` - explicit issuer proof drilldown; paginated raw Company Facts tags; typical public route price $0.03.

Credit packs are not implemented yet. Pricing metadata uses `1 credit = $0.01` of DeltaSignal usage value.

## Daily Monitoring, Evidence, and Export Packaging

DeltaSignal daily activity is packaged as three separate products:

- **Daily Monitoring** answers "what changed today?" through compact MCP and REST responses.
- **Evidence Drilldown** answers "show me why this issuer moved" through `deltasignal_daily_change_evidence` or `GET /v1/daily-changes/evidence`.
- **Bulk Export** is reserved for future artifact-backed full daily evidence packs; full exports should be saved as durable files, not pasted into chat context.

Public REST and payment surfaces:

- `GET /v1/daily-changes/latest` or `GET /mpp/v1/daily-changes/latest` - $0.03 compact monitoring.
- `GET /v1/daily-changes/evidence` or `GET /mpp/v1/daily-changes/evidence` - $0.03 issuer evidence drilldown.
- Future bulk evidence export proposal: small pack $0.15, standard pack $0.30, full daily export $0.75-$1.50.

## Development Modes

- Local: `DELTASIGNAL_PAYMENT_MODE=local`
- Live production: `DELTASIGNAL_PAYMENT_MODE=live`

## Technical

- Bundled STDIO MCP server in `mcp-stdio/`
- Remote MCP endpoint: `https://api.aitrailblazer.net/mcp`
- OpenAPI surface: `https://api.aitrailblazer.net/openapi.json`
- Strict input validation and bounded responses
- Compatible with Codex, Claude, OpenClaw, and other MCP clients
