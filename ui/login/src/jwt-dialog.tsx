import {
  Alert,
  AlertTitle,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  useMediaQuery,
} from "@mui/material"
import { Box } from "@mui/system"
import { useEffect, useState } from "react"
import { decodeBase64url, stringToBuffer } from "./utils"
import type { DecodedJWT } from "./types"

export const JwtDialog = ({
  payload,
  header,
  rawToken,
  onClose,
}: {
  rawToken: string
  onClose(): void
} & DecodedJWT) => {
  const fullScreen = useMediaQuery("(max-width: 600px)")
  const [isValid, setIsValid] = useState<boolean>()
  const [verificationError, setVerificationError] = useState<string>()
  useEffect(() => {
    const verify = async () => {
      try {
        // Fetch the JWKS from the issuer
        const jwks = await fetch(payload.iss + "/.well-known/jwks.json").then(
          p => p.json()
        )

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
            name: "ECDSA",
            namedCurve: "P-256",
          },
          false, // not extractable
          ["verify"]
        )

        // Convert the signature and the payload to ArrayBuffers
        const splitted = rawToken.split(".")
        const signature = splitted[2]
        const encoded = splitted.slice(0, 2).join(".")
        const decodedSign = decodeBase64url(signature)
        const signatureArrayBuffer = stringToBuffer(decodedSign)
        const encodedArrayBuffer = stringToBuffer(encoded)

        // Verify the signature
        const result = await crypto.subtle.verify(
          { name: "ECDSA", hash: "SHA-256" },
          cryptoKeyPublic,
          signatureArrayBuffer,
          encodedArrayBuffer
        )

        if(!result) {
          setVerificationError(`Signature doesn't match`)
        }
        setIsValid(result)
      } catch (error) {
        if (error instanceof Error) {
          setVerificationError(error.message)
          setIsValid(false)
        }
      }
    }

    verify()
  }, [])
  return (
    <Dialog open fullScreen={fullScreen}>
      <DialogTitle>Decoded JWT</DialogTitle>
      <DialogContent>
        {isValid === undefined ? null : isValid ? (
          <Alert variant="outlined" severity="success">
            <AlertTitle>Signature Verified</AlertTitle>
            Issuer: {payload.iss}
          </Alert>
        ) : (
          <Alert variant="outlined" severity="error">
            <AlertTitle>Signature Verification Failed</AlertTitle>
            {verificationError}
          </Alert>
        )}
        <Box
          sx={{
            bgcolor: "action.hover",
            borderRadius: 2,
            overflowX: "auto",
            px: 2,
            mt: 2,
          }}
        >
          <pre>{JSON.stringify(payload, null, 2)}</pre>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="contained">
          Close
        </Button>
      </DialogActions>
    </Dialog>
  )
}
