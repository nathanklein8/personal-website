import { type Handle, type RequestEvent } from '@sveltejs/kit';
import { readFileSync, existsSync } from 'node:fs';

export const handle: Handle = async ({ event, resolve }) => {
    if (event.url.pathname.match(/\/api\/photos\/.*\/image$/)) {
        return serveImageFromVolume(event);
    }
    return resolve(event);
};

async function serveImageFromVolume(event: RequestEvent) {
    // Parse path: /api/photos/{year}/{event}/{filename}/image
    const parts = event.url.pathname.split('/').filter(Boolean);
    // parts: ["api", "photos", "year", "event", "filename", "image"]
    if (parts.length < 6) {
        return new Response('Not found', { status: 404 });
    }

    const year = decodeURIComponent(parts[2]);
    const event_name = decodeURIComponent(parts[3]);
    const filename = decodeURIComponent(parts[4]);

    const size = event.url.searchParams.get('size') || 'med';
    if (!['thumb', 'med'].includes(size)) {
        return new Response('Invalid size parameter', { status: 400 });
    }

    const imageFile = size === 'thumb' ? `${filename}_thumb.jpg` : `${filename}_med.jpg`;
    const volumePath = `/thumbnails/${year}/${event_name}/${imageFile}`; // In prod: /thumbnails (mounted volume). In dev: symlinked or local path.

    if (!existsSync(volumePath)) {
        return new Response('Image not found', { status: 404 });
    }

    const buffer = readFileSync(volumePath);

    return new Response(buffer, {
        headers: {
            'Content-Type': 'image/jpeg',
            'Cache-Control': 'public, max-age=31536000, immutable',
        },
    });
}
