import {
  Alert,
  AlertTitle,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  Divider,
  Stack,
  Typography,
  useMediaQuery,
} from "@mui/material"
import { useEffect, useState } from "react"
import { useVerifyJwt } from "./use-verify-jwt"
import { formatDuration, intervalToDuration } from "date-fns"

export const JwtDialog = ({
  payload,
  token,
  onClose,
}: {
  payload: Record<string, any>
  token: string
  onClose(): void
}) => {
  const fullScreen = useMediaQuery("(max-width: 600px)")
  const [isValid, setIsValid] = useState<boolean>()
  const [error, setError] = useState<string>()
  const [currentTimeInMs, setCurrentTimeInMs] = useState<number>(Date.now())

  useEffect(() => {
    const intervalTimer = setInterval(() => {
      setCurrentTimeInMs(Date.now())
    }, 1000)
    return () => {
      clearInterval(intervalTimer)
    }
  }, [])

  const expInMs = payload.exp * 1000
  const isExpired = expInMs < currentTimeInMs
  const lifetime = isExpired
    ? ""
    : formatDuration(
        intervalToDuration({
          start: currentTimeInMs,
          end: expInMs,
        })
      )

  const verifyJwt = useVerifyJwt()
  useEffect(() => {
    if (!isExpired) {
      verifyJwt({ token })
        .then(isValid => {
          setIsValid(isValid)
          if (!isValid) {
            throw Error(`Signature doesn't match`)
          }
        })
        .catch(err => {
          setIsValid(false)
          if (err instanceof Error) {
            setError(err.message)
          }
        })
    } else {
      setIsValid(false)
      setError("Token has expired")
    }
  }, [token, isExpired])

  return (
    <Dialog open fullScreen={fullScreen}>
      {/* <DialogTitle>Decoded JWT</DialogTitle> */}
      <DialogContent sx={{ mt: 1 }}>
        {isValid === undefined ? null : isValid ? (
          <Alert variant="outlined" severity="success">
            <AlertTitle>Signature Verified</AlertTitle>
            Issuer: {payload.iss}
          </Alert>
        ) : (
          <Alert variant="outlined" severity="error">
            <AlertTitle>Signature Verification Failed</AlertTitle>
            {error}
          </Alert>
        )}
        <Typography variant="h6" mt={3} mb={1}>
          Claims
        </Typography>
        <Stack
          direction="column"
          spacing={0.5}
          divider={<Divider />}
          py={1}
          px={2}
          borderRadius={1}
          overflow="auto"
          bgcolor="action.hover"
        >
          {payload.identity.map((claim: string) => (
            <Typography key={claim}>{claim}</Typography>
          ))}
        </Stack>
        <Typography variant="h6" mt={3} mb={1}>
          Audiences
        </Typography>
        <Stack
          direction="column"
          spacing={0.5}
          divider={<Divider />}
          py={1}
          px={2}
          borderRadius={1}
          overflow="auto"
          bgcolor="action.hover"
        >
          {payload.aud.map((audience: string) => (
            <Typography key={audience}>{audience}</Typography>
          ))}
        </Stack>
        {!isExpired && (
          <DialogContentText mt={3}>
            Token will expire in {lifetime}.
          </DialogContentText>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="contained">
          Close
        </Button>
      </DialogActions>
    </Dialog>
  )
}
