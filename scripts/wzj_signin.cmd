@echo off
setlocal
set "PROJECT_ROOT=%~dp0.."
powershell -NoProfile -ExecutionPolicy Bypass -File "%PROJECT_ROOT%\scripts\run-local.ps1" -Port 8081 %*
