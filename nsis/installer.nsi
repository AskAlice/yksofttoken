; YKSoft Token NSIS Installer Script
; Modern UI 2 installer

!include "MUI2.nsh"
!include "FileFunc.nsh"
!include "EnvVarUpdate.nsh"

; General settings
Name "YKSoft Token"
OutFile "YKSoftToken-Setup.exe"
InstallDir "$PROGRAMFILES64\yksoft"
InstallDirRegKey HKLM "Software\YKSoft Token" "InstallDir"
RequestExecutionLevel admin

; Version information
!define VERSION "1.0.0"
!define PUBLISHER "Arran Cudbard-Bell"
!define WEBSITE "https://github.com/arr2036/yksofttoken"

VIProductVersion "${VERSION}.0"
VIAddVersionKey "ProductName" "YKSoft Token"
VIAddVersionKey "CompanyName" "${PUBLISHER}"
VIAddVersionKey "LegalCopyright" "Copyright (c) 2022-2024 ${PUBLISHER}, modifications (c) 2026 Alice Knag"
VIAddVersionKey "FileDescription" "Yubikey Software Token Emulator"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "ProductVersion" "${VERSION}"

; Modern UI settings
!define MUI_ABORTWARNING
!define MUI_ICON "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"
!define MUI_UNICON "${NSISDIR}\Contrib\Graphics\Icons\modern-uninstall.ico"

; Welcome page
!define MUI_WELCOMEPAGE_TITLE "Welcome to YKSoft Token Setup"
!define MUI_WELCOMEPAGE_TEXT "This wizard will guide you through the installation of YKSoft Token.$\r$\n$\r$\nYKSoft Token is a software Yubikey emulator that generates HOTP One Time Passcodes.$\r$\n$\r$\nClick Next to continue."

; Finish page
!define MUI_FINISHPAGE_RUN "$INSTDIR\yksoft.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Launch YKSoft Token"

; Pages
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "..\LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; Language
!insertmacro MUI_LANGUAGE "English"

; Installation section
Section "Install" SecInstall
    SetOutPath "$INSTDIR"
    
    ; Install files
    File "yksoft.exe"
    
    ; Create uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"
    
    ; Create Start Menu shortcuts
    CreateDirectory "$SMPROGRAMS\YKSoft Token"
    CreateShortcut "$SMPROGRAMS\YKSoft Token\YKSoft Token.lnk" "$INSTDIR\yksoft.exe"
    CreateShortcut "$SMPROGRAMS\YKSoft Token\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
    
    ; Create Desktop shortcut
    CreateShortcut "$DESKTOP\YKSoft Token.lnk" "$INSTDIR\yksoft.exe"
    
    ; Write registry keys for Add/Remove Programs
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "DisplayName" "YKSoft Token"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "UninstallString" "$\"$INSTDIR\Uninstall.exe$\""
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "QuietUninstallString" "$\"$INSTDIR\Uninstall.exe$\" /S"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "InstallLocation" "$\"$INSTDIR$\""
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "DisplayIcon" "$\"$INSTDIR\yksoft.exe$\""
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "Publisher" "${PUBLISHER}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "HelpLink" "${WEBSITE}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "URLInfoAbout" "${WEBSITE}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "DisplayVersion" "${VERSION}"
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "NoModify" 1
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "NoRepair" 1
    
    ; Get installed size
    ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
    IntFmt $0 "0x%08X" $0
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token" \
                     "EstimatedSize" "$0"
    
    ; Save install dir
    WriteRegStr HKLM "Software\YKSoft Token" "InstallDir" "$INSTDIR"
    
    ; Add to PATH
    ${EnvVarUpdate} $0 "PATH" "A" "HKLM" "$INSTDIR"
SectionEnd

; Uninstaller section
Section "Uninstall"
    ; Remove from PATH
    ${un.EnvVarUpdate} $0 "PATH" "R" "HKLM" "$INSTDIR"
    
    ; Remove files
    Delete "$INSTDIR\yksoft.exe"
    Delete "$INSTDIR\Uninstall.exe"
    
    ; Remove shortcuts
    Delete "$SMPROGRAMS\YKSoft Token\YKSoft Token.lnk"
    Delete "$SMPROGRAMS\YKSoft Token\Uninstall.lnk"
    RMDir "$SMPROGRAMS\YKSoft Token"
    Delete "$DESKTOP\YKSoft Token.lnk"
    
    ; Remove installation directory
    RMDir "$INSTDIR"
    
    ; Remove registry keys
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\YKSoft Token"
    DeleteRegKey HKLM "Software\YKSoft Token"
SectionEnd
