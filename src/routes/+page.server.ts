import { getURL } from '$lib/server/backend';
import type { PageServerLoad } from './$types';

function handleError(message: string): any {
    return {
        landingCard: {
            error: message
        }
    };
}

export const load: PageServerLoad = async () => {

    const apiURL = getURL();
    const endpoint = apiURL + '/api/landingcard'

    let res: Response;

    try {
        res = await fetch(endpoint);
    } catch (err) {
        console.error(`Unable to reach backend`, err);
        return handleError("Unable to reach backend")
    }

    if (!res.ok) {
        let bodyText = "";
        try {
            bodyText = await res.text();
        } catch {
            bodyText = "<unreadable>";
        }
        console.error(`HTTP ${res.status}: ${res.statusText}, Body: ${bodyText}`);
        return handleError(`HTTP ${res.status}: ${res.statusText}, Body: ${bodyText}`)
    }

    let json: any;

    try {
        json = await res.json();
    } catch (err) {
        console.error("Backend returned non-JSON response", err);
        return handleError("Backend returned non-JSON response")
    }

    // no error happened, return a json struct to the ui
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