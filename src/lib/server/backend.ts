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