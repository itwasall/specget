package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
		//Background(lipgloss.Color("#505059")).
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

func checkError(error_string string, err error) bool {
	if err != nil {
		fmt.Println(error_string, err)
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

	cpu_name_out, err := c.Output()

	if err != nil {
		fmt.Println("Error getting CPU: ", err)
		return false
	}

	cpu_name := string(cpu_name_out[:])

	if strings.Contains(cpu_name, " M1") {
		return true
	}
	return false
}

func checkFusion() bool {
	c := exec.Command("bash", "-c", "diskutil resetfusion | cat | awk '{print $1}' | sed -n '1 p'")

	fusion_status_out, err := c.Output()

	if err != nil {
		fmt.Println("Error getting fusion drive status: ", err)
	}

	fusion_status := string(fusion_status_out[:])

	if strings.Contains(fusion_status, "s internal disk devices must be solid-state") {
		return false
	}
	return true
}

func getHDD(m model) (tea.Model, tea.Cmd) {
	hdd_fusion_check := exec.Command("bash", "-c", "diskutil info /dev/disk2")
	// dev/disk0 being the default drive macOS is installed on
	hdd_disk0 := exec.Command("bash", "-c", "diskutil info /dev/disk0 | grep \"Disk Size\" | awk '{print $3, $4, \"(\"substr($5, 2, 16)/1000/1000/1000, \"MB)\"}'")
	// dev/disk2 being the default assignment a fusion drive gets. (Comprising of a disk0 SSD and a disk1 HDD for example)
	hdd_disk2 := exec.Command("bash", "-c", "diskutil info /dev/disk2 | grep \"Disk Size\" | awk '{print $3, $4, \"(\"substr($5, 2, 16)/1000/1000/1000, \"MB)\"}'")

	// Boring Go error handling
	fusion_info_bytes, err := hdd_fusion_check.Output()
	if err != nil {
		// Currently doesn't execute as it'll trigger on non-fusion machines
		// fmt.Println("Error fusion! ", err)
	}
	hdd_info, err := hdd_disk0.Output()
	if err != nil {
		fmt.Println("Error disk0!", err)
	}
	hdd_info2, err := hdd_disk2.Output()
	if err != nil {
		fmt.Println("Error disk2!", err)
	}

	// If "Fusion Drive" is found in the output of hdd_fusion_check, then print the disk size of /dev/disk2 (the default allocation
	//  of a fusion drive), otherwise print the disk size of /dev/disk0
	fusion_info_string := string(fusion_info_bytes[:])
	if strings.Contains(fusion_info_string, "Fusion Drive:              Yes") {
		m.command = bodyStyle.Render("\nDEVICE IS USING FUSION DRIVE\n/dev/disk2: ") + commandStyle.Render(string(hdd_info2[:]))
	} else {
		// Otherwise if the drive isn't a fusion drive then whether or not it's an SSD, and what kind of protocol the drive is using
		//  is found and displayed. The latter is important because a SATA SSD is weaker than a PCI-Expess (PCI-E) SSD. Just to confuse
		//  matters further, NVMe drives are PCI-E. No we haven't used enough of the alphabet to convolude this yet. We haven't even
		//  included M.2!

		// too long;didn't give a shit:
		//    SATA SSD = Slow SSD
		//    PCIE SSD aka NVMe SSD = Fast SSD
		hdd_info_string := string(hdd_info[:])
		m.command = bodyStyle.Render("\n   /dev/disk0: ") + commandStyle.Render(hdd_info_string)

		hdd_is_solidstate_cmd := exec.Command("bash", "-c", "diskutil info /dev/disk0 | grep Solid")
		hdd_is_solidstate, err := hdd_is_solidstate_cmd.Output()

		// Boring Go error handling
		if err != nil {
			fmt.Println("Error determing if drive is solid state: ", err)
		}
		is_solidstate := string(hdd_is_solidstate[:])

		// No newline required as calling .Render on a string(<cmd_output>) will automatically add a newline at the
		//   end of itself
		if strings.Contains(is_solidstate, "Yes") {
			m.command += "\n" + bodyStyle.Render("Is SSD: ") + commandStyle.Render("Yes")
		} else {
			m.command += "\n" + bodyStyle.Render("Is SSD: ") + commandStyle.Render("No")
		}

		hdd_protocol_cmd := exec.Command("bash", "-c", "diskutil info /dev/disk0 | grep Protocol | awk '{print $2}'")
		hdd_protocol, err := hdd_protocol_cmd.Output()
		// Boring Go error handling
		if err != nil {
			fmt.Println("Error determing drive protocol: ", err)
		}
		hdd_protocol_string := string(hdd_protocol[:])
		if strings.Contains(hdd_protocol_string, "Apple") {
			hdd_protocol_string = "Apple Fabric"
		}
		m.command += "\n" + bodyStyle.Render("Protocol: ") + commandStyle.Render(hdd_protocol_string)

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
	gpu_cmd := exec.Command("bash", "-c", "ioreg -rc IOPCIDevice | grep \"model\" | sed -n '1 p'")
	gpu_info, err := gpu_cmd.Output()

	if checkError("Error Getting GPU: ", err) {
		return m, nil
	}

	command_raw := string(gpu_info[:])
	m.command = strings.Split(strings.Split(command_raw, `<"`)[1], `">`)[0]

	c2 := exec.Command("bash", "-c", "ioreg -rc IOPCIDevice | grep VRAM | awk '{print $5/1024\"GB\"}'")
	gpuvram_info, err := c2.Output()

	// Boring Go error handling
	if err != nil || len(string(gpuvram_info[:])) == 0 {
		m.command += " (VRAM Unknown)"
	}

	command2_raw := string(gpuvram_info)
	m.command += " " + command2_raw

	// m.command = command_splits[0]
	m.commandType = "GPU"
	return m, nil
}

func getSerial(m model) (tea.Model, tea.Cmd) {
	c := exec.Command("bash", "-c", "ioreg -w0 -l | grep PlatformSerial | awk '{print $4}'")
	serial_no, err := c.Output()

	if err != nil {
		fmt.Println("Error Serial Number:", err)
		return m, nil
	}
	m.command = string(serial_no[:])
	m.commandType = "SERIAL"
	return m, nil
}

func getModel(m model) (tea.Model, tea.Cmd) {
	c := exec.Command("bash", "-c", "ioreg -rk model | grep model | grep ',' | awk '{print $4}'")
	model_no, err := c.Output()

	if err != nil {
		fmt.Println("Error Model Number:", err)
		return m, nil
	}

	m.command = string(model_no[:])
	m.commandType = "MODEL"
	return m, nil
}

func getGPUCore(m model) (tea.Model, tea.Cmd) {
	if checkM1() == false {
		m.command = "This machine doesn't have an M1 chip"
		m.commandType = "NOTM1"
		return m, nil
	}
	gpu_core_cmd := exec.Command("bash", "-c", "ioreg -rc IOGPU | grep core-count | awk '{print $4}'")
	gpu_core_info, err := gpu_core_cmd.Output()

	if checkError("Error Getting GPU Cores: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(gpu_core_info[:]), "GPUCORE"
	return m, nil
}

func getRAM(m model) (tea.Model, tea.Cmd) {
	c := exec.Command("bash", "-c", "sysctl hw.memsize | awk '{print $2/1024/1024/1024 \"GB\"}'")
	ram_info, err := c.Output()

	if checkError("Error Getting RAM:", err) {
		return m, nil
	}

	m.command, m.commandType = (string(ram_info[:])), "RAM"
	return m, nil
}

func getBattery(m model) (tea.Model, tea.Cmd) {
	battery_count_cmd := exec.Command("bash", "-c", "chroot /Volumes/Macintosh\\ HD system_profiler SPPowerDataType | grep Count")
	battery_count_info, err := battery_count_cmd.Output()
	m.commandType = "BATTERY"

	if checkError("Error getting battery count: ", err) {
		return m, nil
	}

	m.command = string(battery_count_info[:])

	battery_condition_cmd := exec.Command("bash", "-c", "chroot /Volumes/Macintosh\\ HD system_profiler SPPowerDataType | grep Condition")
	battery_condition_info, err := battery_condition_cmd.Output()

	if checkError("Error getting battery condition: ", err) {
		return m, nil
	}

	m.command += "\n" + string(battery_condition_info[:])
	return m, nil
}

func getCPU(m model) (tea.Model, tea.Cmd) {
	cpu_cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
	cpu_info, err := cpu_cmd.Output()

	if checkError("Error Getting CPU: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(cpu_info[:]), "CPU"
	return m, nil
}

func getCPUCore(m model) (tea.Model, tea.Cmd) {
	if checkM1() != true {
		m.command = "This machine doesn't have an M1 chip"
		m.commandType = "NOTM1"
		return m, nil
	}
	cpu_core_cmd := exec.Command("sysctl", "-n", "machdep.cpu.core_count")
	cpu_core_info, err := cpu_core_cmd.Output()

	if checkError("Error getting CPU core count: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(cpu_core_info[:]), "CPUCORE"
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
	c.Run()
	install_os_info, err := c.Output()

	if err != nil {
		fmt.Println("Error Installing OS:", err)

		return m, nil
	}

	m.command = string(install_os_info[:])
	m.commandType = "OS"

	return m, nil
}

func demoBatteryNew(m model) (tea.Model, tea.Cmd) {
	c := exec.Command("bash", "-c", "ioreg -rk LegacyBatteryInfo | grep LegacyBatteryInfo | awk '{print $3}'")
	demoBattery_raw, err := c.Output()

	if err != nil {
		fmt.Println("Demo battery fuck up")
		return m, nil
	}

	m.command = string(demoBattery_raw[:])
	m.commandType = "DEMO"
	return m, nil
}
func isFusionDrive(m model) (tea.Model, tea.Cmd) {
	if checkFusion() == true {
		m.command = "No fusion drive"
	} else {
		m.command = "Has fusion drive"
	}
	m.commandType = "FUSION"
	return m, nil
}

func testMenu_HDDWriteTest(m model) (tea.Model, tea.Cmd) {
	write_test_cmd := exec.Command("bash", "-c", "dd if=/dev/zero of=/Volumes/Macintosh HD/tstFile bs=512k count=5000 | grep sec | awk '{print $1/ 1024/ 1024/ $5\"MB/sec\"}'")
	write_test_info, err := write_test_cmd.Output()
	fmt.Println("HDD Write Test engaged")

	if checkError("Error executing write test: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(write_test_info[:]), "TEST_WRITE"
	fmt.Println("HDD Write Test completed")
	return m, nil
}

func testMenu_HDDReadTest(m model) (tea.Model, tea.Cmd) {
	read_test_cmd := exec.Command("bash", "-c", "dd if=/Volumes/Macintosh HD/tstfile of=/dev/null bs=512k count=5000 | grep sec | awk '{print $1/ 1024/ 1024 $5\"MB/sec\"}'")
	read_test_info, err := read_test_cmd.Output()
	fmt.Println("HDD Read Test engaged")

	if checkError("Error executing read test: ", err) {
		return m, nil
	}

	m.command, m.commandType = string(read_test_info[:]), "TEST_READ"
	fmt.Println("HDD Read Test completed")
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
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
				return testMenu_HDDWriteTest(m)
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
				return testMenu_HDDReadTest(m)
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
		case "a":
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

		renderString += "\n" + bodyStyle.Render("'B' for ") + commandStyle.Render("battery info ") + warningStyle.Render("Macbook Only. OS Install Req.")
		renderString += "\n" + bodyStyle.Render("'c' for ") + commandStyle.Render("cpu")
		renderString += "\n" + bodyStyle.Render("'C' for ") + commandStyle.Render("cpu core ") + warningStyle.Render("M1 Only")
		renderString += "\n" + bodyStyle.Render("'r' for ") + commandStyle.Render("ram")
		renderString += "\n" + bodyStyle.Render("'g' for ") + commandStyle.Render("gpu")
		renderString += "\n" + bodyStyle.Render("'G' for ") + commandStyle.Render("gpu core ") + warningStyle.Render("M1 Only")
		renderString += "\n" + bodyStyle.Render("'f' for ") + commandStyle.Render("fusion drive test")
		renderString += "\n" + bodyStyle.Render("'h' for ") + commandStyle.Render("hdd")
		renderString += "\n" + bodyStyle.Render("'s' for ") + commandStyle.Render("serial number")
		renderString += "\n" + bodyStyle.Render("'m' for ") + commandStyle.Render("model reference")
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
		m.commandPresent = !m.commandPresent
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
	m := model{commandPresent: false, disclaimerShow: true, formatMenu: false, testMenu: false}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error! ", err)
		os.Exit(1)
	}
}
