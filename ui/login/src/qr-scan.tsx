import { Box, Typography } from "@mui/material"
import { useEffect, useRef, useState } from "react"
import QrScanner from "qr-scanner"
import jwtDecode from "jwt-decode"
import { JwtDialog } from "./jwt-dialog"
// @ts-ignore
import qrScannerWorkerSource from "qr-scanner/qr-scanner-worker.min.js?raw"
QrScanner.WORKER_PATH = URL.createObjectURL(new Blob([qrScannerWorkerSource]))

export const QrScan = () => {
  const videoRef = useRef<HTMLVideoElement>(null)
  const [jwt, setJwt] = useState<{ payload: Record<string, any>; token: string }>()
  // const [error, setError] = useState<string>()
  useEffect(() => {
    const qrScanner = new QrScanner(
      videoRef.current!,
      result => {
        try {
          const payload = jwtDecode(result) as Record<string, any>
          setJwt({ payload, token: result })
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
    <Box
      sx={{
        width: "90vw",
        "@media screen and (min-width: 600px)": {
          width: "40vw",
        },
      }}
    >
      {jwt ? (
        <JwtDialog
          {...jwt}
          onClose={() => {
            setJwt(undefined)
          }}
        />
      ) : (
        <video
          ref={videoRef}
          style={{
            width: "100%",
          }}
        ></video>
      )}
    </Box>
  )
}
