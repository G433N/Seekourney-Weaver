<script lang="ts">
	let { children } = $props();
	import { page } from '$app/stores';
	import { derived } from 'svelte/store';
</script>

<header>
	<h1>Seekourny <br/>Weaver</h1>
</header>

<!--
<nav>
	<ul>
		<li>
			<a href="/">Home</a>
		</li>
		<li>
			<a href="/about">About</a>
		</li>
	</ul>
</nav>
-->
<nav>
	<a href="/" class:active={$page.url.pathname === '/'}>Home</a>
	<a href="/settings" class:active={$page.url.pathname === '/settings'}>Settings</a>
	<a href="/about" class:active={$page.url.pathname === '/about'}>About</a>
</nav>

{@render children()}

<style>
	ul {
		list-style-type: none;
		margin: 0;
		padding: 0;
		overflow: hidden;
		background-color: #aec6df;
	}
</style>
