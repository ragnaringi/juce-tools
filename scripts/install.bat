@ECHO OFF
set url=https://github.com/ragnaringi/juce-tools/releases/download/0.1.0/juce-tools.exe
set filename=juce-tools.exe
ECHO Downloading binary from %url%
%SYSTEMROOT%\System32\WindowsPowerShell\v1.0\powershell -Command Invoke-WebRequest %url% -OutFile .\%filename%
set destdir="C:\Program Files\juce-tools"
set destfile="%destdir%\%filename%"
echo Installing binary to %destdir%
if not exist %destdir% (
  %SYSTEMROOT%\System32\WindowsPowerShell\v1.0\powershell -Command mkdir '%destdir%'
)
%SYSTEMROOT%\System32\WindowsPowerShell\v1.0\powershell -Command mv -force .\%filename% '%destfile%'
echo Appending %destdir% to PATH
%SYSTEMROOT%\System32\WindowsPowerShell\v1.0\powershell -Command $env:PATH += '%destdir%'