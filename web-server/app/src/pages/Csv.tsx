import { CloudUpload } from "@mui/icons-material";
import useQueryUser from "../hooks/useQueryUser";
import {Card, CardMedia, CardContent, Typography, Container,Button, Snackbar, Alert} from "@mui/material";
import { styled } from '@mui/material/styles';
import Config from "../config"

import { useState } from "react";
import Loading from "../components/Loading";

const VisuallyHiddenInput = styled('input')({
  clip: 'rect(0 0 0 0)',
  clipPath: 'inset(50%)',
  height: 1,
  overflow: 'hidden',
  position: 'absolute',
  bottom: 0,
  left: 0,
  whiteSpace: 'nowrap',
  width: 1,
});

export default function UploadPage() {

  const [file, setFile] = useState<File | null>(null);
  const [openSuccess, setOpenSuccess] = useState(false);
  const [openFail, setOpenFail] = useState(false);
  const [isProgress, setIsProgress] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  const handleCloseSucess = (event?: React.SyntheticEvent | Event, reason?: string) => {
    if (reason === 'clickaway') {
      return;
    }
    setOpenSuccess(false);
  };

  const handleCloseFail= (event?: React.SyntheticEvent | Event, reason?: string) => {
    if (reason === 'clickaway') {
      return;
    }
    setOpenFail(false);
  };

  const handleOnClick = (event: any) => {
    event.target.value = null
  };


  const handleUpload = async () => {
    if (file) {
      const formData = new FormData();
      formData.append("file", file)
      try {
        setIsProgress(true)
        const result = await fetch(Config.uploadCsvUri, {
          credentials: 'include',
          method: 'POST',
          body: formData,
        });
        if (result.status == 200) {
          setOpenSuccess(true)
        } else {
          setOpenFail(true)
        }
      } catch (error) {
        console.error(error);
        setOpenFail(true)
      } finally {
        setIsProgress(false)
        setFile(null)
      }
    }
  };
    const result = useQueryUser()
    if (!result) {
      return null
    }

    return (
      <>
      {!isProgress ? <>{result &&  <Container
        maxWidth="lg"
        sx={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
          height: "100vh",
        }}
      >
        <Card sx={{ maxWidth: 700 }}>
      <CardMedia
        component="img"
        sx={{ width: 300 }}
        
        image={result.data.avatar}
        title="user avatar"
      />
      <CardContent >
        <Typography gutterBottom variant="h5" component="div">
          Email:
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {result.data.email}
        </Typography>
        <Button
  component="label"
  role={undefined}
  variant="contained"
  tabIndex={-1}
  startIcon={<CloudUpload />}
>
  {!!file ? "Change File": "Load file"}
  <VisuallyHiddenInput accept=".csv" type="file" onChange={handleFileChange} onClick={handleOnClick}/>
</Button>
{file && (
        <section>
          File details:
          <ul>
            <li>Name: {file.name}</li>
            <li>Type: {file.type}</li>
            <li>Size: {file.size} bytes</li>
          </ul>
        </section>
      )}
{file && <Button variant="contained" onClick={handleUpload}>Upload File</Button>}
      </CardContent>
    </Card>
    <Snackbar anchorOrigin={{vertical: "top", horizontal: "right"}} open={openSuccess} autoHideDuration={5000} onClose={handleCloseSucess}>
        <Alert
          onClose={handleCloseSucess}
          severity="success"
          sx={{ width: '100%' }}
        >
          The CSV file was processed successfully, and the report was sent via email.
        </Alert>
      </Snackbar>
      <Snackbar anchorOrigin={{vertical: "top", horizontal: "right"}} open={openFail} autoHideDuration={5000} onClose={handleCloseFail}>
        <Alert
          onClose={handleCloseFail}
          severity="error"
          sx={{ width: '100%' }}
        >
          Error processing the CSV file, please try again later.
        </Alert>
      </Snackbar>
    </Container>} </>: <Loading />}
      </>
    );
  }