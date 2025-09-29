<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import '../app.css';

  let mounted = false;

  onMount(async () => {
    await auth.checkAuth();
    mounted = true;
  });

  $: if (mounted && !$auth.isAuthenticated && $page.url.pathname !== '/login') {
    goto('/login');
  }
</script>

{#if !mounted}
  <div class="min-h-screen flex items-center justify-center">
    <div class="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-600"></div>
  </div>
{:else if $auth.isAuthenticated}
  <div class="min-h-screen bg-gray-50">
    <slot />
  </div>
{:else}
  <div class="min-h-screen bg-gray-50">
    <slot />
  </div>
{/if}
