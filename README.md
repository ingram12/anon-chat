# Anonymous Chat

<img width="300" alt="image" src="https://github.com/user-attachments/assets/f23fe618-4b71-4d66-98b2-f093dc074115" />
<img width="300" alt="image" src="https://github.com/user-attachments/assets/5b97783e-16b2-4075-9321-31e09b461502" />

### 🔗 Live Preview

You can try the **preview version** of the chat at <a href="https://anon-chat.top" target="_blank" rel="noopener noreferrer">https://anon-chat.top</a>

> ⚠️ This is a **preview instance** for demonstration purposes. It's not production-ready and may be reset or unavailable at any time.

## Overview

An anonymous chat application with **end-to-end encrypted messages** and **Proof-of-Work (PoW) based protection** against spam and abuse.

- ❌ No cookies  
- ❌ No IP tracking  
- ❌ No email or phone required  
- ❌ No real name (unless you choose to share one)  
- ❌ No ID or passport verification

## How It Works

Upon visiting the chat, users are presented with a computational PoW challenge. Successfully solving it grants access to the system and allows the user to search for a chat partner.

The backend does **not use a database**. All data is stored in-memory only. Messages are **encrypted on the client side**, and the server **cannot decrypt** their contents. Once a message is delivered to the recipient, it is **immediately deleted from server memory** and is **not stored or logged** in any form.

## Features

- Anonymous access without registration
- Proof-of-Work challenge for each session
- Ephemeral in-memory storage
- End-to-end encrypted messages
- Lightweight and privacy-focused

## Tech Stack

- **Backend:** Go (Golang)
- **Frontend:** Svelte
- **Build Tool:** Make
- **API Schema:** OpenAPI 3

## Development Status ⚠️

This project is currently **under active development** and is **not ready for production use**. Expect bugs, missing features, and possible protocol or API changes.

## Requirements

- Go 1.24 or newer
- Node.js 22 or newer
- npm 11 or newer

## Build and Run

Clone the repository and run the following commands:

```bash
make build
cd ./build
./server -port=8080
```
After building, the `build/` directory will contain the compiled binary (`server`)  and the `frontend/` folder with static assets. The backend and frontend are bundled together: the server will serve both the **API** and the **frontend UI**.

The application will be available at [http://localhost:8080](http://localhost:8080).

## Development

To run the application in development mode, run the backend and frontend servers separately. The backend will **proxy** requests to the frontend dev server for seamless integration.

### Backend (Go)

```bash
go mod download
cd ./cmd/server
go run main.go -dev
```

### Frontend (Svelte)

```bash
cd frontend
npm install
npm run dev
```

After starting both development servers, the application will be available at [http://localhost:8080](http://localhost:8080).

## License

MIT © 2025 ingram12

