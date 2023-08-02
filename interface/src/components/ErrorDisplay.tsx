import { Alert, Snackbar } from "@mui/material";
import { useState } from "react";

export default function ErrorDisplay() {
  const [errorMessage, setErrorMessage] = useState("test");

  const handleErrorClose = () => {
    setErrorMessage("");
  };

  return (
    <Snackbar open={errorMessage !== ""} autoHideDuration={6000} onClose={handleErrorClose}>
      <Alert onClose={handleErrorClose} severity="error" sx={{ width: "100%" }}>
        {errorMessage}
      </Alert>
    </Snackbar>
  );
}
