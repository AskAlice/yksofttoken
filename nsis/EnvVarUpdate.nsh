/**
 *  EnvVarUpdate.nsh
 *    : Environmental Variables: append, prepend, and remove entries
 *
 *     WARNING: If you use StrFunc.nsh header then include it before this file
 *              with all required definitions. This is to avoid conflicts
 *
 *  Usage:
 *    ${EnvVarUpdate} "ResultVar" "MYVAR" "A|P|R" "HKLM|HKCU" "PathString"
 *
 *  Credits:
 *  Version 1.0 
 *  * Cal Turney (Wikipedia)
 *  * Wikipedia - Environmental Variables: append, prepend, andடremove entries
 *
 *  Version 1.1 
 *  * LogicLib syntax improvements
 *  
 */

!ifndef ENVVARUPDATE_NSH
!define ENVVARUPDATE_NSH
 
!include "LogicLib.nsh"
!include "WinMessages.nsh"
 
!define EnvVarUpdate '!insertmacro "_EnvVarUpdate"'
!define un.EnvVarUpdate '!insertmacro "_un.EnvVarUpdate"'
 
!macro _EnvVarUpdate _RESULT _NAME _ACTION _REGLOC _PATHNAME
  Push "${_NAME}"
  Push "${_ACTION}"
  Push "${_REGLOC}"
  Push "${_PATHNAME}"
  Call EnvVarUpdate
  Pop ${_RESULT}
!macroend
 
!macro _un.EnvVarUpdate _RESULT _NAME _ACTION _REGLOC _PATHNAME
  Push "${_NAME}"
  Push "${_ACTION}"
  Push "${_REGLOC}"
  Push "${_PATHNAME}"
  Call un.EnvVarUpdate
  Pop ${_RESULT}
!macroend

;------------------------------------------
; EnvVarUpdate function body
;------------------------------------------
!macro EnvVarUpdate_Func un
Function ${un}EnvVarUpdate
 
  Exch $3 ; pathname
  Exch
  Exch $2 ; Action - A=Append, P=Prepend, R=Remove
  Exch 2
  Exch $1 ; HKCU/HKLM
  Exch 3
  Exch $0 ; env variable name
  
  Push $4 ; Current value
  Push $5 ; strlen of $4
  Push $6 ; strlen of $3
  Push $7 ; Temp var
  Push $8 ; Result
  Push $9 ; Registry key
  
  ; Set registry key based on root
  StrCmp $1 "HKCU" 0 +3
    StrCpy $9 "Environment"
    Goto +2
    StrCpy $9 "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"
  
  ; Read current value
  StrCmp $1 "HKCU" 0 +3
    ReadRegStr $4 HKCU $9 $0
    Goto +2
    ReadRegStr $4 HKLM $9 $0
  
  ; Get string lengths  
  StrLen $5 $4
  StrLen $6 $3
  
  ; Handle action type
  StrCmp $2 "R" _Remove
  StrCmp $2 "A" _Append
  StrCmp $2 "P" _Prepend
  Goto _Done
  
_Append:
  ; Check if path already exists
  ${If} $4 == ""
    StrCpy $8 $3
  ${Else}
    Push $4
    Push $3
    Call ${un}EnvVarUpdate_InPath
    Pop $7
    ${If} $7 == 1
      StrCpy $8 $4 ; Already in path, don't add again
    ${Else}
      StrCpy $8 "$4;$3"
    ${EndIf}
  ${EndIf}
  Goto _WriteValue

_Prepend:
  ${If} $4 == ""
    StrCpy $8 $3
  ${Else}
    Push $4
    Push $3
    Call ${un}EnvVarUpdate_InPath
    Pop $7
    ${If} $7 == 1
      StrCpy $8 $4
    ${Else}
      StrCpy $8 "$3;$4"
    ${EndIf}
  ${EndIf}
  Goto _WriteValue
  
_Remove:
  ${If} $4 == ""
    StrCpy $8 ""
    Goto _WriteValue
  ${EndIf}
  
  ; Remove the path from the variable
  Push $4
  Push $3
  Call ${un}EnvVarUpdate_RemovePath
  Pop $8
  Goto _WriteValue
  
_WriteValue:
  ; Write new value to registry
  StrCmp $1 "HKCU" 0 +3
    WriteRegExpandStr HKCU $9 $0 $8
    Goto +2
    WriteRegExpandStr HKLM $9 $0 $8
  
  ; Broadcast WM_SETTINGCHANGE
  SendMessage ${HWND_BROADCAST} ${WM_SETTINGCHANGE} 0 "STR:Environment" /TIMEOUT=5000
  
_Done:
  Pop $9
  Pop $8
  Pop $7
  Pop $6
  Pop $5
  Pop $4
  Pop $0
  Pop $1
  Pop $2
  Pop $3
  Push $8
  
FunctionEnd

; Helper: Check if path is already in PATH variable
Function ${un}EnvVarUpdate_InPath
  Exch $0 ; path to find
  Exch
  Exch $1 ; current PATH
  Push $2 ; temp
  Push $3 ; position
  
  StrCpy $3 0
  
_loop:
  ${If} $3 > 10000
    StrCpy $2 0
    Goto _end
  ${EndIf}
  
  ; Find next semicolon
  Push $1
  Push ";"
  Push $3
  Call ${un}EnvVarUpdate_StrLoc
  Pop $2
  
  ${If} $2 == ""
    ; No more semicolons, check rest of string
    StrCpy $2 $1 "" $3
    ${If} $2 == $0
      StrCpy $2 1
      Goto _end
    ${EndIf}
    StrCpy $2 0
    Goto _end
  ${EndIf}
  
  ; Extract segment
  IntOp $2 $2 - $3
  StrCpy $2 $1 $2 $3
  ${If} $2 == $0
    StrCpy $2 1
    Goto _end
  ${EndIf}
  
  ; Move past semicolon
  StrLen $2 $1
  IntOp $3 $3 + 1
  StrCpy $2 $1 1 $3
  ${If} $2 == ";"
    IntOp $3 $3 + 1
  ${EndIf}
  StrLen $2 $0
  IntOp $3 $3 + $2
  Goto _loop
  
_end:
  Pop $3
  Exch 2
  Pop $1
  Pop $0
  Exch $2
FunctionEnd

; Helper: Remove a path from PATH variable
Function ${un}EnvVarUpdate_RemovePath
  Exch $0 ; path to remove
  Exch
  Exch $1 ; current PATH
  Push $2 ; result
  Push $3 ; temp
  Push $4 ; segment
  Push $5 ; position
  
  StrCpy $2 ""
  StrCpy $5 0
  
_loop:
  ${If} $5 > 10000
    Goto _end
  ${EndIf}
  
  ; Find next semicolon
  Push $1
  Push ";"
  Push $5
  Call ${un}EnvVarUpdate_StrLoc
  Pop $3
  
  ${If} $3 == ""
    ; No more semicolons, check rest of string
    StrCpy $4 $1 "" $5
    ${If} $4 != $0
      ${If} $2 != ""
        StrCpy $2 "$2;$4"
      ${Else}
        StrCpy $2 $4
      ${EndIf}
    ${EndIf}
    Goto _end
  ${EndIf}
  
  ; Extract segment
  IntOp $4 $3 - $5
  StrCpy $4 $1 $4 $5
  
  ${If} $4 != $0
    ${If} $2 != ""
      StrCpy $2 "$2;$4"
    ${Else}
      StrCpy $2 $4
    ${EndIf}
  ${EndIf}
  
  ; Move to next segment
  IntOp $5 $3 + 1
  Goto _loop
  
_end:
  Pop $5
  Pop $4
  Pop $3
  Exch
  Pop $1
  Pop $0
  Exch $2
FunctionEnd

; Helper: Find position of string
Function ${un}EnvVarUpdate_StrLoc
  Exch $0 ; start position
  Exch
  Exch $1 ; search string
  Exch 2
  Exch $2 ; string to search in
  Push $3 ; temp
  Push $4 ; len search
  Push $5 ; current pos
  
  StrLen $4 $1
  StrCpy $5 $0
  
_loop:
  StrCpy $3 $2 $4 $5
  ${If} $3 == ""
    StrCpy $0 ""
    Goto _end
  ${EndIf}
  ${If} $3 == $1
    StrCpy $0 $5
    Goto _end
  ${EndIf}
  IntOp $5 $5 + 1
  Goto _loop
  
_end:
  Pop $5
  Pop $4
  Pop $3
  Pop $2
  Pop $1
  Exch $0
FunctionEnd

!macroend

; Insert the functions
!insertmacro EnvVarUpdate_Func ""
!insertmacro EnvVarUpdate_Func "un."

!endif ; ENVVARUPDATE_NSH
