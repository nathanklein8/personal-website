import { getContent } from '@nk/shared/server/backend';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    const { landingCard, projects, featuredPhotos, apiURL } = await getContent();
    return {
        landingCard,
        projects,
        featuredPhotos,
        apiURL
    };
};
