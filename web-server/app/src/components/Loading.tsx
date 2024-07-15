import CircularProgress from "@mui/material/CircularProgress";
import Backdrop from '@mui/material/Backdrop';

export default function Loading() {
    return (
        <Backdrop open >
            <CircularProgress size={150}/>
        </Backdrop>
    );
  }