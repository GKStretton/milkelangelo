interface GoToBody {
	x: number;
	y: number;
}

interface CollectionBody {
	id: number;
}

export function collectRequest(auth: Twitch.ext.Authorized, vialPos: number) {
	console.log(`collect (${vialPos})`);

	const body: CollectionBody = { id: vialPos };

	fetch("http://localhost:8789/collect", {
		method: "POST",
		body: JSON.stringify(body),
		headers: {
			"Content-Type": "application/json",
			Authorization: `Bearer ${auth.token}`,
			"X-Twitch-Extension-Client-Id": auth.clientId,
		},
	}).catch((e) => console.error(e));
}

export function goToRequest(auth: Twitch.ext.Authorized, x: number, y: number) {
	console.log(`goTo (${x}, ${y})`);

	const body: GoToBody = { x: x, y: y };

	fetch("http://localhost:8789/goto", {
		method: "PUT",
		body: JSON.stringify(body),
		headers: {
			"Content-Type": "application/json",
			Authorization: `Bearer ${auth.token}`,
			"X-Twitch-Extension-Client-Id": auth.clientId,
		},
	}).catch((e) => console.error(e));
}

export function dispenseRequest(auth: Twitch.ext.Authorized) {
	console.log("dispense");

	fetch("http://localhost:8789/dispense", {
		method: "POST",
		body: JSON.stringify({}),
		headers: {
			"Content-Type": "application/json",
			Authorization: `Bearer ${auth.token}`,
			"X-Twitch-Extension-Client-Id": auth.clientId,
		},
	}).catch((e) => console.error(e));
}
