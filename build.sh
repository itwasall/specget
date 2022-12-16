
if [ "$1" == "-w" ] || [ "$1" == "--windows" ]
then
  echo "Building windows amd64"
  # Windows Build
  export GOOS=windows go build

  go build main.go

  mv ./main.exe ./bin/specget_win64.exe

elif [ "$1" == "-l" ] || [ "$1" == "--linux" ]
then
  echo "Building linux amd64"
  # Linux Build
  export GOOS=linux go build

  go build main.go

  mv ./main ./bin/specget_unix64

elif [ "$1" == "-x86" ] || [ "$1" == "-x86-64" ] || [ "$1" == "-mac-intel" ] || [ "$1" == "--mac" ]
then
  if [ "$1" == "--mac" ]
  then
    echo "Please be aware this argument will be depreciated or" 
    echo "redirected to ARM based iMac's and MacBooks in the future"
  fi
  echo "Building darwin amd64"
  # macOS Intel Build
  export GOOS=darwin go build

  go build main.go

  mv ./main ./bin/specget_mac64

elif [ "$1" == "-arm" ] || [ "$1" == "-m1" ] || [ "$1" == "--mac-arm64" ]
then
  echo "Building darwin arm64"
  # macOS M1 Build
  export GOOS=darwin go build
  export GOARCH=arm64 go build

  go build main.go

  mv ./main ./bin/specget_macARM
else
  go build main.go
  
  mv ./main ./bin/main
fi
echo "Complete"
export GOOS=linux go build
export GOARCH=amd64 go build
echo "Envs reset"
