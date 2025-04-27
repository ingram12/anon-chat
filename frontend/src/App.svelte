<script>
  import { onMount } from 'svelte';
  import { GetFirstChallengeRequest } from './proto/users';
  import { UserServiceClient } from './proto/users.client';
  import { solveChallenge } from './lib/challenge';
  import { GrpcWebFetchTransport } from '@protobuf-ts/grpcweb-transport';

  let challenge = '';
  let token = '';
  let userId = '';
  let error = ''; 
  let solution = '';
  const transport = new GrpcWebFetchTransport({
    baseUrl: 'http://localhost:50051'
  });
  const client = new UserServiceClient(transport);

  async function getChallenge() {
    try {
      const request = GetFirstChallengeRequest.create();
      console.log('Sending request...');
      
      const { response } = await client.getFirstChallenge(request);
      
      console.log('Received response:', response);
      
      challenge = response.challenge;
      token = response.token;
      userId = response.userId;
      solution = await solveChallenge(challenge, 3000);
      error = '';
      console.log('Challenge:', challenge);
      console.log('Solution:', solution);
      console.log('Token:', token);
      console.log('UserId:', userId);
    } catch (err) {
      error = 'Error connecting to server';
      console.error('Error:', err);
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
