# yksoft

[![CI build status][BuildStatus]][BuildStatusLink]
[![Cross-Platform Build][CrossBuildStatus]][CrossBuildLink]

## Introduction

Sometimes it's useful to emulate a physical Yubikey token in software, examples of this include:

- Testing, where you don't want to purchase a Yubikey.
- M2M VPN connections where the other party is unable to make an exception to
  their 2FA policy.
- Particularly difficult customers who want to treat everyone in your organisation as
  their employees, and want to issue each of them a hardware token to connect into
  their VPN service.

**NOTE**

yksoft is not intended to be a replacement for a Yubikey in situations that require
high security, this utility is not as secure as a physical Yubikey.

## Features

- **Cross-Platform**: Runs on Windows, macOS (Intel & Apple Silicon), and Linux (x86_64 & ARM64)
- **Modern GUI**: Built with [Fyne](https://fyne.io/) toolkit for a native look and feel
- **Token Management**: Create, manage, and delete multiple software tokens
- **Clipboard Support**: One-click copy of OTPs and registration information
- **Persistent Storage**: Token data is stored securely in `~/.yksoft/`
- **Compatible**: Generates OTPs compatible with standard Yubikey validators

## Screenshots

The application provides a clean, modern interface for managing software Yubikey tokens.

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/arr2036/yksofttoken/releases) page.

#### Windows
- Download `YKSoftToken-Setup.exe` for the installer
- Or download `yksoft-windows-amd64.zip` / `yksoft-windows-arm64.zip` for portable version

#### macOS
Using Homebrew (recommended):
```bash
brew install --cask yksoft
```

Or download the ZIP file for your architecture:
- Intel Mac: `yksoft-darwin-amd64.zip`
- Apple Silicon: `yksoft-darwin-arm64.zip`

#### Linux
Debian/Ubuntu:
```bash
sudo dpkg -i yksofttoken_1.0.0_amd64.deb
# or for ARM64
sudo dpkg -i yksofttoken_1.0.0_arm64.deb
```

Or download the tarball for your architecture.

### Building from Source

#### Prerequisites
- Go 1.21 or later
- Fyne dependencies (see [Fyne Getting Started](https://developer.fyne.io/started/))

#### Build

```bash
# Clone the repository
git clone https://github.com/arr2036/yksofttoken.git
cd yksofttoken

# Build
go build -o yksoft ./cmd/yksoft

# Run
./yksoft
```

#### Cross-Compilation with fyne-cross

```bash
# Install fyne-cross
go install github.com/fyne-io/fyne-cross@latest

# Build for all platforms
fyne-cross linux -arch amd64,arm64 ./cmd/yksoft
fyne-cross darwin -arch amd64,arm64 ./cmd/yksoft
fyne-cross windows -arch amd64,arm64 ./cmd/yksoft
```

## Usage

### GUI Application

1. Launch the application
2. Click "New" to create a new token
3. The registration information (public ID, private ID, AES key) will be displayed
4. Register these values with your authentication server
5. Click "Generate OTP" to create a one-time password
6. Click "Copy" to copy the OTP to clipboard

### Token Storage

Token data is stored in `~/.yksoft/` (or `%USERPROFILE%\.yksoft\` on Windows).

Each token is stored as a plaintext file with the following format:
```
public_id: <modhex>
private_id: <hex>
aes_key: <hex>
counter: <number>
session: <number>
created: <timestamp>
lastuse: <timestamp>
ponrand: <number>
```

**Security Note**: The token files are not encrypted. Ensure appropriate file permissions
are set (the application creates files with mode 0600).

## Registration

When you create a new token or click "Copy Registration Info", you'll get a CSV string:
```
<public_id_modhex>, <private_id_hex>, <aes_key_hex>
```

Use these values to register the token with your authentication server:
- **Public ID (modhex)**: The identifier prepended to each OTP
- **Private ID (hex)**: The secret identifier validated by the server
- **AES Key (hex)**: The encryption key (16 bytes / 32 hex characters)

## Technical Details

### OTP Format

Each OTP consists of:
- 12 modhex characters: Public ID
- 32 modhex characters: Encrypted token block

Total: 44 characters

### Encrypted Token Block Contents

| Field     | Size    | Description                    |
|-----------|---------|--------------------------------|
| uid       | 6 bytes | Private ID                     |
| counter   | 2 bytes | Usage counter                  |
| timestamp | 3 bytes | 8Hz timer value                |
| session   | 1 byte  | Session use counter            |
| random    | 2 bytes | Random value                   |
| crc       | 2 bytes | CRC16 checksum                 |

### Time Simulation

A hardware Yubikey has an 8Hz timer. This software emulates it using:
```
timestamp = ((current_time - created_time) * 8 + ponrand) % 0xFFFFFF
```

## Development

### Project Structure

```
.
├── cmd/yksoft/          # Main application
├── internal/
│   ├── yubikey/         # Yubikey encoding/crypto functions
│   └── token/           # Token management
├── assets/              # Application icons
├── nsis/                # Windows installer script
├── homebrew/            # macOS Homebrew cask
├── debian/              # Debian packaging
└── .github/workflows/   # CI/CD pipelines
```

### Running Tests

```bash
go test ./...
```

### CI/CD

The project uses GitHub Actions for:
- Building binaries for all platforms using `fyne-cross`
- Creating Windows NSIS installers
- Building Debian packages
- Creating GitHub releases

## License

BSD 2-Clause License. See [LICENSE](LICENSE) for details.

## Credits

- Original C implementation by Arran Cudbard-Bell
- Go port using [Fyne](https://fyne.io/) GUI toolkit
- Yubikey protocol compatible with [libyubikey](https://github.com/Yubico/yubico-c)

## Contributing

Contributions are welcome! Please open an issue or pull request.

[BuildStatus]: https://github.com/arr2036/yksofttoken/actions/workflows/ci-linux.yml/badge.svg "CI status"
[BuildStatusLink]: https://github.com/arr2036/yksofttoken/actions/workflows/ci-linux.yml
[CrossBuildStatus]: https://github.com/arr2036/yksofttoken/actions/workflows/build.yml/badge.svg "Cross-Platform Build"
[CrossBuildLink]: https://github.com/arr2036/yksofttoken/actions/workflows/build.yml
