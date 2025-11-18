import { getURL } from '$lib/server/backend';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request }) => {

	const apiURL = getURL();

	const body = await request.json();
	const res = await fetch(apiURL + '/api/test', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});

	const text = await res.text();
	return new Response(text, {
		status: res.status,
		headers: { 'Content-Type': 'application/json' }
	});
};
