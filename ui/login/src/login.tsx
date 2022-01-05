import { useState } from "react"
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
import { useIssue } from "./useIssue"
import { CloakDialog } from "./cloak-dialog"
import jwtDecode from "jwt-decode"

type StringBoolean = "true" | "false"

export const Login = () => {
  const [cloak, setCloak] = useState("__ENABLE_CLOAKING__" as StringBoolean)
  const [sourceType, setSourceType] = useState("MTLS")
  const [sourceNamespace, setSourceNamespace] = useState("/")
  const [sourceName, setSourceName] = useState("")
  const [ldapUsername, setLdapUsername] = useState("")
  const [ldapPassword, setLdapPassword] = useState("")

  const params = new URLSearchParams(window.location.search)
  const OIDCstate = params.get("state")
  const OIDCcode = params.get("code")
  const { issueWithLdap, issueWithMtls, issueWithOidc, issueWithA3s, token } =
    useIssue({
      apiUrl: "__API_URL__",
      // apiUrl: "https://localhost:44443",
      redirectUrl: "__REDIRECT_URL__",
      audience: ["__AUDIENCE__"],
      // audience: ["https://127.0.0.1:44443"],
      saveToken: cloak === "true",
      OIDCstate,
      OIDCcode,
    })

  // Render identities for cloak mode
  let identities: string[] = []
  if (token) {
    try {
      const decoded = jwtDecode(token) as Record<string, any>
      if (Array.isArray(decoded.identity)) {
        // Dedupe
        identities = [...new Set(decoded.identity)]
      }
    } catch (e) {
      console.error(e)
    }
  }

  let mainContent: JSX.Element
  if (OIDCstate && OIDCcode) {
    mainContent = <Typography>Authenticating using OIDC...</Typography>
  } else {
    mainContent = (
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
                  })
                : sourceType === "LDAP"
                ? issueWithLdap({
                    sourceNamespace,
                    sourceName,
                    username: ldapUsername,
                    password: ldapPassword,
                  })
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

  return (
    <>
      {mainContent}
      {!!identities.length && (
        <CloakDialog
          identities={identities}
          onConfirm={cloak => {
            issueWithA3s({ token: token!, cloak })
          }}
        />
      )}
    </>
  )
}
