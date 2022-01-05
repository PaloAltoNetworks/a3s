import { render } from "react-dom"
import { Box } from "@mui/system"
import { Login } from "./login"
import { createTheme, ThemeProvider, useMediaQuery } from "@mui/material"
import { useMemo } from "react"

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

  return (
    <ThemeProvider theme={theme}>
      <Box
        sx={{
          display: "flex",
          height: "100vh",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          bgcolor: 'background.default',
          color: 'text.primary',
        }}
      >
        <Login />
      </Box>
    </ThemeProvider>
  )
}

render(<App />, document.getElementById("root"))
