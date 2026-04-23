import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		alias: {
			'@nk/shared/*': './src/*',
		},
	},
	extensions: ['.svelte'],
};

export default config;
