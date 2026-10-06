const CLAIM_HEADER = "X-Claim-Token";

export interface ApiResult<T = unknown> {
	ok: boolean;
	status: number;
	message: string;
	data?: T;
}

async function request<T>(
	method: string,
	path: string,
	token: string | null,
	body?: unknown,
): Promise<ApiResult<T>> {
	const headers: Record<string, string> = {};
	if (token) headers[CLAIM_HEADER] = token;
	if (body !== undefined) headers["Content-Type"] = "application/json";

	let res: Response;
	try {
		res = await fetch(`/api${path}`, {
			method,
			headers,
			body: body === undefined ? undefined : JSON.stringify(body),
		});
	} catch {
		return { ok: false, status: 0, message: "Can't reach the robot" };
	}

	let data: any;
	try {
		data = await res.json();
	} catch {
		data = undefined;
	}
	return {
		ok: res.ok,
		status: res.status,
		message: typeof data?.message === "string" ? data.message : res.statusText,
		data,
	};
}

export const api = {
	claim: (token: string | null) =>
		request<{ token: string }>("PUT", "/claim", token),
	unclaim: (token: string) => request("PUT", "/unclaim", token),
	collect: (token: string, id: number) =>
		request("POST", "/collection", token, { id }),
	dispense: (token: string, x: number, y: number) =>
		request("POST", "/dispense", token, { x, y }),
	goTo: (token: string, x: number, y: number) =>
		request("PUT", "/goto", token, { x, y }),
};

// Release control when the page closes. fetch with keepalive survives unload.
export function unclaimOnExit(token: string) {
	fetch("/api/unclaim", {
		method: "PUT",
		headers: { [CLAIM_HEADER]: token },
		keepalive: true,
	}).catch(() => {});
}
