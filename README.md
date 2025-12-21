# fauthy

minimal 2fa authenticator with local encrypted storage :D

## install

**arch linux:**
```bash
yay -S fauthy
```

**from source:**
```bash
git clone https://github.com/uIvPuGpT/fauthy.git
cd fauthy
go build -o fauthy
./fauthy
```

## structure
```
fauthy/
├── main.go
├── totp/
│   └── totp.go      # rfc 6238 code generation
├── tui/
│   ├── model.go     # state + data
│   ├── update.go    # input handling
│   ├── view.go      # rendering
│   └── styles.go    # colors + layout
└── storage/
    └── storage.go   # aes-256 encrypted vault
```

## flow
```
start
├── enter passphrase (or create one)
├── decrypt secrets from ~/.config/fauthy/secrets.enc
└── show totp codes

add new
├── press 'a'
├── type label (github, discord, etc)
├── paste secret key
└── auto-saves encrypted

codes refresh every 30s
```

## security
- aes-256-gcm encryption
- pbkdf2 key derivation (100k iterations)
- passphrase-protected
- file stored at `~/.config/fauthy/secrets.enc` with 0600 perms
- 5 wrong attempts → option to reset

## run
```bash
cd fauthy
go build -o fauthy
./fauthy
```

## keys
- `a` = add new secret
- `q` = quit
- `ctrl+c` = force quit
- text selection works (mouse disabled)

paste your totp secrets directly in!
