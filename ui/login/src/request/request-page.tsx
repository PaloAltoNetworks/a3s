import { Button, Card, IconButton, Tooltip, Typography } from "@mui/material"
import { Box } from "@mui/system"
import { useLocalJsonState } from "../utils/use-local-json-state"
import { RequestEntry } from "./types"
import { useState } from "react"
import { EditDialog } from "./edit-dialog"
import { Delete, Edit, Add, QrCode } from "@mui/icons-material"
import { QrCodeDialog } from "../components/qr-code-dialog"

export const RequestPage = () => {
  const [entries, setEntries] = useLocalJsonState<RequestEntry[]>(
    [],
    "requestEntries"
  )
  const [editEntryIndex, setEditEntryIndex] = useState<number>()
  const [qrCodeEntryIndex, setQrCodeEntryIndex] = useState<number>()

  return (
    <Box
      sx={{
        p: 3,
        width: "100%",
        maxWidth: "800px",
      }}
    >
      <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
        <Button
          startIcon={<Add />}
          variant="contained"
          onClick={() => {
            setEditEntryIndex(-1)
          }}
        >
          New entry
        </Button>
      </Box>
      {entries.map((entry, index) => (
        <Card
          key={entry.ID}
          sx={{
            mt: 2,
            p: 2,
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <Box>
            <Typography>{entry.name}</Typography>
            <Typography color="text.secondary">{entry.description}</Typography>
          </Box>
          <Box display="flex">
            <Tooltip title="Request with QR code">
              <IconButton
                aria-label="qr-code"
                onClick={() => {
                  setQrCodeEntryIndex(index)
                }}
              >
                <QrCode />
              </IconButton>
            </Tooltip>
            <Tooltip title="Edit">
              <IconButton
                aria-label="edit"
                onClick={() => {
                  setEditEntryIndex(index)
                }}
              >
                <Edit />
              </IconButton>
            </Tooltip>
            <Tooltip title="Delete">
              <IconButton
                aria-label="delete"
                onClick={() => {
                  setEntries(entries.filter(e => e.ID !== entry.ID))
                }}
              >
                <Delete />
              </IconButton>
            </Tooltip>
          </Box>
        </Card>
      ))}

      {editEntryIndex !== undefined && (
        <EditDialog
          initEntry={editEntryIndex === -1 ? {} : entries[editEntryIndex]}
          onChange={newEntry => {
            const isUpdate = editEntryIndex !== -1
            if (isUpdate) {
              setEntries([
                ...entries.slice(0, editEntryIndex),
                newEntry,
                ...entries.slice(editEntryIndex + 1),
              ])
            } else {
              setEntries([...entries, newEntry])
            }
            setEditEntryIndex(undefined)
          }}
          onCancel={() => {
            setEditEntryIndex(undefined)
          }}
        />
      )}
      {qrCodeEntryIndex !== undefined && (
        <QrCodeDialog
          data={JSON.stringify({
            meta: {
              version: "0.0.0",
            },
            claims: entries[qrCodeEntryIndex].claims,
            issuers: entries[qrCodeEntryIndex].issuers,
            message: entries[qrCodeEntryIndex].message,
          })}
          title={`Requesting ${entries[qrCodeEntryIndex].name}`}
          onClose={() => {
            setQrCodeEntryIndex(undefined)
          }}
        />
      )}
    </Box>
  )
}
