@echo off
REM BackupFinder Windows Setup
REM This batch file helps Windows users set up BackupFinder

echo.
echo =======================================
echo    BackupFinder Windows Setup
echo =======================================
echo.

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% == 0 (
    echo [✓] Go is already installed
    go version
    echo.
    echo [Building BackupFinder...]
    go build -o backupfinder.exe ./cmd/backupfinder
    if %errorlevel% == 0 (
        echo [✓] Build successful!
        echo.
        echo [Testing BackupFinder...]
        backupfinder.exe health-check
        if %errorlevel% == 0 (
            echo [✓] BackupFinder is working correctly!
            echo.
            echo Usage examples:
            echo   backupfinder.exe -u example.com -w
            echo   backupfinder.exe --help
            echo   backupfinder.exe system-check
        ) else (
            echo [!] Build completed but health check failed
            echo     Run: backupfinder.exe system-check
        )
    ) else (
        echo [✗] Build failed. Please check Go installation.
        echo     Run: go version && go env
    )
) else (
    echo [✗] Go is not installed on this system
    echo.
    echo Windows Go Installation Options:
    echo.
    echo 1. Download installer from: https://go.dev/dl/
    echo    - Choose Windows installer (.msi)
    echo    - Follow installation wizard
    echo    - Restart command prompt after installation
    echo.
    echo 2. Using Chocolatey:
    echo    choco install golang
    echo.
    echo 3. Using Scoop:
    echo    scoop install go
    echo.
    echo After installing Go, run this script again.
)

echo.
pause
