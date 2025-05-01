<script lang="ts">
  import { onMount } from 'svelte';
  import { solveChallenge } from './lib/challenge';

  interface ChallengeResponse {
    challenge: string;
    token: string;
    difficulty: number;
  }

  let challenge: string = '';
  let difficulty: number = 0;
  let token: string = '';
  let userId: string = '';
  let error: string = ''; 
  let nonce: string = '';
  let key: string = '';

  async function getChallenge(): Promise<void> {
    try {
      const response = await fetch('http://localhost:8080/challenge/first');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data: ChallengeResponse = await response.json();
      challenge = data.challenge;
      difficulty = data.difficulty;
      token = data.token;

      nonce = await solveChallenge(challenge, difficulty);

      console.log(data);
      error = '';
      setTimeout(() => {
        submitSolution();
      }, 0); // Delay for 1 second before submitting the solution
    } catch (e) {
      error = `Failed to get challenge: ${e instanceof Error ? e.message : String(e)}`;
      console.error('Error getting challenge:', e);
    }
  }

  async function submitSolution(): Promise<void> {
    try {
      const response = await fetch('http://localhost:8080/challenge/solve', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          challenge,
          nonce,
          difficulty,
          token,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log('Solution submitted:', data);
    } catch (e) {
      error = `Failed to submit solution: ${e instanceof Error ? e.message : String(e)}`;
      console.error('Error submitting solution:', e);
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
      <p>Solution: {nonce}</p>
      <p>Key: {token}</p>
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
