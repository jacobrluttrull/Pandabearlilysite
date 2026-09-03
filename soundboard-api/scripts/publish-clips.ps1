<#
.SYNOPSIS
    Publish every new clip in the staging folder to the live soundboard.

.DESCRIPTION
    Drop .mp3 files into "Pandalily Soundbites" at the repo root, then run this.
    It imports anything not already published and audits the result.

    Import is idempotent: clips already stored are skipped, so re-running after an
    interruption is safe and only ever adds what is missing.

    Paths are derived from this script's own location, so the alias works from any
    directory and survives the repo being moved.

.PARAMETER DryRun
    Report what would be published without copying audio or writing rows.

.EXAMPLE
    clips -DryRun
    clips
#>
[CmdletBinding()]
param(
    [switch]$DryRun
)

$ErrorActionPreference = 'Stop'

$apiDir = Split-Path $PSScriptRoot -Parent
$repoDir = Split-Path $apiDir -Parent
$clipsDir = Join-Path $repoDir 'Pandalily Soundbites'
$productionEnv = Join-Path $apiDir '.env.production.local'

if (-not (Test-Path $clipsDir)) {
    Write-Error "Staging folder not found: $clipsDir"
    exit 1
}

$mp3Count = @(Get-ChildItem -Path $clipsDir -Filter *.mp3 -File).Count
if ($mp3Count -eq 0) {
    Write-Host "No .mp3 files in $clipsDir - nothing to publish." -ForegroundColor Yellow
    exit 0
}

if (-not (Test-Path -LiteralPath $productionEnv -PathType Leaf)) {
    Write-Error "Production credentials file not found: $productionEnv`nCopy .env.example to .env.production.local and add freshly rotated publishing credentials."
    exit 1
}

Write-Host "$mp3Count mp3 file(s) staged in $clipsDir" -ForegroundColor Cyan
Write-Host ""

# go run must be invoked from inside the module, and the CLI resolves .env and its
# default paths relative to the working directory. Push/Pop leaves the caller's
# location untouched either way.
Push-Location $apiDir
$previousEnvFile = [Environment]::GetEnvironmentVariable('SOUNDBOARD_ENV_FILE', 'Process')
[Environment]::SetEnvironmentVariable('SOUNDBOARD_ENV_FILE', $productionEnv, 'Process')
try {
    if ($DryRun) {
        go run .\cmd\cli import -dir $clipsDir -dry-run
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

        Write-Host ""
        Write-Host "Dry run only - nothing was published. Re-run without -DryRun to publish." -ForegroundColor Yellow
        exit 0
    }

    go run .\cmd\cli import -dir $clipsDir
    if ($LASTEXITCODE -ne 0) {
        Write-Host ""
        Write-Host "Import failed - skipping the audit." -ForegroundColor Red
        exit $LASTEXITCODE
    }

    # The audit is the part that catches a half-configured run: it compares the
    # database, the clip store and names.json, and says so when they disagree.
    Write-Host ""
    Write-Host "--- audit ---" -ForegroundColor Cyan
    go run .\cmd\cli check
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

    Write-Host ""
    Write-Host "Done. Refresh the soundboard to see the new clips." -ForegroundColor Green
}
finally {
    if ($null -eq $previousEnvFile) {
        Remove-Item Env:SOUNDBOARD_ENV_FILE -ErrorAction SilentlyContinue
    }
    else {
        [Environment]::SetEnvironmentVariable('SOUNDBOARD_ENV_FILE', $previousEnvFile, 'Process')
    }
    Pop-Location
}
