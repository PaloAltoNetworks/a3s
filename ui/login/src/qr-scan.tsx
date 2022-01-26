import { Box, Typography } from "@mui/material"
import { useEffect, useRef, useState } from "react"
import QrScanner from "qr-scanner"

export const QrScan = () => {
  const videoRef = useRef<HTMLVideoElement>(null)
  const [rawString, setRawString] = useState("")
  useEffect(() => {
    const qrScanner = new QrScanner(
      videoRef.current!,
      result => {
        result !== rawString && setRawString(result)
      },
      undefined,
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
    <Box         style={{
      width: "40vw",
    }}>
      <video
        ref={videoRef}
        style={{
          width: '100%',
        }}
      ></video>
      {rawString && <Typography sx={{wordWrap: 'break-word'}}>{rawString}</Typography>}
    </Box>
  )
}
