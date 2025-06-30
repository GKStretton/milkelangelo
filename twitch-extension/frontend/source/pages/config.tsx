import render from "~/app/render";
import { useFollowingStatus } from "~/hooks/twitch";

const App = () => {
	const { followingStatus } = useFollowingStatus();

	return <>Config {followingStatus}</>;
};

render(<App />);
