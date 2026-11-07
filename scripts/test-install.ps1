# Regression tests for install.ps1 helpers (checksum verification and
# zip-slip-safe extraction). Runs under Windows PowerShell 5.1 and PowerShell 7.
$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $PSScriptRoot
$env:GLADIA_INSTALL_FUNCTIONS_ONLY = '1'
. (Join-Path $root 'install.ps1')
Remove-Item Env:\GLADIA_INSTALL_FUNCTIONS_ONLY

Add-Type -AssemblyName System.IO.Compression, System.IO.Compression.FileSystem

$script:failures = 0
function Pass([string]$Name) { Write-Host "ok   - $Name" }
function Fail([string]$Name, [string]$Detail) {
    Write-Host "FAIL - $Name $Detail"
    $script:failures++
}

function New-TestZip {
    param([string]$Path, [string[]]$EntryNames)
    $zip = [System.IO.Compression.ZipFile]::Open($Path, [System.IO.Compression.ZipArchiveMode]::Create)
    try {
        foreach ($name in $EntryNames) {
            $entry = $zip.CreateEntry($name)
            $writer = New-Object System.IO.StreamWriter($entry.Open())
            try { $writer.Write('payload') } finally { $writer.Dispose() }
        }
    }
    finally {
        $zip.Dispose()
    }
}

function Expect-Throw {
    param([string]$Name, [scriptblock]$Block, [string]$Pattern)
    try {
        & $Block
        Fail $Name '(no error)'
    }
    catch {
        if ($_.Exception.Message -match $Pattern) { Pass $Name } else { Fail $Name "(unexpected error: $($_.Exception.Message))" }
    }
}

$work = Join-Path ([System.IO.Path]::GetTempPath()) ("gladia-install-test-{0}" -f [guid]::NewGuid())
New-Item -ItemType Directory -Path $work | Out-Null

try {
    # --- Confirm-ArchiveChecksum ---
    $archive = Join-Path $work 'gladia_9.9.9_windows_amd64.zip'
    New-TestZip -Path $archive -EntryNames @('gladia.exe')
    $hash = (Get-FileHash -Path $archive -Algorithm SHA256).Hash.ToLowerInvariant()
    $sums = Join-Path $work 'checksums.txt'

    Set-Content -Path $sums -Value "$hash  gladia_9.9.9_windows_amd64.zip"
    try { Confirm-ArchiveChecksum -FilePath $archive -ChecksumsPath $sums; Pass 'accepts a matching checksum' }
    catch { Fail 'accepts a matching checksum' $_.Exception.Message }

    Set-Content -Path $sums -Value ("{0}  gladia_9.9.9_windows_amd64.zip" -f ('0' * 64))
    Expect-Throw 'rejects a mismatching checksum' { Confirm-ArchiveChecksum -FilePath $archive -ChecksumsPath $sums } 'Checksum mismatch'

    Set-Content -Path $sums -Value "$hash  other.zip"
    Expect-Throw 'rejects an archive missing from checksums.txt' { Confirm-ArchiveChecksum -FilePath $archive -ChecksumsPath $sums } 'No sha256 entry'

    # --- Expand-ZipArchiveSafe ---
    $good = Join-Path $work 'good.zip'
    New-TestZip -Path $good -EntryNames @('gladia.exe', 'docs/README.md')
    $dest = Join-Path $work 'good-out'
    New-Item -ItemType Directory -Path $dest | Out-Null
    try {
        Expand-ZipArchiveSafe -ZipPath $good -DestinationPath $dest
        if ((Test-Path (Join-Path $dest 'gladia.exe')) -and (Test-Path (Join-Path $dest 'docs\README.md'))) { Pass 'extracts a benign archive' }
        else { Fail 'extracts a benign archive' '(files missing)' }
    }
    catch { Fail 'extracts a benign archive' $_.Exception.Message }

    $cases = @(
        @{ Name = 'rejects a ..\ traversal entry'; Entry = '..\..\gladia-zipslip-a.txt' },
        @{ Name = 'rejects a ../ traversal entry'; Entry = 'sub/../../gladia-zipslip-b.txt' },
        @{ Name = 'rejects a rooted entry'; Entry = '/gladia-zipslip-c.txt' },
        @{ Name = 'rejects a drive-qualified entry'; Entry = 'C:\gladia-zipslip-d.txt' }
    )
    $i = 0
    foreach ($case in $cases) {
        $i++
        $evil = Join-Path $work "evil$i.zip"
        New-TestZip -Path $evil -EntryNames @('gladia.exe', $case.Entry)
        $evilDest = Join-Path $work "evil$i-out\a\b"
        New-Item -ItemType Directory -Path $evilDest -Force | Out-Null
        Expect-Throw $case.Name { Expand-ZipArchiveSafe -ZipPath $evil -DestinationPath $evilDest } 'unsafe path|escapes the destination'
    }
    $escaped = Get-ChildItem -Path $work -Recurse -Filter 'gladia-zipslip-*' -ErrorAction SilentlyContinue
    if ($escaped) { Fail 'no traversal entry written' ($escaped.FullName -join ', ') } else { Pass 'no traversal entry written' }

    # --- Tag validation is applied before use ---
    foreach ($tag in @('v1.0;calc', '..', 'v1/../x')) {
        if ($tag -match '^[A-Za-z0-9][A-Za-z0-9._+-]*$') { Fail "tag pattern rejects '$tag'" '' } else { Pass "tag pattern rejects '$tag'" }
    }
    $src = Get-Content -Path (Join-Path $root 'install.ps1') -Raw
    if ($src.Contains("`$tag -notmatch '^[A-Za-z0-9][A-Za-z0-9._+-]*$'")) { Pass 'installer validates the tag' }
    else { Fail 'installer validates the tag' '(check not found)' }
}
finally {
    Remove-Item -Recurse -Force $work -ErrorAction SilentlyContinue
}

if ($script:failures -ne 0) {
    Write-Host "$($script:failures) installer test(s) failed"
    exit 1
}
Write-Host 'all installer tests passed'
