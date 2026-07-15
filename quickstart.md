---
title: Delta Signal Natural Language Briefs Quickstart
description: Quickstart guide for calling Phase 1 StrategiX Delta Signal Natural Language Brief endpoints.
sidebar_label: Quickstart
sidebar_position: 1
slug: /quickstart
version: ATLAS7_NL_DOCS_V1
status: phase-1
product: Delta Signal Natural Language Briefs
last_updated: 2026-06-19
---

# Delta Signal Natural Language Briefs Quickstart

![Version](https://img.shields.io/badge/version-ATLAS7_NL_DOCS_V1-blue) ![Status](https://img.shields.io/badge/status-phase--1-orange) ![API](https://img.shields.io/badge/API-OpenAPI%203.1-green) ![Safety](https://img.shields.io/badge/investment%20advice-no-red)

Latest public changelog: [`changelog.html`](./changelog.html).

June 19 live contract checkpoint: production exposes OpenAPI with 48 paths, `GET /v1/pricing`, `GET /v1/contract/fields`, public MCP `initialize` and `tools/list` discovery, and public x402 challenges for protected route and MCP `tools/call` execution.

Delta Signal Natural Language Briefs convert canonical Delta Signal evidence into validated, audit-grade Markdown.

They are designed for developer workflows that need human-readable summaries without losing provenance, source dates, caveats, quality flags, null handling, or evidence boundaries.

Compiler model:

```text
Raw Delta Signal evidence
-> Canonical evidence envelope
-> Endpoint-specific Markdown compiler
-> Output validator
-> Natural Language response
```

The Natural Language layer is not a chatbot and not a new analytics source. It renders returned evidence. It does not browse, call tools from inside the renderer, infer missing values, re-rank issuers, invent catalysts, or provide investment advice.

Every successful brief ends with:

```text
Delta Signal outputs are diligence triage and risk-sizing signals, not investment advice.
```

## Standard x402 + Circle-Aware Reviewer Quickstart

Review Delta Signal ATLAS-7 as a paid evidence API for Codex workflows.

Start with discovery, not payment:

1. Use a standard EIP-3009 x402 client, or a Circle-aware buyer that can explicitly select `GatewayWalletBatched`.
2. Install the Codex plugin.
3. Confirm hosted MCP `tools/list` works before paid execution.
4. Inspect OpenAPI and Arazzo metadata.
5. Probe `GET https://api.aitrailblazer.net/v1/readiness`.
6. Confirm HTTP `402` challenge behavior when payment or grant proof is required.
7. Retry through standard `accepts[0]`, explicitly select Circle at `accepts[1]`, or use an active grant context.
8. Confirm returned JSON or Markdown preserves source dates, caveats, quality flags, evidence hashes where available, payload mode, route provenance, and the non-advice boundary.

Reviewer pass conditions:

- OpenAPI declares public route contracts and 402 behavior.
- MCP discovery is available before paid execution.
- x402 challenge includes coherent network, asset, seller, amount, and retry metadata.
- Bazaar metadata includes route template, output example, merchant resources, seller `payTo`, and payment metadata.
- No Agentic Market first-party verification claim is made until Coinbase or the relevant UI confirms it.

Primary public workflows for review:

- Morning Brief
- Company Report
- Pressure Board
- Alpha Sweep
- Quick Ticker Check
- issuer fundamentals
- covenant stress
- peer ranking
- alpha signals
- daily changes
- SPECTRA field maps

## TripCode Research Continuity Quickstart

Public marketplace caveat: TripCode / River subscriber research is currently discoverable in public MCP `tools/list`, but execution is x402/grant-gated and evidence-dependent. If a live tool call cannot resolve the required River / issuer index blobs, keep the missing state explicit and use the supported ATLAS-7 issuer workflow instead.

When a DeltaSignal article subtitle includes `ATLAS-7 TripCode: TF-SUB-...`, an MCP-enabled agent can load the research object behind the article instead of relying only on copied prose.

Recommended flow:

1. Read the article subtitle and extract the current `TF-SUB` TripCode.
2. Call MCP tool `deltasignal_resolve_article_tripcode`.
3. Inspect returned `research_payload`, `continuity`, `publication_targets`, and `azure_blob_paths`.
4. Resolve linked `TF-XBRL` evidence nodes with `deltasignal_resolve_filing_tripcode`.
5. Run `deltasignal_compare_article_to_filing_evidence` when the user asks what changed, what was confirmed, what weakened, or what to monitor next.
6. If prior article nodes are not supplied, run `deltasignal_list_article_tripcodes` with the current TripCode, River TripCode, or issuer symbol to discover prior TF-SUB nodes automatically.
7. Run `deltasignal_resolve_river_tripcode` for the full River graph, or `deltasignal_reverse_search_river` for thesis-lineage reconstruction.
8. Run `deltasignal_article_thesis_map` when the user asks for the full eight-section thesis map across the article, prior River, filing evidence, computed DeltaSignal nodes, and monitor milestones.

Typical MCP prices:

| MCP tool | Price |
| --- | ---: |
| `deltasignal_generate_article_tripcode` | `$0.00` |
| `deltasignal_resolve_article_tripcode` | `$0.02` |
| `deltasignal_list_article_tripcodes` | `$0.02` |
| `deltasignal_resolve_river_tripcode` | `$0.05` |
| `deltasignal_reverse_search_river` | `$0.30` |
| `deltasignal_search_by_claim` | `$0.05` |
| `deltasignal_search_by_issuer` | `$0.05` |
| `deltasignal_compare_claim_to_evidence` | `$0.08` |
| `deltasignal_generate_filing_tripcode` | `$0.00` |
| `deltasignal_resolve_filing_tripcode` | `$0.02` |
| `deltasignal_compare_article_to_filing_evidence` | `$0.08` |
| `deltasignal_article_thesis_map` | `$0.30` |

Do not use Substack post ID, URL, slug, or publication timestamp as canonical TF-SUB identity. Those are publication metadata only. If a resolver returns missing evidence, keep it as missing evidence.

## Base URL

Use your configured StrategiX API base URL:

```text
https://api.example.com
```

Replace this with the production or staging base URL supplied for your environment.

## Authentication And Access

TestFlight and production clients authenticate through the StrategiX backend relay.

The backend holds the internal Delta Signal credential and may call the MCP/x402-capable Delta Signal surface. Client applications should not embed internal Delta Signal credentials.

Public paid routes may use x402 or MPP access depending on deployment. Payment and settlement flows are handled by the backend route infrastructure, not by the StrategiX app. Public buyers should use the x402 `/v1/*` routes.

Payment failure returns:

```text
402 Payment Required
```

Standard x402, Circle-aware buyers, and Agentic Market use the Base route family under `https://api.aitrailblazer.net/v1/*`. Public clients do not need Delta Signal API keys. Before spending USDC, run an x402 discovery probe against `GET /v1/readiness`. A database-managed grant may allow execution when the caller completes the wallet grant-session flow for an enrolled EVM wallet: `GET /v1/grant/challenge`, sign the returned message, `POST /v1/grant/session`, then attach `X-DeltaSignal-Grant-Token`. Raw unsigned wallet headers are not a grant identity. Otherwise expect HTTP `402` and confirm:

- `x402Version=2`
- `network=eip155:8453`
- Base USDC asset `0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913`
- `amount=40000` atomic USDC (`$0.04`)
- seller `payTo=0x6D91ADF2c545047cbbC5b37a5f457cce081B48d3`
- Bazaar metadata includes `extensions.bazaar.routeTemplate=/v1/readiness`
- `accepts[0]` is standard Base x402 for first-entry client compatibility
- `accepts[1]` is Circle `GatewayWalletBatched` when Circle settlement is enabled

Circle-aware buyers must select the Circle requirement explicitly and sign against its Gateway EIP-712 domain using the challenge-provided `extra.verifyingContract`. The signing wallet must have deposited or available Circle Gateway balance. Base ERC-20 USDC alone may return `insufficient_balance`, and that response must not be treated as a charge.

Do not claim Agentic Market first-party verification until the Agentic Market UI or Coinbase/Agentic Market team confirms the branded DeltaSignal ATLAS-7 listing.

## Phase 1 Endpoints

### Top Stressed Natural Language Brief

Returns a compact Markdown brief for the highest-stress issuers returned by Delta Signal.

- Route: `GET /v1/top-stressed/natural`
- MPP mirror: `GET /mpp/v1/top-stressed/natural`
- MCP tool: `deltasignal_top_stressed_natural`
- Price: `$0.95`
- Ticker required: no

Query parameters:

- `style`: `professional`, `concise`, `trader`, or `detailed`
- `limit`: integer, default `10`
- `offset`: integer, default `0`
- `date`: optional `YYYY-MM-DD` selector where supported

Example request:

```bash
curl -sS "$BASE_URL/v1/top-stressed/natural?style=professional&limit=10&offset=0" \
  -H "Authorization: Bearer $STRATEGIX_API_TOKEN"
```

Expected Markdown sections:

```text
Delta Signal Top Stressed Brief

Readiness Summary

Top Stressed Issuers

Main Pattern

Evidence Limits

Suggested Drilldowns
```

### Covenant Stress Natural Language Brief

Returns a ticker-specific Markdown brief for covenant, leverage, liquidity, filing, and issuer stress evidence.

- Route: `GET /v1/covenant-stress/natural`
- MPP mirror: `GET /mpp/v1/covenant-stress/natural`
- MCP tool: `deltasignal_covenant_stress_natural`
- Price: `$1.20`
- Ticker required: yes

Query parameters:

- `ticker`: required issuer ticker
- `style`: `professional`, `concise`, `trader`, or `detailed`
- `date`: optional `YYYY-MM-DD` selector where supported

Example request:

```bash
curl -sS "$BASE_URL/v1/covenant-stress/natural?ticker=RIOT&style=professional" \
  -H "Authorization: Bearer $STRATEGIX_API_TOKEN"
```

Missing ticker response:

```json
{
  "status": "error",
  "code": "ticker_required",
  "message": "deltasignal_covenant_stress_natural requires ticker.",
  "endpoint": "deltasignal_covenant_stress_natural"
}
```

Expected Markdown sections:

```text
Delta Signal Covenant Stress Brief

Issuer Summary

Covenant Pressure Indicators

Filing Evidence

Evidence Limits

Suggested Drilldowns
```

Important constraints:

- Do not substitute a default ticker.
- Do not infer covenant breach unless returned.
- Do not infer default risk unless returned.
- Do not infer insolvency, bankruptcy risk, liquidity crisis, price movement, or trade direction unless returned.
- Do not treat missing covenant data as no covenant risk.

### Morning Brief Natural Language Brief

Returns a backend-composed compact market-wide Markdown brief.

- Route: `GET /v1/morning-brief/natural`
- MPP mirror: `GET /mpp/v1/morning-brief/natural`
- MCP tool: `deltasignal_morning_brief_natural`
- Price: `$0.95`
- Ticker required: no

Query parameters:

- `style`: `professional`, `concise`, `trader`, or `detailed`
- `date`: optional `YYYY-MM-DD` selector where supported
- `limit`: optional compact record limit where supported

Example request:

```bash
curl -sS "$BASE_URL/v1/morning-brief/natural?style=professional" \
  -H "Authorization: Bearer $STRATEGIX_API_TOKEN"
```

Expected Markdown sections:

```text
Delta Signal Morning Brief

Readiness Summary

What Changed

Market-Wide Risk Board

Top Stressed Issuers

Alpha Opportunity Board

Main Pattern

Evidence Limits

Suggested Drilldowns
```

Important constraints:

- Morning Brief uses backend-composed compact evidence.
- Do not perform client-side fan-out.
- Do not call tools from inside the renderer.
- Do not generate X drafts, Substack notes, article candidates, or issuer drilldowns unless the user asks after evidence is returned.
- If daily changes are not returned, say they were not returned.
- If alpha opportunities are not returned, say they were not returned.
- Do not interpret missing evidence as no changes, no opportunities, or no risk.

## Styles

Supported styles:

- `professional`
- `concise`
- `trader`
- `detailed`

Style changes tone and density only. It must not change facts, rankings, caveats, risk interpretation, or investment posture.

Recommended behavior:

- `professional`: default audit-grade summary
- `concise`: shorter summary with the same evidence boundary
- `trader`: direct risk-board framing without trade instructions
- `detailed`: expanded evidence limits and drilldown context

Do not use style to create investment advice.

## Successful Response Shape

All successful Natural Language responses return JSON with a validated Markdown brief and metadata.

Example:

```json
{
  "natural_brief": "# Delta Signal Top Stressed Brief\n\n## Readiness Summary\n\nReturned compact top-stressed evidence...\n\nDelta Signal outputs are diligence triage and risk-sizing signals, not investment advice.",
  "summary_data": {
    "returned_count": 10
  },
  "raw_data_ref": "atlas7://deltasignal/top_stressed/latest?limit=10",
  "cost_usd": 0.95,
  "tokens_used": null,
  "source_date": "2026-05-09",
  "computed_at": "2026-05-09T14:30:00Z",
  "payload_mode": "natural_language",
  "raw_payload_mode": "compact",
  "quality_flags": [],
  "caveats": ["compact payload; full raw artifact arrays may be omitted"],
  "model": "deterministic-atlas7-nl-renderer",
  "renderer_agent": "deterministic-fallback",
  "renderer_mode": "deterministic",
  "validation_status": "passed",
  "fallback_reason": null,
  "prompt_version": "ATLAS7_NL_RENDERER_V1",
  "renderer_version": "ATLAS7_NL_COMPILER_V1",
  "evidence_hash": "sha256...",
  "request_id": "req_...",
  "latency_ms": 184,
  "cache_hit": false,
  "billing_mode": "x402",
  "free_tier_applied": false
}
```

Important response fields:

- `natural_brief`: the validated Markdown brief.
- `summary_data`: compact structured summary values returned alongside the Markdown.
- `raw_data_ref`: reference to raw evidence or raw artifact location when available.
- `cost_usd`: endpoint price for the Natural Language response unless a signed grant session or other live access policy applies.
- `payload_mode`: always `natural_language` for successful Natural Language responses.
- `raw_payload_mode`: the evidence mode used to produce the brief, such as `compact`, `full`, or `evidence_page`.
- `renderer_mode`: `deterministic`, `apex`, or `deterministic_fallback`.
- `validation_status`: `passed`, `failed`, or `skipped`.
- `fallback_reason`: reason deterministic fallback was used, if applicable.
- `evidence_hash`: hash of the canonical evidence envelope when available.

## Renderer Modes

`deterministic`: evidence-preserving baseline renderer. Used when APEX is disabled, during baseline deployment, or when maximum output stability is required.

`apex`: GEN-I-GENESIS-APEX renderer. Used when APEX is enabled and output validation passes.

`deterministic_fallback`: fallback renderer. Used when APEX times out, errors, returns malformed Markdown, fails evidence validation, includes unsafe financial language, or omits required sections or disclaimer.

Example fallback metadata:

```json
{
  "renderer_mode": "deterministic_fallback",
  "validation_status": "failed",
  "fallback_reason": "apex_validation_failed"
}
```

## Error Responses

Unsupported style:

```json
{
  "status": "error",
  "code": "unsupported_style",
  "message": "style must be one of: professional, concise, trader, detailed."
}
```

Invalid date:

```json
{
  "status": "error",
  "code": "invalid_date",
  "message": "date must use YYYY-MM-DD format."
}
```

Payment required:

```json
{
  "status": "error",
  "code": "payment_required",
  "message": "Payment is required to access this Natural Language endpoint."
}
```

Raw evidence missing:

```json
{
  "status": "error",
  "code": "raw_evidence_missing",
  "message": "Raw Delta Signal evidence was not available for the requested selector."
}
```

Evidence service unavailable:

```json
{
  "status": "error",
  "code": "evidence_service_unavailable",
  "message": "Delta Signal evidence is temporarily unavailable. Retry later."
}
```

## Validation Model

Natural Language responses are validated before being returned.

The validator rejects output that:

- Adds unsupported tickers or issuers.
- Changes ranks, scores, dates, caveats, quality flags, filing periods, accession numbers, or null semantics.
- Converts missing values into factual claims.
- Describes compact evidence as full raw evidence.
- Adds external news.
- Invents causality.
- Includes buy, sell, short, long, hold, price-target, entry, exit, stop-loss, guaranteed-move, or trade-recommendation language.
- Omits required Markdown sections.
- Omits or changes the required disclaimer.

If APEX output fails validation, the backend may return deterministic fallback when possible.

## Caching

Natural Language outputs may be cached after validation.

Recommended cache-key fields:

- `endpoint`
- `normalized_request`
- `source_date`
- `style`
- `limit`
- `ticker`, where applicable
- `evidence_hash`
- `renderer_mode`
- `renderer_version`

Cache rules:

- Evidence hash changes force cache miss.
- Renderer version changes force cache miss.
- Style changes force cache miss.
- Source date changes force cache miss.
- APEX output is cached only after validation passes.
- Failed APEX output is not cached as a valid response.
- Cache TTL must not exceed source-data freshness.

## SDK Integration Notes

For client applications:

- Authenticate to the StrategiX backend relay.
- Call the appropriate `*_natural` route.
- Display `natural_brief` as Markdown.
- Surface `source_date`, `renderer_mode`, `validation_status`, and `evidence_hash` in debug or provenance views.
- Treat `deterministic_fallback` as a valid response mode, not as a client failure.
- Do not embed internal Delta Signal credentials in the client.
- Do not execute wallet, settlement, or payment flows inside the app.

### Minimal JavaScript Example

```js
const baseUrl = process.env.STRATEGIX_BASE_URL ?? "https://api.aitrailblazer.net";
const token = process.env.STRATEGIX_API_TOKEN;

async function getTopStressedBrief() {
  const url = new URL("/v1/top-stressed/natural", baseUrl);
  url.searchParams.set("style", "professional");
  url.searchParams.set("limit", "10");

  const response = await fetch(url, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });

  if (!response.ok) {
    throw new Error(`Request failed: ${response.status}`);
  }

  const data = await response.json();

  return {
    markdown: data.natural_brief,
    sourceDate: data.source_date,
    rendererMode: data.renderer_mode,
    validationStatus: data.validation_status,
    evidenceHash: data.evidence_hash,
  };
}
```

### Minimal Swift Example

```swift
import Foundation

struct NaturalLanguageBriefResponse: Decodable {
    let naturalBrief: String
    let sourceDate: String?
    let rendererMode: String
    let validationStatus: String
    let evidenceHash: String?

    enum CodingKeys: String, CodingKey {
        case naturalBrief = "natural_brief"
        case sourceDate = "source_date"
        case rendererMode = "renderer_mode"
        case validationStatus = "validation_status"
        case evidenceHash = "evidence_hash"
    }
}

func fetchTopStressedBrief(baseURL: URL, token: String) async throws -> NaturalLanguageBriefResponse {
    var components = URLComponents(
        url: baseURL.appendingPathComponent("v1/top-stressed/natural"),
        resolvingAgainstBaseURL: false
    )!

    components.queryItems = [
        URLQueryItem(name: "style", value: "professional"),
        URLQueryItem(name: "limit", value: "10")
    ]

    var request = URLRequest(url: components.url!)
    request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

    let (data, response) = try await URLSession.shared.data(for: request)

    guard let http = response as? HTTPURLResponse, 200..<300 ~= http.statusCode else {
        throw URLError(.badServerResponse)
    }

    return try JSONDecoder().decode(NaturalLanguageBriefResponse.self, from: data)
}
```

## Operational Checklist

Before enabling a new Natural Language endpoint in production:

- Endpoint contract exists.
- Route is registered.
- MCP tool is registered.
- MPP route is registered where applicable.
- OpenAPI operation exists.
- Price metadata exists.
- Request validation exists.
- Evidence builder exists.
- Deterministic fallback exists.
- APEX path is feature-flagged.
- Markdown grammar tests exist.
- Exact disclaimer test exists.
- Null preservation test exists.
- Caveat preservation test exists.
- Quality flag preservation test exists.
- Evidence hash preservation test exists.
- No unsupported ticker test exists.
- No changed rank, score, or date test exists.
- No investment advice language test exists.
- Authenticated route success test exists.
- Unpaid x402 route returns HTTP `402`.
- MCP `tools/list` includes the endpoint.
- MCP `tools/call` returns valid ABI output.
- `gofmt` passes.
- No-secrets check passes.
- Deploy audit passes.

## Non-Advice Policy

Delta Signal Natural Language Briefs are diligence triage and risk-sizing summaries. They are not personalized financial advice.

Do not generate or display them as:

- Buy recommendations.
- Sell recommendations.
- Short recommendations.
- Long recommendations.
- Hold recommendations.
- Entry levels.
- Exit levels.
- Stop losses.
- Price targets.
- Guaranteed outcomes.
- High-conviction trade calls.

Use them as evidence-preserving summaries for review, triage, and follow-up drilldowns.
