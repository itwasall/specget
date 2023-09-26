# specget

specget is a portable tool designed for refurbishers to get all required specifications of an iMac or MacBook without having to create a User account. 

## 

## What it might do later
- WiFi test
- Format Drive (APFS/JHFS+/Fusion Formats)
- Install macOS (Either from USB or through Internet Recovery)

# Requirements
- Go

# Building
Use `build.sh -h` to find build options.
Or you can just use `go build main.go`

# Usage
Intended usage is to be copied over to an install USB and called from an instance of terminal either on the USB live system or from the recovery drive
`/Volumes/<YOUR_COOL_USB_NAME_HERE>/main`
