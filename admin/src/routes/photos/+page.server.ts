import type { Actions, PageServerLoad } from './$types';
import { getPhotos, getURL } from '@nk/shared/server/backend';

export const load: PageServerLoad = async () => {
    const { visiblePhotos: photos } = await getPhotos();
    return {
        photos,
        apiURL
    }
};

export const actions: Actions = {
    addPhoto: async ({ request }) => {
        const data = await request.formData();
        
        const payload = {
            title: data.get("title")?.toString(),
            filePath: data.get("filePath")?.toString(),
            altText: data.get("altText")?.toString() || null,
            dateTaken: data.get("dateTaken")?.toString() || null,
            location: data.get("location")?.toString() || null,
            camera: data.get("camera")?.toString() || null,
            lens: data.get("lens")?.toString() || null,
            aperture: data.get("aperture")?.toString() || null,
            shutterSpeed: data.get("shutterSpeed")?.toString() || null,
            iso: data.get("iso")?.toString() || null,
            visible: data.get("visible") === 'on',
            sortOrder: parseInt(data.get("sortOrder")?.toString() || '0'),
        };

        const apiURL = getURL();
        const res = await fetch(`${apiURL}/api/photos`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            return { failure: true, message: `Backend error: ${res.status} ${res.statusText}` };
        }
        return { success: true };
    },

    updatePhoto: async ({ request }) => {
        const data = await request.formData();
        const id = data.get('id');

        const payload = {
            title: data.get("title")?.toString(),
            filePath: data.get("filePath")?.toString(),
            altText: data.get("altText")?.toString() || null,
            dateTaken: data.get("dateTaken")?.toString() || null,
            location: data.get("location")?.toString() || null,
            camera: data.get("camera")?.toString() || null,
            lens: data.get("lens")?.toString() || null,
            aperture: data.get("aperture")?.toString() || null,
            shutterSpeed: data.get("shutterSpeed")?.toString() || null,
            iso: data.get("iso")?.toString() || null,
            visible: data.get("visible") === 'on',
            sortOrder: parseInt(data.get("sortOrder")?.toString() || '0'),
        };

        const apiURL = getURL();
        const res = await fetch(`${apiURL}/api/photos/${id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            return { failure: true, id, message: `Backend error: ${res.status} ${res.statusText}` };
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
            return { failure: true, id, message: `Backend error: ${res.status} ${res.statusText}` };
        }
        return { success: true, id };
    }
};
