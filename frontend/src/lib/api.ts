const API_BASE_URL = 'http://localhost:8080';

export interface ChallengeResponse {
    challenge: string;
    token: string;
    difficulty: number;
}

export interface SolveFirstChallengeResponse {
    userId: string;
    challenge: string;
    token: string;
    difficulty: number;
}

export interface RegisterUserRequest {
    challenge: string;
    token: string;
    difficulty: number;
    userId: string;
    nonce: string;
    nickname: string;
    publicKey: string;
    tags: string[];
}

export async function getFirstChallenge(): Promise<ChallengeResponse> {
    const response = await fetch(`${API_BASE_URL}/challenge/first`);
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
}

export async function solveFirstChallenge(params: {
    challenge: string;
    nonce: string;
    difficulty: number;
    token: string;
}): Promise<SolveFirstChallengeResponse> {
    const response = await fetch(`${API_BASE_URL}/challenge/solve`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(params),
    });

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
}

export async function registerUser(params: RegisterUserRequest): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/users/register`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(params),
    });

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
}

export async function waitChat(userId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/users/waitChat/${userId}`);
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
}

export async function updateChat(userId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/chat/update/${userId}`);
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
}

export async function sendMessage(userId: string, message: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/chat/message/send/${userId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ message }),
    });

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
}

export async function quitChat(userId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/chat/quit/${userId}`);

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
}
