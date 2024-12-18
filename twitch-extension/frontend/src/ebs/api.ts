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

export function collectRequest(
	auth: Twitch.ext.Authorized | undefined,
	vialPos: number,
) {
	console.log(`collect (${vialPos})`);

	const body: CollectionBody = { id: vialPos };

	fetch("http://localhost:8789/collect", {
		method: "POST",
		body: JSON.stringify(body),
		headers: getHeaders(auth),
	}).catch((e) => console.error(e));
}

export function goToRequest(
	auth: Twitch.ext.Authorized | undefined,
	x: number,
	y: number,
) {
	console.log(`goTo (${x}, ${y})`);

	const body: GoToBody = { x: x, y: y };

	fetch("http://localhost:8789/goto", {
		method: "PUT",
		body: JSON.stringify(body),
		headers: getHeaders(auth),
	}).catch((e) => console.error(e));
}

export function dispenseRequest(
	auth: Twitch.ext.Authorized | undefined,
	x: number,
	y: number,
) {
	console.log(`dispense (${x}, ${y})`);

	const body: DispenseBody = { x: x, y: y };

	fetch("http://localhost:8789/dispense", {
		method: "POST",
		body: JSON.stringify(body),
		headers: getHeaders(auth),
	}).catch((e) => console.error(e));
}

export function claimRequest(auth: Twitch.ext.Authorized | undefined) {
	console.log("claim");

	fetch("http://localhost:8789/claim-control", {
		method: "PUT",
		headers: getHeaders(auth),
	}).catch((e) => console.error(e));
}

const getHeaders = (auth: Twitch.ext.Authorized | undefined) => ({
	"Content-Type": "application/json",
	Authorization: `Bearer ${auth?.token}`,
	"X-Twitch-Extension-Client-Id": auth?.clientId ?? "",
});
