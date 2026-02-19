$server = "138.124.101.69"
$port = 50195
Write-Host "Testing TCP connection to $server on port $port..."
try {
    $tcp = Test-NetConnection -ComputerName $server -Port $port
    if ($tcp.TcpTestSucceeded) {
        Write-Host "✅ TCP Connection Successful!" -ForegroundColor Green
    } else {
        Write-Host "❌ TCP Connection Failed! Port is blocked or server is down." -ForegroundColor Red
    }
} catch {
    Write-Host "Error running test: $_"
}
