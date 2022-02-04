import { render } from "react-dom"
import { Box } from "@mui/system"
import { Login } from "./login"
import { createTheme, ThemeProvider, useMediaQuery } from "@mui/material"
import { useMemo } from "react"
import { RequestPage } from "./request/request-page"
import "./main.css"

const App = () => {
  // https://mui.com/customization/dark-mode/#system-preference
  const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)")
  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode: prefersDarkMode ? "dark" : "light",
        },
      }),
    [prefersDarkMode]
  )

  const pathname = window.location.pathname.toLowerCase()
  const isRequestMode =
    pathname.endsWith("request.html") || pathname.endsWith("request")

  return (
    <ThemeProvider theme={theme}>
      <Box
        sx={{
          height: "100vh",
          bgcolor: "background.default",
          color: "text.primary",
          display: "flex",
          flexDirection: "column",
          "@media screen and (min-width: 600px)": {
            alignItems: "center",
          },
        }}
      >
        {isRequestMode ? <RequestPage /> : <Login />}
      </Box>
    </ThemeProvider>
  )
}

render(<App />, document.getElementById("root"))
