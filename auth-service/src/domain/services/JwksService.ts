/**
 * JwksService (ADR-009.1) — SKELETON for the v4 handoff.
 *
 * Target state: auth-service signs with RS256 (kid header), serves
 * GET /.well-known/jwks.json, and api-service verifies locally with a
 * cached JWKS (introspection stays available behind a strict-mode flag).
 *
 * Implementation plan lives in documentation/phases/v4-HANDOFF.md §A6:
 *   - keypair source: PEM via env/secret (init-secrets.sh grows an
 *     RS256 keypair; k8s Secret auth-service-keys)
 *   - kid = first 8 bytes of SHA-256 of the DER SPKI, hex
 *   - TokenService signs RS256 when a keypair is configured (HS256 stays
 *     as a fallback during migration, controlled by JWT_ALGORITHM env)
 *   - JWKS shape: { keys: [{ kty:"RSA", use:"sig", alg:"RS256", kid, n, e }] }
 *   - rotation: KEYS array, newest signs, all published until expiry
 *
 * This stub compiles and is intentionally unused until wired; no new npm
 * deps are required (node:crypto covers RSA + JWK export).
 */
import { createPublicKey, createHash, type KeyObject } from "node:crypto";

export interface JsonWebKey {
  kty: string;
  use: string;
  alg: string;
  kid: string;
  n: string;
  e: string;
}

export class JwksService {
  private readonly publicKeys: KeyObject[];

  constructor(publicKeyPems: string[]) {
    this.publicKeys = publicKeyPems.map((pem) => createPublicKey(pem));
  }

  /** Stable key id derived from the SPKI DER (HANDOFF §A6). */
  static kidFor(key: KeyObject): string {
    const der = key.export({ type: "spki", format: "der" });
    return createHash("sha256").update(der).digest("hex").slice(0, 16);
  }

  /** The /.well-known/jwks.json document. TODO(v4): route + controller. */
  jwks(): { keys: JsonWebKey[] } {
    return {
      keys: this.publicKeys.map((key) => {
        const jwk = key.export({ format: "jwk" }) as { kty: string; n: string; e: string };
        return {
          kty: jwk.kty,
          use: "sig",
          alg: "RS256",
          kid: JwksService.kidFor(key),
          n: jwk.n,
          e: jwk.e,
        };
      }),
    };
  }
}
