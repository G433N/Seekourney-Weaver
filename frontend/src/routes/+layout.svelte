<script lang="ts">
	let { children } = $props();
	import { page } from '$app/stores';
</script>

<header>
	<h1>Seekourny <br />Weaver</h1>
</header>

<nav>
	<a href="/" class:active={$page.url.pathname === '/'}>Home</a>
	<a href="/settings" class:active={$page.url.pathname === '/settings'}>Settings</a>
	<a href="/about" class:active={$page.url.pathname === '/about'}>About</a>
</nav>

{@render children()}

