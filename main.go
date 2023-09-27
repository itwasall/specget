package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/slices"
)

var (
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#505059")).
			Foreground(lipgloss.Color("#ffffff"))
	disclaimerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#505059")).
			Foreground(lipgloss.Color("#ffffff")).
			PaddingLeft(2)
	bodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#404040"))
	warningStyle = lipgloss.NewStyle().
		// Background(lipgloss.Color("#505059")).
		Foreground(lipgloss.Color("#a03000"))

	quitStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a0a0a0")).
			PaddingLeft(3).
			PaddingRight(3)
	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0050c0"))

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#000000")).
			PaddingLeft(3).
			PaddingRight(3)
)

type model struct {
	altscreen      bool
	command        string
	commandPresent bool
	commandType    string
	err            error
	disclaimerShow bool
	testMenu       bool
	formatMenu     bool
}

type commandFinishedMsg struct{ err error }

func (m model) Init() tea.Cmd {
	return nil
}

func checkError(errorString string, err error) bool {
	if err != nil {
		fmt.Println(errorString, err)
		return true
	}
	return false
}

func checkDir(path string, searchterm string) bool {
	// Read items at path
	folderSearch, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Path Error! ", err)
		return false
	}
	var folderItems []string
	// Put the names of the items at path into slice/s
	for _, f := range folderSearch {
		folderItems = append(folderItems, f.Name())
	}
	// Searh slices for search term, returning true if found, false if err or not
	if slices.Contains(folderItems, searchterm) {
		return true
	}
	return false
}

func checkM1() bool {
	c := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")

	cpuNameOut, err := c.Output()
	if err != nil {
		fmt.Println("Error getting CPU: ", err)
		return false
	}

	cpuName := string(cpuNameOut[:])

	if strings.Contains(cpuName, " M1") {
		return true
	}
	return false
}

func checkMacOSVersion() string {
	c := exec.Command("bash", "-c", "sw_vers | sed -n '2 p' | awk '{print $2}'")

	macOsVersionOut, err := c.Output()
	if err != nil {
		fmt.Println("Error getting sw_vers", err)
		return "ERROR"
	}

	macOsVersion := string(macOsVersionOut[:])

	if strings.Contains(macOsVersion, "13") {
		return "Ventura"
	} else if strings.Contains(macOsVersion, "12") {
		return "Monterey"
	} else if strings.Contains(macOsVersion, "11") {
		return "Big Sur"
	} else if strings.Contains(macOsVersion, "10") {
		return "Old as shit boiiii"
	}
	return "Couldn't get OS version"
}

func checkFusion() bool {
	c := exec.Command("bash", "-c", "diskutil resetfusion | cat | awk '{print $1}' | sed -n '1 p'")

	fusionStatusOut, err := c.Output()
	if err != nil {
		fmt.Println("Error getting fusion drive status: ", err)
	}

	fusionStatus := string(fusionStatusOut[:])

	if strings.Contains(fusionStatus, "s internal disk devices must be solid-state") {
		return false
	}
	return true
}

func getHDD(m model) (tea.Model, tea.Cmd) {
	hddFusionCheck := exec.Command("bash", "-c", "diskutil info /dev/disk2")
	// dev/disk0 being the default drive macOS is installed on
	hddDisk0 := exec.Command(
		"bash",
		"-c",
		"diskutil info /dev/disk0 | grep \"Disk Size\" | awk '{print $3, $4, \"(\"substr($5, 2, 16)/1000/1000/1000, \"MB)\"}'",
	)
	// dev/disk2 being the default assignment a fusion drive gets. (Comprising of a disk0 SSD and a disk1 HDD for example)
	hddDisk2 := exec.Command(
		"bash",
		"-c",
		"diskutil info /dev/disk2 | grep \"Disk Size\" | awk '{print $3, $4, \"(\"substr($5, 2, 16)/1000/1000/1000, \"MB)\"}'",
	)

	// Boring Go error handling
	fusionInfoBytes, err := hddFusionCheck.Output()
	if err != nil {
		fmt.Println("Error fusion! ", err)
	}
	hddInfo, err := hddDisk0.Output()
	if err != nil {
		fmt.Println("Error disk0!", err)
	}
	hddInfo2, err := hddDisk2.Output()
	if err != nil {
		fmt.Println("Error disk2!", err)
	}

	// If "Fusion Drive" is found in the output of hdd_fusion_check, then print the disk size of /dev/disk2 (the default allocation
	//  of a fusion drive), otherwise print the disk size of /dev/disk0
	fusionInfoString := string(fusionInfoBytes[:])
	if strings.Contains(fusionInfoString, "Fusion Drive:              Yes") {
		m.command = bodyStyle.Render(
			"\nDEVICE IS USING FUSION DRIVE\n/dev/disk2: ",
		) + commandStyle.Render(
			string(hddInfo2[:]),
		)
	} else {
		//    SATA SSD = Slow SSD
		//    PCIE SSD aka NVMe SSD = Fast SSD
		hddInfoString := string(hddInfo[:])
		m.command = bodyStyle.Render("\n   /dev/disk0: ") + commandStyle.Render(hddInfoString)

		hddIsSolidstateCmd := exec.Command("bash", "-c", "diskutil info /dev/disk0 | grep Solid")
		hddIsSolidstate, err := hddIsSolidstateCmd.Output()
		// Boring Go error handling
		if err != nil {
			fmt.Println("Error determing if drive is solid state: ", err)
		}
		isSolidstate := string(hddIsSolidstate[:])

		// No newline required as calling .Render on a string(<cmd_output>) will automatically add a newline at the
		//   end of itself
		if strings.Contains(isSolidstate, "Yes") {
			m.command += "\n" + bodyStyle.Render("Is SSD: ") + commandStyle.Render("Yes")
		} else {
			m.command += "\n" + bodyStyle.Render("Is SSD: ") + commandStyle.Render("No")
		}

		hddProtocolCmd := exec.Command("bash", "-c", "diskutil info /dev/disk0 | grep Protocol | awk '{print $2}'")
		hddProtocol, err := hddProtocolCmd.Output()
		// Boring Go error handling
		if err != nil {
			fmt.Println("Error determing drive protocol: ", err)
		}
		hddProtocolString := string(hddProtocol[:])
		if strings.Contains(hddProtocolString, "Apple") {
			hddProtocolString = "Apple Fabric"
		}
		m.command += "\n" + bodyStyle.Render("Protocol: ") + commandStyle.Render(hddProtocolString)

	}

	// **IGNORE**
	// m.command = bodyStyle.Render("\n/dev/disk0: ") + commandStyle.Render(string(hdd_info[:])) + bodyStyle.Render("\n/dev/disk1: ")+ commandStyle.Render(string(hdd_info2[:])) + bodyStyle.Render("\n/dev/disk2: ") + commandStyle.Render(string(hdd_info3[:]))
	m.commandType = "HDD"
	return m, nil
}

func getGPU(m model) (tea.Model, tea.Cmd) {
	if checkM1() {
		m.command, m.commandType = "M1 GPU", "GPU"
		return m, nil
	}
	gpuCmd := exec.Command("bash", "-c", "ioreg -rc IOPCIDevice | grep \"model\" | sed -n '1 p'")
	gpuInfo, err := gpuCmd.Output()

	if checkError("Error Getting GPU: ", err) {
		return m, nil
	}

	commandRaw := string(gpuInfo[:])
	m.command = strings.Split(strings.Split(commandRaw, `<"`)[1], `">`)[0]

	c2 := exec.Command(
		"bash",
		"-c",
		"ioreg -rc IOPCIDevice | grep VRAM | awk '{print $5/1024\"GB\"}'",
	)
	gpuVramInfo, err := c2.Output()

	// Boring Go error handling
	if err != nil || len(string(gpuVramInfo[:])) == 0 {
		m.command += " (VRAM Unknown)"
	}

	command2Raw := string(gpuVramInfo)
	m.command += " " + command2Raw

	// m.command = command_splits[0]
	m.commandType = "GPU"
	return m, nil
}

func getSerial(m model) (tea.Model, tea.Cmd) {
	c := exec.Command("bash", "-c", "ioreg -w0 -l | grep PlatformSerial | awk '{print $4}'")
	serialNo, err := c.Output()
	if err != nil {
		fmt.Println("Error Serial Number:", err)
		return m, nil
	}
	m.command = string(serialNo[:])
	m.commandType = "SERIAL"
	return m, nil
}

func getModel(m model) (tea.Model, tea.Cmd) {
	c := exec.Command("bash", "-c", "ioreg -rk model | grep model | grep ',' | awk '{print $4}'")
	modelNo, err := c.Output()
	if err != nil {
		fmt.Println("Error Model Number:", err)
		return m, nil
	}

	m.command = string(modelNo[:])
	m.commandType = "MODEL"
	return m, nil
}

func getGPUCore(m model) (tea.Model, tea.Cmd) {
	if checkM1() == false {
		m.command = "This machine doesn't have an M1 chip"
		m.commandType = "NOTM1"
		return m, nil
	}
	gpuCoreCmd := exec.Command(
		"bash",
		"-c",
		"ioreg -rc IOGPU | grep core-count | awk '{print $4}'",
	)
	gpuCoreInfo, err := gpuCoreCmd.Output()

	if checkError("Error Getting GPU Cores: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(gpuCoreInfo[:]), "GPUCORE"
	return m, nil
}

func getRAM(m model) (tea.Model, tea.Cmd) {
	c := exec.Command("bash", "-c", "sysctl hw.memsize | awk '{print $2/1024/1024/1024 \"GB\"}'")
	ramInfo, err := c.Output()

	if checkError("Error Getting RAM:", err) {
		return m, nil
	}

	m.command, m.commandType = (string(ramInfo[:])), "RAM"
	return m, nil
}

func getBattery(m model) (tea.Model, tea.Cmd) {
	batteryCountCmd := exec.Command(
		"bash",
		"-c",
		"chroot /Volumes/Macintosh\\ HD system_profiler SPPowerDataType | grep Count",
	)
	batteryCountInfo, err := batteryCountCmd.Output()
	m.commandType = "BATTERY"

	if checkError("Error getting battery count: ", err) {
		return m, nil
	}

	m.command = string(batteryCountInfo[:])

	batteryConditionCmd := exec.Command(
		"bash",
		"-c",
		"chroot /Volumes/Macintosh\\ HD system_profiler SPPowerDataType | grep Condition",
	)
	batteryConditionInfo, err := batteryConditionCmd.Output()

	if checkError("Error getting battery condition: ", err) {
		return m, nil
	}

	m.command += "\n" + string(batteryConditionInfo[:])
	return m, nil
}

func getCPU(m model) (tea.Model, tea.Cmd) {
	cpuCmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
	cpuInfo, err := cpuCmd.Output()

	if checkError("Error Getting CPU: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(cpuInfo[:]), "CPU"
	return m, nil
}

func getCPUCore(m model) (tea.Model, tea.Cmd) {
	if !checkM1() {
		m.command = "This machine doesn't have an M1 chip"
		m.commandType = "NOTM1"
		return m, nil
	}
	cpuCoreCmd := exec.Command("sysctl", "-n", "machdep.cpu.core_count")
	cpuCoreInfo, err := cpuCoreCmd.Output()

	if checkError("Error getting CPU core count: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(cpuCoreInfo[:]), "CPUCORE"
	return m, nil
}

func getActivationStatus(m model) (tea.Model, tea.Cmd) {
	activationStatusCmd := exec.Command(
		"bash",
		"-c",
		"chroot /Volumes/Macintosh\\ HD system_profiler SPHardwareDataType | sed -n '20 p'",
	)
	activationStatusInfo, err := activationStatusCmd.Output()

	if checkError("Error getting activation status", err) {
		return m, nil
	}

	m.command, m.commandType = string(activationStatusInfo[:]), "ACTIVATION_STATUS"
	if m.command == "" {
		m.command = "No Activation Label found in SPHardwareDataType"
	}
	return m, nil
}

func disableAuthenticatedRoot(m model) (tea.Model, tea.Cmd) {
	osVersion := checkMacOSVersion()
	if osVersion == "Ventura" {
		m.command = "os confirmed to be ventura.\n if you're seeing this `m.command` failed to update later on in the function"
		authRootDisable := exec.Command("bash", "-c", "csrutil authenticated-root disable")
		authRootDisableGo, err := authRootDisable.Output()
		if err != nil {
			fmt.Println("Error disabling authenticated-root", err, authRootDisableGo)
			return m, nil
		}
		m.command = "os confirmed to be ventura. auth-root is disabled, trying csrutil"
		csrutilDisable := exec.Command("bash", "-c", "csrutil disable")
		csrutilGo, err := csrutilDisable.Output()
		if err != nil {
			fmt.Println("Error disabling authenticated-root", err, csrutilGo)
			return m, nil
		}
		if err != nil {
			fmt.Println("Error disabling authenticated-root", err)
			return m, nil
		}
		m.command = "Authenticated root disabled. Reboot back into USB\n"
		m.commandType = "AUTH_ROOT"
		return m, nil
	}
	m.command = osVersion + "is not ventura. Disabling authenticated root not necessary"
	m.commandType = "AUTH_ROOT"
	return m, nil
}

func formatDrive(m model, fs string) (tea.Model, tea.Cmd) {
	if fs == "APFS" {
		c := exec.Command("bash", "-c", "diskutil erasedisk APFS \"Macintosh HD\" /dev/disk0")
		m.command = "APFS"
		err := c.Run()
		if err != nil {
			fmt.Println("Error formating drive to APFS: ", err)
			return m, nil
		}
	} else if fs == "JHFS+" {
		c := exec.Command("bash", "-c", "diskutil erasedisk JHFS+ \"Macintosh HD\" /dev/disk0")
		m.command = "JHFS+"
		err := c.Run()
		if err != nil {
			fmt.Println("Error formatting drive to JHFS+: ", err)
			return m, nil
		}
	} else {
		c := exec.Command("diskutil", "resetfusion")
		m.command = "FUSION"
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			fmt.Println("Error formatting fusion drive: ", err)
			return m, nil
		}
	}
	m.commandType = "FORMAT"
	return m, nil
}

func installOS(m model) (tea.Model, tea.Cmd) {
	var startOSInstall string
	// Fuck off you try finding a better way
	if checkDir("/", "Install macOS Catalina.app") {
		startOSInstall = `/Install\ macOS\ Catalina.app/Contents/Resources/startosinstall`
	} else if checkDir("/", "Install macOS Big Sur.app") {
		startOSInstall = `/Install\ macOS\ Big\ Sur.app/Contents/Resources/startosinstall`
	} else if checkDir("/", "Install macOS Monterey.app") {
		startOSInstall = `/Install\ macOS\ Monterey.app/Contents/Resources/startosinstall`
	}
	c := exec.Command(startOSInstall, "--agreetolicense", "--volume", `/Volumes/Macintosh\ HD/`)
	err := c.Run()
	if err != nil {
		fmt.Println("Whoops! All fucked os!", err)
	}
	installOsInfo, err := c.Output()
	if err != nil {
		fmt.Println("Error Installing OS:", err)

		return m, nil
	}

	m.command = string(installOsInfo[:])
	m.commandType = "OS"

	return m, nil
}

func demoBatteryNew(m model) (tea.Model, tea.Cmd) {
	c := exec.Command(
		"bash",
		"-c",
		"ioreg -rk LegacyBatteryInfo | grep LegacyBatteryInfo | awk '{print $3}'",
	)
	demoBatteryRaw, err := c.Output()
	if err != nil {
		fmt.Println("Demo battery fuck up")
		return m, nil
	}

	m.command = string(demoBatteryRaw[:])
	m.commandType = "DEMO"
	return m, nil
}

func isFusionDrive(m model) (tea.Model, tea.Cmd) {
	if checkFusion() {
		m.command = "No fusion drive"
	} else {
		m.command = "Has fusion drive"
	}
	m.commandType = "FUSION"
	return m, nil
}

func testMenuHDDWriteTest(m model) (tea.Model, tea.Cmd) {
	writeTestCmd := exec.Command(
		"bash",
		"-c",
		"dd if=/dev/zero of=/Volumes/Macintosh HD/tstFile bs=512k count=5000 | grep sec | awk '{print $1/ 1024/ 1024/ $5\"MB/sec\"}'",
	)
	writeTestInfo, err := writeTestCmd.Output()
	fmt.Println("HDD Write Test engaged")

	if checkError("Error executing write test: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(writeTestInfo[:]), "TEST_WRITE"
	fmt.Println("HDD Write Test completed")
	return m, nil
}

func testMenuHDDReadTest(m model) (tea.Model, tea.Cmd) {
	readTestCmd := exec.Command(
		"bash",
		"-c",
		"dd if=/Volumes/Macintosh HD/tstfile of=/dev/null bs=512k count=5000 | grep sec | awk '{print $1/ 1024/ 1024 $5\"MB/sec\"}'",
	)
	readTestInfo, err := readTestCmd.Output()
	fmt.Println("HDD Read Test engaged")

	if checkError("Error executing read test: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(readTestInfo[:]), "TEST_READ"
	fmt.Println("HDD Read Test completed")
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		// ACTIVATION STATUS
		case "a":
			m.commandPresent = true
			m.disclaimerShow = false
			return getActivationStatus(m)
		// DISABLE AUTHENTICATED ROOT
		case "R":
			m.commandPresent = true
			m.disclaimerShow = false
			return disableAuthenticatedRoot(m)
		// BATTERY
		case "B":
			m.commandPresent = true
			m.disclaimerShow = false
			return getBattery(m)
		// CPU
		case "c":
			m.commandPresent = true
			m.disclaimerShow = false
			return getCPU(m)
		// M1 CPU CORE
		case "C":
			m.commandPresent = true
			m.disclaimerShow = false
			return getCPUCore(m)
		// FUSION DRIVE CHECK
		case "f":
			m.commandPresent = true
			m.disclaimerShow = false
			return isFusionDrive(m)
		case "F":
			m.commandPresent = true
			m.disclaimerShow = false
			return demoBatteryNew(m)
		// GPU
		case "g":
			m.commandPresent = true
			m.disclaimerShow = false
			return getGPU(m)
		// M1 GPU CORE
		case "G":
			m.commandPresent = true
			m.disclaimerShow = false
			return getGPUCore(m)
		// HDD WRITE TEST
		case "h":
			m.commandPresent = true
			m.disclaimerShow = false
			if m.testMenu {
				m.testMenu = !m.testMenu
				return testMenuHDDWriteTest(m)
			}
			return getHDD(m)
		// GET MODEL NUMBER
		case "m":
			m.commandPresent = true
			m.disclaimerShow = false
			return getModel(m)
		// INSTALL OS
		case "o":
			m.commandPresent = true
			m.disclaimerShow = false
			return (installOS(m))
		// RAM
		case "r":
			m.commandPresent = true
			m.disclaimerShow = false
			if m.testMenu {
				m.testMenu = !m.testMenu
				return testMenuHDDReadTest(m)
			}
			return getRAM(m)
		// SERIAL
		case "s":
			m.commandPresent = true
			m.disclaimerShow = false
			if m.testMenu {
				m.testMenu = !m.testMenu
				return m, nil
			}
			return getSerial(m)
		// TEST
		case "t":
			m.commandPresent = true
			m.disclaimerShow = false
			m.commandType = "TEST"
			m.testMenu = true
			return m, nil
		// Format Option 1: APFS
		case "1":
			m.commandPresent = true
			m.disclaimerShow = false
			return formatDrive(m, "APFS")
		// Format Option 2: JHFS+
		case "2":
			m.commandPresent = true
			m.disclaimerShow = false
			return formatDrive(m, "JHFS+")
		// Format Option 3: Fusion
		case "3":
			m.commandPresent = true
			m.disclaimerShow = false
			return formatDrive(m, "Fusion")
		// I have no idea
		case "b":
			m.commandPresent = false
			m.disclaimerShow = false
			return m, nil
		// Show Altscreen
		case "A":
			m.disclaimerShow = false
			m.altscreen = !m.altscreen
			cmd := tea.EnterAltScreen
			if !m.altscreen {
				cmd = tea.ExitAltScreen
			}
			return m, cmd
		// I have no idea
		case "x":
			m.disclaimerShow = false
			return m, nil
		}
	case commandFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var renderString, titleString string
	var displayOptions bool
	displayOptions = true
	renderString = ""
	titleString = ""
	if m.disclaimerShow {
		titleString += disclaimerStyle.Render("Welcome to `specget` 1.1!") + "\n"
	}
	if !m.commandPresent {
		renderString += titleStyle.Render("Please enter a command")

		renderString += "\n" + bodyStyle.Render(
			"'a' for ",
		) + commandStyle.Render(
			"activation status ",
		) + warningStyle.Render(
			"OS Install Req.",
		)
		renderString += "\n" + bodyStyle.Render(
			"'B' for ",
		) + commandStyle.Render(
			"battery info ",
		) + warningStyle.Render(
			"Macbook Only. OS Install Req.",
		)
		renderString += "\n" + bodyStyle.Render(
			"'c' for ",
		) + commandStyle.Render(
			"cpu",
		)
		renderString += "\n" + bodyStyle.Render(
			"'C' for ",
		) + commandStyle.Render(
			"cpu core ",
		) + warningStyle.Render(
			"M1 Only",
		)
		renderString += "\n" + bodyStyle.Render("'r' for ") + commandStyle.Render("ram")
		renderString += "\n" + bodyStyle.Render(
			"'R' for ",
		) + commandStyle.Render(
			"disable authenticated root ",
		) + warningStyle.Render(
			"Ventura Only",
		)
		renderString += "\n" + bodyStyle.Render(
			"'g' for ",
		) + commandStyle.Render(
			"gpu",
		)
		renderString += "\n" + bodyStyle.Render(
			"'G' for ",
		) + commandStyle.Render(
			"gpu core ",
		) + warningStyle.Render(
			"M1 Only",
		)
		renderString += "\n" + bodyStyle.Render(
			"'f' for ",
		) + commandStyle.Render(
			"fusion drive test",
		)
		renderString += "\n" + bodyStyle.Render(
			"'h' for ",
		) + commandStyle.Render(
			"hdd",
		)
		renderString += "\n" + bodyStyle.Render(
			"'s' for ",
		) + commandStyle.Render(
			"serial number",
		)
		renderString += "\n" + bodyStyle.Render(
			"'m' for ",
		) + commandStyle.Render(
			"model reference",
		)
		// renderString += "\n" + bodyStyle.Render("'o' for ") + commandStyle.Render("os install ") + warningStyle.Render("Coming Soon")
		// renderString += "\n" + bodyStyle.Render("'p' for ") + commandStyle.Render("ping test")+ warningStyle.Render("NOT IMPLEMENTED")
		// renderString += "\n" + bodyStyle.Render("'f' for ") + commandStyle.Render("Format Menu")
		// renderString += "\n" + bodyStyle.Render("'t' for ") + commandStyle.Render("Test Menu ") + warningStyle.Render("Coming Not Soon")
		renderString += "\n" + quitStyle.Render("'q' to quit")
	} else {
		renderString = "\n"
		switch m.commandType {
		case "GPU":
			renderString += bodyStyle.Render("GPU is: ") + commandStyle.Render(m.command)
		case "GPUCORE":
			renderString += bodyStyle.Render("GPU Core Count is: ") + commandStyle.Render(m.command)
		case "RAM":
			renderString += bodyStyle.Render("RAM is: ") + commandStyle.Render(m.command)
		case "CPU":
			renderString += bodyStyle.Render("CPU is: ") + commandStyle.Render(m.command)
		case "CPUCORE":
			renderString += bodyStyle.Render("CPU Core Count is: ") + commandStyle.Render(m.command)
		case "BATTERY":
			renderString += commandStyle.Render(m.command)
		case "HDD":
			renderString += bodyStyle.Render("Disk Size: ") + commandStyle.Render(m.command)
		case "SERIAL":
			renderString += bodyStyle.Render("Serial Number: ") + commandStyle.Render(m.command)
		case "FUSION":
			renderString += bodyStyle.Render("Fusion Drive Status: ") + commandStyle.Render(m.command)
		case "MODEL":
			renderString += commandStyle.Render(m.command)
		case "DEMO":
			renderString += commandStyle.Render(m.command)
		case "OS":
			renderString += m.command + "\n\n"
			displayOptions = false
		case "AUTH_ROOT":
			renderString += commandStyle.Render(m.command)
		case "ACTIVATION_STATUS":
			renderString += commandStyle.Render(m.command)
		case "NOTM1":
			renderString += commandStyle.Render(m.command)
		case "PING":
			renderString += bodyStyle.Render("Ping results: \n") + commandStyle.Render(m.command)
		case "TEST":
			renderString += bodyStyle.Render("Yo sorry B but this hasn't been implemented yet. Look out for these tests in the future")
			renderString += "\n" + bodyStyle.Render("'?' for ") + commandStyle.Render("Harddrive read test")

			renderString += "\n" + bodyStyle.Render("'?' for ") + commandStyle.Render("Harddrive write test")
			renderString += "\n" + bodyStyle.Render("'?' for ") + commandStyle.Render("CPU Stress test")
			renderString += "\n" + bodyStyle.Render("'?' for ") + commandStyle.Render("RAM test")
		case "FORMAT":
			renderString += warningStyle.Render("NOT IMPLEMENTED") + "\n" + bodyStyle.Render("Assumed drive is /dev/disk0, fusion assumes /dev/disk0 & /dev/disk1")
			renderString += "\n" + bodyStyle.Render("Use APFS Format for modern OS versions. Use JHFS+ for Catalina & below")
			renderString += "\n" + bodyStyle.Render("'?' for ") + commandStyle.Render("APFS Format")
			renderString += "\n" + bodyStyle.Render("'?' for ") + commandStyle.Render("JHFS+ Format")
			renderString += "\n" + bodyStyle.Render("'?' for ") + commandStyle.Render("Fusion Drive Format")
		case "TEST_READ":
			renderString += "\n" + bodyStyle.Render("Read Test:") + commandStyle.Render(m.command)
		case "TEST_WRITE":
			renderString += "\n" + bodyStyle.Render("Write Test:") + commandStyle.Render(m.command)
		case "OLDOS":
			renderString += "\n" + bodyStyle.Render("The OS you're booting off of is too old for this command to work.") + "\n" + bodyStyle.Render("Please use ") + commandStyle.Render("Big Sur") + bodyStyle.Render(" or later to use this feature")
		default:
			renderString += bodyStyle.Render("Fucking uhhhhhh") + commandStyle.Render(" Idk B")
		}
		m.commandPresent = false
		if displayOptions {
			renderString += "\n" + quitStyle.Render("[c]pu | [g]pu | [r]am | [h]dd | [b]ack to menu | [q]uit")
		}
	}
	/*
	  if m.command.Contains("GB") {
	    m.commandPresent = !m.commandPresent
	  } else if m.command.Contains("Intel") {
	    renderString = "\n"
	    m.commandPresent = !m.command
	  }
	*/

	if m.err != nil {
		return "Error with UI for some reason: " + m.err.Error() + "\n"
	}
	outputString := titleString + borderStyle.Render(renderString)
	titleString, renderString = "", ""
	return outputString
}

func main() {
	m := model{
		commandPresent: false,
		disclaimerShow: true,
		formatMenu:     false,
		testMenu:       false,
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error! ", err)
		os.Exit(1)
	}
}
