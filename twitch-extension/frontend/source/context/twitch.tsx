import React from "react";

import { setAuthorization } from "~api/twitch";

interface TwitchContext {
	authorized: Twitch.ext.Authorized | null;
}

const twitchContext = React.createContext<TwitchContext>(null as any);

const TwitchProvider = (props: React.PropsWithChildren) => {
	const [authorized, setAuthorized] =
		React.useState<Twitch.ext.Authorized | null>(null);

	React.useEffect(() => {
		Twitch.ext.viewer.id;
		Twitch.ext.onAuthorized((authorized) => {
			setAuthorization(authorized.helixToken, authorized.clientId);
			setAuthorized(authorized);
		});
	}, []);

	const context: TwitchContext = {
		authorized,
	};

	return (
		<twitchContext.Provider value={context}>
			{props.children}
		</twitchContext.Provider>
	);
};

const useTwitch = () => {
	const context = React.useContext(twitchContext);

	if (!context) {
		throw new Error("useTwitch must be used within a TwitchProvider");
	}

	return context;
};

export type { TwitchContext };

export { useTwitch };

export default TwitchProvider;
