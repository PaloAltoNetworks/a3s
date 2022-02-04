import { useCallback, useState, useEffect } from "react"
import { Box } from "@mui/system"
import {
  Button,
  TextField,
  Radio,
  RadioGroup,
  FormControlLabel,
  FormControl,
  FormLabel,
  Checkbox,
  Typography,
  useTheme,
} from "@mui/material"
import { useIssue } from "./use-issue"
import { CloakDialog } from "./cloak-dialog"
import jwtDecode from "jwt-decode"
import { useLocalState } from "./utils/use-local-state"
import { QrCodeDialog } from "./qr-code-dialog"
import { QrScan } from "./qr-scan"

type StringBoolean = "true" | "false"
type DialogState =
  | { type: "None" }
  | { type: "Cloak"; token: string }
  | { type: "QrCode"; token: string }

const audience = ["__AUDIENCE__"]
const apiUrl = "__API_URL__"
let redirectUrl = "__REDIRECT_URL__"
// const audience = ["https://127.0.0.1:44443"]
// const apiUrl = "https://localhost:44443"
// let redirectUrl = "https://google.com"
const redirectUrlInLocalStorage = localStorage.getItem("redirectUrl")
if (redirectUrlInLocalStorage) {
  redirectUrl = redirectUrlInLocalStorage
  localStorage.removeItem("redirectUrl")
}
// A temp solution to avoid the issue of the `audience` being empty
if (audience.length === 1 && audience[0] === "") {
  audience[0] = "public"
}
// @ts-ignore
const isQrCodeMode = redirectUrl === ""
const enableCloak = "__ENABLE_CLOAKING__" as StringBoolean

export const Login = () => {
  const theme = useTheme()
  const [cloak, setCloak] = useLocalState(enableCloak, "cloak")
  const [sourceType, setSourceType] = useLocalState<string>(
    "MTLS",
    "sourceType"
  )
  const [sourceNamespace, setSourceNamespace] = useLocalState<string>(
    "/",
    "sourceNamespace"
  )
  const [sourceName, setSourceName] = useLocalState<string>("", "sourceName")
  const [ldapUsername, setLdapUsername] = useState("")
  const [ldapPassword, setLdapPassword] = useState("")
  const [dialogState, setDialogState] = useState<DialogState>({ type: "None" })
  const { issueWithLdap, issueWithMtls, issueWithOidc, issueWithA3s } =
    useIssue({
      apiUrl,
      audience,
    })

  const handleIssueResponse = useCallback(
    (shouldRedirect: boolean) => async (res: Response) => {
      if (res.status === 200) {
        if (shouldRedirect) {
          window.location.href = redirectUrl
        } else {
          return (await res.json()).token as string
        }
      } else {
        throw Error(
          "Request to issue failed. Please check the network tab for details"
        )
      }
    },
    []
  )

  const onToken = useCallback(
    (token?: string) => {
      if (!token) {
        return
      }
      if (cloak === "true") {
        setDialogState({
          type: "Cloak",
          token,
        })
      } else if (isQrCodeMode) {
        setDialogState({
          type: "QrCode",
          token,
        })
      }
    },
    [cloak, isQrCodeMode]
  )

  // Not using `cloak === "false"` in case `__ENABLE_CLOAKING__` is not replaced
  const shouldRedirect = !isQrCodeMode && cloak !== "true"

  // Below for OIDC auto login
  const params = new URLSearchParams(window.location.search)
  const OIDCstate = params.get("state")
  const OIDCcode = params.get("code")
  useEffect(() => {
    if (OIDCstate && OIDCcode) {
      fetch(`${apiUrl}/issue`, {
        method: "POST",
        body: JSON.stringify({
          sourceType: "OIDC",
          sourceNamespace,
          sourceName,
          inputOIDC: {
            state: OIDCstate,
            code: OIDCcode,
          },
          cookie: shouldRedirect,
          cookieDomain: window.location.hostname,
          audience,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then(res => {
          // Clear the state and code
          history.replaceState(null, "", window.location.pathname)
          return handleIssueResponse(shouldRedirect)(res)
        })
        .then(onToken)
    }
  }, [
    OIDCstate,
    OIDCcode,
    sourceNamespace,
    sourceName,
    shouldRedirect,
    handleIssueResponse,
    onToken,
  ])

  if (OIDCstate && OIDCcode) {
    return <Typography>Authenticating using OIDC...</Typography>
  }

  if (dialogState.type === "Cloak") {
    // Render identities for cloak mode
    let identities: string[] = []
    try {
      const decoded = jwtDecode(dialogState.token) as Record<string, any>
      if (Array.isArray(decoded.identity)) {
        // Dedupe
        identities = [...new Set(decoded.identity)].filter(
          // the @source:xxx tags should not be presented
          identity => !identity.startsWith("@source:")
        )
      }
    } catch (e) {
      console.error(e)
      return <Typography>Failed to parse the token: {e}</Typography>
    }

    return (
      <CloakDialog
        identities={identities}
        onConfirm={cloak => {
          issueWithA3s({
            cloak,
            token: dialogState.token,
            cookie: !isQrCodeMode,
          })
            .then(handleIssueResponse(!isQrCodeMode))
            .then(token => {
              token &&
                setDialogState({
                  type: "QrCode",
                  token,
                })
            })
        }}
      />
    )
  }

  if (dialogState.type === "QrCode") {
    return (
      <QrCodeDialog
        token={dialogState.token}
        onClose={() => {
          setDialogState({ type: "None" })
        }}
      />
    )
  }

  return (
    <Box
      sx={{
        "@media screen and (min-width: 600px)": {
          display: "flex",
          // Avoid vertical position shift of the auth sources when switching between them.
          minHeight: "400px",
          alignItems: "flex-start",
        },
        p: 2
      }}
    >
      <FormControl component="fieldset">
        <FormLabel>Authentication Source</FormLabel>
        <RadioGroup
          value={sourceType}
          onChange={e => {
            setSourceType(e.target.value)
          }}
        >
          <FormControlLabel value="MTLS" control={<Radio />} label="MTLS" />
          <FormControlLabel value="LDAP" control={<Radio />} label="LDAP" />
          <FormControlLabel value="OIDC" control={<Radio />} label="OIDC" />
          <FormControlLabel value="QR" control={<Radio />} label="QR Code" />
        </RadioGroup>
      </FormControl>
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          "@media screen and (max-width: 599px)": {
            borderTop: `1px solid ${theme.palette.divider}`,
            pt: 2,
            mt: 2,
          },
          "@media screen and (min-width: 600px)": {
            borderLeft: `1px solid ${theme.palette.divider}`,
            pl: 2,
            ml: 2,
          },
          "& .MuiTextField-root": { mt: 1, mb: 1, width: "32ch" },
        }}
      >
        {sourceType === "QR" && <QrScan />}
        {sourceType !== "QR" && (
          <>
            <TextField
              label="Source Namespace"
              value={sourceNamespace}
              onChange={e => {
                setSourceNamespace(e.target.value)
              }}
            />
            <TextField
              label="Source Name"
              value={sourceName}
              placeholder={`The name of the ${sourceType} source`}
              onChange={e => {
                setSourceName(e.target.value)
              }}
              InputLabelProps={{
                shrink: true,
              }}
            />
            {sourceType === "LDAP" && (
              <>
                <TextField
                  label="LDAP Username"
                  value={ldapUsername}
                  onChange={e => {
                    setLdapUsername(e.target.value)
                  }}
                  InputLabelProps={{
                    shrink: true,
                  }}
                />
                <TextField
                  label="LDAP Password"
                  value={ldapPassword}
                  onChange={e => {
                    setLdapPassword(e.target.value)
                  }}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  type="password"
                />
              </>
            )}
            <FormControlLabel
              control={
                <Checkbox
                  checked={cloak === "true"}
                  onChange={e => {
                    setCloak(e.target.checked ? "true" : "false")
                  }}
                />
              }
              label="Cloak claims"
            />
            <Button
              onClick={() => {
                sourceType === "MTLS"
                  ? issueWithMtls({
                      sourceNamespace,
                      sourceName,
                      cookie: shouldRedirect,
                    })
                      .then(handleIssueResponse(shouldRedirect))
                      .then(onToken)
                  : sourceType === "LDAP"
                  ? issueWithLdap({
                      sourceNamespace,
                      sourceName,
                      username: ldapUsername,
                      password: ldapPassword,
                      cookie: shouldRedirect,
                    })
                      .then(handleIssueResponse(shouldRedirect))
                      .then(onToken)
                  : issueWithOidc({
                      sourceNamespace,
                      sourceName,
                      redirectUrl,
                    })
              }}
              variant="outlined"
              sx={{
                mt: 2,
              }}
            >
              Sign in
            </Button>
          </>
        )}
      </Box>
    </Box>
  )
}
