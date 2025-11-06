<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import CardFooter from '$lib/components/ui/card/card-footer.svelte';
	import {
		ArrowUpRight,
		BadgeQuestionMark,
		BookOpenCheck,
		CircleArrowDown,
		Clipboard,
		ClipboardCheck,
		MessageCircle
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { expoIn } from 'svelte/easing';
	import { fade } from 'svelte/transition';
	import Device from 'svelte-device-info';
	import ProjectCard from '$lib/components/project-card.svelte';
	import { cn } from '$lib/utils';
	import { TypeWriter } from 'svelte-typewrite';
	import { Button } from '$lib/components/ui/button';

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
	let copied = false;

	onMount(() => {
		const isMobile: boolean = Device.isPhone || Device.isTablet;

		if (window.scrollY > window.innerHeight * 0.05) {
			showScrollIcon = false;
			return; // No need to add scroll listener
		} else {
			showScrollIcon = !isMobile;
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

	async function copyEmail() {
		try {
			await navigator.clipboard.writeText('neklein3@gmail.com');
			copied = true;
			setTimeout(() => (copied = false), 3000);
		} catch (e) {
			console.error('Clipboard write failed', e);
		}
	}
</script>

<section
	id="landing"
	class="flex min-h-screen flex-col items-center gap-y-8 bg-gradient-to-b from-green-700 to-background to-55% md:to-70% pt-24 sm:pt-32"
>
	<div class="mx-auto max-w-4xl space-y-8">
		<div class="space-y-4 text-center">
			<div class="mb-6 flex justify-center">
				<div class="h-48 w-48 overflow-hidden rounded-full border-4 border-stone-100 shadow-xl">
					<img src="/me.jpg" alt="me" />
				</div>
			</div>
			<h1 class="font-code text-3xl md:text-4xl">Hi! I'm Nathan</h1>
			<h2 class="text-xl md:text-2xl">
				<TypeWriter typeSpeed={150} deleteSpeed={200} texts={['CS @ RIT', 'New Grad']} />
			</h2>
		</div>
	</div>

	<Card class="max-w-[95vw] shadow-lg md:max-w-[80vw] lg:max-w-[65vw] xl:max-w-[50vw]">
		<CardContent class="grid gap-6 sm:grid-cols-2">
			<div class="space-y-4">
				<h1 class="font-code flex items-center gap-2 text-xl">
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
				<h1 class="font-code flex items-center gap-2 text-xl">
					<MessageCircle size={22} />
					Contact
				</h1>
				<div class="space-y-3">
					<button
						class={cn(
							copied ? 'text-muted-foreground' : '',
							'font-code-wide flex w-full flex-row items-center border-foreground/25 rounded-lg border-1 bg-background/80 p-3 text-sm'
						)}
						onclick={copyEmail}
					>
						$ neklein3@gmail.com
						<span class="flex grow"></span>
						{#if copied}
							<ClipboardCheck size={20} />
						{:else}
							<Clipboard size={20} />
						{/if}
					</button>
					<p class="flex max-w-fit items-center gap-1 underline-offset-4 hover:underline">
						<a
							class="flex items-center gap-0.5"
							href="https://www.linkedin.com/in/neklein"
							target="_blank"
							rel="noopener noreferrer"
						>
							Connect with me on LinkedIn
						</a>
						<ArrowUpRight size={20} />
					</p>
				</div>
			</div>
		</CardContent>
		<CardFooter class="border-t">
			<div class="flex flex-col gap-3">
				<h1 class="font-code flex items-center gap-2 text-xl">
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
			class="invisible sm:visible fixed bottom-0 left-0 mb-6 flex w-full flex-col items-center gap-y-4 text-muted-foreground"
		>
			<p>Check out what I've done!</p>
			<CircleArrowDown class="animate-bounce" size={32} />
		</div>
	{/if}
</section>

<section
	id="projects"
	class="flex min-h-screen flex-col items-center justify-center gap-y-12 py-16"
>
	<h1 class="font-code text-center text-lg md:text-3xl">Some projects I've worked on...</h1>

	<ProjectCard
		icon="Puzzle"
		title="Hangman Web App"
		deploymentLink="https://app.nklein.xyz"
		technologies={[
			'React',
			'Next.js',
			'TypeScript',
			'Prisma',
			'PostgreSQL',
			'tailwindCSS',
			'Docker'
		]}
		description={'Full-stack word-guessing game built with \
			React, Prisma, and PostgreSQL. It features a daily global \
			challenge word and player statistics, all managed server-side for consistent gameplay across \
			users. The application is self-hosted on my home server for complete control over deployment \
			and data.'}
		image="/hangman.jpg"
	/>

	<ProjectCard
		icon="Globe"
		title="This Website!"
		technologies={['Svelte', 'Typescript', 'drizzle', 'PostgreSQL', 'tailwindCSS', 'Docker']}
		description={'Svelte application fed with data from a PostgreSQL database via drizzle, \
			to allow for dynamic content changes without redeployment.  Deployed with Github \
			Actions CI/CD pipeline on my home server.'}
		image=""
	/>
</section>

<section
	id="featured-photography"
	class="flex min-h-[75vh] flex-col items-center justify-center gap-y-12 py-16"
>

	<h1 class="font-code text-center text-lg md:text-3xl">photography blah blah</h1>

	<div class="border-2 border-foreground p-32">
		placeholder
	</div>

	<Button>
		<a href="/photography">See Full Gallery</a>
	</Button>

</section>

<section
	id="featured-hikes"
	class="flex min-h-[75vh] flex-col items-center justify-center gap-y-12 py-16"
>

	<h1 class="font-code text-center text-lg md:text-3xl">hikes blah blah</h1>

	<div class="border-2 border-foreground p-32">
		placeholder
	</div>

	<Button>
		<a href="/photography">See Hikes </a>
	</Button>

</section>