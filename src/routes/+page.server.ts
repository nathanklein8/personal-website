import { getURL } from '$lib/server/backend';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {

    const apiURL = getURL();
    const endpoint = apiURL + '/api/landinggard'

    let res: Response;

    try {
        res = await fetch(endpoint);
    } catch (err) {
        console.error(`Failed to reach backend at ${endpoint}`, err);
        return;
    }

    if (!res.ok) {
        let bodyText = "";
        try {
            bodyText = await res.text();
        } catch {
            bodyText = "<unreadable>";
        }
        console.error(`Backend returned HTTP ${res.status}: ${res.statusText}\nBody: ${bodyText}`);
    }

    let json: any;

    try {
        json = await res.json();
    } catch (err) {
        console.error("Backend returned non-JSON response", err);
        return;
    }

    return {
        landingCard: {
            bio: json.bio,
            email: json.email,
            linkedin: json.linkedin,
            github: json.github,
            skills: json.skills,
        }
    };
    
};