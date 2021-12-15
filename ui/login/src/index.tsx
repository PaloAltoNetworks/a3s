import { useState } from "react"
import { render } from "react-dom"
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
} from "@mui/material"
import { useIssue } from "./useIssue"
import { CloakDialog } from "./cloak-dialog"
import jwtDecode from "jwt-decode"

type StringBoolean = "true" | "false"

const App = () => {
  const [cloak, setCloak] = useState("__ENABLE_CLOAKING__" as StringBoolean)
  const [sourceType, setSourceType] = useState("MTLS")
  const [sourceNamespace, setSourceNamespace] = useState("/")
  const [sourceName, setSourceName] = useState("")
  const [ldapUsername, setLdapUsername] = useState("")
  const [ldapPassword, setLdapPassword] = useState("")
  const { issueWithLdap, issueWithMtls, issueWithOidc, issueWithA3s, token } =
    useIssue({
      apiUrl: "__API_URL__",
      // apiUrl: "https://localhost:44443",
      redirectUrl: "__REDIRECT_URL__",
      audience: ["__AUDIENCE__"],
      // audience: ["https://127.0.0.1:44443"],
      saveToken: cloak === "true",
    })

  let identities: string[] = []
  if (token) {
    try {
      const decoded = jwtDecode(token) as Record<string, any>
      if (Array.isArray(decoded.identity)) {
        identities = [...new Set(decoded.identity)]
      }
      console.log(identities)
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <Box
      sx={{
        display: "flex",
        height: "100vh",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Box
        sx={{
          display: "flex",
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
      {!!identities.length && (
        <CloakDialog
          identities={identities}
          onConfirm={cloak => {
            issueWithA3s({ token: token!, cloak })
          }}
        />
      )}
    </Box>
  )
}

render(<App />, document.getElementById("root"))
