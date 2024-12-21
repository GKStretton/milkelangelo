import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { toast } from "sonner";
import { EbsState } from "../types";

// Create axios instance with default config
const api = axios.create({
	baseURL:
		process.env.NODE_ENV === "development"
			? "http://localhost:8789"
			: process.env.EBS_URL,
});

export function setupApiAuth(auth: Twitch.ext.Authorized) {
	api.interceptors.request.clear();
	api.interceptors.request.use((config) => {
		config.headers.Authorization = `Bearer ${auth.token}`;
		config.headers["X-Twitch-Extension-Client-Id"] = auth.clientId;
		return config;
	});
}

// API functions
interface GoToBody {
	x: number;
	y: number;
}

interface DispenseBody {
	x: number;
	y: number;
}

interface CollectionBody {
	id: number;
}

// Mutation hooks
export function useCollect() {
	return useMutation({
		mutationFn: (vialPos: number) => api.post("/collect", { id: vialPos }),
		onSuccess: () => toast.success("Collection successful"),
		onError: (error) => toast.error(`Collection failed: ${error.message}`),
	});
}

export function useGoTo() {
	return useMutation({
		mutationFn: ({ x, y }: GoToBody) => api.put("/goto", { x, y }),
		onSuccess: () => toast.success("Movement successful"),
		onError: (error) => toast.error(`Movement failed: ${error.message}`),
	});
}

export function useDispense() {
	return useMutation({
		mutationFn: ({ x, y }: DispenseBody) => api.post("/dispense", { x, y }),
		onSuccess: () => toast.success("Dispense successful"),
		onError: (error) => toast.error(`Dispense failed: ${error.message}`),
	});
}

export function useClaim() {
	return useMutation({
		mutationFn: () => api.put("/claim-control"),
		onSuccess: () => toast.success("Control claimed successfully"),
		onError: (error) => toast.error(`Claim failed: ${error.message}`),
	});
}

export function useDirectEbsState(enabled: boolean) {
	return useQuery({
		queryKey: ["state"],
		queryFn: async () => {
			const { data } = await api.get<EbsState>("/direct-state");
			return data;
		},
		enabled,
		refetchInterval: 1000,
	});
}
