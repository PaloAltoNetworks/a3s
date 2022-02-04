import { Button, Card, IconButton, Typography } from "@mui/material"
import { Box } from "@mui/system"
import { useLocalJsonState } from "../utils/use-local-json-state"
import { RequestEntry } from "./types"
import { useState } from "react"
import { EditDialog } from "./edit-dialog"
import { Delete, Edit, Add } from "@mui/icons-material"

export const RequestPage = () => {
  const [entries, setEntries] = useLocalJsonState<RequestEntry[]>(
    [],
    "requestEntries"
  )
  const [editEntryIndex, setEditEntryIndex] = useState<number>()
  return (
    <Box
      sx={{
        height: "calc(100vh - 64px)",
        "@media screen and (min-width: 600px)": {
          width: 800,
        },
        width: "100%",
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
            <IconButton
              aria-label="edit"
              onClick={() => {
                setEditEntryIndex(index)
              }}
            >
              <Edit />
            </IconButton>
            <IconButton
              aria-label="delete"
              onClick={() => {
                setEntries(entries.filter(e => e.ID !== entry.ID))
              }}
            >
              <Delete />
            </IconButton>
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
    </Box>
  )
}
