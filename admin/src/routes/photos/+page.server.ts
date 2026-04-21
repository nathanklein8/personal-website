import type { Actions, PageServerLoad } from './$types';
import { getAllPhotos, getURL } from '@nk/shared/server/backend';

export const load: PageServerLoad = async () => {
    const photos = await getAllPhotos();
    return {
        photos
    };
};

export const actions: Actions = {
    addPhoto: async ({ request }) => {
        const data = await request.formData();
        
        const payload = {
            filename: data.get("filename")?.toString() || '',
            title: data.get("title")?.toString() || '',
            sortOrder: parseInt(data.get("sortOrder")?.toString() || '0'),
        };

        const apiURL = getURL();
        const res = await fetch(`${apiURL}/api/photos`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            const body = await res.text();
            return { failure: true, message: body || `Backend error: ${res.status}` };
        }
        const result = await res.json();
        return { success: true, id: result.id };
    },

    updatePhoto: async ({ request }) => {
        const data = await request.formData();
        const id = data.get('id');

        const payload: any = {};
        const title = data.get("title")?.toString();
        const sortOrder = data.get("sortOrder")?.toString();
        const visible = data.get("visible") === 'on';

        if (title !== undefined) payload.title = title;
        if (sortOrder !== undefined) payload.sortOrder = parseInt(sortOrder);
        payload.visible = visible;

        const apiURL = getURL();
        const res = await fetch(`${apiURL}/api/photos/${id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            const body = await res.text();
            return { failure: true, id, message: body || `Backend error: ${res.status}` };
        }
        return { success: true, id };
    },

    deletePhoto: async ({ request }) => {
        const data = await request.formData();
        const id = data.get('id');

        const apiURL = getURL();
        const res = await fetch(`${apiURL}/api/photos/${id}`, {
            method: "DELETE"
        });

        if (!res.ok) {
            const body = await res.text();
            return { failure: true, id, message: body || `Backend error: ${res.status}` };
        }
        return { success: true, id };
    },

    regenerateThumbnails: async () => {
        const apiURL = getURL();
        const res = await fetch(`${apiURL}/api/photos/regenerate-thumbnails`, {
            method: "POST"
        });

        if (!res.ok) {
            const body = await res.text();
            return { failure: true, message: body || `Backend error: ${res.status}` };
        }
        return { success: true };
    },

    availableYears: async () => {
        const apiURL = getURL();
        const res = await fetch(`${apiURL}/api/photos/available`);
        if (!res.ok) {
            return { failure: true, message: `Failed to load years: ${res.status}` };
        }
        const years = await res.json();
        return { success: true, years };
    },

    availableEvents: async ({ request }) => {
        const apiURL = getURL();
        const formData = await request.formData();
        const year = formData.get('year')?.toString();
        if (!year) {
            return { failure: true, message: 'No year specified' };
        }
        const res = await fetch(`${apiURL}/api/photos/available/${encodeURIComponent(year)}`);
        if (!res.ok) {
            return { failure: true, message: `Failed to load events: ${res.status}` };
        }
        const events = await res.json();
        return { success: true, events };
    },

    availablePhotos: async ({ request }) => {
        const apiURL = getURL();
        const formData = await request.formData();
        const year = formData.get('year')?.toString();
        const event = formData.get('event')?.toString();
        if (!year || !event) {
            return { failure: true, message: 'Year and event required' };
        }
        const res = await fetch(`${apiURL}/api/photos/available/${encodeURIComponent(year)}/${encodeURIComponent(event)}`);
        if (!res.ok) {
            return { failure: true, message: `Failed to load photos: ${res.status}` };
        }
        const photos = await res.json();
        return { success: true, photos };
    },

    selectPhoto: async ({ request }) => {
        const apiURL = getURL();
        const formData = await request.formData();
        const year = formData.get('year')?.toString();
        const event = formData.get('event')?.toString();
        const filename = formData.get('filename')?.toString();
        if (!year || !event || !filename) {
            return { failure: true, message: 'Year, event, and filename required' };
        }
        return { success: true, year, event, filename };
    }
};
