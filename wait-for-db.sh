# Script to wait for db container to be ready for connections before running the bot

#!/bin/bash

param(
    [string]$host,
    [int]$port
)


while ($true) {
    try {
        $socket = New-Object System.Net.Sockets.TcpClient($host, $port)
        $socket.Close()
        Write-Host "Database is up - executing command"
        # Add your bot start command here
        break
    } catch {
        Write-Host "Database is unavailable - sleeping"
        Start-Sleep -Seconds 1
    }
}

while ($true) {
    try {
        $socket = New-Object System.Net.Sockets.TcpClient($host, $port)
        $socket.Close()
        Write-Host "Database is up - executing command"
        # Add your bot start command here
        break
    } catch {
        Write-Host "Database is unavailable - sleeping"
        Start-Sleep -Seconds 1
    }
}