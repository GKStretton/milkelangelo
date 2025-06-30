import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import ReactDOM from "react-dom/client";
import "./render.css";
import { StateProvider } from "~/helpers/State";

const queryClient = new QueryClient();

const render = (node: React.ReactNode) => {
	const app = document.getElementById("app");

	return ReactDOM.createRoot(app).render(
		<React.StrictMode>
			<StateProvider>
				{/* <TwitchProvider> */}
				<QueryClientProvider client={queryClient}>{node}</QueryClientProvider>
				{/* </TwitchProvider> */}
			</StateProvider>
		</React.StrictMode>,
	);
};

export default render;
