# 🏴‍☠️ CHEST

#### Don't lose your jewels, keep them in a chest!

A minimalist CLI secret manager to store credentials, keys, and configurations — and use them directly from the terminal.

---

## 📦 Installation

### Requirements

- `xclip`, `xsel`, or `wl-copy` for clipboard support

```bash
# Install clipboard tool (pick one)
sudo apt install xclip
```

### Download binary and install

**Linux (amd64)**
```bash
wget -qO chest https://github.com/tommyliverani/chest/releases/latest/download/chest-linux-amd64 && chmod +x chest && sudo mv chest /usr/local/bin/chest
```

**Windows (amd64)**

Download [`chest-windows-amd64.exe`](https://github.com/tommyliverani/chest/releases/latest/download/chest-windows-amd64.exe) from the latest release, rename it to `chest.exe` and move it to a folder included in your `PATH` (e.g. `C:\Windows\System32` or any custom directory).

### Build from source

Requires Go 1.25+.

```bash
git clone https://github.com/tommyliverani/chest.git
cd chest
make build
sudo mv chest /usr/local/bin/chest
```


---

## 🏗️ Concepts

| Concept | Description |
|---------|-------------|
| **Jewel** | Evrything that can contains secrets data like simple credentisals, kubeconfigs, login command, ecc... |
| **Chest** | A container of jewels |


Chests are stored as encrypted JSON files. When you **open** a chest, its decrypted contents are temporarily saved in `/run/user/<uid>/` and are available to all jewel commands until you **close** it.

---

## 🗝️ Jewel kinds

| Emoji | Kind | Command | Description |
|-------|------|---------|-------------|
| 💎 | `up` | `chest up ...` | Username + password pair |
| 🗝️ | `key` | `chest key ...` | Secret key or password |
| ☁️ | `aws` | `chest aws ...` | AWS Access Key ID + Secret Access Key |
| 🧭 | `ssh` | `chest ssh ...` | SSH credentials (username, host, password) |
| ☸️ | `kube` | `chest kube ...` | Kubeconfig — merges into `~/.kube/config` |
| ♦️ | `oc` | `chest oc ...` | OpenShift API URL + token for `oc login` |
| 🐳 | `docker` | `chest docker ...` | Docker registry credentials for `docker login` |
| 🌐 | `curl` | `chest curl ...` | Saved curl request (URL, method, options, auth) |
| 📜 | `file` | `chest file ...` | Secret file stored as base64 |

---

## 🚀 Usage

### Chest commands

```bash
chest create              # create a new chest (prompts for name, type, password)
chest ls                  # list all chests
chest open <name>         # open a chest (prompts for password)
chest close               # close all open chests
chest edit <name>         # edit name or description of a chest
chest rm <name>           # delete a chest
chest jewels              # list all jewels across open chests (alias: chest js)
```

### Jewel commands

Each jewel kind supports the same set of sub-commands:

```bash
chest <kind> add              # add a new jewel (interactive prompt)
chest <kind> ls               # list all jewels of that kind in open chests
chest <kind> edit [name]      # edit a jewel's fields
chest <kind> rm [name]        # remove a jewel
chest <kind> print [name]     # print the jewel's sensitive value (with confirmation)
chest <kind> copy [name]      # copy the relevant value/command to the clipboard
chest <kind> [name]           # use the jewel (connect, login, apply, etc.)
chest <kind> help             # show help for that jewel kind
```

---

## 🔍 Jewel details

### 💎 `up` — Username & Password
Stores a username and password pair.
- **copy** → copies the password to the clipboard
- **use** → copies the password to the clipboard

### 🗝️ `key` — Secret Key
Stores a secret key or password.
- **copy** → copies the key value to the clipboard
- **use** → copies the key value to the clipboard

### ☁️ `aws` — AWS Credentials
Stores AWS Access Key ID and Secret Access Key.
- **print** → shows the `export AWS_*` commands with masked secret
- **copy** → copies the full `export AWS_ACCESS_KEY_ID=... && export AWS_SECRET_ACCESS_KEY=...` to the clipboard
- **use** → writes credentials to `~/.aws/credentials` and verifies with `sts get-caller-identity`

### 🧭 `ssh` — SSH
Stores SSH credentials (username, host, password).
- **print** → shows username, host and password (with confirmation)
- **copy** → copies `ssh user@host` to the clipboard
- **use** → opens an interactive SSH session

### ☸️ `kube` — Kubeconfig
Stores a full kubeconfig (paste via stdin, `Ctrl+D` to finish).
- **print** → shows the raw kubeconfig YAML (with confirmation)
- **copy** → copies the kubeconfig YAML to the clipboard
- **use** → merges the kubeconfig into `~/.kube/config` and sets the current context

### ♦️ `oc` — OpenShift
Stores an OpenShift API URL and token.
- **print** → shows `oc login <url> --token=***`
- **copy** → copies the full `oc login <url> --token=<token>` command to the clipboard
- **use** → runs `oc login` with the stored credentials

### 🐳 `docker` — Docker Registry
Stores Docker registry credentials.
- **print** → shows `docker login` command with masked password
- **copy** → copies the full `docker login` command to the clipboard
- **use** → runs `docker login` with the stored credentials

### 🌐 `curl` — Curl Request
Stores a full curl request (URL, HTTP method, options, optional basic auth).
- **print** → shows the curl command with masked password
- **copy** → copies the full curl command to the clipboard
- **use** → executes the curl request

### 📜 `file` — Secret File
Stores a secret file as base64 (paste content via stdin, `Ctrl+D` to finish).
- **print** → shows the file content (with confirmation)
- **copy** → saves the file to a path chosen interactively

---

## 🔒 Security

- Chests are encrypted with **AES-256-GCM**.
- The encryption key is derived from the chest password + a random salt using **SHA-256**.
- Open chest sessions are stored in `/run/user/<uid>/` (tmpfs, not persisted to disk across reboots).
- Secrets are never printed without explicit confirmation.

---

## 🛠️ Development

```bash
make test      # run tests with race detector
make quality   # run golangci-lint + govulncheck
make build     # compile the binary
make all       # quality + test + build
```

---

## 📋 Quick example

```bash
# Create a chest
chest create
# > name: mychest | type: aes | password: ****

# Open it
chest open mychest

# Add an SSH credential
chest ssh add
# > name: prod-server | description: production box | username: admin | host: 192.168.1.10 | password: ****

# Connect
chest ssh prod-server

# Copy the ssh command to clipboard
chest ssh copy prod-server

# Close when done
chest close
```