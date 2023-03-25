import { Button } from "@mui/material";
import { useEffect } from "react";

let number = 0;
function increment() {
  number++;
}
export default function Test() {
  useEffect(() => {
    console.log("number changed");
  }, [number]);

  return (
    <>
    {number}
    <Button onClick={increment}>Increment</Button>
    </>
  )
}
