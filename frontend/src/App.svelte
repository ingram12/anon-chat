<script lang="ts">
  import { onMount } from 'svelte';
  import { solveChallenge } from './lib/challenge';

  interface ChallengeResponse {
    challenge: string;
    token: string;
    difficulty: number;
  }

  interface RegisterUserRequest {
    challenge: string;
    token: string;
    difficulty: number;
    userId: string;
    nonce: string;
    nickname: string;
    publicKey: string;
    tags: string[];
  }

  interface SolveFirstChallengeResponse {
    userId: string;
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
  let nickname: string = 'anon228';
  let publicKey: string = 'key228';
  let tags: string[] = [];

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

      const data: SolveFirstChallengeResponse = await response.json();
      userId = data.userId;
      challenge = data.challenge;
      difficulty = data.difficulty;
      token = data.token;
      nonce = await solveChallenge(challenge, difficulty);

      console.log('Solution submitted:', data);

      setTimeout(() => {
        registrationUser();
      }, 0);
    } catch (e) {
      error = `Failed to submit solution: ${e instanceof Error ? e.message : String(e)}`;
      console.error('Error submitting solution:', e);
    }
  }

  async function registrationUser(): Promise<void> {
    try {
      const response = await fetch('http://localhost:8080/users/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          challenge,
          nonce,
          difficulty,
          token,
          userId,
          nickname,
          publicKey,
          tags,
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
  :global(body) {
    background-color: #121212;
    color: #e0e0e0;
    font-family: system-ui, sans-serif;
    margin: 0;
    padding: 0;
  }

  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    background-color: #1e1e1e;
    border-radius: 3px;
    box-shadow: 0 0 12px rgba(0, 0, 0, 0.6);
  }

  h1 {
    text-align: center;
    margin-bottom: 2rem;
    color: #ffffff;
  }

  .error {
    color: #ff6b6b;
    background-color: #2e2e2e;
    padding: 1rem;
    border: 1px solid #ff6b6b;
    border-radius: 3px;
    margin-bottom: 1rem;
  }

  .challenge {
    padding: 1rem;
    border: 1px solid #333;
    border-radius: 3px;
    background-color: #2a2a2a;
  }

  .challenge p {
    margin: 0.5rem 0;
    color: #ccc;
    font-family: monospace;
  }

  .loading {
    text-align: center;
    padding: 2rem;
    color: #aaa;
  }

  a {
    color: #90caf9;
    text-decoration: none;
  }

  a:hover {
    text-decoration: underline;
  }
</style>
