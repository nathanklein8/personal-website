import type { RequestHandler } from './$types';
import { API_URL } from '$env/static/private';

export const POST: RequestHandler = async ({ params }) => {
    console.log('API_URL ' + API_URL)
	const { id } = params;

	const res = await fetch(`${API_URL}/api/test/${id}/increment`, {
		method: 'POST'
	});

	const text = await res.text();
	return new Response(text, {
		status: res.status,
		headers: { 'Content-Type': 'application/json' }
	});
};
