import React, { useState } from "react"
import { render } from "react-dom"
import { Box } from "@mui/system"
import {
  Button,
  TextField,
  Typography,
  Autocomplete,
  Alert,
  Radio,
  RadioGroup,
  FormControlLabel,
  FormControl,
  FormLabel,
} from "@mui/material"
import { useAuthz } from "./useAuthz"
import { useIssue } from "./useIssue"

const resourceOptions = ["/secret", "/topsecret"]

const App = () => {
  const { token, issueWithLdap, issueWithMtls, issueWithOidc } = useIssue({
    baseUrl: "https://127.0.0.1:44443",
  })
  const [resource, setResource] = useState(resourceOptions[0])
  const [authed, setAuthed] = useState<boolean>()
  const [sourceType, setSourceType] = useState("MTLS")
  const [sourceNamespace, setSourceNamespace] = useState("/")
  const [authzNs, setAuthzNs] = useState("/")
  const [sourceName, setSourceName] = useState("")
  const [ldapUsername, setLdapUsername] = useState("")
  const [ldapPassword, setLdapPassword] = useState("")
  const authz = useAuthz({ baseUrl: "https://127.0.0.1:44443" })

  return (
    <Box
      height="100vh"
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
    >
      {/* <Typography variant="h4" sx={{mb: 2}}>A3S</Typography> */}
      {token ? (
        <Box minHeight="300px">
          <Typography>
            You are authenticated with A3S. Now you can check your privileges.
          </Typography>
          <Box
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "start",
              "& .MuiTextField-root": { mt: 1, mb: 1, width: "32ch" },
              my: 2,
            }}
          >
            <TextField
              label="Namespace"
              value={authzNs}
              onChange={e => {
                setAuthzNs(e.target.value)
              }}
            />
            <Autocomplete
              freeSolo
              options={resourceOptions}
              sx={{ flexGrow: 1, mr: 2 }}
              value={resource}
              renderInput={params => <TextField {...params} label="Resource" />}
              onChange={(event, newValue) => {
                setResource(newValue || '')
                setAuthed(undefined)
              }}
            />
            <Button
              onClick={() => {
                authz({ resource, token, namespace: authzNs }).then(setAuthed)
              }}
            >
              Check
            </Button>
          </Box>
          {authed === true && (
            <Alert severity="success">
              You are authorized to access {resource}
            </Alert>
          )}
          {authed === false && (
            <Alert severity="error">
              You are not authorized to access {resource}
            </Alert>
          )}
        </Box>
      ) : (
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
      )}
    </Box>
  )
}

render(<App />, document.getElementById("root"))
