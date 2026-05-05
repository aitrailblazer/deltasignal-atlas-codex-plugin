# DeltaSignal ATLAS-7 - Native Codex MCP Server

**Real-time intelligence signals** as a **native MCP server** for Codex CLI and other MCP agents.

DeltaSignal ATLAS-7 gives your agent structured, validated tools for market signals, risk analysis, peer ranking, and fundamentals, with automatic x402 micropayments in USDC. No subscription required.

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
@DeltaSignal alpha_signals ticker:NVDA
@DeltaSignal top_stressed tickers:AAPL,TSLA
@DeltaSignal company_fundamentals ticker:MSFT
```

## Access Model

- First 5 calls are free per user.
- After the free tier, supported clients receive automatic x402 payment requirements in USDC.
- Internal testing can use a pre-authorized API key without payment.
- All tools are read-only, schema-validated, and closed-world.

## Available MCP Tools

- `deltasignal_readiness`
- `deltasignal_top_stressed`
- `deltasignal_covenant_stress`
- `deltasignal_peer_ranking`
- `deltasignal_alpha_signals`
- `deltasignal_company_fundamentals`
- `deltasignal_risk_distribution`
- `deltasignal_daily_changes`

## Development Modes

- Local: `DELTASIGNAL_PAYMENT_MODE=local`
- Internal testing: `DELTASIGNAL_PAYMENT_MODE=internal` plus `DELTASIGNAL_API_KEY`
- Live production: `DELTASIGNAL_PAYMENT_MODE=live`

## Technical

- Bundled STDIO MCP server in `mcp-stdio/`
- Strict input validation and bounded responses
- Compatible with Codex, Claude, OpenClaw, and other MCP clients
