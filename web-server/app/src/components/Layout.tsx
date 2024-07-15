import Header from "./Header"
import { Outlet, useLocation } from "react-router-dom";

import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { useState, useMemo, useEffect } from "react";
import { PaletteMode } from "@mui/material";

const key = "mode"

export default () => {
  const modeLoaded = (localStorage.getItem(key) || 'dark' ) as PaletteMode
  const [mode, setMode] = useState<PaletteMode>(modeLoaded)
  const switchMode = () => setMode(prev => prev === 'light' ? 'dark' : 'light')
  
  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode,
        },
      }),
    [mode],
  );
  const location = useLocation();
  
  useEffect(() => {
    localStorage.setItem(key, mode)
  }, [mode])

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      { !["/", "/home"].includes(location.pathname) ? <Header mode={mode} switchMode={switchMode}/> : null }
      <Outlet />
  </ThemeProvider>
  )
}
  