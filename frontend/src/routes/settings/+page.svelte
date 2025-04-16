<script>
    let showFiles = $state(true);
    let showWebpages = $state(true);
    let showAllResults = $state(true);
    let maxResults = $state(100);

    let cpuDefault = $state(true);
    let maxCores = $state(8);
    let cpuCores = $state(4);
</script>

<main style="max-width: 600px;">
    <h2>Filter</h2>
    <div class="box">
        <h3>Show:</h3>
        <label class="toggle">
            <input type="checkbox" bind:checked={showFiles} />
            Files
        </label>
        <label class="toggle">
            <input type="checkbox" bind:checked={showWebpages} />
            Webpages
        </label>
    </div>

    <div class="box column">
        <label class="toggle">
            <input type="checkbox" bind:checked={showAllResults} />
            Show all results
        </label>

        {#if !showAllResults}
            <label class="max-label">
                Max results shown:
                <input type="number" bind:value={maxResults} />
            </label>
        {/if}
    </div>

    <h2>Searching</h2>
    <div class="box column">
        <label class="toggle">
            <input type="checkbox" bind:checked={cpuDefault} />
            default
        </label>
        {#if !cpuDefault}
            <div class="slider">
                <label for="cpuSlider">CPU cores used:</label>
                <input id="cpuSlider" type="range" min="1" max={maxCores} bind:value={cpuCores} />
                <span>{cpuCores}</span>
            </div>
        {/if}
    </div>

    <button class="custom-scraper">Add a custom scraper</button>
</main>


<style>

    .box {
        background: white;
        border: 2px solid #F0EEEE;
        border-radius: 12px;
        padding: 1rem 1.5rem;
        margin-bottom: 1.5rem;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        align-items: center;
        gap: 1rem;
    }

    .box.column {
        flex-direction: column;
        align-items: flex-start;
    }

    .toggle {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 1rem;
    }

    h3 {
        font-family: "Jost", sans-serif;
	    font-optical-sizing: auto;
	    font-weight: 500;
	    font-style: normal;
    }

    .slider {
        width: 100%;
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-top: 0.5rem;
    }

    .max-label {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
        margin-top: 0.5rem;
    }

    input[type="number"] {
        padding: 0.3rem;
        border: 1px solid #ccc;
        border-radius: 4px;
        width: 100px;
    }

    input[type="range"] {
        flex-grow: 1;
    }

    .custom-scraper {
        display: block;
        margin: 0 auto;
        background-color: #92acc2;
        border: none;
        border-radius: 6px;
        padding: 0.75rem 1.5rem;
        font-size: 1rem;
        color: white;
        cursor: pointer;
    }

    .custom-scraper:hover {
        background-color: #7d98b2;
    }
  
</style>
