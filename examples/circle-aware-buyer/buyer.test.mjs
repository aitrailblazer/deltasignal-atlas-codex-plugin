import assert from "node:assert/strict";
import test from "node:test";
import { Wallet, verifyTypedData } from "ethers";

import {
  createCirclePaymentPayload,
  pickReceipt,
  selectCircleRequirement,
  summarizeRequirement,
  validateCircleRequirement
} from "./buyer.mjs";

const standard = {
  scheme: "exact",
  network: "eip155:8453",
  asset: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
  amount: "40000",
  payTo: "0x6D91ADF2c545047cbbC5b37a5f457cce081B48d3",
  maxTimeoutSeconds: 60,
  extra: { name: "USD Coin", version: "2" }
};
const circle = {
  ...standard,
  maxTimeoutSeconds: 604800,
  extra: {
    name: "GatewayWalletBatched",
    version: "1",
    verifyingContract: "0x77777777dcc4d5a8b6e418fd04d8997ef11000ee",
    minValiditySeconds: 604800
  }
};
const challenge = {
  x402Version: 2,
  resource: { url: "https://api.aitrailblazer.net/v1/readiness" },
  accepts: [standard, circle]
};

test("selects Circle by capability instead of array position", () => {
  assert.deepEqual(selectCircleRequirement(challenge), circle);
  assert.equal(summarizeRequirement(circle).gatewayDomain.name, "GatewayWalletBatched");
});

test("fails closed on altered seller or missing Gateway contract", () => {
  assert.throws(
    () => validateCircleRequirement({ ...circle, payTo: "0x0000000000000000000000000000000000000001" }),
    /seller payTo/
  );
  assert.throws(
    () => validateCircleRequirement({ ...circle, extra: { ...circle.extra, verifyingContract: "" } }),
    /missing extra.verifyingContract/
  );
});

test("creates a verifiable Gateway-domain authorization without leaking it through the summary", async () => {
  // Public test-only key. Never fund or reuse it.
  const signer = new Wallet("0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d");
  const nonce = `0x${"11".repeat(32)}`;
  const { payload, domain, types } = await createCirclePaymentPayload({
    challenge,
    requirement: circle,
    signer,
    nowSeconds: 1_800_000_000,
    nonce
  });

  assert.deepEqual(payload.accepted, circle);
  assert.equal(payload.payload.authorization.nonce, nonce);
  assert.equal(domain.name, "GatewayWalletBatched");
  assert.equal(domain.verifyingContract, circle.extra.verifyingContract);
  assert.equal(
    verifyTypedData(domain, types, payload.payload.authorization, payload.payload.signature),
    signer.address
  );
  assert.equal(JSON.stringify(summarizeRequirement(circle)).includes(payload.payload.signature), false);
});

test("returns only parser-stable receipt fields", () => {
  const receipt = pickReceipt({
    billing: {
      settlement_receipt: {
        receipt_id: "dsr_test",
        settlement_mode: "circle_gateway_batched",
        charge_status: "submitted",
        reconciliation_key: "dsrec_test",
        secret_internal_field: "must-not-leak"
      }
    }
  });
  assert.deepEqual(receipt, {
    receipt_id: "dsr_test",
    settlement_mode: "circle_gateway_batched",
    reconciliation_key: "dsrec_test",
    charge_status: "submitted"
  });
});
