# macOS Hardware Detection Program

Ever just like, recieve 150 MacBooks and need to verify the specs of each one? Ya me neither, but I made this program to do it anyway

## What it do
- CPU
- RAM
- GPU
- HDD (Size/SSD/SSD Protocol/Fusion Status)
- Battery Information (Cycle Count/Condition) (Requires macOS to be installed. Doesn't require the set-up process to be completed though)

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
