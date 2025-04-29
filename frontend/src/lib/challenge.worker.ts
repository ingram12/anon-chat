interface WorkerRequest {
  challenge: string;
  difficulty: number;
}

interface WorkerResponse {
  nonce?: string;
  error?: string;
}

self.onmessage = async (e: MessageEvent<WorkerRequest>) => {
  const { challenge, difficulty } = e.data;
  const maxAttempts = 1000000;

  const enc = new TextEncoder();

  for (let i = 0; i < maxAttempts; i++) {
    let buf = enc.encode(challenge + i);

    for (let j = 0; j < difficulty; j++) {
      buf = new Uint8Array(await crypto.subtle.digest('SHA-256', buf));
    }

    const hashBytes = buf;

    if (hashBytes[0] === 0) {
      self.postMessage({ nonce: i.toString() } as WorkerResponse);
      return;
    }
  }

  self.postMessage({ error: "Can't find nonce" } as WorkerResponse);
}; 