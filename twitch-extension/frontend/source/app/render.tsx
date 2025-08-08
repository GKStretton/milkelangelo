import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import ReactDOM from "react-dom/client";
import "./render.css";
import TwitchProvider from "~/context/twitch";
import { StateProvider } from "~/helpers/State";

const queryClient = new QueryClient();

const render = (node: React.ReactNode) => {
	const app = document.getElementById("app");

	return ReactDOM.createRoot(app).render(
		<React.StrictMode>
			<TwitchProvider>
				<StateProvider>
					<QueryClientProvider client={queryClient}>{node}</QueryClientProvider>
				</StateProvider>
			</TwitchProvider>
		</React.StrictMode>,
	);
};

export default render;
