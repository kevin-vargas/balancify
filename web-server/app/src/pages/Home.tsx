import { Box, Button, Container, Typography, useTheme } from "@mui/material";
import { useEffect, useState } from "react";
import Config from "../config"
import Cookies from 'js-cookie';
import { Navigate } from "react-router-dom";
import { IconButton } from '@mui/material';
import { GitHub } from "@mui/icons-material";

const makeRef = (url: string) => {
  const from = encodeURI(window.location.href)
  const to = encodeURI(url)
  //return `https://github.com/login/oauth/authorize?client_id=Ov23litmboG5YxceI3Mv&redirect_uri=${to}&scope=read:user,user:email&state=${from}`
  return `https://github.com/login/oauth/authorize?client_id=Ov23litbrjOBo0urgKmX&redirect_uri=${to}&scope=read:user,user:email&state=${from}`
}

export default function HomePage() {
    const theme = useTheme();
    const isLogged = Cookies.get("logged")
      if (isLogged && isLogged == "true" ) {
        return <Navigate to="/csv" />
      }
    const [href, setHref] = useState("")
    useEffect(() => {
      setHref(makeRef(Config.authorizeUri))    
    },[])
    return (
        <Container
        maxWidth="lg"
        sx={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          height: "100vh",
        }}
      >
        <Box sx={{ mb: 5, mt: -10 }}>
        </Box>
        <Typography
          sx={{
            textAlign: "center",
            marginTop: "-4rem",
            fontSize: "5rem",
            fontWeight: 700,
            letterSpacing: "-0.5rem",
            display: "inline-block",
            whiteSpace: "nowrap",
            [theme.breakpoints.down("sm")]: {
              fontSize: "4rem",
              letterSpacing: "-0.4rem",
            },
          }}
          gutterBottom
        >
          Balancify
        </Typography>
        
          <IconButton size="large" href={href}>
            Github Login
              <GitHub/>
            </IconButton>            
      </Container>
    );
  }