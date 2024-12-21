import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import { StateProvider } from "./helpers/State";
import "./index.css";

const queryClient = new QueryClient();

const root = ReactDOM.createRoot(
	document.getElementById("root") as HTMLElement,
);
root.render(
	<React.StrictMode>
		<QueryClientProvider client={queryClient}>
			<StateProvider>
				<App />
			</StateProvider>
		</QueryClientProvider>
	</React.StrictMode>,
);
