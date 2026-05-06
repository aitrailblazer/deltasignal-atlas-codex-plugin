# DeltaSignal ATLAS-7 - Native Codex MCP Server

[![deltasignal-atlas-codex-plugin MCP server](https://glama.ai/mcp/servers/aitrailblazer/deltasignal-atlas-codex-plugin/badges/score.svg)](https://glama.ai/mcp/servers/aitrailblazer/deltasignal-atlas-codex-plugin)
[![deltasignal-atlas-codex-plugin MCP server](https://glama.ai/mcp/servers/aitrailblazer/deltasignal-atlas-codex-plugin/badges/card.svg)](https://glama.ai/mcp/servers/aitrailblazer/deltasignal-atlas-codex-plugin)

**Real-time intelligence signals** as a **native MCP server** for Codex CLI and other MCP agents.

DeltaSignal ATLAS-7 gives your agent structured, validated tools for market signals, risk analysis, peer ranking, and fundamentals, with automatic x402 micropayments in USDC. No subscription required.

Also available on [Smithery](https://smithery.ai/servers/aitrailblazer/deltasignal-atlas-7).

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
- Glama release container entrypoint in `Dockerfile`
- OpenAPI surface: `https://api.aitrailblazer.net/openapi.json`
- Strict input validation and bounded responses
- Compatible with Codex, Claude, OpenClaw, and other MCP clients
