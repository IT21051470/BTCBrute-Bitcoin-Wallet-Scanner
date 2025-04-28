🚀 BTCBrute – Bitcoin Address Scanner
BTCBrute is a high-performance Bitcoin private key and address scanner focused on single-signature (P2PKH) Bitcoin addresses.

It randomly generates Bitcoin private keys, derives corresponding addresses, and checks them against a known list of Bitcoin addresses (richlist, database, leaked addresses, etc.).
Matches are logged and sent via Telegram notification.

Built for educational and ethical research purposes only.

✨ Features
🔍 Random private key generation

📬 Bitcoin address derivation (P2PKH format, addresses starting with 1)

📑 Match checking against large address datasets

📈 Real-time progress display (speed, total scanned, matches)

📲 Telegram Bot notifications on match found

🛡 Efficient multithreaded scanning

🧠 Memory-efficient (streamed processing)

⚙️ Requirements
Go 1.18+ installed (https://go.dev/dl/)

Basic understanding of command-line usage

📦 Installation
Clone or download this repository:

git clone https://github.com/yourusername/btcbrute.git
cd btcbrute

Initialize Go module:

go mod init btcbrute

Install dependencies:

go get github.com/btcsuite/btcd/btcec/v2
go get github.com/btcsuite/btcutil/base58
go get golang.org/x/crypto/ripemd160

Build the executable:

go build -o btcbrute btcbrute.go

✅ Now you will have btcbrute.exe ready.

🚀 Usage
Basic command to start scanning:

btcbrute.exe -threads 8 -output matches.txt -data btc-data-file.txt


Argument	Description
-threads	Number of concurrent threads (default: 4)
-output	File to save matches (default: matches.txt)
-data	Path to your Bitcoin address database (one address per line)
Example:

btcbrute.exe -threads 16 -output hits.txt -data my-richlist.txt

🛎 Telegram Notification Setup
To receive a Telegram message when a match is found:

Create a Telegram Bot using @BotFather

Get your Bot Token

Get your Telegram Chat ID (use @userinfobot)

Edit these constants in btcbrute.go:

go
Copy
Edit
const (
    botToken = "YOUR_BOT_TOKEN"
    chatID   = "YOUR_CHAT_ID"
)
✅ Now you’ll receive an instant Telegram alert if a private key matches an address!

📈 Live Progress
During scanning, you will see real-time stats:

[BTCBrute] - Scanning | Checked: 5,482,129 | Matches: 0 | Speed: 85,322/s


Stat	Meaning
Checked	Total private keys generated and tested
Matches	Successful address matches found
Speed	Current speed (keys/second)
⚠️ Disclaimer
BTCBrute is designed for educational and ethical research purposes only.
Unauthorized access to cryptocurrency wallets or assets you do not own is illegal.
The author is not responsible for any misuse or damages caused by this software.

🔥 License
This project is open-source and released under the MIT License.

📬 Contact
Questions? Contributions?
Open an issue or pull request!

✨ Final Notes
📚 Use it to learn about Bitcoin key generation and address security.

🔐 Never use weak brainwallets or predictable private keys.

🧠 Real scanning success needs smart targeting, not blind luck.

⭐ Support
If you find this tool useful, please consider giving a ⭐ star to support future updates and new tools!
