📚 SapPress-DL (sap-press downloader)

Download True EPUB Books from [SAP PRESS](https://www.sap-press.com/)
Built with ❤️ in Golang by [@Caliginous_0](https://t.me/Caliginous_0)

---

✨ Features

- 📥 Download real EPUB books from sap-press.com
- ⚡ Multi-threaded downloading for speed (default: 16 threads)
- 🎯 Minimalistic and fast
- 🧩 Easy-to-use CLI interface

---

🛠️ Installation

Clone the repository and run the project using Go:

```bash
git clone https://github.com/mohamedhani1/SAP-Press.git
cd sap-press-dl
go run main.go
````

Or compile it:

```bash
go build -o SapPress-DL main.go
./SapPress-DL
```

---

🚀 Usage

General Help

```bash
SapPress-DL --help
```

Output:

```
NAME:
   SapPress-DL - Download True EPUB Books
                 Created by t.me/@Caliginous_0

USAGE:
   SapPress-DL [global options] command [command options]

COMMANDS:
   download  Download a new book by ID
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

---

📘 Download a Book

```bash
SapPress-DL download --bookid <BOOK_ID>
```

Optional: customize number of threads (default is 16)

```bash
SapPress-DL download --bookid <BOOK_ID> --threads 32
```

Example:

```bash
SapPress-DL download --bookid 12345 --threads 8
```

Help for `download` command

```bash
SapPress-DL download --help
```

```
NAME:
   SapPress-DL download - Download a new book by ID

USAGE:
   SapPress-DL download [command options]

OPTIONS:
   --bookid value   ID of the book to download
   --threads value  Number of concurrent threads (default: 16)
   --help, -h       show help
```

---

🧑‍💻 Author

* Telegram: [@Caliginous\_0](https://t.me/Caliginous_0)

---

⭐️ Star This Repo

If you find this project useful, please consider giving it a ⭐️!

```

Let me know if you'd like a matching `LICENSE` file or a Go-based `main.go` skeleton.
```
