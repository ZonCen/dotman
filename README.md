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

# Remove the file and restore it
dotman remove ~/.zshrc
# Successfully removed file from path
```

After this workflow:  
- `~/.zshrc` has been restored to its original location.  
- The repo file `~/dotfiles/.zshrc` is now gone (moved back).  

---

## 🛠 Roadmap
- [x] Add support for initializing config automatically if missing.  
- [ ] Add `sync` command to push/pull dotfiles via Git.  
- [ ] Add flags for overwrite/force when adding files.  
- [ ] Add verbose mode (show actions instead of silent ops).  

---

## 📜 License
This project is licensed under the MIT License.  
