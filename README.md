# specget

specget is a portable tool designed for refurbishers to get all required specifications of an iMac or MacBook without having to create a User account. 

# For acquiring battery information for Venutra/Sonoma, or for completely negating this applications existence please check the wiki!!

## Usage
From terminal on a MacOS install USB:

`/Volumes/Image\ Volume/specget`

The exact command may vary depending on if you gave your USB a name, just check the `/Volumes/` directory


## specget can currently get
- CPU name & clock speed
- GPU name & VRAM (AMD only)
- Hard Drive Size, Type, Protocol
- Serial Number
- Model ID (i.e. iMac 18,3)
- MacBook Battery Cycle Count (*Install required)
- MacBook Battery Condition (*Install required)

### "Install required" Commands
*Whilst no user account needs to be made, an installation of MacOS with partition name "Macintosh HD" is required for these commands to work. In order to achieve this:*

- Boot into your install USB
- Install MacOS as normal, making sure to call the hard-drive partition **Macintosh HD**
- Once installation is complete (i.e. you are at the "Welcome" screen), reboot back into the install USB
- Run the desired command
## Usable on:
- High Seierra
- Mojave
- Catalina
- Big Sur
- Monterey
- Ventura (partial)

**Ventura workaround is in progress. Currently only the commands that require a `chroot` are affected, which are the *MacBook Battery Cycle Count* and *MacBook Battery Condition* commands**

## Building from Source
1. Don't.

2. Install `go` and run `go build main.go` from the root repo directory.

    or
 
2. Run the `build.sh` script. I haven't looked at it in ages, it might be a shambles.

## Dependancies
**Powered by Bubbletea** https://github.com/charmbracelet/bubbletea
