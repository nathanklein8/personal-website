import { getContent, getLandingCard, getProjects, getFeaturedPhotos } from '@nk/shared/server/backend';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    const { apiURL, results } = await getContent(getLandingCard, getProjects, getFeaturedPhotos);
    const [landingCard, projects, featuredPhotos] = results;
    return {
        landingCard,
        projects,
        featuredPhotos,
        apiURL
    };
};
