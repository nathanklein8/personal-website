<script lang="ts">
	import { ArrowUpRight, CircleArrowDown } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { expoIn } from 'svelte/easing';
	import { fade } from 'svelte/transition';
	import Device from 'svelte-device-info';
	import ProjectCard from '$lib/components/project-card.svelte';
	import { TypeWriter } from 'svelte-typewrite';
	import { Button } from '$lib/components/ui/button';
	import type { PageProps } from './$types';
	import ErrorCard from '$lib/components/error-card.svelte';
	import LandingCard from '$lib/components/landing-card.svelte';

	let { data }: PageProps = $props();

	let showScrollIcon = $state(false);

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
</script>

<section
	id="landing"
	class="flex min-h-screen flex-col items-center justify-center gap-y-12 bg-gradient-to-b from-green-700 to-background to-55% pt-24 pb-12 md:to-70%"
>
	<!-- heading that displays my name and a typewriter effect tagline -->
	<div class="font-code flex flex-col items-center space-y-4">
		<!-- <div class="h-48 w-48 overflow-hidden rounded-full border-4 border-stone-100 shadow-xl">
				<img src="/me.jpg" alt="me" />
			</div> -->
		<h1 class="text-2xl sm:text-3xl">Hi! I'm Nathan</h1>
		<h2 class="text-lg sm:text-xl">
			<TypeWriter typeSpeed={150} deleteSpeed={200} texts={['CS @ RIT', 'New Grad']} />
		</h2>
	</div>

	<!-- display landing 'about me' card, when data is present, and there is no error field on the json -->
	{#if data.landingCard && !data.landingCard.error}
		<LandingCard
			bio={data.landingCard.bio}
			email={data.landingCard.email}
			linkedin={data.landingCard.linkedin}
			github={data.landingCard.github}
			skills={data.landingCard.skills}
		/>
	{/if}

	<!-- display an error alert if there is an error field on the landing card json -->
	{#if data.landingCard.error}
		<ErrorCard description={data.landingCard.error} />
	{/if}

	<!-- show a little fixed position hint scroll -->
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

<section
	id="projects"
	class="flex min-h-screen flex-col items-center justify-center gap-y-12 py-12"
>
	<h1 class="font-code mx-3 max-w-2xl text-center text-lg sm:text-xl">
		Some projects I've worked on...
	</h1>

	<div class="flex flex-col gap-y-8">
		<ProjectCard
			icon="Puzzle"
			title="Hangman Web App"
			description={'Full-stack word-guessing game built with \
			React, Prisma, and PostgreSQL. It features a daily global \
			challenge word and player statistics, all managed server-side for consistent gameplay across \
			users. The application is self-hosted on my home server for complete control over deployment \
			and data.'}
			technologies={[
				'React',
				'Next.js',
				'TypeScript',
				'Prisma',
				'PostgreSQL',
				'tailwindCSS',
				'Docker'
			]}
			image="/hangman.jpg"
			deploymentLink="https://app.nklein.xyz"
		/>

		<ProjectCard
			icon="Globe"
			title="This Website!"
			description={'Svelte C C++ C application fed with data from a PostgreSQL database via a GoLang backend, \
			to allow for dynamic content changes without redeployment.  Deployed with Github \
			Actions CI/CD pipeline on my home server.'}
      technologies={['Svelte', 'C', 'Typescript', 'Go', 'PostgreSQL', 'tailwindCSS', 'Docker']}
		/>
	</div>

	<h1 class="font-code mx-3 max-w-2xl text-center text-lg text-balance sm:text-xl">
		All of my source code & more projects are on
		<span class="underline-offset-4 hover:underline">
			<a
				class="inline-flex items-center gap-0.5 underline hover:underline-offset-6"
				href={data.landingCard?.github}
				target="_blank"
				rel="noopener noreferrer"
			>
				Github
			</a>
		</span>
	</h1>
</section>

<section
	id="featured-photography"
	class="flex min-h-[85vh] flex-col items-center justify-center gap-y-12 bg-secondary py-16"
>
	<h1 class="font-code text-center text-lg sm:text-xl">Photography</h1>

	<div
		class="border-2 border-foreground bg-background p-24 text-center text-muted-foreground lg:p-32"
	>
		under construction...
	</div>

	<Button class="font-code">
		<a href="/photography">My Full Gallery</a>
	</Button>
</section>

<section
	id="featured-hike"
	class="flex min-h-[85vh] flex-col items-center justify-center gap-y-12 py-16"
>
	<h1 class="font-code text-center text-lg sm:text-xl">Hiking</h1>

	<div
		class="border-2 border-foreground bg-background p-24 text-center text-muted-foreground lg:p-32"
	>
		under construction...
	</div>

	<Button class="font-code">
		<a href="/photography">Places I've Gone</a>
	</Button>
</section>

<section id="footer" class="flex h-16 flex-col justify-center bg-secondary">
	<p class="text-center text-muted-foreground italic">Created by Nathan Klein</p>
</section>
