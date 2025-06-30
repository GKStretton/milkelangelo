import React, { createContext, useContext, useState, ReactNode } from "react";
import { Coords, EbsState } from "../types";

// 1. Define the shape of the global state
interface GlobalStateType {
	auth: Twitch.ext.Authorized | undefined;
	ebsState: EbsState | null;
	isDebugMode: boolean;
	isLocalMode: boolean;

	setAuth: (value: Twitch.ext.Authorized) => void;
	setEbsState: (value: EbsState) => void;
	setDebugMode: (value: boolean) => void;
	setLocalMode: (value: boolean) => void;
}

// 2. Create the context with a default value of `undefined`
const StateContext = createContext<GlobalStateType | undefined>(undefined);

// 3. Define the provider component
interface StateProviderProps {
	children: ReactNode;
}

export const StateProvider: React.FC<StateProviderProps> = ({ children }) => {
	const [isDebugMode, setDebugMode] = useState<boolean>(false);
	const [isLocalMode, setLocalMode] = useState<boolean>(false);
	const [auth, setAuth] = useState<Twitch.ext.Authorized>();
	const [ebsState, setEbsState] = useState<EbsState | null>(null);

	return (
		<StateContext.Provider
			value={{
				auth,
				ebsState,
				isDebugMode,
				isLocalMode,
				setAuth,
				setEbsState,
				setDebugMode,
				setLocalMode,
			}}
		>
			{children}
		</StateContext.Provider>
	);
};

// 4. Custom hook for accessing the state context
export const useGlobalState = (): GlobalStateType => {
	const context = useContext(StateContext);
	if (!context) {
		throw new Error("useGlobalState must be used within a StateProvider");
	}
	return context;
};
