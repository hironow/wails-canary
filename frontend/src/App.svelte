<script lang="ts">
  import {
    GetBreedList,
    GetImageUrlsByBreed,
    GetRandomImageUrl,
    SelectFolder,
    SelectFile,
    Search,
  } from "../wailsjs/go/main/App.js";

  let path;
  let pattern;
  let caseInsensitive = false;
  let wholeWord = false;
  let wholeLine = false;
  let filenameOnly = false;
  let filesWoMatches = false;
  let options = {
    caseInsensitive: caseInsensitive,
    wholeWord: wholeWord,
    wholeLine: wholeLine,
    filenameOnly: filenameOnly,
    filesWoMatches: filesWoMatches,
  };
  let results = "";

  let randomImageUrl = "";
  let breeds = [];
  let photos = [];
  let selectedBreed: string;
  let showRandomPhoto = false;
  let showBreedPhotos = false;

  function init() {
    getBreedList();
  }

  init();

  function processOptions(event) {
    options = {
      caseInsensitive: caseInsensitive,
      wholeWord: wholeWord,
      wholeLine: wholeLine,
      filenameOnly: filenameOnly,
      filesWoMatches: filesWoMatches,
    };
  }

  function selectFolder() {
    SelectFolder().then((result) => (path = result));
  }

  function selectFile() {
    SelectFile().then((result) => (path = result));
  }

  function search() {
    results = "Searching...";
    Search(path, pattern, options).then((result) => (results = result));
  }

  function getRandomImageUrl() {
    showRandomPhoto = false;
    showBreedPhotos = false;
    GetRandomImageUrl().then((result) => (randomImageUrl = result));
    showRandomPhoto = true;
  }

  function getBreedList() {
    GetBreedList().then((result) => (breeds = result));
  }

  function getImageUrlsByBreed() {
    init();
    showRandomPhoto = false;
    showBreedPhotos = false;
    GetImageUrlsByBreed(selectedBreed).then((result) => (photos = result));
    showBreedPhotos = true;
  }
</script>

<main>
  <h3>SIMPLE GREP</h3>
  <div class="input-box" id="input">
    <label for="path">Path:</label>
    &nbsp;
    <input
      autocomplete="off"
      bind:value={path}
      class="input"
      id="path"
      type="text"
    />
    &nbsp; &nbsp;
    <label for="pattern">Pattern:</label>
    &nbsp;
    <input
      autocomplete="off"
      bind:value={pattern}
      class="input"
      id="pattern"
      type="text"
    />
    &nbsp;
  </div>
  <div class="checkboxes">
    <span class="checkboxes-with-labels">
      <input
        class="checkbox"
        type="checkbox"
        id="caseInsensitive"
        value="caseInsensitive"
        bind:checked={caseInsensitive}
        on:change={processOptions}
      />
      <label for="caseInsensitive">Case Insensitive</label>
    </span>
    <span class="checkboxes-with-labels">
      <input
        class="checkbox"
        type="checkbox"
        id="wholeWord"
        value="wholeWord"
        bind:checked={wholeWord}
        on:change={processOptions}
      />
      <label for="wholeWord">Whole Word</label>
    </span>
    <span class="checkboxes-with-labels">
      <input
        class="checkbox"
        type="checkbox"
        id="wholeLine"
        value="wholeLine"
        bind:checked={wholeLine}
        on:change={processOptions}
      />
      <label for="wholeLine">Whole Line</label>
    </span>
    <span class="checkboxes-with-labels">
      <input
        class="checkbox"
        type="checkbox"
        id="fileNameOnly"
        value="fileNameOnly"
        bind:checked={filenameOnly}
        on:change={processOptions}
      />
      <label for="fileNameOnly">File Name Only</label>
    </span>
    <span class="checkboxes-with-labels">
      <input
        class="checkbox"
        type="checkbox"
        id="filesWoMatches"
        value="filesWoMatches"
        bind:checked={filesWoMatches}
        on:change={processOptions}
      />
      <label for="filesWoMatches">Files Without Matches</label>
    </span>
  </div>

  <div class="buttons">
    <button class="button" id="select-folder-btn" on:click={selectFolder}>
      Select a Folder
    </button>
    <button class="button" id="select-file-btn" on:click={selectFile}>
      Select a File
    </button>
    <button class="button" id="search-btn" on:click={search}>Search</button>
  </div>
  <div>
    <textarea value={results} />
  </div>

  <h3>Dogs API</h3>
  <div>
    <button class="btn" on:click={getRandomImageUrl}>
      Fetch a dog randomly
    </button>
    Click on down arrow to select a breed
    <select bind:value={selectedBreed}>
      {#each breeds as breed}
        <option value={breed}>
          {breed}
        </option>
      {/each}
    </select>
    <button class="btn" on:click={getImageUrlsByBreed}>
      Fetch by this breed
    </button>
  </div>
  <br />
  {#if showRandomPhoto}
    <img id="random-photo" src={randomImageUrl} alt="No dog found" />
  {/if}
  {#if showBreedPhotos}
    {#each photos as photo}
      <img id="breed-photos" src={photo} alt="No dog found" />
    {/each}
  {/if}
</main>

<style>
  #random-photo {
    width: 600px;
    height: auto;
  }

  #breed-photos {
    width: 300px;
    height: auto;
  }

  .btn:focus {
    border-width: 3px;
  }

  .buttons {
    display: flex;
    flex-direction: row;
    justify-content: space-around;
  }

  .button {
    margin-left: 10px;
    margin-right: 10px;
  }

  .button:focus {
    border-width: 3px;
  }

  .checkbox:focus {
    outline: 2px solid white;
  }

  .checkboxes {
    display: flex;
    flex-direction: row;
    justify-content: space-around;
  }

  .checkboxes-with-labels {
    margin-left: 10px;
    margin-right: 10px;
  }

  textarea {
    width: 100%;
    height: 550px;
    background-color: gray;
  }

  #pattern {
    width: 250px;
  }

  .input-box {
    display: flex;
    justify-content: flex-start;
    padding: 0px 10px;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    width: 450px;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }
</style>
