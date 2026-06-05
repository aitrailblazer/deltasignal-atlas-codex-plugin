---
name: deltasignal-atlas-7
description: Use DeltaSignal ATLAS-7 for SEC-grounded issuer intelligence on crypto public companies, including composite MCP presets, readiness, covenant stress, top-stressed issuers, peer ranking, alpha signals, fundamentals, compact daily SEC monitoring, paginated evidence drilldowns, and Azure-native audit status over the free starter tier, Tempo MPP, or Base x402.
---

# DeltaSignal ATLAS-7

Use this skill when the user asks for DeltaSignal, ATLAS-7, crypto public issuer stress, SEC-grounded crypto equity intelligence, covenant stress, top-stressed issuers, issuer fundamentals, or DeltaSignal paid API access.

## Access Model

DeltaSignal is a usage-based intelligence API with a free starter tier. Do not imply there is a subscription requirement.

Free tier:

- The Base x402 `/v1/*` routes allow 5 free calls per user before payment fallback.
- If the free tier is exhausted, the same route returns x402 payment requirements for compatible clients.
- If a caller includes an x402 payment header, treat it as an intentional paid request rather than consuming free quota.

Primary rail:

- Tempo MPP routes under `https://api.aitrailblazer.net/mpp/v1/*`

Agentic Market compatibility rail:

- Base x402 routes under `https://api.aitrailblazer.net/v1/*`
- Base network: `eip155:8453`
- Payment asset: USDC

Bundled MCP server:

- Treat DeltaSignal as an MCP-first plugin. Prefer named MCP tool calls over hand-written HTTP routes whenever the plugin is installed.
- Public MCP `initialize` and `tools/list` are free discovery methods. Public MCP `tools/call` remains protected by x402 payment unless a free-tier grant applies.
- Prefer composite MCP presets for common workflows:
  - `deltasignal_morning_brief` for daily scans.
  - `deltasignal_company_report` for full ticker reports.
  - `deltasignal_pressure_board` for risk monitoring.
  - `deltasignal_alpha_sweep` for opportunity screening.
  - `deltasignal_quick_ticker_check` for fast ticker checks.
  - `deltasignal_daily_changes` for compact Daily Monitoring.
  - `deltasignal_daily_change_evidence` for explicit issuer proof after a monitoring result.
  - `deltasignal_atlas7_audit_status` for operator checks that the Azure-native 215-issuer regression audit is healthy and fresh.
- Use granular tools only for custom drilldowns or when a composite preset is unavailable.
- Prefer the `deltasignal_*` MCP tools when this plugin is installed and the request maps to a supported DeltaSignal route.
- MCP tools are read-only, idempotent, closed-world tools with strict argument validation.
- In live mode, the MCP wrapper calls the deployed Base x402 `/v1/*` routes and returns free-tier headers or x402 payment requirements in the tool result.
- In internal mode, the MCP wrapper uses a pre-authorized `DELTASIGNAL_API_KEY` against the deployed API for no-payment smoke testing.
- In local mode, the MCP wrapper calls `http://127.0.0.1:8080/v1/*` without x402 payment so DeltaSignal can be used during development work.
- If the MCP wrapper returns `payment_required`, do not claim payment succeeded. Tell the user the free tier is exhausted or unavailable and switch to an x402/MPP-capable payment flow.

Do not translate between route families unless the user explicitly switches payment rails. If the active client is Agentic Market or Coinbase x402, use `/v1/*`. If the active client is MPPScan, mpp.dev, or asks for Tempo MPP, use `/mpp/v1/*`.

## Composite Presets

Composite presets are discounted server-enforced workflows. Public Builder pricing uses `1 credit = $0.01` of DeltaSignal usage value. Credit packs are not implemented yet.

| MCP tool | Builder price | Internal calls |
| --- | ---: | --- |
| `deltasignal_morning_brief` | 18 credits / `$0.18` | readiness, daily_changes, risk_distribution, top_stressed(limit=10), alpha_opportunities(limit=10) |
| `deltasignal_company_report` | 30 credits / `$0.30` | readiness, company_fundamentals(ticker), alpha_signals(ticker), peer_ranking(ticker), covenant_stress(ticker) |
| `deltasignal_pressure_board` | 14 credits / `$0.14` | readiness, top_stressed(limit=15), risk_distribution |
| `deltasignal_alpha_sweep` | 14 credits / `$0.14` | readiness, alpha_opportunities(limit=15), daily_changes |
| `deltasignal_quick_ticker_check` | 18 credits / `$0.18` | readiness, covenant_stress(ticker), alpha_signals(ticker) |

Daily Monitoring and Evidence are separate products. `deltasignal_daily_changes` is compact by default and should not be used as a bulk evidence export. Use `deltasignal_daily_change_evidence` only when the user asks to inspect raw proof for a named issuer or CIK.

## Operator Audit Status

- Use `deltasignal_atlas7_audit_status` when the user asks whether ATLAS-7 is healthy, current, monitored, Azure-native, or fully regression-tested.
- Treat it as an operational readiness/audit surface, not as issuer intelligence or investment evidence.
- Preserve `status`, `stale`, `artifact.prefix`, `finished_at_utc`, `issuer_count`, `operation_count`, `failed_count`, `historical_failed_count`, and `composite_failed_count` when summarizing.
- A healthy result currently means the latest scheduled Go audit covered the 215-issuer universe, current routes, historical routes, and composite MCP workflows with zero failures inside the freshness window.
- If the route returns unauthorized from the bundled STDIO wrapper, explain that this operator endpoint requires internal/pre-authorized access and tell the user to use the hosted MCP/tooling path that already has the DeltaSignal credential.

## Paid Routes

Use readiness before higher-cost calls when freshness matters.

| Route | Typical price | Purpose |
| --- | ---: | --- |
| `GET /mpp/v1/readiness` or `GET /v1/readiness` | `$0.04` | Latest service/data readiness and coverage snapshot |
| `GET /mpp/v1/daily-changes/latest` or `GET /v1/daily-changes/latest` | `$0.03` | Compact Daily Monitoring: freshness, counts, compact changed-company rows, and evidence refs |
| `GET /mpp/v1/daily-changes/evidence` or `GET /v1/daily-changes/evidence` | `$0.03` | Explicit issuer evidence drilldown with paginated raw Company Facts tags |
| `GET /mpp/v1/risk-distribution` or `GET /v1/risk-distribution` | `$0.04` | Current risk-tier distribution |
| `GET /mpp/v1/top-stressed` or `GET /v1/top-stressed` | `$0.05` | Most stressed crypto public issuers |
| `GET /mpp/v1/peer-ranking/{ticker}` or `GET /v1/peer-ranking/{ticker}` | `$0.06` | Peer covenant ranking for one issuer |
| `GET /mpp/v1/alpha-signals/{ticker}` or `GET /v1/alpha-signals/{ticker}` | `$0.07` | Phase 1 alpha, resilience, treasury, and regime-fit signals |
| `GET /mpp/v1/covenant-stress` or `GET /v1/covenant-stress` | `$0.08` | Filter or list the active covenant stress slice |
| `GET /mpp/v1/spectra-field-map/{ticker}` or `GET /v1/spectra-field-map/{ticker}` | `$0.08` | Historical field-map pressure and filing choreography |
| `GET /mpp/v1/covenant-stress/{ticker}` or `GET /v1/covenant-stress/{ticker}` | `$0.10` | Detailed ATLAS-7 covenant stress for one issuer |

Supported list filters include `limit`, `offset`, `risk_tier`, `min_stress`, `linkbase_only`, and debt coverage status filters where exposed by the route.

Future artifact-backed bulk daily evidence exports are not inline chat responses. Proposed prices: small pack $0.15, standard pack $0.30, full daily evidence export $0.75-$1.50.

## Operating Rules

- Confirm whether the request should use remaining free calls or a paid rail before initiating payment.
- Confirm the route family and approximate cost before initiating any paid request.
- Never invent a payment result. Report the route, status, payer/payment metadata if available, and the response summary.
- Treat DeltaSignal outputs as issuer intelligence, not investment advice.
- Preserve source dates, row counts, readiness fields, risk tiers, and debt coverage status when summarizing.
- For named issuers, normalize tickers to uppercase before constructing routes.
- For large endpoints such as company-report and daily-changes, do not pipe `npx agentcash fetch` stdout directly into another process; write the response to a file or use compact/paginated routes to avoid Node/pipe truncation around 64 KB.
- If payment tooling is unavailable, give the exact route and expected cost, then explain that payment execution needs an x402 or MPP-capable client.

## Common Prompts

- "Give me a DeltaSignal morning brief" -> call `deltasignal_morning_brief`.
- "Run a company report for RIOT" -> call `deltasignal_company_report` with `ticker=RIOT`.
- "Show the pressure board" -> call `deltasignal_pressure_board`.
- "Find alpha opportunities" -> call `deltasignal_alpha_sweep`.
- "What changed today?" -> call `deltasignal_daily_changes`.
- "Show me why ARKB moved" -> call `deltasignal_daily_change_evidence` with `ticker=ARKB` and the source date from Daily Monitoring when available.
- "Is the full ATLAS-7 audit healthy?" -> call `deltasignal_atlas7_audit_status`.
- "Quick check MARA" -> call `deltasignal_quick_ticker_check` with `ticker=MARA`.
- "Is DeltaSignal current?" -> call readiness.
- "Which issuers are most stressed?" -> call top-stressed.
- "Run covenant stress on MARA" -> call covenant-stress detail for `MARA`.
- "Show high-risk linkbase-backed issuers" -> call covenant-stress list with `risk_tier=HIGH`, `min_stress`, and `linkbase_only=true`.
- "Compare MARA to peers" -> call peer-ranking for `MARA`.
