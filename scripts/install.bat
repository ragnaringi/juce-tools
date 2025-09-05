@echo off
setlocal

:: ----------------------------
:: Configuration
:: ----------------------------
set "REPO=ragnaringi/juce-tools"
set "BINARY=juce-tools.exe"
set "DEST=C:\Program Files\juce-tools"

:: Full PowerShell path
set "POWERSHELL=%SYSTEMROOT%\System32\WindowsPowerShell\v1.0\powershell.exe"

:: ----------------------------
:: Fetch latest release tag
:: ----------------------------
echo Fetching latest release info from GitHub...
%POWERSHELL% -NoProfile -Command "Invoke-RestMethod -Uri 'https://api.github.com/repos/%REPO%/releases/latest' -Headers @{ 'User-Agent'='batch' } | Select-Object -ExpandProperty tag_name | Out-File -Encoding ascii release_tag.txt"
for /f "usebackq delims=" %%V in ("release_tag.txt") do set "VERSION=%%V"
del release_tag.txt

if "%VERSION%"=="" (
    echo Could not determine latest release.
    exit /b 1
)
echo Latest release: %VERSION%

:: ----------------------------
:: Fetch download URL for Windows .exe
:: ----------------------------
%POWERSHELL% -NoProfile -Command "($release = Invoke-RestMethod -Uri 'https://api.github.com/repos/%REPO%/releases/latest' -Headers @{ 'User-Agent'='batch' }).assets | Where-Object { $_.name -imatch '\.exe$' } | Select-Object -First 1 -ExpandProperty browser_download_url | Out-File -Encoding ascii download_url.txt"
for /f "usebackq delims=" %%U in ("download_url.txt") do set "URL=%%U"
del download_url.txt

if "%URL%"=="" (
    echo Could not find a Windows .exe binary in the latest release.
    exit /b 1
)
echo Download URL: %URL%

:: ----------------------------
:: Download binary
:: ----------------------------
echo Downloading %BINARY% ...
curl -L -o "%BINARY%" "%URL%"
if errorlevel 1 (
    echo Download failed.
    exit /b 1
)

:: ----------------------------
:: Install binary
:: ----------------------------
echo Installing to %DEST%...
if not exist "%DEST%" (
    mkdir "%DEST%"
)
move /Y "%BINARY%" "%DEST%\%BINARY%" >nul
if errorlevel 1 (
    echo Failed to move binary to %DEST%.
    exit /b 1
)

:: ----------------------------
:: Add to PATH if missing
:: ----------------------------
echo %PATH% | find /i "%DEST%" >nul
if %errorlevel% neq 0 (
    setx PATH "%PATH%;%DEST%"
    echo Added %DEST% to system PATH. You may need to restart your terminal.
) else (
    echo %DEST% is already in PATH.
)

echo Installation complete.
echo You can now run: juce-tools --help
endlocal
