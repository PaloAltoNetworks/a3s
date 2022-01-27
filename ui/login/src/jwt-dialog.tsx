import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  useMediaQuery,
} from "@mui/material"
import { Box } from "@mui/system"
import type { DecodedJWT } from "./types"

export const JwtDialog = ({
  payload,
  header,
  onClose,
}: {
  onClose(): void
} & DecodedJWT) => {
  const fullScreen = useMediaQuery("(max-width: 600px)")
  return (
    <Dialog open fullScreen={fullScreen}>
      <DialogTitle>Decoded JWT</DialogTitle>
      <DialogContent>
        <DialogContentText>Payload</DialogContentText>
        <Box
          sx={{
            bgcolor: "action.hover",
            borderRadius: 2,
            overflowX: "auto",
            px: 2,
            mt: 1
          }}
        ><pre>
          {JSON.stringify(payload, null, 2)}
          </pre></Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="contained">Close</Button>
      </DialogActions>
    </Dialog>
  )
}
