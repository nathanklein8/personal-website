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
	import ErrorCard from '@nk/shared/components/error-card';
	import LandingCard from '$lib/components/landing-card.svelte';
	import { Carousel, CarouselContent, CarouselItem, CarouselNext, CarouselPrevious } from '$lib/components/ui/carousel';

	let { data }: PageProps = $props();

	let showScrollIcon = $state(false);

	onMount(() => {
		const isMobile: boolean = Device.isPhone || Device.isTablet;

		if (window.scrollY > window.innerHeight * 0.05) {
			return;
		} else {
			showScrollIcon = !isMobile;
		}
		const handleScroll = () => {
			const scrollY = window.scrollY;
			const threshold = window.innerHeight * 0.05;
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
	<div class="font-code flex flex-col items-center space-y-4">
		<h1 class="text-2xl sm:text-3xl">Hi! I'm Nathan</h1>
		<h2 class="text-lg sm:text-xl">
			<TypeWriter typeSpeed={150} deleteSpeed={200} texts={['CS @ RIT', 'New Grad']} />
		</h2>
	</div>

	{#if data.landingCard && !data.landingCard.error}
		<LandingCard {...data.landingCard} />
	{/if}

	{#if data.landingCard.error}
		<ErrorCard description={data.landingCard.error} />
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

<section
	id="projects"
	class="flex min-h-screen flex-col items-center justify-center gap-y-12 py-12"
>
	<h1 class="font-code mx-3 max-w-2xl text-center text-lg sm:text-xl">
		Some projects I've worked on...
	</h1>

	<div class="flex flex-col gap-y-8">
		{#each data.projects as proj}
			<ProjectCard {...proj} />
		{/each}
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
				<ArrowUpRight class="h-4 w-4" />
			</a>
		</span>
	</h1>
</section>

<section
	id="featured-photography"
	class="flex min-h-[85vh] flex-col items-center justify-center gap-y-12 bg-secondary py-16"
>
	<h1 class="font-code text-center text-lg sm:text-xl">Featured Photography</h1>

	<div class="grid gap-8 justify-items-center px-4 w-full">
		{#if data.featuredPhotos && data.featuredPhotos.length > 0}
			<Carousel class="max-w-[85vw]" opts={{ loop: true }}>
				<CarouselContent class="">
					{#each data.featuredPhotos as photo, index}
						<CarouselItem class="basis-auto">
							<div class="w-[50vw] md:w-[75vw] bg-muted flex items-center justify-center overflow-hidden">
								<img
									src={`${data.apiURL}/api/photos/${photo.id}/image`}
									alt={photo.altText || photo.title}
									class="aspect-3/2 object-contain w-full h-full"
								/>
							</div>
						</CarouselItem>
					{/each}
				</CarouselContent>
				<CarouselPrevious />
				<CarouselNext />
			</Carousel>
		{/if}

		{#if !(data.featuredPhotos && data.featuredPhotos.length > 0)}
			<div class="text-center border-2 border-foreground bg-background p-24 text-center text-muted-foreground lg:p-32">
				<p>No featured photos yet. Check back later!</p>
			</div>
		{/if}
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
