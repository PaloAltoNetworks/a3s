export interface RequestQrJson {
  claims: string[]
  issuers: string[]
  message: string
  meta: {
    version: string
  }
}

export interface RequestEntry {
  name: string
  description: string
  claims: string[]
  issuers: string[]
  message: string
  ID: string
}
