# # make sure we start in the right directory
Push-Location $PSScriptRoot

# # load the environment variables
& ./env.ps1

$env:BUILD_WEBUI = "true"

# # pull these vars from the environment
$remoteServer = $env:REMOTE_SERVER
$remoteUser = $env:REMOTE_USER
$dockerUser = $env:DOCKER_USER

# print the vars
Write-Host "Remote Server: $remoteServer"
Write-Host "Remote User: $remoteUser"
Write-Host "Docker User: $dockerUser"

# build and push the docker images
Push-Location ../api-server/go
docker build -t ${dockerUser}/shorturls-api-go:latest .
docker push ${dockerUser}/shorturls-api-go:latest
Pop-Location

if ($env:BUILD_WEBUI -eq "true") {
    Write-Host "Building webui"
    # build and deploy the webui
    Push-Location ../webui
    yarn
    #yarn release
    #git push
    yarn prod
    # zip the dist/ folder
    $date = Get-Date -Format "yyyy-MM-dd-HH-mm-ss"
    $commit = git rev-parse --short HEAD
    $zipFile = "webui_${date}_${commit}.zip"
    mkdir -ErrorAction Ignore -Path ./artifacts
    Compress-Archive -Path ./dist/* -DestinationPath ./artifacts/${zipFile}

    # deploy to remote server
    # deploy the server
    pscp -P 22 ./artifacts/${zipFile} ${remoteUser}@${remoteServer}:/var/www/
    plink -batch -P 22 -ssh ${remoteUser}@${remoteServer} "cd /var/www/ && rm -rf ./html && unzip -o ./${zipFile} -d ./html && rm -f ./${zipFile}"
    Pop-Location
}

plink -batch -P 22 -ssh ${remoteUser}@${remoteServer} "docker pull caellach/shorturls-api-go:latest && docker stop shorturls-api-go || true && docker rm shorturls-api-go || true && docker run -d --restart unless-stopped -v /etc/shorturls-api-go/config.json:/app/config.json -v /etc/shorturls-api-go/wordlist.json:/app/wordlist.json --name shorturls-api-go -p 8080:3000 ${dockerUser}/shorturls-api-go:latest"
