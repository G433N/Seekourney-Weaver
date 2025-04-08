<script lang="ts">
    interface SearchResult {
        title: string;
        desc: string;   /*fix how interface should look to correspond with backend*/
    }

    let query: string = '';
    let results: SearchResult[] = [];

    async function search(): Promise<void> {
        const res = await fetch(`http://localhost:8080/search?q=${query}`); /*change to url we are using for backend*/
        results = await res.json() as SearchResult[];
    }

</script>

<h1>Seekourny Weaver</h1>

<input bind:value={query} type="text" placeholder="Write your search here!">
<button on:click={search}> Search </button>

{#if results.length > 0}
    <ul>
        {#each results as res}
            <li> {res.title} </li>
        {/each}
    </ul>
{/if}