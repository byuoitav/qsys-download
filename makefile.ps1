$COMMAND = $args[0]

$NAME = "qsys-download"

function Test {
    Write-Output "Test"
    Invoke-Expression "go test -v $PKG_LIST"
}

function Deps {
    Write-Output "Downloading Dependencies"
    Invoke-Expression "go mod download"
}

function Build {
    Write-Output "Build"

    New-Item -Path dist -ItemType Directory
    # $location = Get-Location
    # Write-Output $location\deps

    if (Test-Path "cmd") {
        Set-Location "cmd"
        Write-Output "Entering \cmd"
    
        Write-Output "*****************************************"
        Write-Output "Building for linux-amd64"
        Set-Item -Path env:CGO_ENABLED -Value 0
        Set-Item -Path env:GOOS -Value "linux"
        Set-Item -Path env:GOARCH -Value "amd64"
        Invoke-Expression "go build -v -o ../dist/${NAME}-linux-amd64"

        Write-Output "*****************************************"
        Write-Output "Building for linux-arm"
        Set-Item -Path env:CGO_ENABLED -Value 0
        Set-Item -Path env:GOOS -Value "linux"
        Set-Item -Path env:GOARCH -Value "arm"
        Invoke-Expression "go build -v -o ../dist/${NAME}-linux-arm"

        Write-Output "*****************************************"
        Write-Output "Building for linux-arm"
        Set-Item -Path env:GOOS -Value "windows"
        Set-Item -Path env:GOARCH -Value "amd64"
        Invoke-Expression "go build -v -o ../dist/${NAME}-windows-arm"

        Write-Output "Build output is located in ./dist/."
        Invoke-Expression "cd .."
    }
}

function Cleanup {
    Write-Output "Clean"
    Invoke-Expression "go clean"
    if (Test-Path -Path "dist") {
    Remove-Item dist -recurse
    Write-Output "Recursively deleted dist/"
    } else {
        Write-Output "No dist directory to delete"
    }
}

if ($COMMAND -eq "All") {
    Cleanup
    Deps
    Build 
}
elseif ($COMMAND -eq "Test") {
    Deps
    Test
}

elseif ($COMMAND -eq "Deps") {
    Deps
}
elseif ($COMMAND -eq "Build") {
    Cleanup
    Deps
    Build
}
elseif ($COMMAND -eq "Clean") {
    Cleanup
}

else {
    Write-Output "Please include a valid command parameter"
}