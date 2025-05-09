# smtp2chat

**A lightweight SMTP server that relays email messages to Telegram or Microsoft Teams.**

---

## 🚀 Features

* Converts incoming SMTP emails into chat messages
* Supports **Telegram bots** and **Microsoft Teams webhooks**
* Simple binary, minimal configuration
* Optional verbose logging

---

## 💾 Installation

```bash
go build -o smtp2chat main.go
```

---

## ⚙️ Usage

```bash
./smtp2chat [OPTIONS]
```

### Available Options:

| Flag            | Description                             |
| --------------- | --------------------------------------- |
| `-s, --service` | Message service: `telegram` or `teams`  |
| `-t, --token`   | Telegram bot token                      |
| `-c, --chatid`  | Telegram chat ID                        |
| `-w, --webhook` | Teams incoming webhook URL              |
| `-b, --bind`    | Bind address for SMTP (default `:2525`) |
| `-v, --verbose` | Enable verbose logging                  |

### Example (Telegram):

```bash
./smtp2chat \
  --service telegram \
  --token 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11 \
  --chatid -1001234567890
```

### Example (Teams):

```bash
./smtp2chat \
  --service teams \
  --webhook https://outlook.office.com/webhook/...
```

---

## 🔧 How It Works

1. Listens for SMTP connections on a configurable port
2. Parses incoming email
3. Relays content to your chosen service:

   * **Telegram** via bot API
   * **Teams** via webhook

---

## ✉️ Email Format Support

* Extracts and forwards:

  * Subject
  * From / To addresses
  * Body (plain text preferred)

---

## 🛡️ Security Notes

* Ensure your SMTP server is only exposed to trusted sources (e.g., use firewalls or network policies).
* Avoid exposing Telegram bot tokens or Teams webhooks.

---

## 🚜 Contributing

PRs are welcome! Please include a test plan or validation example.

---

## 🚀 License

MIT License

---

## 🔗 Related Projects

* [Telegram Bot API](https://core.telegram.org/bots/api)
* [Microsoft Teams Webhooks](https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook)
