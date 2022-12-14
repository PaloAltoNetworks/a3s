import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
} from "@mui/material"
import { ReactNode } from "react"
import QrCode from "../utils/qr-code-gen"

/**
 * Encode data into QR code
 */
export const QrCodeDialog = ({
  data,
  title,
  children,
  onClose,
}: {
  data: string
  title?: string
  children?: ReactNode
  onClose(): void
}) => {
  return (
    <Dialog open>
      {title && <DialogTitle>{title}</DialogTitle>}
      <DialogContent sx={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
        {children}
        <canvas
          ref={canvas => {
            if (canvas) {
              const dpi = window.devicePixelRatio
              const qr = QrCode.encodeText(data, QrCode.Ecc.LOW)
              const border = 3
              const scale = 3
              const width: number = (qr.size + border * 2) * scale
              canvas.width = width * dpi
              canvas.height = width * dpi
              canvas.style.width = `${width}px`
              canvas.style.height = `${width}px`
              let ctx = canvas.getContext("2d") as CanvasRenderingContext2D
              ctx.scale(dpi, dpi)
              for (let y = -border; y < qr.size + border; y++) {
                for (let x = -border; x < qr.size + border; x++) {
                  ctx.fillStyle = qr.getModule(x, y) ? "#000000" : "#FFFFFF"
                  ctx.fillRect(
                    (x + border) * scale,
                    (y + border) * scale,
                    scale,
                    scale
                  )
                }
              }
            }
          }}
        />
      </DialogContent>
      <DialogActions>
        <Button variant="outlined" onClick={onClose}>
          Close
        </Button>
      </DialogActions>
    </Dialog>
  )
}
