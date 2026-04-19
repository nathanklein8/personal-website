import type { PageServerLoad } from './$types';
import { getVisiblePhotos, getURL } from '@nk/shared/server/backend';

export const load: PageServerLoad = async () => {
    const photos = await getVisiblePhotos();
    const apiURL = getURL();
    return { 
        photos,
        apiURL
    };
};
