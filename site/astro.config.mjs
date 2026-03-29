// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://thinkupfront.dev',
	integrations: [
		starlight({
			title: 'Upfront',
			tagline: 'Force thinking before code.',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/brennhill/upfront' }],
			customCss: ['./src/styles/custom.css'],
			sidebar: [
				{ label: 'Why Upfront', slug: 'why' },
				{ label: 'Human-First Development', slug: 'human-first' },
				{
					label: 'Commands',
					items: [
						{ label: 'Overview', slug: 'commands/overview' },
						{ label: '/ideate', slug: 'commands/ideate' },
						{ label: '/explore', slug: 'commands/explore' },
						{ label: '/feature', slug: 'commands/feature' },
						{ label: '/refine', slug: 'commands/refine' },
						{ label: '/plan', slug: 'commands/plan' },
						{ label: '/build', slug: 'commands/build' },
						{ label: '/patch', slug: 'commands/patch' },
						{ label: '/quick', slug: 'commands/quick' },
						{ label: '/debug', slug: 'commands/debug' },
						{ label: '/ship', slug: 'commands/ship' },
						{ label: '/retro', slug: 'commands/retro' },
						{ label: '/teach', slug: 'commands/teach' },
					],
				},
				{
					label: 'Support Commands',
					items: [
						{ label: '/note', slug: 'commands/note' },
						{ label: '/pause & /resume', slug: 'commands/pause-resume' },
					],
				},
				{ label: 'Audit Trail', slug: 'audit-trail' },
				{ label: 'Research', slug: 'research' },
				{ label: 'Install', slug: 'install' },
			],
		}),
	],
});
