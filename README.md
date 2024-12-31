# 🤖 Diamond Hand Discord Bot
Diamond Hand is a feature-rich Discord bot built with Go, offering AI responses and server utilities.

## ✨ Features

### 🧠 AI-Powered Responses
- **Command:** `/ai <query>`
  - 🤔 Generates intelligent responses using Gemini API
  - 💡 Example: `/ai How is the weather today?`

### 🛠️ Server Utilities

#### 📡 Network & Performance
- **Command:** `/dh latence`
  - 🟢 Shows bot latency with color indicators
    - ✅ Green: Excellent connection
    - ⚠️ Yellow: Poor connection
    - ❌ Red: Critical latency

#### 🔍 Server Navigation
- **Command:** `/dh ls`
  - 📂 Lists all server channels and categories
  - 🗂️ Organized hierarchical view

#### 📍 Location Tracking
- **Command:** `/dh pwd`
  - 🌐 Shows your current server "path":
  ```
  ./<Username>/<Server Name>/<Channel Name>
  ```

#### 🔊 Voice Commands
- **Command:** `/dh dwi [text]`
  - 🗣️ Text-to-Speech functionality
  - 🎯 Clear voice output

#### 🎵 Music Features
- **Commands:**
  - 🎵 `/dh play [url]` - Play music
  - 📑 `/dh playlist [url]` - Load playlist
  - ⏭️ `/dh skip` - Skip track
  - ⏹️ `/dh stop` - Stop playback
  - 📋 `/dh queue` - View playlist

#### 🎮 Gaming Features
- **Valorant Server Status:**
  - 🎯 `/dh valo-ping` - Check server latency
  - 🌍 Shows ping to EU servers

#### 🔒 Security Scanner
- **Command:** `/dh scan [url]`
  - 🔍 Website security analysis
  - ℹ️ Technology detection
  - 🌐 Server information

## 🚀 Setup

### Prerequisites
- Go 1.22 or higher
- Discord Bot Token
- Gemini API Key

### Installation

1. 📥 Clone the repository:
```bash
git clone https://github.com/01000001BDO/bot-discord.git && cd bot-discord
```

2. ⚙️ Create a `.env` file:
```bash
TOKEN=your_discord_bot_token
GEMINI_API_KEY=your_gemini_api_key
```

3. 🏃‍♂️ Run the bot:
```bash
go run .
```

### 🐳 Docker Deployment
```bash
# Build the image
docker build -t discord-bot .

# Run the container
docker run -d discord-bot
```

## 🤝 Contributing
Feel free to contribute! Open issues and pull requests are welcome.

## 📝 License
This project is open source and available under the [MIT License](LICENSE).

## 💖 Support
If you enjoy using Diamond Hand, consider giving it a star on GitHub!
