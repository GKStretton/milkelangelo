import _ from "lodash";
import { toast } from "sonner";
import { useCollect, useDispense } from "../ebs/api";
import { useGlobalState } from "../helpers/State";
import { Status } from "../types";
import "./ControlPanel.css";

export default function ControlPanel() {
	const gs = useGlobalState();

	const gooState = gs.ebsState?.GooState;

	const enableCollectionButtons =
		gooState?.Status === Status.StatusDecidingCollection;
	const enableDispenseButton =
		gooState?.Status === Status.StatusDecidingDispense;
	// const enableCollectionButtons = gooState?.DispenseState?.Completed === true; // finished placing
	// const enableDispenseButton = gooState?.DispenseState?.Completed === false;

	const { mutate: collect, isPending: collectPending } = useCollect();
	const { mutate: dispense, isPending: dispensePending } = useDispense();

	const collectionHandler =
		(auth: Twitch.ext.Authorized | undefined, vialPos: number) => () => {
			if (!enableCollectionButtons) {
				return;
			}
			collect(vialPos);
		};
	const dispenseHandler = () => {
		if (!enableDispenseButton) {
			return;
		}

		if (gooState) {
			dispense({ x: gooState.X, y: gooState.Y });
		} else {
			toast.error("No goo state available");
		}
	};
	return (
		<div className="color-vote-area">
			{_.times(5, (i) => {
				const vialPos = i + 2;
				return (
					<div
						className={`color-option ${
							!enableCollectionButtons ? "disabled" : ""
						}`}
						onClick={collectionHandler(gs.auth, vialPos)}
						onKeyDown={collectionHandler(gs.auth, vialPos)}
						style={
							enableCollectionButtons &&
							gs.ebsState?.GooState?.VialProfiles[vialPos]?.Colour
								? {
										color: gs.ebsState.GooState.VialProfiles[vialPos].Colour,
								  }
								: {}
						}
					>
						{gs.ebsState?.GooState?.VialProfiles[vialPos]?.Name ?? vialPos}
					</div>
				);
			})}
			<div
				className={`dispense-button ${!enableDispenseButton ? "disabled" : ""}`}
				onClick={dispenseHandler}
				onKeyDown={dispenseHandler}
			>
				Dispense
			</div>
		</div>
	);
}
