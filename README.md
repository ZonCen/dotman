# dotman

A simple, CLI-based **dotfiles manager** written in Go.  
`dotman` helps you keep your configuration files (dotfiles) organized in a central repository, while automatically creating symlinks so your system keeps working as expected.  

---

## ✨ Features (MVP)
- **Add**: move a file into the dotfiles repo and replace it with a symlink.  
- **List**: show all files currently tracked in the repo.  
- **Remove**: restore a file from the repo by removing the symlink and moving the file back.  

---

## 📦 Installation
Clone the repository and build from source:

```bash
git clone https://github.com/ZonCen/dotman.git
cd dotman
go build -o dotman
```

Optionally, move it into your `$PATH`:

```bash
mv dotman ~/.local/bin/
```

Note: You also need to create the dotfiles repository folder you want to use (example: /Users/yourname/dotfiles) and make sure it is initialized as a Git repository if you plan to use the sync command.


---

## ⚡ Usage

### 1. Initialize Config
The first time you run `dotman`, make sure you have a config file at:

```
~/dotfiles/.dotconfig
```

Example `.dotconfig`:

```yaml
repo_path: /Users/yourname/dotfiles
```

This tells `dotman` where to store your dotfiles.  

---

### 2. Add a File
```bash
dotman add ~/.zshrc
```
- Moves `~/.zshrc` → `~/dotfiles/.zshrc`.  
- Creates a symlink: `~/.zshrc → ~/dotfiles/.zshrc`.  

---

### 3. List Files
```bash
dotman list
```
Shows all tracked files in the repo.  

---

### 4. Remove a File
```bash
dotman remove ~/.zshrc
```
- Removes the symlink `~/.zshrc`.  
- Moves `~/dotfiles/.zshrc` back to `~/.zshrc`.  

---

### 5. Sync your repo (experimental)
- Stages new/modified files
- Commit them with the message `dotman sync`.
- Pushes your configured Github repository
- Pulls changes from remote (fast-forward only)

## 🔄 Full Example Workflow

Here’s a typical session:

```bash
# Check your config (should exist at ~/dotfiles/.dotconfig)
cat ~/dotfiles/.dotconfig
# repo_path: /Users/yourname/dotfiles

# Add your zsh config
dotman add ~/.zshrc
# Successfully added .zshrc to repository

# List tracked files
dotman list
# .zshrc

# Sync with GitHub
dotman sync
# dotman sync -> git commit + push

# Remove the file and restore it
dotman remove ~/.zshrc
# Successfully removed file from path
```

After this workflow:  
- `~/.zshrc` has been restored to its original location.  
- The repo file `~/dotfiles/.zshrc` is now gone (moved back).
- Your repo can be synced with GitHub.

---

## 🛠 Roadmap
- [x] Add support for initializing config automatically if missing.  
- [x] Add `sync` command to push/pull dotfiles via Git.
- [ ] Improvements to `sync` command such as dry-run functionality 
- [ ] Add possibility to automatically download your dotfiles folder and make it ready for `sync`
- [ ] Add flags for overwrite/force when adding files.  
- [ ] Add verbose mode (show actions instead of silent ops).  

---

## 📜 License
This project is licensed under the MIT License.  
