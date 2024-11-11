import _ from "lodash";
import { collectRequest, dispenseRequest } from "../ebs/api";
import { Coords } from "../types";

export default function ControlPanel({
	auth,
	coords,
}: { auth: Twitch.ext.Authorized | undefined; coords: Coords }) {
	const collectionHandler =
		(auth: Twitch.ext.Authorized | undefined, vialPos: number) => () => {
			collectRequest(auth, vialPos);
		};
	const dispenseHandler =
		(auth: Twitch.ext.Authorized | undefined, x: number, y: number) => () => {
			dispenseRequest(auth, x, y);
		};
	return (
		<div id="color-vote-area">
			{_.times(5, (i) => {
				const vialPos = i + 2;
				return (
					<div
						className="color-option"
						onClick={collectionHandler(auth, vialPos)}
						onKeyDown={collectionHandler(auth, vialPos)}
					>
						{vialPos}
					</div>
				);
			})}
			<div
				className="dispense-button"
				onClick={dispenseHandler(auth, coords.x, coords.y)}
				onKeyDown={dispenseHandler(auth, coords.x, coords.y)}
			>
				Dispense
			</div>
		</div>
	);
}
