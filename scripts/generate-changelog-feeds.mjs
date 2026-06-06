#!/usr/bin/env node
import { readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const repoRoot = dirname(dirname(fileURLToPath(import.meta.url)));
const siteURL = "https://aitrailblazer.github.io/deltasignal-atlas-codex-plugin";
const changelogURL = `${siteURL}/changelog.html`;
const llmChangelogURL = `${siteURL}/changelog-llm.html`;
const feedLimit = 30;

const html = readFileSync(join(repoRoot, "changelog-llm.html"), "utf8");
const jsonMatch = html.match(/<script type="application\/json" id="changelog-llm-json">\s*([\s\S]*?)\s*<\/script>/);
if (!jsonMatch) {
  throw new Error("Could not find changelog-llm-json script block");
}

const contract = JSON.parse(jsonMatch[1]);
const entries = [...(contract.timeline_entries || [])]
  .filter((entry) => isCalendarDate(entry.date))
  .sort((a, b) => Number(b.river_index || 0) - Number(a.river_index || 0))
  .slice(0, feedLimit);

if (entries.length === 0) {
  throw new Error("No dated changelog entries found for feed generation");
}

const latest = entries[0];
const updated = dateToISOString(latest.date);

writeFileSync(join(repoRoot, "rss.xml"), buildRSS(entries, updated));
writeFileSync(join(repoRoot, "atom.xml"), buildAtom(entries, updated));

function buildRSS(feedEntries, updatedISO) {
  const items = feedEntries.map((entry) => `    <item>
      <title>${xml(entry.title)}</title>
      <link>${xml(entryURL(entry))}</link>
      <guid isPermaLink="false">${xml(entryGUID(entry))}</guid>
      <pubDate>${new Date(dateToISOString(entry.date)).toUTCString()}</pubDate>
      <category>${xml(entry.tripcode || "DeltaSignal")}</category>
      <description>${xml(entryDescription(entry))}</description>
    </item>`).join("\n");

  return `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Delta Signal ATLAS-7 Changelog</title>
    <link>${changelogURL}</link>
    <atom:link href="${siteURL}/rss.xml" rel="self" type="application/rss+xml" />
    <description>Public changelog for Delta Signal ATLAS-7 MCP, OpenAPI, x402, notification, SPECTRA, and evidence-contract changes.</description>
    <language>en-us</language>
    <lastBuildDate>${new Date(updatedISO).toUTCString()}</lastBuildDate>
    <ttl>720</ttl>
${items}
  </channel>
</rss>
`;
}

function buildAtom(feedEntries, updatedISO) {
  const items = feedEntries.map((entry) => `  <entry>
    <title>${xml(entry.title)}</title>
    <link href="${xml(entryURL(entry))}" />
    <id>${xml(entryGUID(entry))}</id>
    <updated>${dateToISOString(entry.date)}</updated>
    <category term="${xml(entry.tripcode || "DeltaSignal")}" />
    <summary>${xml(entryDescription(entry))}</summary>
  </entry>`).join("\n");

  return `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Delta Signal ATLAS-7 Changelog</title>
  <subtitle>Public changelog for Delta Signal ATLAS-7 MCP, OpenAPI, x402, notification, SPECTRA, and evidence-contract changes.</subtitle>
  <link href="${changelogURL}" />
  <link href="${siteURL}/atom.xml" rel="self" type="application/atom+xml" />
  <id>${siteURL}/changelog</id>
  <updated>${updatedISO}</updated>
  <author><name>AITrailblazer</name></author>
${items}
</feed>
`;
}

function entryDescription(entry) {
  const parts = [
    entry.summary || "",
    entry.bread ? `Breadcrumb: ${entry.bread}` : "",
    entry.track ? `Track: ${entry.track}` : "",
    entry.commit ? `Commit: ${entry.commit}` : "",
  ].filter(Boolean);
  return parts.join(" ");
}

function entryURL(entry) {
  const anchor = entry.tripcode ? `#${entry.tripcode}` : "#tracking-river";
  return `${llmChangelogURL}${anchor}`;
}

function entryGUID(entry) {
  return `${siteURL}/changelog/${entry.tripcode || entry.river_index}`;
}

function isCalendarDate(value) {
  return /^(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2},\s+\d{4}$/.test(String(value || ""));
}

function dateToISOString(value) {
  const parsed = new Date(`${value} 00:00:00 UTC`);
  if (Number.isNaN(parsed.valueOf())) {
    throw new Error(`Invalid changelog date: ${value}`);
  }
  return parsed.toISOString();
}

function xml(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&apos;");
}
