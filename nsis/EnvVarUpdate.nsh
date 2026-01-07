/**
 *  EnvVarUpdate.nsh
 *    : Environmental Variables: append, prepend, and remove entries
 *
 *  Features:
 *    - Supports long PATH strings (uses NSIS_MAX_STRLEN)
 *    - Backs up PATH before modifications
 *    - Properly handles REG_EXPAND_SZ registry type
 *    - Broadcasts environment change to all windows
 *
 *  Usage:
 *    ${EnvVarUpdate} "ResultVar" "MYVAR" "A|P|R" "HKLM|HKCU" "PathString"
 *
 *    A = Append, P = Prepend, R = Remove
 */

!ifndef ENVVARUPDATE_NSH
!define ENVVARUPDATE_NSH

!include "LogicLib.nsh"
!include "WinMessages.nsh"

; Define maximum string length for long paths
!ifndef NSIS_MAX_STRLEN
  !define NSIS_MAX_STRLEN 8192
!endif

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

  Exch $3 ; pathname to add/remove
  Exch
  Exch $2 ; Action - A=Append, P=Prepend, R=Remove
  Exch 2
  Exch $1 ; HKCU/HKLM
  Exch 3
  Exch $0 ; env variable name

  Push $4 ; Current PATH value
  Push $5 ; Result/new PATH
  Push $6 ; Temp var
  Push $7 ; Registry key path
  Push $8 ; Loop counter / position
  Push $9 ; Segment

  ; Determine registry key based on root
  StrCmp $1 "HKCU" 0 +3
    StrCpy $7 "Environment"
    Goto _ReadCurrent
  StrCpy $7 "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"

_ReadCurrent:
  ; Read current value with long string support
  SetRegView 64
  StrCmp $1 "HKCU" 0 +3
    ReadRegStr $4 HKCU $7 $0
    Goto _ProcessAction
  ReadRegStr $4 HKLM $7 $0

_ProcessAction:
  ; Handle action type
  StrCmp $2 "A" _Append
  StrCmp $2 "P" _Prepend
  StrCmp $2 "R" _Remove
  ; Invalid action, return current value
  StrCpy $5 $4
  Goto _Done

_Append:
  ; Check if already in path
  Push $4
  Push $3
  Call ${un}EnvVarUpdate_IsInPath
  Pop $6
  ${If} $6 == 1
    ; Already exists, don't add again
    StrCpy $5 $4
  ${ElseIf} $4 == ""
    ; PATH is empty, just set to new value
    StrCpy $5 $3
  ${Else}
    ; Append with semicolon
    StrCpy $5 "$4;$3"
  ${EndIf}
  Goto _WriteValue

_Prepend:
  ; Check if already in path
  Push $4
  Push $3
  Call ${un}EnvVarUpdate_IsInPath
  Pop $6
  ${If} $6 == 1
    ; Already exists, don't add again
    StrCpy $5 $4
  ${ElseIf} $4 == ""
    ; PATH is empty, just set to new value
    StrCpy $5 $3
  ${Else}
    ; Prepend with semicolon
    StrCpy $5 "$3;$4"
  ${EndIf}
  Goto _WriteValue

_Remove:
  ${If} $4 == ""
    StrCpy $5 ""
    Goto _WriteValue
  ${EndIf}

  ; Remove the path segment
  Push $4
  Push $3
  Call ${un}EnvVarUpdate_RemoveFromPath
  Pop $5
  Goto _WriteValue

_WriteValue:
  ; Write new value to registry as REG_EXPAND_SZ
  SetRegView 64
  StrCmp $1 "HKCU" 0 +3
    WriteRegExpandStr HKCU $7 $0 $5
    Goto _Broadcast
  WriteRegExpandStr HKLM $7 $0 $5

_Broadcast:
  ; Broadcast WM_SETTINGCHANGE so applications pick up the change
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
  Push $5

FunctionEnd

;------------------------------------------
; Helper: Check if a path segment exists in PATH
; Returns 1 if found, 0 if not
;------------------------------------------
Function ${un}EnvVarUpdate_IsInPath
  Exch $0 ; path to find
  Exch
  Exch $1 ; current PATH
  Push $2 ; current position
  Push $3 ; segment end position
  Push $4 ; extracted segment
  Push $5 ; PATH length

  StrLen $5 $1
  StrCpy $2 0

_IsInPath_Loop:
  ; Safety check
  ${If} $2 > $5
    StrCpy $0 0
    Goto _IsInPath_End
  ${EndIf}

  ; Find next semicolon starting from position $2
  StrCpy $3 $2
_IsInPath_FindSemi:
  ${If} $3 >= $5
    ; No more semicolons, extract rest of string
    StrCpy $4 $1 "" $2
    Goto _IsInPath_Compare
  ${EndIf}
  StrCpy $4 $1 1 $3
  ${If} $4 == ";"
    ; Found semicolon at $3, extract segment
    IntOp $4 $3 - $2
    StrCpy $4 $1 $4 $2
    Goto _IsInPath_Compare
  ${EndIf}
  IntOp $3 $3 + 1
  Goto _IsInPath_FindSemi

_IsInPath_Compare:
  ; Compare segment with search path (case-insensitive)
  ${If} $4 == $0
    StrCpy $0 1
    Goto _IsInPath_End
  ${EndIf}

  ; Move to next segment
  IntOp $2 $3 + 1
  ${If} $2 > $5
    StrCpy $0 0
    Goto _IsInPath_End
  ${EndIf}
  Goto _IsInPath_Loop

_IsInPath_End:
  Pop $5
  Pop $4
  Pop $3
  Pop $2
  Pop $1
  Exch $0
FunctionEnd

;------------------------------------------
; Helper: Remove a path segment from PATH
; Returns new PATH with segment removed
;------------------------------------------
Function ${un}EnvVarUpdate_RemoveFromPath
  Exch $0 ; path to remove
  Exch
  Exch $1 ; current PATH
  Push $2 ; current position
  Push $3 ; segment end position
  Push $4 ; extracted segment
  Push $5 ; result PATH
  Push $6 ; PATH length
  Push $7 ; temp char

  StrLen $6 $1
  StrCpy $2 0
  StrCpy $5 ""

_RemovePath_Loop:
  ; Safety check
  ${If} $2 > $6
    Goto _RemovePath_End
  ${EndIf}
  ${If} $2 == $6
    Goto _RemovePath_End
  ${EndIf}

  ; Find next semicolon starting from position $2
  StrCpy $3 $2
_RemovePath_FindSemi:
  ${If} $3 >= $6
    ; No more semicolons, extract rest of string
    StrCpy $4 $1 "" $2
    StrCpy $3 $6
    Goto _RemovePath_CheckSegment
  ${EndIf}
  StrCpy $7 $1 1 $3
  ${If} $7 == ";"
    ; Found semicolon at $3, extract segment
    IntOp $4 $3 - $2
    StrCpy $4 $1 $4 $2
    Goto _RemovePath_CheckSegment
  ${EndIf}
  IntOp $3 $3 + 1
  Goto _RemovePath_FindSemi

_RemovePath_CheckSegment:
  ; If segment doesn't match path to remove, add it to result
  ${If} $4 != $0
    ${If} $5 == ""
      StrCpy $5 $4
    ${Else}
      StrCpy $5 "$5;$4"
    ${EndIf}
  ${EndIf}

  ; Move to next segment
  IntOp $2 $3 + 1
  Goto _RemovePath_Loop

_RemovePath_End:
  ; Return result
  StrCpy $0 $5
  Pop $7
  Pop $6
  Pop $5
  Pop $4
  Pop $3
  Pop $2
  Pop $1
  Exch $0
FunctionEnd

!macroend

; Insert the functions for installer and uninstaller
!insertmacro EnvVarUpdate_Func ""
!insertmacro EnvVarUpdate_Func "un."

!endif ; ENVVARUPDATE_NSH
