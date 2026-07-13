param(
    [string]$Root = "",
    [string]$MockgenVersion = "v0.6.0",
    [switch]$UseInstalledMockgen
)

$ErrorActionPreference = "Stop"

function Find-ModuleRoot {
    param([string]$StartPath)

    $current = Resolve-Path $StartPath
    while ($null -ne $current) {
        if (Test-Path (Join-Path $current "go.mod")) {
            return $current.Path
        }

        $parent = Split-Path $current -Parent
        if ([string]::IsNullOrWhiteSpace($parent) -or $parent -eq $current) {
            break
        }

        $current = Resolve-Path $parent
    }

    throw "go.mod not found from $StartPath"
}

function Convert-ToGoPath {
    param(
        [string]$RootPath,
        [string]$Path
    )

    $rootFullPath = [System.IO.Path]::GetFullPath($RootPath).TrimEnd("\", "/")
    $fullPath = [System.IO.Path]::GetFullPath($Path)
    $relative = $fullPath.Substring($rootFullPath.Length).TrimStart("\", "/")
    return $relative.Replace("\", "/")
}

function Get-GoPackage {
    param([string]$Path)

    $packageLine = Select-String -Path $Path -Pattern '^\s*package\s+([A-Za-z_][A-Za-z0-9_]*)\s*$' -List
    if ($null -eq $packageLine) {
        throw "package declaration not found in $Path"
    }

    return $packageLine.Matches[0].Groups[1].Value
}

function Test-HasInterface {
    param([string]$Path)

    return $null -ne (Select-String -Path $Path -Pattern '^\s*type\s+[A-Za-z_][A-Za-z0-9_]*\s+interface\s*\{' -List)
}

function Invoke-Mockgen {
    param([string[]]$MockgenArgs)

    if ($UseInstalledMockgen -and (Get-Command mockgen -ErrorAction SilentlyContinue)) {
        & mockgen @MockgenArgs
        return
    }

    & go run "go.uber.org/mock/mockgen@$MockgenVersion" @MockgenArgs
}

function New-Mock {
    param(
        [string]$RootPath,
        [string]$SourcePath
    )

    $destinationPath = $SourcePath -replace '\.go$', '_mock.go'
    $packageName = Get-GoPackage $SourcePath
    $source = Convert-ToGoPath $RootPath $SourcePath
    $destination = Convert-ToGoPath $RootPath $destinationPath

    Write-Host "mockgen $source -> $destination"

    Invoke-Mockgen @(
        "-source=$source",
        "-destination=$destination",
        "-package=$packageName"
    )
}

if ([string]::IsNullOrWhiteSpace($Root)) {
    $Root = Find-ModuleRoot $PSScriptRoot
}

$Root = (Resolve-Path $Root).Path
Set-Location $Root

$sources = @()
$baseUseCase = Join-Path $Root "pkg/base/usecase.go"
if (Test-Path $baseUseCase) {
    $sources += (Get-Item $baseUseCase)
}

$appsPath = Join-Path $Root "apps"
if (Test-Path $appsPath) {
    $sources += Get-ChildItem -Path $appsPath -Recurse -File -Filter "*.go" |
        Where-Object {
            $_.FullName -match '\\internal\\application\\port\\' -and
            $_.Name -notmatch '_mock\.go$' -and
            $_.Name -notmatch '_test\.go$' -and
            (Test-HasInterface $_.FullName)
        } |
        Sort-Object FullName
}

foreach ($source in $sources) {
    New-Mock -RootPath $Root -SourcePath $source.FullName
}

Write-Host "Mocks generated: $($sources.Count)"
