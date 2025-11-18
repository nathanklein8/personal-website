import { getURL } from '$lib/server/backend';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ params }) => {

	const apiURL = getURL();

	console.log("API_URL: " + apiURL);
    
	const { id } = params;
	const res = await fetch(apiURL + `/api/test/${id}/increment`, {
		method: 'POST'
	});

	const text = await res.text();
	return new Response(text, {
		status: res.status,
		headers: { 'Content-Type': 'application/json' }
	});
};
