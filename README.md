# Anonymous Chat

An anonymous chat application with **end-to-end encrypted messages** and **Proof-of-Work (PoW) based protection** against spam and abuse.

- ❌ No cookies  
- ❌ No IP tracking  
- ❌ No email or phone required  

## How It Works

Upon visiting the chat, users are presented with a computational PoW challenge. Successfully solving it grants access to the system and allows the user to search for a chat partner.

The backend does **not use a database**. All data is stored in-memory, messages are encrypted, and immediately deleted after delivery to the recipient.

## Features

- Anonymous access without registration
- Proof-of-Work challenge for each session
- Ephemeral in-memory storage
- End-to-end encrypted messages
- Lightweight and privacy-focused

## Requirements

- Go 1.24 or newer
- Node.js and npm

## Installation

Clone the repository and run the following commands:

```bash
make build
cd ./build
./server -port=8080
```

The application will be available at http://localhost:8080.

