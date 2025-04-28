<script>
  import { onMount } from 'svelte';
  import { solveChallenge } from './lib/challenge';

  let challenge = '';
  let difficulty = 0;
  let token = '';
  let userId = '';
  let error = ''; 
  let solution = '';

  async function getChallenge() {
    try {
      const response = await fetch('http://localhost:8080/challenge/first');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      challenge = data.challenge;
      token = data.token;
      userId = data.userId;
      difficulty = data.difficulty;

      solution = await solveChallenge(challenge, 5000);

      console.log(data);
      error = '';
    } catch (e) {
      error = `Failed to get challenge: ${e.message}`;
      console.error('Error getting challenge:', e);
    }
  }

  onMount(() => {
    getChallenge();
  });
</script>

<main>
  <h1>Anonymous Chat</h1>

  {#if error}
    <div class="error">{error}</div>
  {/if}

  {#if challenge}
    <div class="challenge">
      <h2>Challenge Received</h2>
      <p>Challenge: {challenge}</p>
      <p>Solution: {solution}</p>
      <p>Token: {token}</p>
      <p>UserId: {userId}</p>
    </div>
  {:else}
    <div class="loading">Loading challenge...</div>
  {/if}
</main>

<style>
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
  }

  h1 {
    text-align: center;
    margin-bottom: 2rem;
  }

  .error {
    color: red;
    padding: 1rem;
    border: 1px solid red;
    border-radius: 4px;
    margin-bottom: 1rem;
  }

  .challenge {
    padding: 1rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    background-color: #f5f5f5;
  }

  .loading {
    text-align: center;
    padding: 2rem;
  }
</style>
