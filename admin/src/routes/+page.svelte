<script lang="ts">
	import { enhance } from '$app/forms';
	import { buttonVariants } from '$lib/components/ui/button';
	import { Card, CardContent, CardTitle } from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
	// import ErrorCard from '$lib/components/error-card.svelte';
	import type { PageProps } from './$types';
	import { Check, X } from '@lucide/svelte';
	import Header from '$lib/components/Header.svelte';

	let { data, form }: PageProps = $props();
</script>

<section id="landing-card-editor" class="flex flex-col items-center gap-4 py-5 mt-16">
	<h1 class="font-code text-lg">Edit Landing Card</h1>
	<Card class="mx-3 w-2xl shadow-lg">
		<CardContent>
			{#if data.landingCard.error}
				<!-- <ErrorCard description={data.landingCard.error} /> -->
                 uh oh
			{:else}
				<form
					class="flex flex-col items-center space-y-4"
					method="POST"
					action="?/updateLanding"
					use:enhance
				>
					<div class="flex w-full flex-col gap-4">
						<label>
							Bio
							<textarea
								required
								name="bio"
								rows="4"
								defaultValue={data.landingCard.bio}
								class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
							></textarea>
						</label>

						<label>
							Email
							<input
								required
								name="email"
								class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
								defaultValue={data.landingCard.email}
							/>
						</label>

						<label>
							LinkedIn
							<input
								required
								name="linkedin"
								class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
								defaultValue={data.landingCard.linkedin}
							/>
						</label>

						<label>
							GitHub
							<input
								required
								name="github"
								class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
								defaultValue={data.landingCard.github}
							/>
						</label>

						<label>
							Skills
							<span class="text-xs text-muted-foreground">
								(comma separated list, separate categories w/ newline)
							</span>
							<textarea
								required
								name="skills"
								rows="4"
								class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
								defaultValue={data.landingCard.skills.join('\n')}
							></textarea>
						</label>
					</div>
					<div class="flex items-center gap-2">
						<button class={buttonVariants({ size: 'sm' })}> Submit </button>
						{#if form?.success && form?.id == null}
							<Check color={'green'} />
						{:else if form?.failure && form?.id == null}
							<X color={'red'} />
							{form.message}
						{/if}
					</div>
				</form>
			{/if}
		</CardContent>
	</Card>
</section>

<section id="projects-editor" class="flex flex-col items-center gap-4 py-5">
	<h1 class="font-code text-lg">Edit Projects</h1>
	{#each data.projects as project}
		<Card class="mx-3 w-2xl shadow-lg">
			<CardContent>
				<form
					class="flex flex-col items-center space-y-4"
					method="POST"
					action={'?/updateProject'}
					use:enhance
				>
					<div class="flex w-full flex-col gap-2">
						<input
							type="hidden"
							required
							name="id"
							class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
							value={project.id}
						/>
						<div class="flex flex-row gap-4">
							<label>
								Title
								<input
									required
									name="title"
									class="mt-1 w-sm rounded-md border p-2 text-muted-foreground"
									defaultValue={project.title}
								/>
							</label>
							<label>
								Icon
								<input
									required
									name="icon"
									class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
									defaultValue={project.icon}
								/>
							</label>
						</div>
						<label>
							Description
							<textarea
								required
								name="description"
								rows="4"
								defaultValue={project.description}
								class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
							></textarea>
						</label>
						<label>
							Technologies
							<input
								required
								name="technologies"
								defaultValue={project.technologies.join(',')}
								class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
							/>
						</label>
						<label>
							Deployment Link
							<input
								name="deploymentLink"
								class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
								defaultValue={project.deploymentLink ?? ''}
							/>
						</label>
						<div class="flex flex-row gap-4">
							<label>
								Image
								<input
									name="image"
									class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
									defaultValue={project.image ?? ''}
								/>
							</label>
							<label>
								Alt Text
								<input
									name="altText"
									class="mt-1 w-full rounded-md border p-2 text-muted-foreground"
									defaultValue={project.altText ?? ''}
								/>
							</label>
						</div>
					</div>
					<div class="flex items-center gap-2">
						<button class={buttonVariants({ size: 'sm' })}> Submit </button>
						{#if form?.success && form?.id == project.id}
							<Check color={'green'} />
						{:else if form?.failure && form?.id == project.id}
							<X color={'red'} />
							{form.message}
						{/if}
					</div>
				</form>
			</CardContent>
		</Card>
	{/each}
</section>
