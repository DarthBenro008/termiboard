package memory

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"runtime"
	"sort"
	"termiboard/utils"
)

func getReadableSize(sizeInBytes uint64) (readableSizeString string) {
	var (
		units = []string{"B", "KB", "MB", "GB", "TB", "PB"}
		size  = float64(sizeInBytes)
		i     = 0
	)
	for ; i < len(units) && size >= 1024; i++ {
		size = size / 1024
	}
	readableSizeString = fmt.Sprintf("%.2f %s", size, units[i])
	return
}

func GetRamUsage() {
	//Implement the function to fetch ram usage - Total Ram used , free ram available and used percentage
	m, err := mem.VirtualMemory()
	if err != nil {
		utils.StandardPrinter(utils.ErrorRedColor, "Could not retrieve RAM details.")
		utils.PrintVerbosePanic(err) //Exit upon error, below code must not be executed
	}
	usedMessage := fmt.Sprintf(
		"%s (%.2f%%)",
		getReadableSize(m.Used),
		m.UsedPercent,
	)
	utils.ResultPrinter("Ram Total: ", getReadableSize(m.Total))
	utils.ResultPrinter("Ram Used: ", usedMessage)
	utils.ResultPrinter("Ram Available: ", getReadableSize(m.Available))
	utils.ResultPrinter("Ram Free: ", getReadableSize(m.Free))

}

// GetTopProcesses print out the top 5 process that are consuming most RAM
func GetTopProcesses() {
	strOutput := ""
	processes, err := process.Processes()
	if err != nil {
		utils.StandardPrinter(utils.ErrorRedColor, "Could not retrieve running process list.")
		utils.PrintVerbosePanic(err)
	}

	sort.Slice(processes, func(i, j int) bool {
		memoryPercentOfIthProcess, err := processes[i].MemoryPercent()
		if err != nil {
			utils.StandardPrinter(utils.ErrorRedColor, "Could not retrieve memory usage details.")
			utils.PrintVerbosePanic(err)
		}
		memoryPercentOfJthProcess, err := processes[j].MemoryPercent()
		if err != nil {
			utils.StandardPrinter(utils.ErrorRedColor, "Could not retrieve memory usage details.")
			utils.PrintVerbosePanic(err)
		}
		return memoryPercentOfIthProcess > memoryPercentOfJthProcess
	})

	//Windows Compatibility(Leave out the System Idle Process), for more info refer #30

	if runtime.GOOS == "windows" {
		if processes[0].Pid == 0 {
			processes = processes[1:]
		}
	}

	sort.Slice(processes, func(i, j int) bool {
		memoryPercentOfIthProcess, err := processes[i].MemoryPercent()
		if err != nil {
			utils.StandardPrinter(utils.ErrorRedColor, "Process memory usage access denied.") //More descriptive error
			utils.PrintVerbosePanic(err)
		}
		memoryPercentOfJthProcess, err := processes[j].MemoryPercent()
		if err != nil {
			utils.StandardPrinter(utils.ErrorRedColor, "Process memory usage access denied.") //More descriptive error
			utils.PrintVerbosePanic(err)
		}
		return memoryPercentOfIthProcess > memoryPercentOfJthProcess
	})

	for i := 0; i < 5; i++ {
		memoryPercentOfIthProcess, err := processes[i].MemoryPercent()
		if err != nil {
			utils.StandardPrinter(utils.ErrorRedColor, "Process memory usage access denied.") //More descriptive error
			utils.PrintVerbosePanic(err)
		}
		strOutput += fmt.Sprintf("PID: %5d, memory %%: %2.1f\n", processes[i].Pid, memoryPercentOfIthProcess)
	}
	utils.ResultPrinter("Top 5 processes by memory usage: \n", strOutput)
}

func GetDiskUsage() {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		utils.StandardPrinter(utils.ErrorRedColor, "Could not retrieve disk usage details.")
		utils.PrintVerbosePanic(err) //Exit upon error, below code must not be executed
	}
	usedPercent := fmt.Sprintf("%.2f", diskUsage.UsedPercent)
	utils.ResultPrinter("Disk Usage: ", usedPercent+"%")
	utils.ResultPrinter("Disk Space Available: ", getReadableSize(diskUsage.Free))
}
