<script lang="ts">
  import { onMount } from 'svelte';
  import { solveChallenge } from './lib/challenge';
  import { getFirstChallenge, solveFirstChallenge, registerUser, waitChat, updateChat, sendMessage } from './lib/api';
  import { E2ECryptoHandler } from './lib/e2e-crypto';

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
    const data = await getFirstChallenge();
    challenge = data.challenge;
    difficulty = data.difficulty;
    token = data.token;
    nonce = await solveChallenge(challenge, difficulty);

    console.log('Challenge received:', data);
    error = '';
  }

  async function submitSolution(): Promise<void> {
    const data = await solveFirstChallenge({
      challenge,
      nonce,
      difficulty,
      token,
    });

    userId = data.userId;
    challenge = data.challenge;
    difficulty = data.difficulty;
    token = data.token;
    nonce = await solveChallenge(challenge, difficulty);

    console.log('Solution submitted:', data);
  }

  async function registrationUser(): Promise<void> {
    const keys = new E2ECryptoHandler();
    await keys.generateKeyPair();

    const data = await registerUser({
      challenge,
      nonce,
      difficulty,
      token,
      userId,
      nickname,
      publicKey: keys.publicKeyBase64,
      tags,
    });

    console.log('Registration submitted:', data);
  }

  async function run(): Promise<void> {
    try {
      // Step 1: Get and solve initial challenge
      await getChallenge();
      
      // Step 2: Submit solution and get user ID
      await submitSolution();
      
      // Step 3: Register user and wait for chat
      await registrationUser();

      const chatData = await waitChat(userId);
      console.log('Chat connection:', chatData);

      let upData = await updateChat(userId);
      console.log('Chat connection:', upData);

      upData = await updateChat(userId);
      console.log('Chat connection:', upData);

      const chat = await sendMessage(userId, 'Hello, world!');

      function sleep(ms: number): Promise<void> {
        return new Promise(resolve => setTimeout(resolve, ms));
      }

      await sleep(1000);

      upData = await updateChat(userId);
      console.log('Chat connection:', upData);
      await sleep(1000);
      
      upData = await updateChat(userId);
      console.log('Chat connection:', upData);
      upData = await updateChat(userId);
      console.log('Chat connection:', upData);

      upData = await updateChat(userId);
      console.log('Chat connection:', upData);

      upData = await updateChat(userId);
      console.log('Chat connection:', upData); 

    } catch (e) {
      error = `Error in authentication flow: ${e instanceof Error ? e.message : String(e)}`;
      console.error('Error:', e);
    }
  }

  onMount(() => {
    run();
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
