export interface DecodedJWT {
  payload: Record<string, any>
  header: { alg: string; typ: string; kid: string }
}
