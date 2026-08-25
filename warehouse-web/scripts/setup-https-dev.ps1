param(
  [string[]]$ExtraHosts = @(),
  [string]$CertDir = "certs"
)

$ErrorActionPreference = "Stop"

$mkcert = Get-Command mkcert -ErrorAction SilentlyContinue
if (-not $mkcert) {
  Write-Host ""
  Write-Host "mkcert is not installed."
  Write-Host ""
  Write-Host "Install mkcert first, then run this command again:"
  Write-Host "  npm run cert:https"
  Write-Host ""
  Write-Host "Common install options:"
  Write-Host "  winget install FiloSottile.mkcert"
  Write-Host "  choco install mkcert"
  Write-Host ""
  Write-Host "After creating certs, copy certs\dev-ca.pem to each phone/tablet and trust it as a CA."
  exit 1
}

$targetDir = Join-Path $PSScriptRoot "..\$CertDir"
New-Item -ItemType Directory -Force -Path $targetDir | Out-Null

$ipv4Hosts = Get-NetIPAddress -AddressFamily IPv4 |
  Where-Object {
    $_.IPAddress -notlike "127.*" -and
    $_.IPAddress -notlike "169.254.*" -and
    $_.PrefixOrigin -ne "WellKnown"
  } |
  Select-Object -ExpandProperty IPAddress -Unique

$hosts = @("localhost", "127.0.0.1", $env:COMPUTERNAME) + $ipv4Hosts + $ExtraHosts |
  Where-Object { $_ } |
  Select-Object -Unique

$certPath = Join-Path $targetDir "dev-server-cert.pem"
$keyPath = Join-Path $targetDir "dev-server-key.pem"
$caPath = Join-Path $targetDir "dev-ca.pem"

& mkcert -install
& mkcert -cert-file $certPath -key-file $keyPath @hosts

$caRoot = (& mkcert -CAROOT).Trim()
$rootCA = Join-Path $caRoot "rootCA.pem"
if (Test-Path $rootCA) {
  Copy-Item -LiteralPath $rootCA -Destination $caPath -Force
}

$primaryIp = $ipv4Hosts | Select-Object -First 1
$openHost = if ($primaryIp) { $primaryIp } else { "localhost" }

Write-Host ""
Write-Host "HTTPS dev certificate created:"
Write-Host "  Cert: $certPath"
Write-Host "  Key:  $keyPath"
Write-Host "  CA:   $caPath"
Write-Host ""
Write-Host "Included hosts/IPs:"
foreach ($hostItem in $hosts) {
  Write-Host "  - $hostItem"
}
Write-Host ""
Write-Host "Restart Vite and open:"
Write-Host "  https://$openHost`:8084/warehouse-ktp/move-order"
Write-Host ""
Write-Host "For phones/tablets: install certs\dev-ca.pem and trust it as a CA."
