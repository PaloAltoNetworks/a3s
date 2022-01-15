import { useCallback, useEffect } from "react"

interface IssueParams {
  sourceNamespace: string
  sourceName: string
  /**
   * The redirect url after we've authenticated the user.
   * If empty, the token will be returned in the promise.
   */
  redirectUrl?: string
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
   * The audience for the JWT.
   */
  audience: string[]
  /**
   * Used to immediatly redirct the user after OIDC login is successful.
   * If empty, `onOidcSuccess` will be called with the token.
   */
  oidcRedirectUrl?: string
  onOidcSuccess?: (token: string) => void
}

/**
 * TODO: Support custom fetch function.
 * TODO: Add error handling
 */
export function useIssue({
  apiUrl,
  audience,
  oidcRedirectUrl,
  onOidcSuccess,
}: UseIssueOptions) {
  const issueUrl = `${apiUrl}/issue`

  const handleIssueResponse = useCallback(
    (redirectUrl?: string) => async (res: Response) => {
      if (res.status === 200) {
        if (redirectUrl) {
          window.location.href = redirectUrl
        } else {
          return (await res.json()).token as string
        }
      } else {
        throw Error("Request to issue failed. Please check the network tab for details")
      }
    },
    []
  )

  const issueWithLdap = useCallback(
    ({
      sourceNamespace,
      sourceName,
      username,
      password,
      redirectUrl,
    }: IssueLdapParams) =>
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
          cookie: !!redirectUrl,
          cookieDomain: window.location.hostname,
          audience,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      }).then(handleIssueResponse(redirectUrl)),
    [issueUrl, audience, handleIssueResponse]
  )

  const issueWithMtls = useCallback(
    ({ sourceNamespace, sourceName, redirectUrl }: IssueParams) =>
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "MTLS",
          sourceNamespace,
          sourceName,
          cookie: !!redirectUrl,
          cookieDomain: window.location.hostname,
          audience,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      }).then(handleIssueResponse(redirectUrl)),
    [issueUrl, audience, handleIssueResponse]
  )

  const issueWithOidc = useCallback(
    ({ sourceNamespace, sourceName }: IssueParams) => {
      // Remove the trailing slash
      const currentUrl = (
        window.location.origin + window.location.pathname
      ).replace(/\/$/, "")
      return fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "OIDC",
          sourceNamespace,
          sourceName,
          inputOIDC: {
            redirectURL: currentUrl,
            redirectErrorURL: currentUrl,
            noAuthRedirect: true,
          },
        }),
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(res => res.json())
        .then(obj => {
          window.location.href = obj.inputOIDC.authURL
        })
    },
    [issueUrl]
  )

  // Cloak an existing token
  const issueWithA3s = useCallback(
    ({
      cloak,
      token,
      redirectUrl,
    }: {
      cloak: string[]
      token: string
      redirectUrl?: string
    }) =>
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "A3S",
          inputA3S: {
            token,
          },
          cloak,
          cookie: !!redirectUrl,
          cookieDomain: window.location.hostname,
          audience,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      }).then(handleIssueResponse(redirectUrl)),
    [audience, issueUrl]
  )

  // Auto login after OIDC redirection
  const params = new URLSearchParams(window.location.search)
  const OIDCstate = params.get("state")
  const OIDCcode = params.get("code")
  useEffect(() => {
    const sourceNamespace = localStorage.getItem("sourceNamespace")
    const sourceName = localStorage.getItem("sourceName")
    const saveTokenStorage = localStorage.getItem("saveToken") === "true"
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
      }).then(res => {
        // Clear the state and code
        history.replaceState(null, "", window.location.pathname)
        handleIssueResponse(oidcRedirectUrl)(res).then(token => {
          token && onOidcSuccess && onOidcSuccess(token)
        })
      })
    }
  }, [
    OIDCstate,
    OIDCcode,
    issueUrl,
    audience,
    handleIssueResponse,
    onOidcSuccess,
    oidcRedirectUrl,
  ])

  return { issueWithLdap, issueWithOidc, issueWithMtls, issueWithA3s, oidcIssuing: OIDCstate && OIDCcode }
}
