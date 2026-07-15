#!/usr/bin/env node

import { Wallet, getAddress, hexlify, randomBytes } from "ethers";
import { pathToFileURL } from "node:url";

export const DEFAULTS = Object.freeze({
  url: "https://api.aitrailblazer.net/v1/readiness",
  network: "eip155:8453",
  chainId: 8453,
  asset: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
  payTo: "0x6D91ADF2c545047cbbC5b37a5f457cce081B48d3",
  maxAmountAtomic: 40000n,
  gatewayName: "GatewayWalletBatched",
  gatewayVersion: "1",
  minValiditySeconds: 604800,
  clockSkewSeconds: 600
});

const receiptFields = [
  "schema_id",
  "receipt_id",
  "response_id",
  "route_uri",
  "quoted_cost_usd",
  "charged_amount_usdc",
  "settlement_mode",
  "payer_identity",
  "reconciliation_key",
  "circle_x402_transfer_id",
  "batch_id",
  "settlement_tx_hash",
  "created_at_utc",
  "updated_at_utc",
  "charge_status",
  "quality_flags",
  "source_boundary"
];

export function decodeBase64JSON(value) {
  if (!value) throw new Error("missing base64 JSON value");
  return JSON.parse(Buffer.from(value, "base64").toString("utf8"));
}

export function encodeBase64JSON(value) {
  return Buffer.from(JSON.stringify(value), "utf8").toString("base64");
}

export function parsePaymentRequired(response, bodyText = "") {
  const header = response.headers.get("payment-required");
  if (header) return decodeBase64JSON(header);
  if (bodyText) {
    const body = JSON.parse(bodyText);
    if (body?.x402Version && Array.isArray(body.accepts)) return body;
  }
  throw new Error("402 response did not include a decodable PAYMENT-REQUIRED challenge");
}

export function selectCircleRequirement(challenge) {
  const matches = (challenge?.accepts ?? []).filter(
    (requirement) =>
      requirement?.scheme === "exact" &&
      requirement?.extra?.name === DEFAULTS.gatewayName
  );
  if (matches.length !== 1) {
    throw new Error(`expected exactly one ${DEFAULTS.gatewayName} requirement; found ${matches.length}`);
  }
  return matches[0];
}

export function validateCircleRequirement(requirement, policy = DEFAULTS) {
  if (requirement.scheme !== "exact") throw new Error("Circle requirement scheme must be exact");
  if (requirement.network !== policy.network) throw new Error(`unexpected network: ${requirement.network}`);
  if (getAddress(requirement.asset) !== getAddress(policy.asset)) throw new Error("unexpected payment asset");
  if (getAddress(requirement.payTo) !== getAddress(policy.payTo)) throw new Error("unexpected seller payTo");

  const amount = BigInt(requirement.amount);
  if (amount <= 0n || amount > policy.maxAmountAtomic) {
    throw new Error(`amount ${amount} exceeds policy maximum ${policy.maxAmountAtomic}`);
  }
  if (requirement.extra?.name !== policy.gatewayName) throw new Error("unexpected Gateway domain name");
  if ((requirement.extra?.version ?? "") !== policy.gatewayVersion) throw new Error("unexpected Gateway domain version");
  if (!requirement.extra?.verifyingContract) throw new Error("Circle requirement is missing extra.verifyingContract");
  getAddress(requirement.extra.verifyingContract);
  return requirement;
}

export async function createCirclePaymentPayload({
  challenge,
  requirement,
  signer,
  nowSeconds = Math.floor(Date.now() / 1000),
  nonce = hexlify(randomBytes(32))
}) {
  validateCircleRequirement(requirement);

  const minimum = Number(requirement.extra?.minValiditySeconds ?? DEFAULTS.minValiditySeconds);
  const timeout = Number(requirement.maxTimeoutSeconds ?? 0);
  const validity = Math.max(minimum, timeout) + DEFAULTS.clockSkewSeconds;
  const authorization = {
    from: await signer.getAddress(),
    to: requirement.payTo,
    value: String(requirement.amount),
    validAfter: String(nowSeconds - DEFAULTS.clockSkewSeconds),
    validBefore: String(nowSeconds + validity),
    nonce
  };
  const domain = {
    name: requirement.extra.name,
    version: requirement.extra.version,
    chainId: DEFAULTS.chainId,
    verifyingContract: requirement.extra.verifyingContract
  };
  const types = {
    TransferWithAuthorization: [
      { name: "from", type: "address" },
      { name: "to", type: "address" },
      { name: "value", type: "uint256" },
      { name: "validAfter", type: "uint256" },
      { name: "validBefore", type: "uint256" },
      { name: "nonce", type: "bytes32" }
    ]
  };
  const signature = await signer.signTypedData(domain, types, authorization);
  const payload = {
    x402Version: 2,
    accepted: requirement,
    payload: { signature, authorization }
  };
  if (challenge.resource) payload.resource = challenge.resource;
  if (challenge.extensions) payload.extensions = challenge.extensions;
  return { payload, domain, types };
}

export function decodeOptionalResponseHeader(value) {
  if (!value) return null;
  try {
    return decodeBase64JSON(value);
  } catch {
    try {
      return JSON.parse(value);
    } catch {
      return { decode_error: true };
    }
  }
}

export function pickReceipt(body) {
  const receipt = body?.settlement_receipt ?? body?.billing?.settlement_receipt;
  if (!receipt || typeof receipt !== "object") return null;
  return Object.fromEntries(receiptFields.filter((key) => receipt[key] !== undefined).map((key) => [key, receipt[key]]));
}

export function summarizeRequirement(requirement) {
  return {
    scheme: requirement.scheme,
    network: requirement.network,
    asset: requirement.asset,
    amount: requirement.amount,
    payTo: requirement.payTo,
    maxTimeoutSeconds: requirement.maxTimeoutSeconds,
    gatewayDomain: {
      name: requirement.extra?.name,
      version: requirement.extra?.version,
      verifyingContract: requirement.extra?.verifyingContract
    }
  };
}

export async function run({
  url = DEFAULTS.url,
  pay = false,
  privateKey = process.env.EVM_PRIVATE_KEY,
  freeTierUser = process.env.DELTASIGNAL_X402_FREE_TIER_USER,
  fetchImpl = fetch
} = {}) {
  const requestHeaders = { accept: "application/json" };
  if (freeTierUser) requestHeaders["x-codex-user"] = freeTierUser;

  const challengeResponse = await fetchImpl(url, { headers: requestHeaders });
  const challengeBody = await challengeResponse.text();
  if (challengeResponse.status !== 402) {
    return {
      mode: "no-payment-required",
      status: challengeResponse.status,
      note: "Use a pre-exhausted free-tier identity to exercise the 402 flow."
    };
  }

  const challenge = parsePaymentRequired(challengeResponse, challengeBody);
  const requirement = validateCircleRequirement(selectCircleRequirement(challenge));
  const selected = summarizeRequirement(requirement);
  if (!pay) {
    return {
      mode: "dry-run",
      status: 402,
      selected,
      safeToSign: true,
      charged: false,
      next: "Re-run with --pay only after the wallet has available Circle Gateway balance."
    };
  }
  if (!privateKey) throw new Error("EVM_PRIVATE_KEY is required with --pay");

  const signer = new Wallet(privateKey);
  const { payload } = await createCirclePaymentPayload({ challenge, requirement, signer });
  const paymentHeader = encodeBase64JSON(payload);
  const paidResponse = await fetchImpl(url, {
    headers: {
      ...requestHeaders,
      "payment-signature": paymentHeader
    }
  });
  const paidText = await paidResponse.text();
  let paidBody = null;
  try {
    paidBody = paidText ? JSON.parse(paidText) : null;
  } catch {
    paidBody = { non_json_response: true };
  }

  return {
    mode: "circle-payment",
    selected,
    payer: signer.address,
    status: paidResponse.status,
    settlement: decodeOptionalResponseHeader(paidResponse.headers.get("payment-response")),
    receipt: pickReceipt(paidBody),
    response: paidBody
  };
}

function parseArgs(argv) {
  const options = {};
  for (let i = 0; i < argv.length; i += 1) {
    const arg = argv[i];
    if (arg === "--pay") options.pay = true;
    else if (arg === "--dry-run") options.pay = false;
    else if (arg === "--url") options.url = argv[++i];
    else if (arg === "--free-tier-user") options.freeTierUser = argv[++i];
    else if (arg === "--help") options.help = true;
    else throw new Error(`unknown argument: ${arg}`);
  }
  return options;
}

function usage() {
  return `Circle-aware DeltaSignal x402 buyer

Dry-run challenge selection (never signs or pays):
  node buyer.mjs --dry-run --free-tier-user <pre-exhausted-id>

Explicit Circle payment:
  EVM_PRIVATE_KEY=... node buyer.mjs --pay --free-tier-user <pre-exhausted-id>

Options:
  --url <url>              Seller route (default: ${DEFAULTS.url})
  --free-tier-user <id>    Optional pre-exhausted identity used to force a 402
  --dry-run                Inspect and validate Circle without signing (default)
  --pay                    Sign and submit the explicit Circle requirement
`;
}

if (import.meta.url === pathToFileURL(process.argv[1]).href) {
  try {
    const options = parseArgs(process.argv.slice(2));
    if (options.help) {
      process.stdout.write(usage());
    } else {
      const result = await run(options);
      process.stdout.write(`${JSON.stringify(result, null, 2)}\n`);
      if (result.mode === "circle-payment" && result.status !== 200) process.exitCode = 1;
    }
  } catch (error) {
    process.stderr.write(`circle-aware-buyer: ${error.message}\n`);
    process.exitCode = 1;
  }
}
