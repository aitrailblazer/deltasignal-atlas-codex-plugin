# Executable Circle-aware x402 buyer

This example performs the complete buyer-side flow Jeff requested:

1. fetch `GET /v1/readiness`;
2. decode the x402 v2 `PAYMENT-REQUIRED` challenge;
3. select `GatewayWalletBatched` by capability rather than array position;
4. validate Base, USDC, amount, seller, Gateway domain, and verifying contract;
5. sign `TransferWithAuthorization` against the challenge-provided Gateway EIP-712 domain;
6. submit `PAYMENT-SIGNATURE`; and
7. print a sanitized settlement and receipt-reconciliation summary.

It never prints the private key, raw signature, or encoded payment header.

## Prerequisites

- Node.js 20 or newer.
- A Circle-aware EVM wallet private key supplied at runtime through `EVM_PRIVATE_KEY`.
- **Deposited or available Circle Gateway balance.** Base ERC-20 USDC in the wallet is not sufficient by itself.
- A caller identity whose free starter allowance is exhausted when you need to force the `402` flow.

## Install and test

```bash
cd examples/circle-aware-buyer
npm ci
npm test
```

## Safe dry run

Dry run is the default. It validates and displays the selected Circle requirement but never signs or pays:

```bash
node buyer.mjs \
  --dry-run \
  --free-tier-user <pre-exhausted-free-tier-id>
```

Expected result includes:

```json
{
  "mode": "dry-run",
  "status": 402,
  "safeToSign": true,
  "charged": false
}
```

## Execute one Circle payment

Only run this after confirming the wallet has available Gateway balance:

```bash
export EVM_PRIVATE_KEY='0x...'
node buyer.mjs \
  --pay \
  --free-tier-user <pre-exhausted-free-tier-id>
unset EVM_PRIVATE_KEY
```

The output includes the selected requirement, safe payer address, HTTP status, decoded `PAYMENT-RESPONSE`, and these receipt fields when returned:

- `route_uri`
- `quoted_cost_usd`
- `charged_amount_usdc`
- `settlement_mode`
- `payer_identity`
- `reconciliation_key`
- `circle_x402_transfer_id`
- `batch_id`
- `settlement_tx_hash`
- timestamps
- `charge_status`
- `quality_flags`
- `source_boundary`

For `circle_gateway_batched`, the response receipt is the per-call attribution authority until reconciliation adds batch or transaction references. An `insufficient_balance` response means the Circle selection and signature reached Gateway but no payment settled.

## Safety boundaries

- The program fails closed if network, asset, seller, amount, Gateway name/version, or verifying contract violates policy.
- `--pay` is required before a signature is created.
- There is no automatic fallback to another payment requirement.
- The example is payment plumbing, not investment advice or trading software.
