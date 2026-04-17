import type { PageServerLoad } from './$types';
import { getContent } from '@nk/shared/server/backend';

export const load: PageServerLoad = async () => {
    const { photos, apiURL } = await getContent();
    return { 
        photos,
        apiURL
    };
};
