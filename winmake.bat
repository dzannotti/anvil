@echo off
setlocal enabledelayedexpansion

:: Function to ensure bin directory exists
if not exist bin mkdir bin

:: Handle arguments
if "%1"=="run" (
    call :build
    bin\anvil.exe
    exit /b
)

if "%1"=="build" (
    go build -o bin/anvil.exe cmd/hello/main.go
    exit /b
)

if "%1"=="release" (
    go build -trimpath, -ldflags="-w -s" -o bin/anvil.exe cmd/hello/main.go
    exit /b
)

if "%1"=="test" (
    go test -json ./...
    exit /b
)

if "%1"=="watch" (
    air
    exit /b
)

if "%1"=="tdd" (
    gotestsum --format testname --watch ./internal/...
    exit /b
)

if "%1"=="lint" (
    golangci-lint run
    exit /b
)
:: Default case (show usage)
echo Usage: make [run | build | test | watch | tdd | lint]
exit /b

:build
go build -o bin/main.exe cmd/main.go
exit /b