import adapter from '@sveltejs/adapter-node';
import { optimizeImports } from "carbon-preprocess-svelte";

const production = process.env.NODE_ENV === 'production';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    preprocess: [optimizeImports()],
	kit: {
		adapter: adapter(),
		vite: {
			optimizeDeps: {
				include: ['@carbon/charts'],
			},
			ssr: {
				noExternal: [production && '@carbon/charts'].filter(Boolean),
			},
		},
	},
};

export default config;