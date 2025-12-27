import type { Actions, PageServerLoad } from './$types';
import { getContent, getURL } from '$lib/server/backend';

export const load: PageServerLoad = async () => {
    return getContent();
};

export const actions: Actions = {
    updateLanding: async ({ request }) => {
        const data = await request.formData();

        const rawSkills = (data.get("skills") as string) ?? "";

        // split on newline to get categories
        const lines = rawSkills
            .split("\n")
            .map(l => l.trim())
            .filter(l => l.length > 0);

        // split each line on comma to get skills per category
        const parsedSkills = lines.map(line =>
            line
                .split(",")
                .map(s => s.trim())
                .filter(Boolean)
        );

        const payload = {
            bio: data.get("bio")?.toString(),
            email: data.get("email")?.toString(),
            linkedin: data.get("linkedin")?.toString(),
            github: data.get("github")?.toString(),
            skills: parsedSkills,
        };

        const apiURL = getURL();

        const res = await fetch(`${apiURL}/api/landingcard`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            return {
                failure: true,
                message: `Backend error: ${res.status} ${res.statusText}`
            };
        }

        return { success: true };
    },
    updateProject: async ({ request }) => {

        const data = await request.formData();

        // pull project ID out of form data
        const proj_id = data.get('id');

        const payload = {
            title: data.get("title")?.toString(),
            icon: data.get("icon")?.toString(),
            description: data.get("description")?.toString(),
            technologies: data.get("technologies")?.toString().split(','),
            deploymentLink: data.get("deploymentLink")?.toString(),
            image: data.get("image")?.toString(),
            altText: data.get("altText")?.toString(),
        }

        const apiURL = getURL();

        const res = await fetch(`${apiURL}/api/projects/${proj_id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            return {
                failure: true,
                id: proj_id,
                message: `Backend error: ${res.status} ${res.statusText}`
            };
        }

        return { success: true, id: proj_id, };

    }
};
