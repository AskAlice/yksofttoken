# Homebrew Cask formula for YKSoft Token
# To install: brew install --cask yksoft
#
# Note: This is a template. The actual cask should be submitted to homebrew-cask
# or hosted in a custom tap.

cask "yksoft" do
  version "1.0.0"
  
  # Determine architecture
  on_intel do
    sha256 "PLACEHOLDER_SHA256_AMD64"
    url "https://github.com/arr2036/yksofttoken/releases/download/v#{version}/yksoft-darwin-amd64.zip"
  end
  on_arm do
    sha256 "PLACEHOLDER_SHA256_ARM64"
    url "https://github.com/arr2036/yksofttoken/releases/download/v#{version}/yksoft-darwin-arm64.zip"
  end

  name "YKSoft Token"
  desc "Software Yubikey token emulator for HOTP One Time Passcodes"
  homepage "https://github.com/arr2036/yksofttoken"

  # Application
  app "yksoft.app", target: "YKSoft Token.app"

  # Uninstall configuration
  uninstall quit: "org.freeradius.yksoft"

  # Cleanup token data on uninstall (optional - user must confirm)
  zap trash: [
    "~/.yksoft",
  ]

  caveats <<~EOS
    YKSoft Token stores token data in ~/.yksoft/
    
    This is a software Yubikey emulator and is NOT as secure as a
    physical Yubikey. Use only for testing or M2M VPN connections.
    
    To generate a new token, launch the app and click "New".
    The registration information will be displayed for you to
    register with your authentication server.
  EOS
end
