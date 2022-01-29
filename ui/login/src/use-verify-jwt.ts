import jwtDecode from "jwt-decode"
import { useCallback } from "react"
import { DecodedJWT } from "./types"
import { decodeBase64url, stringToBuffer } from "./utils"

export const useVerifyJwt = () =>
  useCallback(async ({ token }: { token: string }) => {
    // Decode the token.
    // Though we can also pass in the decoded object as arguments,
    // decoding in place again simplifies the API and reduces the chance of bugs.
    const payload = jwtDecode(token) as DecodedJWT["payload"]
    const header = jwtDecode(token, {
      header: true,
    }) as DecodedJWT["header"]

    // Fetch the JWKS from the issuer
    const jwksRes = await fetch(payload.iss + "/.well-known/jwks.json")
    const jwks = await jwksRes.json()

    // Get the JWK to use
    const jwk = jwks.keys.find((k: { kid: string }) => k.kid === header.kid)
    if (!jwk) {
      throw Error(`No matching JWK from the issuer: ${payload.iss}`)
    }

    // Remove extra fields which are not supported by `crypto.subtle.importKey`
    const jwkToUse = {
      crv: jwk.crv,
      kty: jwk.kty,
      x: jwk.x,
      y: jwk.y,
    }

    // Import the public key with the `CryptoKey` interface of the Web Crypto API
    const cryptoKeyPublic = await crypto.subtle.importKey(
      "jwk",
      jwkToUse,
      {
        // TODO: Support other algorithms
        name: "ECDSA",
        namedCurve: "P-256",
      },
      false, // not extractable
      ["verify"]
    )

    // Convert the signature and the payload to ArrayBuffers
    const splitted = token.split(".")
    const signature = splitted[2]
    const encoded = splitted.slice(0, 2).join(".")
    const decodedSign = decodeBase64url(signature)
    const signatureArrayBuffer = stringToBuffer(decodedSign)
    const encodedArrayBuffer = stringToBuffer(encoded)

    // Verify the signature
    return crypto.subtle.verify(
      { name: "ECDSA", hash: "SHA-256" },
      cryptoKeyPublic,
      signatureArrayBuffer,
      encodedArrayBuffer
    )
  }, [])
