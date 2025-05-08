<script lang="ts">
  import { onMount } from 'svelte';
  import { solveChallenge } from './lib/challenge';
  import { getFirstChallenge, solveFirstChallenge, registerUser, waitChat, updateChat, sendMessage, quitChat } from './lib/api';
  import { E2ECryptoHandler } from './lib/e2e-crypto';

  // State
  let state: 'initial' | 'solving' | 'registration' | 'waiting' | 'chatting' | 'disconnected' = 'initial';
  let error: string = '';
  let loading: boolean = false;

  // Challenge data
  let challenge: string = '';
  let difficulty: number = 0;
  let token: string = '';
  let userId: string = '';
  let nonce: string = '';

  // User data
  let nickname: string = '';
  let tags: string[] = [];
  let cryptoHandler: E2ECryptoHandler;

  // Chat data
  let messages: { text: string; fromPeer: boolean; timestamp: Date }[] = [];
  let messageInput: string = '';
  let peerNickname: string | null = null;
  let chatUpdateInterval: ReturnType<typeof setInterval>;
  let messagesContainer: HTMLDivElement;

  function scrollToBottom() {
    if (messagesContainer) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }

  onMount(() => {
    cryptoHandler = new E2ECryptoHandler();
    return () => {
      if (chatUpdateInterval) {
        clearInterval(chatUpdateInterval);
      }
    };
  });

  async function startChat() {
    state = 'solving';
    loading = true;
    try {
      // Step 1: Get and solve initial challenge
      const data = await getFirstChallenge();
      challenge = data.challenge;
      difficulty = data.difficulty;
      token = data.token;
      nonce = await solveChallenge(challenge, difficulty);

      // Step 2: Submit solution and get user ID
      const solveResp = await solveFirstChallenge({
        challenge,
        nonce,
        difficulty,
        token,
      });

      challenge = solveResp.challenge;
      difficulty = solveResp.difficulty;
      userId = solveResp.userId;
      nonce = await solveChallenge(challenge, difficulty);
      state = 'registration';
    } catch (e) {
      error = `Error in authentication flow: ${e instanceof Error ? e.message : String(e)}`;
    } finally {
      loading = false;
    }
  }

  async function submitRegistration() {
    if (!nickname) {
      error = 'Please enter a nickname';
      return;
    }

    loading = true;
    try {
      await cryptoHandler.generateKeyPair();

      await registerUser({
        challenge,
        nonce,
        difficulty,
        userId,
        nickname,
        publicKey: cryptoHandler.publicKeyBase64,
        tags,
      });

      state = 'waiting';
      waitForPeer();
    } catch (e) {
      state = 'initial';
    } finally {
      loading = false;
    }
  }

  async function waitForPeer() {
    try {
      while (state === 'waiting') {
        const response = await waitChat(userId);
        
        if (response.status === 'assigned' && response.peerPublicKey && response.nickname) {
          await cryptoHandler.setPeerPublicKey(response.peerPublicKey);
          peerNickname = response.nickname;
          state = 'chatting';
          messages = [];
          startChatUpdates();
          break;
        }
        
        await new Promise(resolve => setTimeout(resolve, 2000));
      }
    } catch (e) {
      state = 'initial';
    }
  }

  function startChatUpdates() {
    chatUpdateInterval = setInterval(checkForMessages, 1000);
  }

  async function checkForMessages() {
    const response = await updateChat(userId);

    if (response?.error === 'User not found') {
      clearInterval(chatUpdateInterval);
      messages = [{ text: 'Logout :(', fromPeer: true, timestamp: new Date() }, ...messages];
      return;
    }

    if (response.ok === false) {
      return;
    }
    
    if (response.status === 'closed') {
      clearInterval(chatUpdateInterval);
      messages = [{ text: 'Your peer live chat :(', fromPeer: true, timestamp: new Date() }, ...messages];
      return;
    }

    for (const msg of response.messages) {
      const decrypted = await cryptoHandler.decrypt(msg.message);
      messages = [{ text: decrypted, fromPeer: true, timestamp: new Date(msg.timestamp) }, ...messages];
    }
    if (response.messages.length > 0) {
      scrollToBottom();
    }
  }

  async function sendChatMessage() {
    if (!messageInput.trim()) return;

    try {
      const encrypted = await cryptoHandler.encrypt(messageInput);
      await sendMessage(userId, encrypted);
      
      messages = [{ text: messageInput, fromPeer: false, timestamp: new Date() }, ...messages];
      messageInput = '';
      scrollToBottom();
    } catch (e) {
      error = `Failed to send message: ${e instanceof Error ? e.message : String(e)}`;
    }
  }

  async function handleQuit() {
    try {
      await quitChat(userId);
      clearInterval(chatUpdateInterval);
      state = 'waiting';
      waitForPeer();
    } catch (e) {
      error = `Error quitting chat: ${e instanceof Error ? e.message : String(e)}`;
    }
  }

  function startNewChat() {
    messages = [];
    error = '';
    state = 'initial';
    peerNickname = null;
  }

  // Handle submit with Enter key in chat input
  function handleKeyPress(event: KeyboardEvent) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      sendChatMessage();
    }
  }
</script>

<main>
  {#if error}
    <div class="error">{error}</div>
  {/if}

  {#if state === 'initial'}
    <div class="start-screen">
      <h1>Anonymous Chat</h1>
      <button on:click={startChat} disabled={loading}>
        {loading ? 'Starting...' : 'Start Chat'}
      </button>
    </div>
  {:else if state === 'solving'}
    <div class="solving">
      <h1>Anonymous Chat</h1>
      <h2>Solving challenge...</h2>
      <div class="loading-spinner"></div>
    </div>
  {:else if state === 'registration'}
    <div class="registration">
      <h1>Anonymous Chat</h1>
      <h2>Enter Your Details</h2>
      <form on:submit|preventDefault={submitRegistration}>
        <div class="form-group">
          <label for="nickname">Nickname:</label>
          <input
            type="text"
            id="nickname"
            bind:value={nickname}
            placeholder="Enter a nickname"
            required
          />
        </div>
        <button type="submit" disabled={loading}>
          {loading ? 'Registering...' : 'Join Chat'}
        </button>
      </form>
    </div>
  {:else if state === 'waiting'}
    <div class="waiting">
      <h1>Anonymous Chat</h1>
      <h2>Waiting for a chat partner...</h2>
      <div class="loading-spinner"></div>
    </div>
  {:else if state === 'chatting'}
    <div class="chat">
      <div class="chat-header">
        <h2>Chatting with <span style='color:#4b806d'>{peerNickname}</span></h2>
        <button class="quit-button" on:click={handleQuit}>Leave Chat</button>
      </div>
      
      <div class="messages" bind:this={messagesContainer}>
        {#each messages as message}
          <div class="message {message.fromPeer ? 'peer' : 'self'}">
            <div class="message-time">{message.timestamp.toLocaleTimeString()} | {message.fromPeer ? 'peer' : 'me'}</div>
            <div class="message-content">{message.text}</div>
          </div>
        {/each}
      </div>

      <div class="message-input">
        <textarea
          bind:value={messageInput}
          on:keypress={handleKeyPress}
          placeholder="Type a message..."
          rows="2"
        ></textarea>
        <button on:click={sendChatMessage}>Send</button>
      </div>
    </div>
  {:else if state === 'disconnected'}
    <div class="disconnected">
      <h1>Anonymous Chat</h1>
      <h2>Chat Ended</h2>
      <p>Your chat partner has disconnected.</p>
      <button on:click={startNewChat}>Start New Chat</button>
    </div>
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
    margin: 0 auto;
    padding: 0;
    background-color: #1e1e1e;
    box-shadow: 0 0 12px rgba(0, 0, 0, 0.6);
    min-height: 100vh;
  }

  h1 {
    text-align: center;
    margin-bottom: 2rem;
    color: #ffffff;
  }

  h2 {
    margin: 0;
    font-weight: normal;
  }

  .error {
    color: #ff6b6b;
    background-color: #2e2e2e;
    padding: 1rem;
    border: 1px solid #ff6b6b;
    border-radius: 4px;
    margin-bottom: 1rem;
  }

  .start-screen, .solving, .registration, .waiting, .disconnected {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding-top: 22vh;
  }

  .form-group {
    margin-bottom: 1rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
  }

  input, textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #333;
    border-radius: 4px;
    background: #2a2a2a;
    color: #fff;
    font-size: 1rem;
  }

  button {
    background-color: #4a4a4a;
    color: white;
    border: none;
    padding: 0.8rem 1.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.2s;
  }

  button:hover {
    background-color: #5a5a5a;
  }

  button:disabled {
    background-color: #3a3a3a;
    cursor: not-allowed;
  }

  .chat {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }

  .chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background-color: #2a2a2a;
    border-radius: 4px 4px 0 0;
  }

  .messages {
    flex-grow: 1;
    overflow-y: auto;
    padding: 1rem;
    display: flex;
    flex-direction: column-reverse;
    gap: 0.5rem;
    background-color: #252525;
  }

  /* Scrollbar styling for Webkit browsers (Chrome, Safari, etc.) */
  .messages::-webkit-scrollbar {
    width: 8px;
  }

  .messages::-webkit-scrollbar-track {
    background: #1a1a1a;
    border-radius: 4px;
  }

  .messages::-webkit-scrollbar-thumb {
    background: #4a4a4a;
    border-radius: 4px;
  }

  .messages::-webkit-scrollbar-thumb:hover {
    background: #5a5a5a;
  }

  /* Scrollbar styling for Firefox */
  .messages {
    scrollbar-width: thin;
    scrollbar-color: #4a4a4a #1a1a1a;
  }

  .message {
    max-width: 96%;
    padding: 0.4rem;
    border-radius: 4px;
    position: relative;
  }

  .message.peer {
    background-color: #283531;
    align-self: flex-start;
  }

  .message.self {
    background-color: #2b3044;
    align-self: flex-start;
  }

  .message-time {
    font-size: 0.5rem;
    color: #999;
  }

  .message-input {
    display: flex;
    gap: 1rem;
    padding: 1rem 1rem 2rem 1rem;
    background-color: #2a2a2a;
    border-radius: 0 0 4px 4px;
  }

  .message-input textarea {
    flex-grow: 1;
    resize: none;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #333;
    border-top: 4px solid #fff;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .quit-button {
    background-color: #442827;
    color: #634a4b;
  }

  .quit-button:hover {
    background-color: #51201d;
    color: #907676;
  }
</style>
