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
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/ThinkUpfront/Upfront' },
				{ icon: 'linkedin', label: 'LinkedIn', href: 'https://www.linkedin.com/in/brennhill/' },
			],
			customCss: ['./src/styles/custom.css'],
			head: [
				{
					tag: 'script',
					attrs: {
						async: true,
						src: 'https://www.googletagmanager.com/gtag/js?id=G-EBNJWXFM9Z',
					},
				},
				{
					tag: 'script',
					content: "window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments)}gtag('js',new Date());gtag('config','G-EBNJWXFM9Z',{storage:'none',anonymize_ip:true});",
				},
				{
					tag: 'link',
					attrs: {
						rel: 'preconnect',
						href: 'https://fonts.googleapis.com',
					},
				},
				{
					tag: 'link',
					attrs: {
						rel: 'preconnect',
						href: 'https://fonts.gstatic.com',
						crossorigin: '',
					},
				},
			],
			sidebar: [
				{ label: 'Why Upfront', slug: 'why' },
				{ label: 'Get Started', slug: 'get-started' },
				{ label: 'Human-First Development', slug: 'human-first' },
				{
					label: 'Skills',
					items: [
						{ label: 'Overview', slug: 'commands/overview' },
						{ label: '/upfront:vision', slug: 'commands/vision' },
						{ label: '/upfront:increment', slug: 'commands/increment' },
						{ label: '/upfront:assess', slug: 'commands/assess' },
						{ label: '/upfront:ideate', slug: 'commands/ideate' },
						{ label: '/upfront:explore', slug: 'commands/explore' },
						{ label: '/upfront:enlighten', slug: 'commands/enlighten' },
						{ label: '/upfront:spike', slug: 'commands/spike' },
						{ label: '/upfront:feature', slug: 'commands/feature' },
						{ label: '/upfront:refine', slug: 'commands/refine' },
						{ label: '/upfront:plan', slug: 'commands/plan' },
						{ label: '/upfront:build', slug: 'commands/build' },
						{ label: '/upfront:patch', slug: 'commands/patch' },
						{ label: '/upfront:quick', slug: 'commands/quick' },
						{ label: '/upfront:debug', slug: 'commands/debug' },
						{ label: '/upfront:ship', slug: 'commands/ship' },
						{ label: '/upfront:retro', slug: 'commands/retro' },
						{ label: '/upfront:teach', slug: 'commands/teach' },
						{ label: '/upfront:architect', slug: 'commands/architect' },
						{ label: '/upfront:re-architect', slug: 'commands/re-architect' },
					],
				},
				{
					label: 'Support',
					items: [
						{ label: '/upfront:up', slug: 'commands/up' },
						{ label: '/upfront:upgrade', slug: 'commands/upgrade' },
						{ label: '/upfront:note', slug: 'commands/note' },
						{ label: '/upfront:pause & /upfront:resume', slug: 'commands/pause-resume' },
					],
				},
				{ label: 'Audit Trail', slug: 'audit-trail' },
				{ label: 'Research', slug: 'research' },
				{
					label: 'Compare',
					items: [
						{ label: 'vs Superpowers', slug: 'compare/vs-superpowers' },
						{ label: 'vs GSD', slug: 'compare/vs-gsd' },
					],
				},
				{ label: 'Install', slug: 'install' },
			],
		}),
	],
});
