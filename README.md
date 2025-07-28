# BackupFinder ❤️

Good Day!

I truly hope everything is awesome on your side of the screen! 😊

**BackupFinder** discovers backup files on web servers by generating intelligent patterns.  
It creates thousands of potential backup file names based on your target domain.  
Perfect for penetration testing, bug bounty hunting, and security audits.

---

## Some Results 

<img width="1512" height="982" alt="Screenshot 2025-07-26 at 3 20 51 AM" src="https://github.com/user-attachments/assets/2d12d61f-7747-4319-bd3a-885f68e39fb7" />

## Hackerone Disclose Report 

For reference, here is an H1 report showing a backup file named after the domain:
[https://hackerone.com/reports/1516520](https://hackerone.com/reports/1516520)
Backup file: `https://mtn.co.rw/mtn.zip`

<img width="1512" height="982" alt="image" src="https://github.com/user-attachments/assets/2b0d12a6-6c4d-4add-bddb-c27e050d2d2e" />




## Special Thanks ❤️



Thanks to everyone who made this tool possible:


- [@infosec_au](https://x.com/infosec_au) for wordlists ([wordlists.assetnote.io](https://wordlists.assetnote.io/))
- [@coffinxp7](https://x.com/coffinxp7) for inspiration and support
- [@GodfatherOrwa](https://x.com/GodfatherOrwa) for support
- [Nishantbhagat57/backup-gen](https://github.com/Nishantbhagat57/backup-gen) — inspired by and credit to this project
- [ffuf](https://github.com/ffuf/ffuf.git) — Fast web fuzzer
- [ProjectDiscovery](https://github.com/projectdiscovery) — Security automation tools
  - [subfinder](https://github.com/projectdiscovery/subfinder.git) — Subdomain discovery
  - [httpx](https://github.com/projectdiscovery/httpx.git) — HTTP probing

Thank you to the authors, maintainers, and the security community for your contributions and inspiration!

---


## Requirements

- Go 1.21 or higher
- Internet connection for installation

---


##  Latest Install Command (Recommended)

```bash
go install github.com/MuhammadWaseem29/BackupFinder/cmd/backupfinder@v1.0.2
```

## Installation

```bash
go install github.com/MuhammadWaseem29/BackupFinder/cmd/backupfinder@v1.0.2
```

### Manual Installation

```bash
git clone https://github.com/MuhammadWaseem29/BackupFinder.git
cd BackupFinder
go build -o backupfinder ./cmd/backupfinder/
sudo mv backupfinder /usr/local/bin/
```

### Verify Installation

```bash
# Check version
backupfinder version

# Verify assets are embedded  
backupfinder health-check

# Quick test
backupfinder -u https://example.com --silent | head -3
```

---

## Usage

```
Usage:
  backupfinder [flags]

Flags:
INPUT:
  -u, -target string       target URL/domain to scan
  -l, -list string         file containing list of targets

PATTERNS:
  -w, -wordlist            use wordlist mode (comprehensive 1900+ patterns)
  -e, -extensions string   custom extensions file

OUTPUT:
  -o, -output string       file to write output to
      -je string           export to JSON file
      -json                JSON output format
      -silent              show only results
  -v, -verbose             verbose mode

PERFORMANCE:
  -c, -concurrency int     number of concurrent workers (default 10)
      -rate-limit int      rate limit for requests (default 50)
      -timeout int         request timeout in seconds (default 30)
      -retries int         maximum number of retries (default 3)

CONFIGURATION:
      -no-color            disable colored output
      -timestamp           add timestamps to output
      -stats               show statistics (default true)
      -store-resp          store responses
      -store-resp-dir      response storage directory (default "responses")

COMMANDS:
  version                  show version information
  health-check             verify installation and assets
  templates                list available pattern templates
  help                     show this help message
```

---

## Examples

<img width="1512" height="982" alt="image" src="https://github.com/user-attachments/assets/9d208507-f1b5-42d9-8bb5-859902523505" />


### Basic Usage
```bash
# Basic scan (92 extension patterns)
backupfinder -u https://admin.microsoft.com


# Comprehensive scan (1907 wordlist patterns)  
backupfinder -u https://admin.microsoft.com -w


# Multiple targets
backupfinder -l targets.txt
```

##  Dont miss use -w flag

<img width="1512" height="982" alt="image" src="https://github.com/user-attachments/assets/2c42bdc6-b318-4c34-a977-2d6613e8724d" />


### Pattern Generation
```bash
# Generate patterns silently for piping
backupfinder -u https://admin.microsoft.com --silent
backupfinder -u https://admin.microsoft.com -w --silent 

# Save to file
backupfinder -u https://admin.microsoft.com -w -o patterns.txt
```

### Output Formats
```bash
# Verbose mode with statistics
backupfinder -u https://admin.microsoft.com -w -v

# JSON export
backupfinder -u https://admin.microsoft.com -w --json -o results.json

# Silent mode (perfect for automation)
backupfinder -u https://admin.microsoft.com -w --silent
```

### FFUF Integration
```bash
backupfinder -u https://admin.microsoft.com --silent | ffuf -w /dev/stdin -u https://admin.microsoft.com/FUZZ -mc 200,403,500 -t 50
ffuf -w patterns.txt -u https://admin.microsoft.com/FUZZ -mc 200,403,500 -fc 404 -t 50 -o results.txt
```

<img width="1510" height="732" alt="image" src="https://github.com/user-attachments/assets/9a66801b-80e7-49b7-a8ae-fda677a50114" />


---

## Integration 😊

### Complete Bug Bounty Workflow
```bash
# Find subdomains
subfinder -d microsoft.com -silent > subdomains.txt

# Check live targets
cat subdomains.txt | httpx -silent > live_subdomains.txt

# Generate patterns for all subdomains
cat live_subdomains.txt | while read url; do 
    backupfinder -u "$url" --silent >> all_patterns.txt
done

# Scan with ffuf
cat live_subdomains.txt | while read url; do 
    backupfinder -u "$url" --silent | ffuf -w /dev/stdin -u "$url/FUZZ" -mc 200,403,500 -fc 404 -t 50 > results.txt
done
```

### Direct Piping
```bash
backupfinder -u https://admin.microsoft.com --silent | ffuf -w /dev/stdin -u https://admin.microsoft.com/FUZZ
backupfinder -u https://admin.microsoft.com --silent | httpx -status-code
```

### Automation Pipeline
```bash
subfinder -d microsoft.com -silent | httpx -silent | head -5 | while read url; do 
    backupfinder -u "$url" -w --silent | ffuf -w /dev/stdin -u "$url/FUZZ" -mc 200,403,500 -fc 404 -t 50
done
```

### Multiple Targets
```bash
echo -e "https://admin.microsoft.com\nhttps://api.microsoft.com" | while read url; do 
    backupfinder -u "$url" --silent | ffuf -w /dev/stdin -u "$url/FUZZ" -mc 200,403,500 -t 30
done
```

---

<img width="1510" height="982" alt="image" src="https://github.com/user-attachments/assets/ee07492e-20ad-4c4a-8690-9c19de9226be" />


---

## Features

- 9000+ backup patterns in wordlist mode
- Smart subdomain handling (admin.example.com → admin.zip, admin-example.sql)  
- Professional JSON export for automation
- Real-time statistics with performance metrics
- Concurrent processing for fast pattern generation
- Custom wordlists support
- Silent mode for integration with other tools

<img width="1510" height="732" alt="image" src="https://github.com/user-attachments/assets/ac02078f-e228-4a11-ba33-5a10f9e08aa5" />



---

## Pattern Types

### Extension Mode (Default - 92 patterns)
Common backup extensions: `.bak`, `.backup`, `.old`, `.sql`, `.zip`, etc.

### Wordlist Mode (-w - 1900+ patterns)  
Comprehensive patterns for database dumps, configuration backups, archive variants

### Custom Extensions (-e)
Use your own pattern file (one pattern per line, supports `#` comments)




---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---


### Getting Help

-  **Issues**: [GitHub Issues](https://github.com/MuhammadWaseem29/BackupFinder/issues)
-  **Documentation**: This README
-  **Discussions**: [GitHub Discussions](https://github.com/MuhammadWaseem29/BackupFinder/discussions)


## Author

**MuhammadWaseem**  
GitHub: [@MuhammadWaseem29](https://github.com/MuhammadWaseem29)  
Tool: BackupFinder v1.0.0

---

 Thank you for using BackupFinder! We appreciate your support! ❤️

---

---



May you be well on your side of the screen :)

---



## Star History


<a href="https://star-history.com/#MuhammadWaseem29/BackupFinder&Date">
  <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=MuhammadWaseem29/BackupFinder&type=Date&theme=light" style="background:white;" />
</a>

---

