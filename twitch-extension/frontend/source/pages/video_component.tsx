import render from "~/app/render";
import { useFollowingStatus } from "~/hooks/twitch";

const App = () => {
	const { followingStatus } = useFollowingStatus();

	return <>Video Component {followingStatus}</>;
};

render(<App />);
