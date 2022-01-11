import { useCallback, useEffect, useState } from "react"

interface IssueParams {
  sourceNamespace: string
  sourceName: string
}

interface IssueLdapParams extends IssueParams {
  username: string
  password: string
}

interface UseIssueOptions {
  /**
   * The base URL for the a3s backend. Shouldn't include the `/` in the end.
   * Example: `https://127.0.0.1:44443`
   */
  apiUrl: string
  /**
   * The redirect url after we've authenticated the user.
   */
  redirectUrl: string
  /**
   * The audience for the JWT.
   */
  audience: string[]
  /**
   * Save the token in the `token` state, instead of instant redirect.
   */
  saveToken?: boolean
  /**
   * If both state and code exists, OIDC login will be automatically triggered.
   */
  OIDCstate?: string | null
  /**
   * If both state and code exists, OIDC login will be automatically triggered.
   */
  OIDCcode?: string | null
}

/**
 * TODO: Support custom fetch function.
 * TODO: Add error handling
 */
export function useIssue({
  apiUrl,
  redirectUrl,
  audience,
  saveToken: saveTokenDefault = false,
  OIDCstate,
  OIDCcode,
}: UseIssueOptions) {
  const issueUrl = `${apiUrl}/issue`
  /**
   * Token is only stored when `saveToken` is true.
   */
  const [token, setToken] = useState<string | null>(null)

  const handleIssueResponse = useCallback(
    (saveToken: boolean) => async (res: Response) => {
      if (res.status === 200) {
        if (saveToken) {
          setToken((await res.json()).token)
        } else {
          window.location.href = redirectUrl
        }
      } else {
        console.error(
          "Request to issue failed. Please check the network tab for details"
        )
      }
    },
    [setToken, redirectUrl]
  )

  const issueWithLdap = useCallback(
    ({ sourceNamespace, sourceName, username, password }: IssueLdapParams) =>
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "LDAP",
          sourceNamespace,
          sourceName,
          inputLDAP: {
            username,
            password,
          },
          cookie: !saveTokenDefault,
          cookieDomain: window.location.hostname,
          audience,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      }).then(handleIssueResponse(saveTokenDefault)),
    [issueUrl, audience, handleIssueResponse, saveTokenDefault]
  )

  const issueWithMtls = useCallback(
    ({ sourceNamespace, sourceName }: IssueParams) =>
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "MTLS",
          sourceNamespace,
          sourceName,
          cookie: !saveTokenDefault,
          cookieDomain: window.location.hostname,
          audience,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      }).then(handleIssueResponse(saveTokenDefault)),
    [issueUrl, audience, handleIssueResponse, saveTokenDefault]
  )

  const issueWithOidc = useCallback(
    ({ sourceNamespace, sourceName }: IssueParams) => {
      // Remove the trailing slash
      const oidcRedirectURL = (
        window.location.origin + window.location.pathname
      ).replace(/\/$/, "")
      return fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "OIDC",
          sourceNamespace,
          sourceName,
          inputOIDC: {
            redirectURL: oidcRedirectURL,
            redirectErrorURL: oidcRedirectURL,
            noAuthRedirect: true,
          },
        }),
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(res => res.json())
        .then(obj => {
          localStorage.setItem("sourceNamespace", sourceNamespace)
          localStorage.setItem("sourceName", sourceName)
          localStorage.setItem("saveToken", saveTokenDefault ? "true" : "false")
          window.location.href = obj.inputOIDC.authURL
        })
    },
    [issueUrl, saveTokenDefault]
  )

  const issueWithA3s = useCallback(
    ({ token, cloak }: { token: string; cloak: string[] }) =>
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "A3S",
          inputA3S: {
            token,
          },
          cloak,
          cookie: true,
          cookieDomain: window.location.hostname,
          audience,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      }).then(handleIssueResponse(false)),
    []
  )

  // Below for handling auto login after OIDC redirection
  useEffect(() => {
    const sourceNamespace = localStorage.getItem("sourceNamespace")
    const sourceName = localStorage.getItem("sourceName")
    const saveTokenStorage = localStorage.getItem("saveToken") === "true"
    // Clear local storage immediately
    localStorage.removeItem("sourceNamespace")
    localStorage.removeItem("sourceName")
    localStorage.removeItem("saveToken")
    if (OIDCstate && OIDCcode && sourceNamespace && sourceName) {
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "OIDC",
          sourceNamespace,
          sourceName,
          inputOIDC: {
            state: OIDCstate,
            code: OIDCcode,
          },
          cookie: !saveTokenStorage,
          cookieDomain: window.location.hostname,
          audience,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      }).then(handleIssueResponse(saveTokenStorage))
    }
  }, [OIDCstate, OIDCcode, issueUrl, audience, handleIssueResponse])

  return { token, issueWithLdap, issueWithOidc, issueWithMtls, issueWithA3s }
}
