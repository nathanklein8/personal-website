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

export async function getLandingCard() {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/landingcard');

  if (res.ok) {
    return await res.json();
  }
  console.error(`Unable to fetch landing card content, HTTP ${res.status}: ${res.statusText}`);
  return { error: `Unable to fetch landing card content, HTTP ${res.status}: ${res.statusText}` };
}

export async function getProjects() {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/projects');

  if (res.ok) {
    return await res.json();
  }
  console.error(`Unable to fetch projects, HTTP ${res.status}: ${res.statusText}`);
  return [{ error: `Unable to fetch projects, HTTP ${res.status}: ${res.statusText}` }];
}

export async function getPhotos() {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/photos');

  let photos: any[] = [];
  if (res.ok) {
    const data = await res.json();
    photos = data ?? [];
  } else {
    console.error(`Unable to fetch photos, HTTP ${res.status}: ${res.statusText}`);
  }

  const visiblePhotos = photos.filter((p: any) => p.visible !== false);
  return { visiblePhotos, featuredPhotos: visiblePhotos.slice(0, 3) };
}

export async function getContent<T extends readonly (() => Promise<any>)[]>(...getters: T) {
  const apiURL = getURL();
  const results = await Promise.all(getters.map(fn => fn()));
  return { apiURL, results };
}
