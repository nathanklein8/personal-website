import { getLandingCard, getProjects, getPhotos, getURL } from '@nk/shared/server/backend';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    const landingCard = await getLandingCard();
    const projects = await getProjects();
    const { featuredPhotos } = await getPhotos();
    const apiURL = getURL();
    return {
        landingCard,
        projects,
        featuredPhotos,
        apiURL
    };
};
