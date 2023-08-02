import { Alert, Snackbar } from "@mui/material";
import { ReactNode, createContext, useCallback, useContext, useState } from "react";

const ErrorContext = createContext<ErrorContextType>((msg) => {});

type ErrorContextType = (msg: string) => void;

export function useError(): ErrorContextType {
  return useContext(ErrorContext);
}

interface ErrorManagerProps {
  children: ReactNode;
}
export function ErrorManager({ children }: ErrorManagerProps) {
  const [errorMessage, setErrorMessage] = useState("test");
  const [errorShowing, setErrorShowing] = useState(false);

  const error = useCallback((msg: string) => {
    if (errorMessage === msg && errorShowing) return;
    setErrorMessage(msg);
    setErrorShowing(true);
    console.error(msg);
  }, []);

  const handleErrorClose = () => {
    setErrorShowing(false);
  };

  return (
    <>
      <ErrorContext.Provider value={error}>{children}</ErrorContext.Provider>
      <Snackbar open={errorShowing} autoHideDuration={6000} onClose={handleErrorClose}>
        <Alert onClose={handleErrorClose} severity="error" sx={{ width: "100%" }}>
          {errorMessage}
        </Alert>
      </Snackbar>
    </>
  );
}
