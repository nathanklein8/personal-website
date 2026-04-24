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

export async function getAllPhotos() {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/photos?type=all');

  let photos: any[] = [];
  if (res.ok) {
    const data = await res.json();
    photos = data ?? [];
  } else {
    console.error(`Unable to fetch all photos, HTTP ${res.status}: ${res.statusText}`);
  }

  return photos;
}

export async function getVisiblePhotos() {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/photos?type=visible');

  let photos: any[] = [];
  if (res.ok) {
    const data = await res.json();
    photos = data ?? [];
  } else {
    console.error(`Unable to fetch visible photos, HTTP ${res.status}: ${res.statusText}`);
  }

  return photos;
}

export async function getFeaturedPhotos() {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/photos?type=featured');

  let photos: any[] = [];
  if (res.ok) {
    const data = await res.json();
    photos = data ?? [];
  } else {
    console.error(`Unable to fetch featured photos, HTTP ${res.status}: ${res.statusText}`);
  }

  return photos;
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

  return { photos, featuredPhotos: photos.filter((p: any) => p.visible !== false).slice(0, 3) };
}

export async function getAvailableYears() {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/photos/available');

  let years: string[] = [];
  if (res.ok) {
    years = await res.json();
  } else {
    console.error(`Unable to fetch available years: HTTP ${res.status}: ${res.statusText}`);
  }

  return years;
}

export async function getAvailableEvents(year: string) {
  const apiURL = getURL();
  const res = await fetch(apiURL + `/api/photos/available/${encodeURIComponent(year)}`);

  let events: string[] = [];
  if (res.ok) {
    events = await res.json();
  } else {
    console.error(`Unable to fetch available events for ${year}: HTTP ${res.status}: ${res.statusText}`);
  }

  return events;
}

export async function getAvailablePhotos(year: string, event: string) {
  const apiURL = getURL();
  const res = await fetch(apiURL + `/api/photos/available/${encodeURIComponent(year)}/${encodeURIComponent(event)}`);

  let photos: string[] = [];
  if (res.ok) {
    photos = await res.json();
  } else {
    console.error(`Unable to fetch available photos for ${year}/${event}: HTTP ${res.status}: ${res.statusText}`);
  }

  return photos;
}

export async function addPhoto(filename: string, title: string, sortOrder: number) {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/photos', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ filename, title, sortOrder })
  });

  if (!res.ok) {
    const body = await res.text();
    return { success: false, message: body || `Backend error: ${res.status}` };
  }

  const result = await res.json();
  return { success: true, id: result.id, photo: result };
}

export async function updatePhoto(id: string | number, title: string, sortOrder: number, visible: boolean) {
  const apiURL = getURL();
  const res = await fetch(apiURL + `/api/photos/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title, sortOrder, visible })
  });

  if (!res.ok) {
    const body = await res.text();
    return { success: false, id, message: body || `Backend error: ${res.status}` };
  }

  return { success: true, id };
}

export async function deletePhoto(id: string | number) {
  const apiURL = getURL();
  const res = await fetch(apiURL + `/api/photos/${id}`, {
    method: 'DELETE'
  });

  if (!res.ok) {
    const body = await res.text();
    return { success: false, id, message: body || `Backend error: ${res.status}` };
  }

  return { success: true, id };
}

export async function regenerateThumbnails() {
  const apiURL = getURL();
  const res = await fetch(apiURL + '/api/photos/regenerate-thumbnails', {
    method: 'POST'
  });

  if (!res.ok) {
    const body = await res.text();
    return { success: false, message: body || `Backend error: ${res.status}` };
  }

  return { success: true };
}

export async function getContent<T extends readonly (() => Promise<any>)[]>(...getters: T) {
  const apiURL = getURL();
  const results = await Promise.all(getters.map(fn => fn()));
  return { apiURL, results };
}
