<script lang="ts">
 import { SunIcon, MoonIcon } from "@lucide/svelte";
 
 import { resetMode, setMode, mode } from "mode-watcher";
 import * as DropdownMenu from "@nk/shared/components/ui/dropdown-menu";
 import { buttonVariants } from "@nk/shared/components/ui/button";

 const Icon = $derived($mode === "dark" ? MoonIcon : SunIcon);
</script>

<DropdownMenu.Root>
 <DropdownMenu.Trigger
  class={buttonVariants({ variant: "ghost", size: "icon-sm" })}
 >
  {@render Icon({ class: "h-[1.2rem] w-[1.2rem] !transition-all" })}
  <span class="sr-only">Toggle theme</span>
 </DropdownMenu.Trigger>
 <DropdownMenu.Content align="end">
  <DropdownMenu.Item onclick={() => setMode("light")}>Light</DropdownMenu.Item
  >
  <DropdownMenu.Item onclick={() => setMode("dark")}>Dark</DropdownMenu.Item>
  <DropdownMenu.Item onclick={() => resetMode()}>System</DropdownMenu.Item>
 </DropdownMenu.Content>
</DropdownMenu.Root>