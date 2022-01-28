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
          "@media screen and (min-width: 600px)": {
            alignItems: "center",
            justifyContent: "center",
            display: "flex",
          },
          height: "100vh",
          bgcolor: "background.default",
          color: "text.primary",
        }}
      >
        <Login />
      </Box>
    </ThemeProvider>
  )
}

render(<App />, document.getElementById("root"))
