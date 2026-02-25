$ErrorActionPreference = "Stop"

$Repo = "frostyeti/osv"
$InstallDir = "$env:LOCALAPPDATA\Programs\bin"

# Determine Architecture
$Arch = $env:PROCESSOR_ARCHITECTURE.ToLower()
$ArchName = ""
if ($Arch -match "amd64" -or $Arch -match "x64") {
    $ArchName = "amd64"
} elseif ($Arch -match "arm64") {
    $ArchName = "arm64"
} else {
    Write-Error "Unsupported architecture: $Arch"
    exit 1
}

Write-Host "Detecting latest release for windows_$ArchName..."
$ReleaseUrl = "https://api.github.com/repos/$Repo/releases/latest"
$Release = Invoke-RestMethod -Uri $ReleaseUrl

$Asset = $Release.assets | Where-Object { $_.name -like "*windows_$ArchName.zip" }
if (-not $Asset) {
    Write-Error "Could not find a release for windows $ArchName."
    exit 1
}

$AssetUrl = $Asset.browser_download_url
Write-Host "Downloading $AssetUrl..."

$TmpDir = New-TemporaryFile | Select-Object -ExpandProperty DirectoryName
$TmpPath = Join-Path $TmpDir "osv_temp_install_$([guid]::NewGuid().ToString().Substring(0,8))"
New-Item -ItemType Directory -Path $TmpPath | Out-Null
$ZipFile = Join-Path $TmpPath "osv.zip"

# Windows ships with curl (curl.exe) natively since build 17063
curl.exe -sL $AssetUrl -o $ZipFile

Write-Host "Extracting..."
# Windows also ships with an inbox version of tar natively capable of extracting zips
Set-Location -Path $TmpPath
tar.exe -xf $ZipFile

Write-Host "Installing to $InstallDir..."
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Move-Item -Path (Join-Path $TmpPath "osv.exe") -Destination (Join-Path $InstallDir "osv.exe") -Force

Remove-Item -Path $TmpPath -Recurse -Force

Write-Host "========================================================="
Write-Host "osv was successfully installed to $InstallDir\osv.exe"
Write-Host ""
$PathEnv = [Environment]::GetEnvironmentVariable("Path", "User")
if ($PathEnv -notmatch [regex]::Escape($InstallDir)) {
    Write-Host "WARNING: $InstallDir is not in your PATH."
    Write-Host "Please update your Path environment variable by running the following command:"
    Write-Host ""
    Write-Host "    [Environment]::SetEnvironmentVariable('Path', `"`$env:Path;$InstallDir`", 'User')"
    Write-Host ""
}
Write-Host "Run 'osv --help' to get started!"
Write-Host "========================================================="