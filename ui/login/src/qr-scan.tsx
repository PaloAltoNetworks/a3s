import { Box, Typography } from "@mui/material"
import { useEffect, useRef, useState } from "react"
import QrScanner from "qr-scanner"

import { JwtDialog } from "./jwt-dialog"
// @ts-ignore
import qrScannerWorkerSource from "qr-scanner/qr-scanner-worker.min.js?raw"
QrScanner.WORKER_PATH = URL.createObjectURL(new Blob([qrScannerWorkerSource]))

export const QrScan = ({
  onResult,
}: {
  /**
   * Must throw if the result is not usable so the scan can keep going
   */
  onResult(result: string): void
}) => {
  const videoRef = useRef<HTMLVideoElement>(null)

  // const [error, setError] = useState<string>()
  useEffect(() => {
    const qrScanner = new QrScanner(
      videoRef.current!,
      result => {
        try {
          onResult(result)
          qrScanner.stop()
        } catch (e) {
          console.error(e)
          // setError(`${e}`)
        }
      },
      err => {
        console.error(err)
        // setError(err)
      },
      video => {
        return {
          x: 0,
          y: 0,
          width: video.videoWidth,
          height: video.videoHeight,
          downScaledWidth: video.videoWidth,
          downScaledHeight: video.videoHeight,
        }
      }
    )
    qrScanner.start()
    return () => {
      qrScanner.stop()
    }
  }, [])
  return (
    <video
      ref={videoRef}
      style={{
        width: "100%",
      }}
    ></video>
  )
}
