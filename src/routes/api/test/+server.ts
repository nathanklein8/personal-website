import type { RequestHandler } from './$types';
import { API_URL } from '$env/static/private';

export const POST: RequestHandler = async ({ request }) => {
    console.log('API_URL ' + API_URL)
	const body = await request.json();
	const res = await fetch(`${API_URL}/api/test`, {
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
