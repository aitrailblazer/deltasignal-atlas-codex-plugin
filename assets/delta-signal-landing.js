import { LitElement, css, html } from "https://cdn.jsdelivr.net/npm/lit@3/+esm";

const rawRoutes = [
  ["morning_brief", "$0.18", "Daily scan: readiness, daily changes, risk distribution, top stressed, and alpha opportunities."],
  ["company_report", "$0.30", "Full issuer report across fundamentals, alpha signals, peer ranking, and covenant stress."],
  ["pressure_board", "$0.14", "Risk board with top stressed issuers and risk distribution."],
  ["alpha_sweep", "$0.14", "Opportunity screen with alpha opportunities and daily changes."],
  ["quick_ticker_check", "$0.18", "Fast issuer triage using covenant stress and alpha signals."],
  ["daily_changes", "$0.03", "Compact daily monitoring without raw tag arrays."],
];

const naturalRoutes = [
  ["top_stressed_natural", "$0.95", "Live evidence-preserving Markdown brief for highest-stress issuers."],
  ["morning_brief_natural", "$1.80", "Phase 1 backend-composed market-wide Natural Language brief."],
  ["covenant_stress_natural", "$1.20", "Phase 1 ticker-specific covenant, leverage, liquidity, and filing brief."],
  ["peer_ranking_natural", "$1.10", "Planned Phase 2 peer-positioning brief."],
  ["alpha_signals_natural", "$1.30", "Planned Phase 2 alpha signal brief with strict non-advice validation."],
  ["company_fundamentals_natural", "$1.40", "Planned Phase 2 fundamentals brief over SEC/XBRL evidence."],
];

const installTabs = {
  smithery: {
    label: "Smithery",
    code: "smithery mcp add aitrailblazer/deltasignal-atlas-7",
  },
  codex: {
    label: "Codex Plugin",
    code: "codex plugin install aitrailblazer/deltasignal-atlas-codex-plugin",
  },
  marketplace: {
    label: "Marketplace",
    code: "npx codex-marketplace add aitrailblazer/deltasignal-atlas-codex-plugin --plugin --project",
  },
  mcp: {
    label: "Direct MCP",
    code: `{
  "mcpServers": {
    "deltasignal-atlas-7": {
      "url": "https://api.aitrailblazer.net/mcp"
    }
  }
}`,
  },
  openapi: {
    label: "Open API",
    code: "https://api.aitrailblazer.net/openapi.json",
  },
};

const x402Steps = [
  "Use Coinbase Wallet, Coinbase Smart Wallet, or another Base-compatible wallet.",
  "Switch the wallet network to Base.",
  "Add USDC on Base.",
  "Keep a small amount of ETH on Base if your wallet requires gas.",
  "Re-run the agent query after funding.",
];

class DSCopyPrompt extends LitElement {
  static properties = {
    prompt: { type: String },
    label: { type: String },
    copied: { state: true },
  };

  static styles = css`
    :host { display: block; }
    .card {
      min-height: 116px;
      display: grid;
      align-content: space-between;
      gap: 1rem;
      padding: 1.25rem;
      border: 1px solid rgba(255, 117, 34, 0.24);
      border-radius: 8px;
      background: linear-gradient(135deg, rgba(255,154,26,.09), rgba(227,0,85,.08)), #11100f;
      color: #ffe8d9;
      box-shadow: 0 24px 60px rgba(0, 0, 0, 0.22);
    }
    p {
      margin: 0;
      font: 600 1.02rem/1.45 Inter, system-ui, sans-serif;
    }
    button {
      justify-self: start;
      border: 1px solid rgba(255, 154, 26, 0.42);
      border-radius: 999px;
      background: rgba(255, 154, 26, 0.1);
      color: #ffd1a3;
      padding: 0.55rem 0.8rem;
      font: 800 0.78rem/1 Inter, system-ui, sans-serif;
      cursor: pointer;
    }
    button:hover { background: rgba(255, 154, 26, 0.18); }
    .copied { color: #86efac; font-size: 0.78rem; margin-left: 0.55rem; }
  `;

  constructor() {
    super();
    this.prompt = "";
    this.label = "Copy prompt";
    this.copied = false;
  }

  async copyPrompt() {
    const text = this.prompt || this.textContent?.trim() || "";
    if (!text) return;

    try {
      await navigator.clipboard.writeText(text);
      this.copied = true;
      window.setTimeout(() => {
        this.copied = false;
      }, 1400);
    } catch {
      this.dispatchEvent(new CustomEvent("copy-unavailable", { bubbles: true, composed: true }));
    }
  }

  render() {
    const text = this.prompt || this.textContent?.trim() || "";
    return html`
      <div class="card">
        <p>"${text}"</p>
        <div>
          <button type="button" @click=${this.copyPrompt}>${this.copied ? "Copied" : this.label}</button>
          ${this.copied ? html`<span class="copied">Ready for your agent.</span>` : null}
        </div>
      </div>
    `;
  }
}

class DSInstallTabs extends LitElement {
  static properties = {
    active: { state: true },
    copied: { state: true },
  };

  static styles = css`
    :host { display: block; }
    .tabs { display: flex; flex-wrap: wrap; gap: 0.55rem; margin-bottom: 0.85rem; }
    button {
      border: 1px solid rgba(255, 154, 26, 0.28);
      border-radius: 999px;
      background: rgba(17, 16, 15, 0.88);
      color: #ffe4cf;
      padding: 0.58rem 0.85rem;
      cursor: pointer;
      font: 800 0.84rem/1 Inter, system-ui, sans-serif;
    }
    button[aria-selected="true"], .copy {
      color: #170806;
      background: linear-gradient(90deg, #ff9a1a, #ff3028 65%, #e30055);
      border-color: transparent;
    }
    pre {
      margin: 0;
      padding: 1.3rem;
      overflow-x: auto;
      border: 1px solid rgba(255, 117, 34, 0.24);
      border-radius: 8px;
      background: #0c0a09;
      color: #ffe4cf;
      font: 0.98rem/1.8 ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace;
      white-space: pre-wrap;
    }
    .copy { margin-top: 0.8rem; }
  `;

  constructor() {
    super();
    this.active = "smithery";
    this.copied = false;
  }

  async copy() {
    try {
      await navigator.clipboard.writeText(installTabs[this.active].code);
      this.copied = true;
      window.setTimeout(() => {
        this.copied = false;
      }, 1400);
    } catch {
      this.dispatchEvent(new CustomEvent("copy-unavailable", { bubbles: true, composed: true }));
    }
  }

  render() {
    const current = installTabs[this.active];
    return html`
      <div class="tabs" role="tablist" aria-label="Install options">
        ${Object.entries(installTabs).map(([key, tab]) => html`
          <button type="button" role="tab" aria-selected=${this.active === key} @click=${() => { this.active = key; }}>
            ${tab.label}
          </button>
        `)}
      </div>
      <pre><code>${current.code}</code></pre>
      <button class="copy" type="button" @click=${this.copy}>${this.copied ? "Copied" : "Copy"}</button>
    `;
  }
}

class DSPricingToggle extends LitElement {
  static properties = {
    tier: { state: true },
  };

  static styles = css`
    :host { display: block; }
    .shell {
      border: 1px solid rgba(255, 117, 34, 0.24);
      border-radius: 10px;
      padding: 1.15rem;
      background: rgba(12, 10, 9, 0.94);
      color: #ffe4cf;
    }
    .tabs { display: flex; flex-wrap: wrap; gap: 0.55rem; margin-bottom: 1rem; }
    button {
      border: 1px solid rgba(255, 154, 26, 0.32);
      border-radius: 999px;
      background: transparent;
      color: #ffe4cf;
      padding: 0.58rem 0.85rem;
      cursor: pointer;
      font-weight: 800;
    }
    button[aria-pressed="true"] {
      color: #170806;
      background: linear-gradient(90deg, #ff9a1a, #ff3028 65%, #e30055);
      border-color: transparent;
    }
    .row {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 1rem;
      padding: 0.85rem 0;
      border-top: 1px solid rgba(255, 117, 34, 0.16);
    }
    .name { font-weight: 900; color: #fff6ee; }
    .desc { margin-top: 0.25rem; color: #cbb9ae; font-size: 0.92rem; }
    .price { color: #ffca76; font-weight: 900; white-space: nowrap; }
    .note { margin: 1rem 0 0; color: #cbb9ae; font-size: 0.92rem; }
    @media (max-width: 640px) {
      .row { grid-template-columns: 1fr; }
    }
  `;

  constructor() {
    super();
    this.tier = "raw";
  }

  render() {
    const rows = this.tier === "raw" ? rawRoutes : naturalRoutes;
    return html`
      <section class="shell" aria-label="Delta Signal pricing tiers">
        <div class="tabs">
          <button type="button" aria-pressed=${this.tier === "raw"} @click=${() => { this.tier = "raw"; }}>Raw + Composite Routes</button>
          <button type="button" aria-pressed=${this.tier === "natural"} @click=${() => { this.tier = "natural"; }}>Natural Language Routes</button>
        </div>
        ${rows.map(([name, price, desc]) => html`
          <div class="row">
            <div>
              <div class="name">${name}</div>
              <div class="desc">${desc}</div>
            </div>
            <div class="price">${price}</div>
          </div>
        `)}
        <p class="note">Natural Language routes are priced separately because they compile evidence into validated Markdown with provenance metadata and non-advice disclosure.</p>
      </section>
    `;
  }
}

class DSX402Checklist extends LitElement {
  static properties = {
    checked: { state: true },
  };

  static styles = css`
    :host { display: block; }
    .box {
      border: 1px solid rgba(255, 117, 34, 0.24);
      border-radius: 8px;
      padding: 1.2rem;
      background: #11100f;
      color: #ffe4cf;
    }
    label {
      display: flex;
      align-items: flex-start;
      gap: 0.65rem;
      padding: 0.65rem 0;
      color: #cbb9ae;
      cursor: pointer;
    }
    input { margin-top: 0.25rem; accent-color: #ff9a1a; }
    .note { margin: 1rem 0 0; color: #cbb9ae; font-size: 0.92rem; }
  `;

  constructor() {
    super();
    this.checked = new Set();
  }

  toggle(index) {
    const next = new Set(this.checked);
    if (next.has(index)) next.delete(index);
    else next.add(index);
    this.checked = next;
  }

  render() {
    return html`
      <div class="box">
        ${x402Steps.map((step, index) => html`
          <label>
            <input type="checkbox" .checked=${this.checked.has(index)} @change=${() => this.toggle(index)}>
            <span>${step}</span>
          </label>
        `)}
        <p class="note">StrategiX app clients authenticate through the backend relay. Wallet, settlement, and payment flows should not execute inside the app.</p>
      </div>
    `;
  }
}

class DSRouteExample extends LitElement {
  static properties = {
    base: { type: String },
    path: { type: String },
    ticker: { state: true },
    limit: { state: true },
    style: { state: true },
  };

  static styles = css`
    :host { display: block; }
    .box {
      border: 1px solid rgba(255, 117, 34, 0.24);
      border-radius: 8px;
      padding: 1.2rem;
      background: #0c0a09;
      color: #ffe4cf;
    }
    .grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 0.8rem; margin-bottom: 1rem; }
    label { display: grid; gap: 0.4rem; color: #cbb9ae; font-size: 0.86rem; }
    input, select {
      border: 1px solid rgba(255, 154, 26, 0.28);
      border-radius: 8px;
      background: #050505;
      color: #fff6ee;
      padding: 0.68rem;
    }
    code {
      display: block;
      padding: 0.9rem;
      border-radius: 8px;
      background: #050505;
      color: #ffca76;
      overflow-wrap: anywhere;
      font: 0.92rem/1.5 ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace;
    }
    @media (max-width: 720px) { .grid { grid-template-columns: 1fr; } }
  `;

  constructor() {
    super();
    this.base = "https://api.aitrailblazer.net";
    this.path = "/v1/top-stressed/natural";
    this.ticker = "";
    this.limit = "10";
    this.style = "professional";
  }

  get url() {
    const params = new URLSearchParams();
    if (this.ticker) params.set("ticker", this.ticker.toUpperCase());
    if (this.limit) params.set("limit", this.limit);
    if (this.style) params.set("style", this.style);
    const query = params.toString();
    return `${this.base}${this.path}${query ? `?${query}` : ""}`;
  }

  render() {
    return html`
      <section class="box">
        <div class="grid">
          <label>
            Ticker
            <input .value=${this.ticker} placeholder="RIOT" @input=${(event) => { this.ticker = event.target.value; }}>
          </label>
          <label>
            Limit
            <input .value=${this.limit} type="number" min="1" max="75" @input=${(event) => { this.limit = event.target.value; }}>
          </label>
          <label>
            Style
            <select .value=${this.style} @change=${(event) => { this.style = event.target.value; }}>
              <option value="professional">professional</option>
              <option value="concise">concise</option>
              <option value="trader">trader</option>
              <option value="detailed">detailed</option>
            </select>
          </label>
        </div>
        <code>${this.url}</code>
      </section>
    `;
  }
}

class DSNaturalTierBanner extends LitElement {
  static styles = css`
    :host { display: block; }
    .banner {
      border: 1px solid rgba(255, 154, 26, 0.32);
      border-radius: 12px;
      padding: 1.55rem;
      background:
        radial-gradient(circle at top left, rgba(255,154,26,.14), transparent 20rem),
        linear-gradient(135deg, rgba(255,48,40,.1), rgba(227,0,85,.08)),
        #0c0a09;
      color: #ffe4cf;
    }
    h2 { margin: 0 0 0.65rem; font: 700 clamp(1.55rem, 3vw, 2.5rem)/1.06 "Space Grotesk", Inter, sans-serif; }
    p { margin: 0 0 1rem; color: #cbb9ae; line-height: 1.55; max-width: 74ch; }
    ul { margin: 0; padding-left: 1.2rem; color: #cbb9ae; line-height: 1.65; }
    code { color: #ffca76; }
  `;

  render() {
    return html`
      <section class="banner">
        <h2>Natural Language briefs compile evidence into validated Markdown.</h2>
        <p>Raw and composite routes return structured evidence. Natural Language routes package that evidence for human workflows while preserving source dates, caveats, quality flags, nulls, evidence hashes, and non-advice disclaimers.</p>
        <ul>
          <li><code>deltasignal_top_stressed_natural</code> is live at $0.95.</li>
          <li><code>deltasignal_morning_brief_natural</code> and <code>deltasignal_covenant_stress_natural</code> are Phase 1 targets.</li>
          <li>The brief is a renderer output, not a new analytics source.</li>
        </ul>
      </section>
    `;
  }
}

customElements.define("ds-copy-prompt", DSCopyPrompt);
customElements.define("ds-install-tabs", DSInstallTabs);
customElements.define("ds-pricing-toggle", DSPricingToggle);
customElements.define("ds-x402-checklist", DSX402Checklist);
customElements.define("ds-route-example", DSRouteExample);
customElements.define("ds-natural-tier-banner", DSNaturalTierBanner);
