import _ from "lodash";
import { toast } from "sonner";
import { useCollect, useDispense } from "../ebs/api";
import { useGlobalState } from "../helpers/State";
import "./ControlPanel.css";

export default function ControlPanel() {
	const gs = useGlobalState();

	const gooState = gs.ebsState?.GooState;

	const { mutate: collect, isPending: collectPending } = useCollect();
	const { mutate: dispense, isPending: dispensePending } = useDispense();

	const collectionHandler =
		(auth: Twitch.ext.Authorized | undefined, vialPos: number) => () => {
			collect(vialPos);
		};
	const dispenseHandler = () => {
		if (gooState) {
			dispense({ x: gooState.X, y: gooState.Y });
		} else {
			toast.error("No goo state available");
		}
	};
	return (
		<div id="color-vote-area">
			{_.times(5, (i) => {
				const vialPos = i + 2;
				return (
					<div
						className="color-option"
						onClick={collectionHandler(gs.auth, vialPos)}
						onKeyDown={collectionHandler(gs.auth, vialPos)}
						style={
							gs.ebsState?.GooState?.VialProfiles[vialPos]?.Colour
								? {
										backgroundColor:
											gs.ebsState.GooState.VialProfiles[vialPos].Colour,
										color:
											gs.ebsState.GooState.VialProfiles[vialPos].Colour ===
											"#0000ff"
												? "#ffffff"
												: "#000000",
								  }
								: {}
						}
					>
						{gs.ebsState?.GooState?.VialProfiles[vialPos]?.Name ?? vialPos}
					</div>
				);
			})}
			<div
				className="dispense-button"
				onClick={dispenseHandler}
				onKeyDown={dispenseHandler}
			>
				Dispense
			</div>
		</div>
	);
}
