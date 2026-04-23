import path from 'path';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	resolve: {
		extensions: ['.svelte', '.ts', '.js', '.mjs'],
		alias: {
			'@nk/shared': path.resolve(__dirname, '../packages/shared/src'),
		},
	},
});
