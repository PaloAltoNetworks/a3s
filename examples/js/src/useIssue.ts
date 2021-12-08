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
}

/**
 * TODO: Support custom fetch function.
 */
export function useIssue({ apiUrl }: UseIssueOptions) {
  const [token, setToken] = useState()
  const issueUrl = `${apiUrl}/issue`

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
          cookie: true,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(res => res.json())
        .then(obj => {
          setToken(obj.token)
        }),
    [issueUrl]
  )

  const issueWithMtls = useCallback(
    ({ sourceNamespace, sourceName }: IssueParams) =>
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "MTLS",
          sourceNamespace,
          sourceName,
          cookie: true,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(res => res.json())
        .then(obj => {
          setToken(obj.token)
        }),
    [issueUrl]
  )

  const issueWithOidc = useCallback(
    ({ sourceNamespace, sourceName }: IssueParams) =>
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "OIDC",
          sourceNamespace,
          sourceName,
          inputOIDC: {
            redirectURL: window.location.origin,
            redirectErrorURL: window.location.origin,
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
          window.location.href = obj.inputOIDC.authURL
        }),
    [issueUrl]
  )

  // Below for handling auto login after OIDC redirection
  const params = new URLSearchParams(window.location.search)
  const state = params.get("state")
  const code = params.get("code")
  useEffect(() => {
    const sourceNamespace = localStorage.getItem("sourceNamespace")
    const sourceName = localStorage.getItem("sourceName")
    // Clear local storage immediately
    localStorage.removeItem("sourceNamespace")
    localStorage.removeItem("sourceName")
    if (state && code && sourceNamespace && sourceName) {
      fetch(issueUrl, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "OIDC",
          sourceNamespace,
          sourceName,
          inputOIDC: {
            state,
            code,
          },
          cookie: true,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(res => res.json())
        .then(obj => {
          setToken(obj.token)
          // Clear the search params in the URL
          history.replaceState(null, "", "/")
        })
    }
  }, [state, code, issueUrl])

  return { issueWithLdap, issueWithOidc, issueWithMtls, token }
}
