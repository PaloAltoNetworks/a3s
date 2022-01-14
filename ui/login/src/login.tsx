import { useCallback, useState } from "react"
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
} from "@mui/material"
import { useIssue } from "./use-issue"
import { CloakDialog } from "./cloak-dialog"
import jwtDecode from "jwt-decode"
import { useLocalState } from "./use-local-state"
import { QrCodeDialog } from "./qr-code-dialog"

type StringBoolean = "true" | "false"
type DialogState =
  | { type: "None" }
  | { type: "Cloak"; token: string }
  | { type: "QrCode"; token: string }

const redirectUrl = "__REDIRECT_URL__"
// const redirectUrl = ""
const enableCloak = "__ENABLE_CLOAKING__" as StringBoolean
const audience = ["__AUDIENCE__"]
// const audience = ["https://127.0.0.1:44443"]

export const Login = () => {
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

  // @ts-ignore
  const isQrCodeMode = redirectUrl === ""
  const shouldRedirectImmediately = !isQrCodeMode && cloak === "false"
  const onToken = useCallback(
    token => {
      if (shouldRedirectImmediately) {
        return
      }
      if (cloak === "true") {
        setDialogState({ type: "Cloak", token })
      } else if (isQrCodeMode) {
        setDialogState({ type: "QrCode", token })
      }
    },
    [cloak, isQrCodeMode]
  )

  const {
    issueWithLdap,
    issueWithMtls,
    issueWithOidc,
    issueWithA3s,
    oidcIssuing,
  } = useIssue({
    apiUrl: "__API_URL__",
    // apiUrl: "https://localhost:44443",
    audience,
    oidcRedirectUrl: shouldRedirectImmediately ? redirectUrl : undefined,
    onOidcSuccess: onToken,
  })

  if (oidcIssuing) {
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
            redirectUrl: isQrCodeMode ? undefined : redirectUrl,
          }).then(token => {
            token && setDialogState({ type: "QrCode", token })
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
        display: "flex",
        // Avoid vertical position shift of the auth sources when switching between them.
        minHeight: "400px",
        alignItems: "flex-start",
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
        </RadioGroup>
      </FormControl>
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          "& .MuiTextField-root": { mt: 1, mb: 1, width: "32ch" },
          borderLeft: "1px solid #ccc",
          pl: 2,
          ml: 2,
        }}
      >
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
                  redirectUrl: shouldRedirectImmediately
                    ? redirectUrl
                    : undefined,
                }).then(onToken)
              : sourceType === "LDAP"
              ? issueWithLdap({
                  sourceNamespace,
                  sourceName,
                  username: ldapUsername,
                  password: ldapPassword,
                  redirectUrl: shouldRedirectImmediately
                    ? redirectUrl
                    : undefined,
                }).then(onToken)
              : issueWithOidc({
                  sourceNamespace,
                  sourceName,
                })
          }}
          variant="outlined"
          sx={{
            mt: 2,
          }}
        >
          Sign in
        </Button>
      </Box>
    </Box>
  )
}
