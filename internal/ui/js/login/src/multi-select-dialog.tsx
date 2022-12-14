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
  Divider,
} from "@mui/material"
import { useState, ReactNode } from "react"

export const MultiSelectDialog = ({
  options,
  title,
  description,
  onConfirm,
  onCancel,
  children
}: {
  title: string
  // Should be unique
  options: string[]
  description?: string
  onConfirm(selected: string[]): void
  onCancel?: () => void
  children?: ReactNode
}) => {
  const [selected, setSelected] = useState<string[]>(options)
  const allSelected = selected.length === options.length
  return (
    <Dialog open={!!options.length}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        {description && <DialogContentText>{description}</DialogContentText>}
        {children}
        <FormGroup sx={{ mt: 2 }}>
          <FormControlLabel
            control={
              <Checkbox
                checked={allSelected}
                indeterminate={
                  !allSelected && options.some(id => selected.includes(id))
                }
                onChange={e => {
                  if (e.target.checked) {
                    setSelected([...options])
                  } else {
                    setSelected([])
                  }
                }}
              />
            }
            label="Select All"
          />
          <Divider />
          {options.map(option => (
            <FormControlLabel
              control={
                <Checkbox
                  checked={selected.includes(option)}
                  onChange={e => {
                    if (e.target.checked) {
                      setSelected([...selected, option])
                    } else {
                      setSelected(selected.filter(i => i !== option))
                    }
                  }}
                />
              }
              label={option}
              key={option}
            />
          ))}
        </FormGroup>
      </DialogContent>
      <DialogActions>
        {onCancel && <Button onClick={onCancel}>Cancel</Button>}
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
