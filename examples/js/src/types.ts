// TODO: Auto generate this file from go spec

export type IssueSourceType = "LDAP" | "MTLS" | "OIDC"

export interface IssueParams {
  sourceNamespace: string
  sourceName: string
}

export interface IssueLdapParams extends IssueParams {
  username: string
  password: string
}
