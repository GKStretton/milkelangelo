import React, { createContext, useContext, useState, ReactNode } from "react";

// 1. Define the shape of the global state
interface GlobalStateType {
	isDebugMode: boolean;
	isLocalMode: boolean;
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
	const [isDebugMode, setDebugMode] = useState<boolean>(
		process.env.NODE_ENV === "development",
	);
	const [isLocalMode, setLocalMode] = useState<boolean>(false);

	return (
		<StateContext.Provider
			value={{ isDebugMode, isLocalMode, setDebugMode, setLocalMode }}
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
