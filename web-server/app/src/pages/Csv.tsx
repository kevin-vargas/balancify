import { CloudUpload } from "@mui/icons-material";
import useQueryUser from "../hooks/useQueryUser";
import {Card, CardMedia, CardContent, Typography, Container,Button} from "@mui/material";
import { styled } from '@mui/material/styles';
import Config from "../config"

import { useState } from "react";

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

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (file) {
      const formData = new FormData();
      formData.append("file", file)
      try {
        const result = await fetch(Config.uploadCsvUri, {
          credentials: 'include',
          method: 'POST',
          body: formData,
        });

        const data = await result.json();

        console.log(data);
      } catch (error) {
        console.error(error);
      }
    }
  };
    const result = useQueryUser()
    if (!result) {
      return null
    }

    return (
      <>
      {result &&  <Container
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
  Load file
  <VisuallyHiddenInput type="file" onChange={handleFileChange} />
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
{file && <Button variant="contained" onClick={handleUpload}>Upload</Button>}
      </CardContent>
    </Card>
    </Container>}
      </>
    );
  }