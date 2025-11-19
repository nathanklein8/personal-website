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
		Github,
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
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let showScrollIcon = $state(false);
	let copied = $state(false);

	onMount(() => {
		const isMobile: boolean = Device.isPhone || Device.isTablet;

		if (window.scrollY > window.innerHeight * 0.05) {
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
	class="flex min-h-screen flex-col justify-center items-center gap-y-16 bg-gradient-to-b from-green-700 to-background to-55% py-12 pt-24 sm:pt-32 md:to-70%"
>
	<div class="space-y-8">
		<div class="space-y-4 flex flex-col items-center">
			<!-- <div class="h-48 w-48 overflow-hidden rounded-full border-4 border-stone-100 shadow-xl">
				<img src="/me.jpg" alt="me" />
			</div> -->
			<h1 class="font-code text-3xl md:text-4xl">Hi! I'm Nathan</h1>
			<h2 class="text-xl md:text-2xl text-primary/80">
				<TypeWriter typeSpeed={150} deleteSpeed={200} texts={['CS @ RIT', 'New Grad']} />
			</h2>
		</div>
	</div>

	{#if data.landingCard}
		<Card class="max-w-[95vw] shadow-lg md:max-w-[80vw] lg:max-w-[65vw] xl:max-w-[50vw]">
			<CardContent class="grid gap-6 sm:grid-cols-2">
				<div class="space-y-4">
					<h1 class="font-code flex items-center gap-2 text-xl">
						<BadgeQuestionMark size={24} />
						About Me
					</h1>
					<p>
						{data.landingCard.bio}
					</p>
				</div>

				<div class="space-y-4">
					<h1 class="font-code flex items-start gap-2 text-xl">
						<MessageCircle size={24} />
						Connect
					</h1>
					<div class="space-y-3">
						<div class="flex flex-row gap-4">
							<p class="flex max-w-fit items-center gap-0.5 underline-offset-4 hover:underline">
								<a href={data.landingCard?.linkedin} target="_blank" rel="noopener noreferrer">
									LinkedIn
								</a>
								<ArrowUpRight size={20} />
							</p>
							<p class="flex max-w-fit items-center gap-0.5 underline-offset-4 hover:underline">
								<a href={data.landingCard?.github} target="_blank" rel="noopener noreferrer">
									Github
								</a>
								<ArrowUpRight size={20} />
							</p>
						</div>
						<button
							class={cn(
								copied ? 'text-muted-foreground' : '',
								'font-code-wide flex w-full flex-row items-center rounded-lg border-1 border-foreground/25 bg-background/80 p-3 text-sm'
							)}
							onclick={copyEmail}
						>
							$ {data.landingCard.email}
							<span class="flex grow"></span>
							{#if copied}
								<ClipboardCheck size={20} />
							{:else}
								<Clipboard size={20} />
							{/if}
						</button>
					</div>
				</div>
			</CardContent>
			<CardFooter class="border-t">
				<div class="flex flex-col gap-3">
					<h1 class="font-code flex items-center gap-2 text-xl">
						<BookOpenCheck size={24} />
						I have experience with
					</h1>
					{#each data.landingCard.skills as skillClass, i}
						<div class="flex flex-wrap gap-2">
							{#each skillClass as skill}
								<Badge variant={i % 2 === 0 ? 'outline' : 'secondary'}>
									{skill}
								</Badge>
							{/each}
						</div>
					{/each}
				</div>
			</CardFooter>
		</Card>
	{/if}

	{#if showScrollIcon}
		<div
			in:fade={{ duration: 1200, easing: expoIn }}
			out:fade={{ duration: 200 }}
			class="invisible fixed bottom-0 left-0 mb-6 flex w-full flex-col items-center gap-y-4 text-muted-foreground sm:visible"
		>
			<p>Check out what I've done!</p>
			<CircleArrowDown class="animate-bounce" size={32} />
		</div>
	{/if}
</section>

<section id="projects" class="flex min-h-screen flex-col items-center justify-center gap-y-8 py-12">
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
	class="flex min-h-[75vh] flex-col items-center justify-center gap-y-12 bg-secondary py-16"
>
	<h1 class="font-code text-center text-lg md:text-3xl">Photography</h1>

	<div class="border-2 border-foreground bg-background p-32">under construction</div>

	<Button>
		<a href="/photography">My Full Gallery</a>
	</Button>
</section>

<section
	id="featured-hike"
	class="flex min-h-[75vh] flex-col items-center justify-center gap-y-12 py-16"
>
	<h1 class="font-code text-center text-lg md:text-3xl">Hiking</h1>

	<div class="border-2 border-foreground bg-background p-32">under construction</div>

	<Button>
		<a href="/photography">Places I've Gone</a>
	</Button>
</section>

<section id="footer" class="h-16 bg-secondary flex flex-col justify-center">
	<p class="text-muted-foreground text-center italic">
		Created by Nathan Klein
	</p>
</section>
