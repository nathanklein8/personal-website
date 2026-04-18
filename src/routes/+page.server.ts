import { getContent, getLandingCard, getProjects, getPhotos } from '@nk/shared/server/backend';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    const { apiURL, results } = await getContent(getLandingCard, getProjects, getPhotos);
    const [landingCard, projects, { featuredPhotos }] = results;
    return {
        landingCard,
        projects,
        featuredPhotos,
        apiURL
    };
};
