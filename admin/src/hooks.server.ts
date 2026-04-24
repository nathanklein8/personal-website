import { type Handle, type RequestEvent } from '@sveltejs/kit';
import { readFileSync, existsSync } from 'node:fs';

export const handle: Handle = async ({ event, resolve }) => {
    if (event.url.pathname.match(/\/api\/photos\/available\/.*\/preview$/)) {
        return servePreviewFromVolume(event);
    }
    return resolve(event);
};

async function servePreviewFromVolume(event: RequestEvent) {
    // Parse path: /api/photos/available/{year}/{event}/{filename}/preview
    const parts = event.url.pathname.split('/').filter(Boolean);
    // parts: ["api", "photos", "available", "year", "event", "filename", "preview"]
    if (parts.length < 7) {
        return new Response('Not found', { status: 404 });
    }

    const year = decodeURIComponent(parts[3]);
    const event_name = decodeURIComponent(parts[4]);
    const filename = decodeURIComponent(parts[5]);

    // Determine size from filename
    let imageFile = filename;
    if (filename.endsWith('_thumb.jpg')) {
        imageFile = filename;
    } else if (filename.endsWith('_med.jpg')) {
        imageFile = filename;
    } else {
        imageFile = `${filename}_med.jpg`;
    }

    const volumePath = `/thumbnails/${year}/${event_name}/${imageFile}`;

    if (!existsSync(volumePath)) {
        return new Response('Preview not found', { status: 404 });
    }

    const buffer = readFileSync(volumePath);

    return new Response(buffer, {
        headers: {
            'Content-Type': 'image/jpeg',
            'Cache-Control': 'public, max-age=31536000, immutable',
        },
    });
}
