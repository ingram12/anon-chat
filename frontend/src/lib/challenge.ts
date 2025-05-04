let worker: Worker | null = null;

interface WorkerMessage {
  error?: string;
  nonce?: string;
}

interface WorkerRequest {
  challenge: string;
  difficulty: number;
}

export async function solveChallenge(challenge: string, difficulty: number): Promise<string> {
  return new Promise((resolve, reject) => {
    if (!worker) {
      worker = new Worker(new URL('./challenge.worker.ts', import.meta.url), { type: 'module' });
    }

    const messageHandler = (e: MessageEvent<WorkerMessage>) => {
      if (e.data.error) {
        reject(new Error(e.data.error));
      } else if (e.data.nonce) {
        resolve(e.data.nonce);
      } else {
        reject(new Error('Invalid worker response'));
      }
      worker?.removeEventListener('message', messageHandler);
    };

    const errorHandler = (error: ErrorEvent) => {
      reject(error);
      worker?.removeEventListener('error', errorHandler);
    };

    worker.addEventListener('message', messageHandler);
    worker.addEventListener('error', errorHandler);
    worker.postMessage({ challenge, difficulty } as WorkerRequest);
  });
}

/**
 * Terminates the worker
 */
export function terminateWorker(): void {
  if (worker) {
    worker.terminate();
    worker = null;
  }
} 