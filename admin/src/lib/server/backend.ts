let url: string | null = null;

export function getURL() {
  if (!url) {
    const env_var = process.env.API_URL;
    if (!env_var) {
      throw new Error('API_URL environment variable is not set!');
    }
    url = env_var;
  }
  return url;
}

export async function getContent() {
  const apiURL = getURL();

  // endpoints to fetch
  const landingCardEndpoint = apiURL + '/api/landingcard';
  const projectsEndpoint = apiURL + '/api/projects';

  // get promises for all endpoints
  const [landingRes, projectsRes] = await Promise.allSettled([
    fetch(landingCardEndpoint),
    fetch(projectsEndpoint)
  ]);

  // handle landing card promise
  let landingCard: any = {};
  if (landingRes.status === 'fulfilled') {
    const res = landingRes.value;
    if (res.ok) {
      landingCard = await res.json();
    } else {
      console.error(`Unable to fetch landing card content, HTTP ${res.status}: ${res.statusText}`);
      landingCard = { error: `Unable to fetch landing card content, HTTP ${res.status}: ${res.statusText}` }
    }
  } else {
    console.error("Unable to reach backend:", landingRes.reason);
    landingCard = { error: `Unable to reach backend: ${landingRes.reason}` };
  }

  // handle projects promise
  let projects: any[] = [];
  if (projectsRes.status === 'fulfilled') {
    const res = projectsRes.value;
    if (res.ok) {
      projects = await res.json();
    } else {
      console.error(`Unable to fetch projects, HTTP ${res.status}: ${res.statusText}`);
      projects = [{ error: `Unable to fetch projects, HTTP ${res.status}: ${res.statusText}` }]
    }
  } else {
    console.error(`Unable to reach backend: ${projectsRes.reason}`);
    projects = [{ error: `Unable to reach backend: ${projectsRes.reason}` }]
  }

  return {
    landingCard,
    projects
  };
}