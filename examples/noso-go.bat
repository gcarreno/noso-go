:loop
@echo off

Rem #################################
Rem ## Begin of user-editable part ##
Rem #################################


REM Valid Pools:
REM   DevNoso
REM   dukedog.io
REM   mining.moe
REM   russiapool

set "POOL=DevNoso"
set "CPU=2"
set "WALLET=N2RUEEpVEyF9fgmQg6HKcrwkm4MERDx"

Rem #################################
Rem ##  End of user-editable part  ##
Rem #################################
setlocal enableDelayedExpansion

tasklist /FI "IMAGENAME eq noso-go.exe" 2>NUL | find /I /N "noso-go.exe">NUL
if "%ERRORLEVEL%"=="0" taskkill /F /im noso-go.exe

noso-go.exe mine pool !POOL! --wallet !WALLET! --cpu !CPU!
timeout 10
goto loop
