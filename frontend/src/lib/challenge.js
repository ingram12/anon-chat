let worker = null;

/**
 * Решает челлендж, полученный от сервера
 * @param {string} challenge - Челендж для решения
 * @param {number} difficulty - Сложность челленджа
 * @returns {Promise<string>} - Решение челленджа
 */
export async function solveChallenge(challenge, difficulty) {
  return new Promise((resolve, reject) => {
    if (!worker) {
      worker = new Worker(new URL('./challenge.worker.js', import.meta.url), { type: 'module' });
    }

    const messageHandler = (e) => {
      if (e.data.error) {
        reject(new Error(e.data.error));
      } else {
        resolve(e.data.nonce);
      }
      worker.removeEventListener('message', messageHandler);
    };

    const errorHandler = (error) => {
      reject(error);
      worker.removeEventListener('error', errorHandler);
    };

    worker.addEventListener('message', messageHandler);
    worker.addEventListener('error', errorHandler);
    worker.postMessage({ challenge, difficulty });
  });
}

/**
 * Закрывает воркер
 */
export function terminateWorker() {
  if (worker) {
    worker.terminate();
    worker = null;
  }
}