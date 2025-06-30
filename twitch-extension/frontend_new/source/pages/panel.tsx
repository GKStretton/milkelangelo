import render from "~/app/render";
import { useFollowingStatus } from "~/hooks/twitch";

const App = () => {
	const { followingStatus } = useFollowingStatus();

	return <>Panel {followingStatus}</>;
};

render(<App />);
