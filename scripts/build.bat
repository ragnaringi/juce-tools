@echo off
setlocal enabledelayedexpansion

:: ----------------------------
:: Configuration
:: ----------------------------
:: First argument is TAG, default to "dev"
set "TAG=%~1"
if "%TAG%"=="" set "TAG=dev"

:: Go CLI directory (relative to script)
set "CLIDIR=%~dp0\..\cli"
set "OUTDIR=%~dp0\.."

:: Detect host architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set "HOSTARCH=amd64"
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set "HOSTARCH=arm64"
) else (
    echo Unsupported architecture: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)

echo Detected host architecture: %HOSTARCH%

:: Change to CLI directory
cd /d "%CLIDIR%"

echo Fetching Go module dependencies...
go mod tidy
go mod download

echo Building juce-tools for Windows %HOSTARCH%...
set "GOOS=windows"
set "GOARCH=%HOSTARCH%"

go build -o "%OUTDIR%\juce-tools_%TAG%_windows_%HOSTARCH%.exe" ./...

echo Built juce-tools_%TAG%_windows_%HOSTARCH%.exe
endlocal
