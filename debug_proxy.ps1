$regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings"
Write-Host "--- Current Web Proxy Settings ---"
Get-ItemProperty -Path $regPath | Select-Object ProxyEnable, ProxyServer, ProxyOverride

Write-Host "`n--- Testing Local Proxy Connection (if enabled) ---"
$proxy = Get-ItemProperty -Path $regPath -Name ProxyServer -ErrorAction SilentlyContinue
if ($proxy -and $proxy.ProxyServer) {
    $server = $proxy.ProxyServer
    Write-Host "Proxy Configured: $server"
    try {
        $response = Invoke-WebRequest -Uri "http://checkip.amazonaws.com" -Proxy "http://$server" -TimeoutSec 5
        Write-Host "Success! IP via Proxy: $($response.Content)"
    } catch {
        Write-Host "Failed to connect via proxy: $_"
    }
} else {
    Write-Host "No Proxy Server configured."
}
