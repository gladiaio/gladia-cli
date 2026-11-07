# Gladia CLI installer for Windows
# Downloads the latest GitHub release and installs the binary to a local directory.
# Requires: PowerShell 5+ (ships with Windows 10/11), internet access.
#
# Usage (PowerShell):
#   irm https://raw.githubusercontent.com/gladiaio/gladia-cli/main/install.ps1 | iex
#
# Environment variables:
#   GLADIA_INSTALL_DIR           Override install directory (default: $env:LOCALAPPDATA\Programs\gladia-cli\bin)
#   GITHUB_TOKEN                 Optional, for GitHub API rate limits (only sent to api.github.com)
#   GLADIA_INSTALL_API_URL       Optional override for the release API endpoint (testing)
#   GLADIA_INSTALL_DOWNLOAD_BASE Optional override for the release download base URL (testing)
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
    $defaultApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
    $apiUrl = if ($env:GLADIA_INSTALL_API_URL) { $env:GLADIA_INSTALL_API_URL } else { $defaultApiUrl }

    $headers = @{ Accept = 'application/vnd.github+json' }
    # Only send GITHUB_TOKEN to the real GitHub API, never to an override URL.
    if ($env:GITHUB_TOKEN -and $apiUrl -eq $defaultApiUrl) {
        $headers['Authorization'] = "Bearer $($env:GITHUB_TOKEN)"
    }

    $release = Invoke-RestMethod -Uri $apiUrl -Headers $headers
    return $release.tag_name
}

function Confirm-ArchiveChecksum {
    param(
        [string]$FilePath,
        [string]$ChecksumsPath
    )

    $name = Split-Path $FilePath -Leaf
    $expected = $null
    foreach ($line in (Get-Content -Path $ChecksumsPath)) {
        if ($line -match "^([0-9a-fA-F]{64})[ \t]+\*?$([regex]::Escape($name))[ \t]*$") {
            $expected = $Matches[1]
            break
        }
    }
    if (-not $expected) {
        throw "No sha256 entry for $name in checksums.txt"
    }

    $actual = (Get-FileHash -Path $FilePath -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $expected.ToLowerInvariant()) {
        throw "Checksum mismatch for ${name}: expected $expected, got $actual"
    }
    Write-Host "Checksum verified: $name"
}

# Extracts a zip archive while confining every entry to $DestinationPath.
# Windows PowerShell 5.1's built-in Expand-Archive joins entry names to the
# destination with no validation (zip-slip), so entries with '..' segments,
# absolute paths, or drive/UNC prefixes would be written outside the temp
# directory. This mirrors the containment check added upstream in
# Microsoft.PowerShell.Archive v1.2.6.
function Expand-ZipArchiveSafe {
    param(
        [string]$ZipPath,
        [string]$DestinationPath
    )

    Add-Type -AssemblyName System.IO.Compression, System.IO.Compression.FileSystem
    $zip = [System.IO.Compression.ZipFile]::OpenRead($ZipPath)
    try {
        $destRoot = (Resolve-Path -LiteralPath $DestinationPath).ProviderPath
        if (-not $destRoot.EndsWith([System.IO.Path]::DirectorySeparatorChar)) {
            $destRoot += [System.IO.Path]::DirectorySeparatorChar
        }

        foreach ($entry in $zip.Entries) {
            $rawName = $entry.FullName
            # Skip empty / root entries.
            $segments = @($rawName -split '[/\\]' | Where-Object { $_ -ne '' })
            if (@($segments).Count -eq 0) {
                continue
            }

            # Reject drive-letter (C:), UNC/device (\\.\, \\?\), rooted
            # (leading slash), and '..' entries before any path resolution.
            if ($rawName -match '^[A-Za-z]:' -or $rawName -match '^[\\/]' -or @($segments) -contains '..') {
                throw "Archive entry has an unsafe path: $rawName"
            }

            $targetPath = [System.IO.Path]::GetFullPath((Join-Path $DestinationPath $rawName))
            if (-not $targetPath.StartsWith($destRoot, [System.StringComparison]::OrdinalIgnoreCase)) {
                throw "Archive entry escapes the destination directory: $rawName"
            }

            if ($rawName.EndsWith('/') -or $rawName.EndsWith('\')) {
                if (-not (Test-Path -LiteralPath $targetPath -PathType Container)) {
                    New-Item -ItemType Directory -Path $targetPath -Force | Out-Null
                }
            }
            else {
                $parent = Split-Path -Parent $targetPath
                if (-not (Test-Path -LiteralPath $parent)) {
                    New-Item -ItemType Directory -Path $parent -Force | Out-Null
                }
                [System.IO.Compression.ZipFileExtensions]::ExtractToFile($entry, $targetPath, $true)
            }
        }
    }
    finally {
        if ($null -ne $zip) {
            $zip.Dispose()
        }
    }
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

function Write-CompletionHints {
    Write-Host ''
    Write-Host 'Shell tab completion is available. To set up manually:'
    Write-Host '  gladia completion --help'
    Write-Host '  gladia completion powershell'
}

function Install-PowerShellCompletions {
    param(
        [string]$GladiaExe
    )

    $profilePath = $PROFILE
    $profileDir = Split-Path -Parent $profilePath
    if (-not (Test-Path $profileDir)) {
        New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
    }
    if (-not (Test-Path $profilePath)) {
        New-Item -ItemType File -Path $profilePath -Force | Out-Null
    }

    $marker = '# gladia-cli completions'
    if (Test-Path $profilePath) {
        $existing = Get-Content -Path $profilePath -Raw -ErrorAction SilentlyContinue
        if ($existing -and $existing.Contains($marker)) {
            Write-Host "PowerShell completion already configured in $profilePath"
            Write-Host 'Restart your terminal for completion to take effect.'
            return
        }
    }

    $completion = & $GladiaExe completion powershell
    Add-Content -Path $profilePath -Value "`n$marker"
    Add-Content -Path $profilePath -Value $completion
    Write-Host "Wrote PowerShell completion to $profilePath"
    Write-Host 'Restart your terminal for completion to take effect.'
}

function Maybe-PromptCompletions {
    param(
        [string]$GladiaExe
    )

    if ($env:GLADIA_NO_COMPLETION_PROMPT) {
        Write-CompletionHints
        return
    }

    if ([Console]::IsInputRedirected -or [Console]::IsOutputRedirected) {
        Write-CompletionHints
        return
    }

    $reply = Read-Host 'Install shell tab completion? [y/N]'
    if ($reply -match '^[yY]') {
        Install-PowerShellCompletions -GladiaExe $GladiaExe
    }
    else {
        Write-CompletionHints
    }
}

# Lets tests dot-source this script to load the functions above without
# running the installation.
if ($env:GLADIA_INSTALL_FUNCTIONS_ONLY) {
    return
}

$arch = Get-Arch
$tag = Get-LatestTag

if (-not $tag) {
    throw 'Could not determine latest release'
}

# Sanity-check the tag before interpolating it into URLs or filenames.
if ($tag -notmatch '^[A-Za-z0-9][A-Za-z0-9._+-]*$') {
    throw "Unexpected tag format from GitHub API: $tag"
}

$version = $tag.TrimStart('v')
$archive = "gladia_${version}_windows_${arch}.zip"
$downloadBase = if ($env:GLADIA_INSTALL_DOWNLOAD_BASE) { $env:GLADIA_INSTALL_DOWNLOAD_BASE.TrimEnd('/','\') } else { "https://github.com/$Repo/releases/download" }
$url = "$downloadBase/$tag/$archive"
$checksumsUrl = "$downloadBase/$tag/checksums.txt"
$installDir = if ($env:GLADIA_INSTALL_DIR) { $env:GLADIA_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA 'Programs\gladia-cli\bin' }

$tmpdir = Join-Path $env:TEMP ("gladia-install-{0}" -f [guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $tmpdir -Force | Out-Null

try {
    Write-Host "Installing gladia $tag (windows/$arch)..."

    $zipPath = Join-Path $tmpdir $archive
    Invoke-WebRequest -Uri $url -OutFile $zipPath -UseBasicParsing
    $checksumsPath = Join-Path $tmpdir 'checksums.txt'
    Invoke-WebRequest -Uri $checksumsUrl -OutFile $checksumsPath -UseBasicParsing
    Confirm-ArchiveChecksum -FilePath $zipPath -ChecksumsPath $checksumsPath
    Expand-ZipArchiveSafe -ZipPath $zipPath -DestinationPath $tmpdir

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

    Maybe-PromptCompletions -GladiaExe $installed
}
finally {
    Remove-Item -Recurse -Force $tmpdir -ErrorAction SilentlyContinue
}
