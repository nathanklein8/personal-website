import { getContent } from '$lib/server/backend';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    return getContent();
};
