import type { PageServerLoad } from './$types';
import { getContent } from '$lib/server/backend';

export const load: PageServerLoad = async () => {
    const { photos, apiURL } = await getContent();
    return { 
        photos,
        apiURL
    };
};
