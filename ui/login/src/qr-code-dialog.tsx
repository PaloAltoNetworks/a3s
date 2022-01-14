import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
} from "@mui/material"

export const QrCodeDialog = ({ token, onClose }: { token: string, onClose(): void }) => {
  return (
    <Dialog open={!!token}>
      <DialogTitle>QR Code</DialogTitle>
      <DialogContent>
        <DialogContentText>Below is your QR Code</DialogContentText>
        {token}
      </DialogContent>
      <DialogActions>
        <Button variant="contained" autoFocus>Save Image</Button>
        <Button variant="outlined" onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  )
}
