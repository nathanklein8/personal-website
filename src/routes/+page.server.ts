import { getURL } from '$lib/server/backend';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {

    const apiURL = getURL();

    const res = await fetch(apiURL + '/api/landingcard', {
        method: 'GET',
    });

    if (res.ok) {
        const json = await res.json();
        return {
            landingCard: {
                bio: json.bio,
                email: json.email,
                linkedin: json.linkedin,
                skills: json.skills,
            }
        };
    } else {
        console.log('uh oh');
    }

};