import { getURL } from '$lib/server/backend';
import type { PageServerLoad } from './$types';


export const load: PageServerLoad = async () => {
    const apiURL = getURL();

    // Define endpoints
    const landingCardEndpoint = apiURL + '/api/landingcard';
    const projectsEndpoint = apiURL + '/api/projects';

    // Fetch both concurrently
    const [landingRes, projectsRes] = await Promise.allSettled([
        fetch(landingCardEndpoint),
        fetch(projectsEndpoint)
    ]);

    // Landing card
    let landingCard: any = {};
    if (landingRes.status === 'fulfilled') {
        const res = landingRes.value;
        if (res.ok) {
            try {
                const json = await res.json();
                landingCard = {
                    bio: json.bio,
                    email: json.email,
                    linkedin: json.linkedin,
                    github: json.github,
                    skills: json.skills
                };
            } catch (err) {
                console.error("Landing card JSON parse error:", err);
                landingCard = { error: "Invalid landing card JSON" };
            }
        } else {
            console.error(`Landing card HTTP ${res.status}: ${res.statusText}`);
            landingCard = { error: `HTTP ${res.status}: ${res.statusText}`};
        }
    } else {
        console.error("Unable to fetch landing card:", landingRes.reason);
        landingCard = { error: `Unable to fetch landing card, ${landingRes.reason}` };
    }

    // Projects
    let projects: any[] = [];
    if (projectsRes.status === 'fulfilled') {
        const res = projectsRes.value;
        if (res.ok) {
            try {
                projects = await res.json();
            } catch (err) {
                console.error("Projects JSON parse error:", err);
            }
        } else {
            console.error(`Projects HTTP ${res.status}: ${res.statusText}`);
        }
    } else {
        console.error("Unable to fetch projects:", projectsRes.reason);
    }

    return {
        landingCard,
        projects
    };
};
