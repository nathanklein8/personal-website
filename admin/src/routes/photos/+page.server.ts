import type { Actions, PageServerLoad } from './$types';
import type { Photo } from '@nk/shared/types/photo';
import {
	getAllPhotos,
	getAvailableYears,
	getAvailableEvents,
	getAvailablePhotos,
	addPhoto,
	updatePhoto,
	deletePhoto,
	regenerateThumbnails
} from '@nk/shared/server/backend';

export const load: PageServerLoad = async ({ url }) => {
	const year = url.searchParams.get('year');
	const event = url.searchParams.get('event');

	const [photos, years] = await Promise.all([getAllPhotos(), getAvailableYears()]);

	let events: string[] = [];
	let availablePhotos: string[] = [];

	if (year && event) {
		availablePhotos = await getAvailablePhotos(year, event);
	} else if (year) {
		events = await getAvailableEvents(year);
	}

	return {
		photos: photos as Photo[],
		years,
		events,
		availablePhotos
	};
};

export const actions: Actions = {
	addPhoto: async ({ request }) => {
		const data = await request.formData();

		const filename = data.get('filename')?.toString() || '';
		const title = data.get('title')?.toString() || '';
		const sortOrder = parseInt(data.get('sortOrder')?.toString() || '0');

		return await addPhoto(filename, title, sortOrder);
	},

	updatePhoto: async ({ request }) => {
		const data = await request.formData();
		const id = data.get('id')?.toString() || 'unknown';
		const title = data.get('title')?.toString() || '';
		const sortOrder = parseInt(data.get('sortOrder')?.toString() || '0');
		const visible = data.get('visible') === 'on';

		return await updatePhoto(id, title, sortOrder, visible);
	},

	deletePhoto: async ({ request }) => {
		const data = await request.formData();
		const id = data.get('id')?.toString() || 'unknown';

		return await deletePhoto(id);
	},

	regenerateThumbnails: async () => {
		return await regenerateThumbnails();
	}
};
