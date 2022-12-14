import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField,
} from "@mui/material"
import { useState } from "react"
import { RequestEntry } from "./types"

export const EditDialog = ({
  initEntry,
  onChange,
  onCancel,
}: {
  initEntry: Partial<RequestEntry>
  onChange(newEntry: RequestEntry): void
  onCancel(): void
}) => {
  const [name, setName] = useState(initEntry.name)
  const [claims, setClaims] = useState(initEntry.claims)
  const [description, setDescription] = useState(initEntry.description)
  const [message, setMessage] = useState(initEntry.message)
  const [issuers, setIssuers] = useState(initEntry.issuers)
  const isUpdate = initEntry.ID !== undefined

  return (
    <Dialog open>
      <DialogTitle>{isUpdate ? "Edit" : "New"} Entry</DialogTitle>
      <DialogContent>
        <DialogContentText mb={2}>An entry is an entry.</DialogContentText>
        <TextField
          required
          label="Name"
          value={name}
          onChange={e => {
            setName(e.target.value)
          }}
          fullWidth
          margin="dense"
        />
        <TextField
          label="Description"
          value={description}
          onChange={e => {
            setDescription(e.target.value)
          }}
          fullWidth
          margin="dense"
        />
        <TextField
          required
          label="Claims"
          multiline
          maxRows={6}
          value={claims?.join("\n")}
          onChange={e => {
            const value = e.target.value
            const splitted = value.split("\n")
            setClaims(splitted)
          }}
          helperText="Required claim prefixes, separated by newlines."
          fullWidth
          margin="dense"
        />
        <TextField
          label="Issuers"
          multiline
          maxRows={4}
          value={issuers?.join("\n")}
          onChange={e => {
            const value = e.target.value
            const splitted = value.split("\n")
            setIssuers(splitted)
          }}
          helperText="Trusted issuers, separated by newlines."
          fullWidth
          margin="dense"
        />
        <TextField
          label="Message"
          value={message}
          onChange={e => {
            setMessage(e.target.value)
          }}
          helperText="An optional message to be displayed to the client."
          fullWidth
          margin="dense"
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onCancel}>Cancel</Button>
        <Button
          disabled={!name || !claims?.length}
          variant="contained"
          onClick={() => {
            onChange({
              name: name!,
              claims: claims!.map(s => s.trim()).filter(s => s.length > 0),
              description: description || "",
              message: message || "",
              issuers: issuers || [],
              ID: isUpdate ? initEntry.ID! : `${Date.now()}`,
            })
          }}
        >
          {isUpdate ? "Update" : "Create"}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
