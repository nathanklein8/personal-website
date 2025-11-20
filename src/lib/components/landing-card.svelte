<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import CardFooter from '$lib/components/ui/card/card-footer.svelte';
	import {
		ArrowUpRight,
		BadgeQuestionMark,
		BookOpenCheck,
		Clipboard,
		ClipboardCheck,
		MessageCircle
	} from '@lucide/svelte';
	import { cn } from '$lib/utils';

    interface Props {
        bio: string,
        email: string,
        linkedin: string,
        github: string,
        skills: string[][],
    }

	let { bio, email, linkedin, github, skills }: Props = $props();

	let copied = $state(false);

	async function copyEmail() {
		try {
			await navigator.clipboard.writeText(email);
			copied = true;
			setTimeout(() => (copied = false), 3000);
		} catch (e) {
			console.error('Clipboard write failed', e);
		}
	}
</script>

<Card class="mx-3 max-w-3xl shadow-lg">
	<CardContent class="grid gap-6 sm:grid-cols-2">
		<div class="space-y-4">
			<h1 class="font-code flex items-center gap-2 text-xl">
				<BadgeQuestionMark size={24} class="mb-0.5" />
				About Me
			</h1>
			<p class="text-foreground/80">
				{bio}
			</p>
		</div>

		<div class="space-y-4">
			<h1 class="font-code flex items-start gap-2 text-xl">
				<MessageCircle size={24} class="mb-0.5" />
				Connect
			</h1>
			<div class="space-y-3">
				<div class="flex flex-row gap-4">
					<p class="flex max-w-fit items-center gap-0.5 underline-offset-4 hover:underline">
						<a href={linkedin} target="_blank" rel="noopener noreferrer">
							LinkedIn
						</a>
						<ArrowUpRight size={20} />
					</p>
					<p class="flex max-w-fit items-center gap-0.5 underline-offset-4 hover:underline">
						<a href={github} target="_blank" rel="noopener noreferrer">
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
					$ {email}
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
			{#each skills as skillClass, i}
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
