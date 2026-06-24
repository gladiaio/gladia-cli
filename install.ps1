#Requires -Version 5.1
$ErrorActionPreference = 'Stop'

$Repo = 'gladiaio/gladia-cli'
$BinaryName = 'gladia.exe'

function Get-Arch {
    if ([Environment]::Is64BitOperatingSystem) {
        return 'amd64'
    }
    return '386'
}

function Get-LatestTag {
    $headers = @{ Accept = 'application/vnd.github+json' }
    if ($env:GITHUB_TOKEN) {
        $headers['Authorization'] = "Bearer $($env:GITHUB_TOKEN)"
    }

    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Headers $headers
    return $release.tag_name
}

function Add-ToUserPath {
    param([string]$Directory)

    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if ($null -eq $userPath) {
        $userPath = ''
    }

    $parts = $userPath -split ';' | Where-Object { $_ -and $_ -ne $Directory }
    $newPath = ($parts + $Directory) -join ';'
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')

    if ($env:Path -notlike "*$Directory*") {
        $env:Path = "$env:Path;$Directory"
    }
}

$arch = Get-Arch
$tag = Get-LatestTag

if (-not $tag) {
    throw 'Could not determine latest release'
}

$version = $tag.TrimStart('v')
$archive = "gladia_${version}_windows_${arch}.zip"
$url = "https://github.com/$Repo/releases/download/$tag/$archive"
$installDir = if ($env:GLADIA_INSTALL_DIR) { $env:GLADIA_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA 'Programs\gladia-cli\bin' }

$tmpdir = Join-Path $env:TEMP ("gladia-install-{0}" -f [guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $tmpdir -Force | Out-Null

try {
    Write-Host "Installing gladia $tag (windows/$arch)..."

    $zipPath = Join-Path $tmpdir $archive
    Invoke-WebRequest -Uri $url -OutFile $zipPath -UseBasicParsing
    Expand-Archive -Path $zipPath -DestinationPath $tmpdir -Force

    $binaryPath = Join-Path $tmpdir $BinaryName
    if (-not (Test-Path $binaryPath)) {
        throw "$BinaryName not found in archive"
    }

    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    Copy-Item -Path $binaryPath -Destination (Join-Path $installDir $BinaryName) -Force
    Add-ToUserPath -Directory $installDir

    $installed = Join-Path $installDir $BinaryName
    Write-Host "Installed gladia to $installed"
    Write-Host "Restart your terminal if gladia is not found"
}
finally {
    Remove-Item -Recurse -Force $tmpdir -ErrorAction SilentlyContinue
}
