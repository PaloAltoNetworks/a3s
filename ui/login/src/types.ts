import type { JwtHeader } from "jwt-decode"

export interface DecodedJWT {
  payload: Record<string, any>;
  header: JwtHeader & { kid: string }
}