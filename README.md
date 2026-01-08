# QuickDo
A Terminal ToDo List app to practice Golang coding

## How to Run it

### 1. Download
Go to the [Releases](https://github.com/fapinter/QuickDo/releases) page and download the version for your system:
- **Windows**: `quickdo.exe`
- **Linux**: `quickdo`
- **Mac**: `quickdo-macos-arm`

### 2. Setup (Run from any folder)

####  Windows
1. Move `quickdo.exe` to a permanent folder (e.g., `C:\Tools\`).
2. Search for **"Edit the system environment variables"** in the Start Menu.
3. Click **Environment Variables** > Find **Path** in 'User variables' > Click **Edit**.
4. Click **New** and paste the folder path `C:\Tools\`.
5. Open a new Terminal and just type `quickdo`.

#### 🍎 macOS & 🐧 Linux
1. Open your terminal in the folder where you downloaded the file.
2. Move it to your system's bin folder and make it executable:
   ```bash
   # Move and rename to just 'quickdo'
   sudo mv quickdo* /usr/local/bin/quickdo
   
   # Give permission to run
   sudo chmod +x /usr/local/bin/quickdo

   sudo chown $USER /usr/local/bin/quickdo
