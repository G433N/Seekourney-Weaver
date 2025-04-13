<script lang="ts">
	let { children } = $props();
</script>

<header>
	<h1 style="font-size: 2.5rem; font-weight: bold; line-height: 1.2;">Seekourny <br/>Weaver</h1>
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
<nav style="background-color: #AEC6DF; display: flex; justify-content: center; gap: 6rem; padding: 1rem;">
	<a href="/" class="active">Home</a>
	<a href="/settings">Settings</a>
	<a href="/about">About</a>
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
