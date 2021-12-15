import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  FormGroup,
  FormControlLabel,
  Checkbox,
} from "@mui/material"
import { useState } from "react"

export const CloakDialog = ({
  identities,
  onConfirm,
}: {
  // Unique list of identities
  identities: string[]
  onConfirm(selected: string[]): void
}) => {
  const [selected, setSelected] = useState<string[]>(identities)
  return (
    <Dialog open={!!identities.length}>
      <DialogTitle>Select Claims</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Please select the claims that you want to include.
        </DialogContentText>
        <FormGroup>
          {identities.map(identity => (
            <FormControlLabel
              control={
                <Checkbox
                  checked={selected.includes(identity)}
                  onChange={e => {
                    if (e.target.checked) {
                      setSelected([...selected, identity])
                    } else {
                      setSelected(selected.filter(i => i !== identity))
                    }
                  }}
                />
              }
              label={identity}
              key={identity}
            />
          ))}
        </FormGroup>
      </DialogContent>
      <DialogActions>
        <Button
          variant="outlined"
          onClick={() => {
            onConfirm(selected)
          }}
          autoFocus
        >
          Confirm
        </Button>
      </DialogActions>
    </Dialog>
  )
}
