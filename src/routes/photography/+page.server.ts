import type { PageServerLoad } from './$types';
import { getPhotos, getURL } from '@nk/shared/server/backend';

export const load: PageServerLoad = async () => {
    const { visiblePhotos: photos } = await getPhotos();
    const apiURL = getURL();
    return { 
        photos,
        apiURL
    };
};
