<script>
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import CardFooter from '$lib/components/ui/card/card-footer.svelte';
	import CardTitle from '$lib/components/ui/card/card-title.svelte';
	import {
		ArrowUpRight,
		BadgeQuestionMark,
		BookOpenCheck,
		CircleArrowDown,
		MessageCircle,
		Puzzle
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { expoIn } from 'svelte/easing';
	import { fade } from 'svelte/transition';

	const languages = ['TypeScript', 'Python', 'Java', 'C', 'SQL'];
	const fullstacktechs = [
		'React',
		'Svelte',
		'Prisma',
		'drizzle',
		'PostgreSQL',
		'Docker',
		'tailwindCSS',
		'shadcn/ui'
	];
	const mltechs = ['Pytorch', 'scikit-learn', 'mlflow'];

	let showScrollIcon = false;

	onMount(() => {
		if (window.scrollY > window.innerHeight * 0.05) {
			showScrollIcon = false;
			return; // No need to add scroll listener
		} else {
			showScrollIcon = true;
		}
		const handleScroll = () => {
			const scrollY = window.scrollY;
			const threshold = window.innerHeight * 0.05;
			// turn off scroll hint after user scrolls over 5% of screen height
			// and don't turn scroll hint back on, even if they scroll back up
			showScrollIcon = showScrollIcon && scrollY < threshold;
		};

		window.addEventListener('scroll', handleScroll);
		return () => window.removeEventListener('scroll', handleScroll);
	});
</script>

<section
	id="landing"
	class="flex min-h-[calc(100vh-64px)] flex-col items-center justify-center gap-y-8 bg-gradient-to-b from-green-700 to-background py-12"
>
	<div class="mx-auto max-w-4xl space-y-8">
		<div class="space-y-4 text-center">
			<div class="mb-6 flex justify-center">
				<div class="h-48 w-48 overflow-hidden rounded-full border-4 border-stone-100 shadow-xl">
					<img src="/me.jpg" alt="me" />
				</div>
			</div>

			<h2 class="font-pixel text-2xl text-5xl">Hi! I'm Nathan</h2>
			<p class="mx-auto max-w-2xl">New Grad | CS @ RIT</p>
		</div>
	</div>

	<Card
		class="max-w-[95vw] shadow-lg md:max-w-[80vw] lg:max-w-[65vw] xl:max-w-[50vw] dark:shadow-stone-800"
	>
		<CardContent class="grid gap-6 md:grid-cols-2">
			<div class="space-y-4">
				<h1 class="font-pixel flex items-center gap-2 text-2xl">
					<BadgeQuestionMark size={22} />
					About Me
				</h1>
				<p>
					I am a Senior at RIT studying Computer Science. I have internship experience in software
					engineering and machine learning. In my free time I enjoy hiking, photography, and playing
					volleyball!
				</p>
			</div>

			<div class="space-y-4">
				<h1 class="font-pixel flex items-center gap-2 text-2xl">
					<MessageCircle size={22} />
					Contact
				</h1>
				<div class="space-y-2">
					<p class="flex items-center gap-2">neklein3@gmail.com</p>
					<p class="flex max-w-fit items-center gap-2 border-foreground hover:border-b">
						<a
							class="flex items-center gap-0.5"
							href="https://www.linkedin.com/in/neklein"
							target="_blank"
							rel="noopener noreferrer"
						>
							Connect with me on LinkedIn
							<ArrowUpRight size={20} />
						</a>
					</p>
				</div>
			</div>
		</CardContent>
		<CardFooter class="border-t">
			<div class="flex flex-col gap-3">
				<h1 class="font-pixel flex items-center gap-2 text-2xl">
					<BookOpenCheck size={22} />
					I have experience with
				</h1>
				<div class="flex flex-wrap gap-2">
					{#each languages as lang}
						<Badge variant="outline">
							{lang}
						</Badge>
					{/each}
				</div>
				<div class="flex flex-wrap gap-2">
					{#each fullstacktechs as thing}
						<Badge variant="secondary">
							{thing}
						</Badge>
					{/each}
				</div>
				<div class="flex flex-wrap gap-2">
					{#each mltechs as thing}
						<Badge variant="outline">
							{thing}
						</Badge>
					{/each}
				</div>
			</div>
		</CardFooter>
	</Card>

	{#if showScrollIcon}
		<div
			in:fade={{ duration: 1200, easing: expoIn }}
			out:fade={{ duration: 200 }}
			class="fixed bottom-0 left-0 mb-4 flex w-full flex-col items-center gap-y-4"
		>
			<p>Check out what I've done!</p>
			<CircleArrowDown class="animate-bounce" size={32} />
		</div>
	{/if}
</section>

<section id="projects" class="flex flex-col items-center justify-center gap-y-8 py-16">
	<h1 class="font-pixel text-center text-4xl">Some projects I've worked on...</h1>

	<Card
		class="max-w-[95vw] shadow-lg md:max-w-[80vw] lg:max-w-[65vw] xl:max-w-[50vw] dark:shadow-stone-800"
	>
		<CardTitle class="font-pixel flex flex-row items-center gap-2 px-6 text-2xl">
			<Puzzle size={22} />
			<h1>Hangman Web App</h1>
			<span class="flex grow"></span>

			<p class="flex flex-row max-w-fit items-center gap-2 border-foreground hover:border-b">
				<a href="https://app.nklein.xyz" target="_blank" rel="noopener noreferrer">
					check it out
				</a>
        <ArrowUpRight size={32} />
			</p>
		</CardTitle>
		<CardContent class="grid gap-6 md:grid-cols-2">
			<img
				src="/hangman.jpg"
				alt="Screenshot of my Hangman Webapp"
				class="w-full max-w-80 rounded-lg object-cover md:max-w-full"
			/>
			<p class="text-md text-muted-foreground">
				Full-stack word-guessing game built with <span class="font-medium text-foreground"
					>React</span
				>,
				<span class="font-medium text-foreground">Prisma</span>, and
				<span class="font-medium text-foreground">PostgreSQL</span>. It features a daily global
				challenge word and player statistics, all managed server-side for consistent gameplay across
				users. The application is self-hosted on my home server for complete control over deployment
				and data.
			</p>
		</CardContent>
		<CardFooter class="flex flex-wrap gap-2 border-t">
			<Badge variant="secondary">React</Badge>
			<Badge variant="secondary">Next.js</Badge>
			<Badge variant="secondary">TypeScript</Badge>
			<Badge variant="secondary">Prisma</Badge>
			<Badge variant="secondary">PostgreSQL</Badge>
			<Badge variant="secondary">tailwindCSS</Badge>
			<Badge variant="secondary">Docker</Badge>
		</CardFooter>
	</Card>
</section>
