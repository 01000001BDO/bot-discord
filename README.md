# Diamond Hand Discord Bot

Diamond Hand  is a Discord bot built with Go, featuring AI responses and server utilities.

## Features

### AI-Powered Responses
- **Command:** `/ai <query>`
  - Generates AI-driven responses using the Gemini API.
  - Example: `/ai How is the weather today?`

### Server Utilities
1. **Latency Check**
   - **Command:** `/dh latence`
   - Displays the bot's latency in milliseconds.
   - Latency is color-coded:
     - **Green:** Low latency (good connection).
     - **Yellow:** High latency (poor connection).

2. **Channel Listing**
   - **Command:** `/dh ls`
   - Lists all server channels categorized by their types.

3. **Current Path**
   - **Command:** `/dh pwd`
   - Provides the user’s "path" in the server hierarchy:
     ```
     ./<Username>/<Server Name>/<Channel Name>
     ```

4. **Memes (Planned)**
   - **Command:** `/dh memes`
   - Placeholder for future meme-sharing functionality.

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/01000001BDO/bot-discord.git && cd bot-discord
   ```
2. Create a .env file with the following:
   ```bash
   TOKEN=your_discord_bot_token
   GEMINI_API_KEY=your_gemini_api_key
   ```
3. Clone the repository:
   ```bash
   go run .
   ```
