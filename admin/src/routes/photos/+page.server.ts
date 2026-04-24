import type { Actions, PageServerLoad } from './$types';
import type { Photo } from '@nk/shared/types/photo';
import { getAllPhotos, getAvailableYears, getAvailableEvents, getAvailablePhotos, addPhoto, updatePhoto, deletePhoto, regenerateThumbnails } from '@nk/shared/server/backend';

export const load: PageServerLoad = async () => {
    const [photos, years] = await Promise.all([
        getAllPhotos(),
        getAvailableYears()
    ]);

    return {
        photos: photos as Photo[],
        years
    };
};

export const actions: Actions = {
    addPhoto: async ({ request }) => {
        const data = await request.formData();

        const filename = data.get("filename")?.toString() || '';
        const title = data.get("title")?.toString() || '';
        const sortOrder = parseInt(data.get("sortOrder")?.toString() || '0');

        return await addPhoto(filename, title, sortOrder);
    },

    updatePhoto: async ({ request }) => {
        const data = await request.formData();
        const id = data.get('id')?.toString() || 'unknown';
        const title = data.get("title")?.toString() || '';
        const sortOrder = parseInt(data.get("sortOrder")?.toString() || '0');
        const visible = data.get("visible") === 'on';

        return await updatePhoto(id, title, sortOrder, visible);
    },

    deletePhoto: async ({ request }) => {
        const data = await request.formData();
        const id = data.get('id')?.toString() || 'unknown';

        return await deletePhoto(id);
    },

    regenerateThumbnails: async () => {
        return await regenerateThumbnails();
    },

    availableYears: async () => {
        const years = await getAvailableYears();
        return { success: true, years };
    },

    availableEvents: async ({ request }) => {
        const formData = await request.formData();
        const year = formData.get('year')?.toString();
        if (!year) {
            return { failure: true, message: 'No year specified' };
        }
        const events = await getAvailableEvents(year);
        return { success: true, events };
    },

    availablePhotos: async ({ request }) => {
        const formData = await request.formData();
        const year = formData.get('year')?.toString();
        const event = formData.get('event')?.toString();
        if (!year || !event) {
            return { failure: true, message: 'Year and event required' };
        }
        const photos = await getAvailablePhotos(year, event);
        return { success: true, photos };
    },

    selectPhoto: async ({ request }) => {
        const formData = await request.formData();
        const year = formData.get('year')?.toString();
        const event = formData.get('event')?.toString();
        const filename = formData.get('filename')?.toString();
        if (!year || !event || !filename) {
            return { failure: true, message: 'Year, event, and filename required' };
        }
        return { success: true, year, event, filename };
    },
};
