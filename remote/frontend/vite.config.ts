import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

const goo = process.env.GOO_API ?? "http://127.0.0.1:8789";

// In dev, /api and /video are proxied to goo so everything shares an origin,
// as it does when goo serves the built page (PUBLIC_UI_DIR).
export default defineConfig({
	plugins: [react()],
	server: {
		port: 3000,
		strictPort: true,
		proxy: {
			"/api": goo,
			"/video": { target: goo, ws: true },
		},
	},
});
