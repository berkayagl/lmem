package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/berkayagl/cmlib"
)

type Memory struct {
	TotalMem          int
	FreeMem           int
	AvailableMem      int
	Buffers           int
	Cached            int
	SwapCached        int
	Active            int
	Inactive          int
	Unevictable       int
	Mlocked           int
	SwapTotalMem      int
	SwapFreeMem       int
	Zswap             int
	Zswapped          int
	Dirty             int
	Writeback         int
	AnonPages         int
	Mapped            int
	Shmem             int
	KReclaimable      int
	Slab              int
	SReclaimable      int
	SUnreclaim        int
	KernelStack       int
	PageTables        int
	SecPageTables     int
	NFS_Unstable      int
	Bounce            int
	WritebackTmp      int
	CommitLimit       int
	Committed_AS      int
	VmallocTotal      int
	VmallocUsed       int
	VmallocChunk      int
	Percpu            int
	HardwareCorrupted int
	AnonHugePages     int
	ShmemHugePages    int
	ShmemPmdMapped    int
	FileHugePages     int
	FilePmdMapped     int
	HugePages_Total   int
	HugePages_Free    int
	HugePages_Rsvd    int
	HugePages_Surp    int
	Hugepagesize      int
	Hugetlb           int
	DirectMap4k       int
	DirectMap2M       int
	DirectMap1G       int
}

const (
	MemoryInfo = "/proc/meminfo"
)

// This function will take a string (s) and return an integer and an error.
func parseMemoryVal(s string) (int, error) {

	// The given string s is split into pieces and the result is assigned to fields.
	fields := strings.Fields(s)

	// If the length of the fragmented string is less than 2 : Invalid input format
	if len(fields) < 2 {
		return 0, errors.New("Invalid input format")
	}

	// The second element in the fields slice, the one with index 1, is converted to an integer using strconv.ParseInt.
	// In this case, the string is converted to a 64-bit signed integer type.
	i, err := strconv.ParseInt(fields[1], 10, 64)

	// errors control
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// This function is used to collect RAM and swap memory information.
func GetRam() (Memory, error) {

	// Read 16 lines from the MemoryInfo file.
	lines, err := cmlib.ReadLines(MemoryInfo, 54)

	// errors control
	if err != nil {
		return Memory{}, err
	}

	// 54 < errors control
	if len(lines) < 54 {
		return Memory{}, errors.New("Not enough lines in MemoryInfo!")
	}

	// Parallel parsing of memory values
	ch := make(chan memResult, 54)

	// A separate goroutine is called for each line with the range lines loop.
	// These goroutines parse the memory values using the parseMemoryVal function and write the results to the channel in the memResult structure.
	for i := range lines {
		go func(line string) {
			value, err := parseMemoryVal(line)
			ch <- memResult{value, err}
		}(lines[i])
	}

	// An instance of ram is created from the memory structure.
	ram := Memory{}

	// one result from the channel for each index. 0-54
	for i := 0; i < 54; i++ {
		result := <-ch
		if result.err != nil {
			return Memory{}, result.err
		}

		// In switch i, according to a given index, allocated memory values are assigned to the corresponding areas of the ram structure.
		switch i {
		case 0:
			ram.TotalMem = result.value
		case 1:
			ram.FreeMem = result.value
		case 2:
			ram.AvailableMem = result.value
		case 3:
			ram.Buffers = result.value
		case 4:
			ram.Cached = result.value
		case 5:
			ram.SwapCached = result.value
		case 6:
			ram.Active = result.value
		case 7:
			ram.Inactive = result.value
		case 12:
			ram.Unevictable = result.value
		case 13:
			ram.Mlocked = result.value
		case 14:
			ram.SwapTotalMem = result.value
		case 15:
			ram.SwapFreeMem = result.value
		case 16:
			ram.Zswap = result.value
		case 17:
			ram.Zswapped = result.value
		case 18:
			ram.Dirty = result.value
		case 19:
			ram.Writeback = result.value
		case 20:
			ram.AnonPages = result.value
		case 21:
			ram.Mapped = result.value
		case 22:
			ram.Shmem = result.value
		case 23:
			ram.KReclaimable = result.value
		case 24:
			ram.Slab = result.value
		case 25:
			ram.SReclaimable = result.value
		case 26:
			ram.SUnreclaim = result.value
		case 27:
			ram.KernelStack = result.value
		case 28:
			ram.PageTables = result.value
		case 29:
			ram.SecPageTables = result.value
		case 30:
			ram.NFS_Unstable = result.value
		case 31:
			ram.Bounce = result.value
		case 32:
			ram.WritebackTmp = result.value
		case 33:
			ram.CommitLimit = result.value
		case 34:
			ram.Committed_AS = result.value
		case 35:
			ram.VmallocTotal = result.value
		case 36:
			ram.VmallocUsed = result.value
		case 37:
			ram.VmallocChunk = result.value
		case 38:
			ram.Percpu = result.value
		case 39:
			ram.HardwareCorrupted = result.value
		case 40:
			ram.AnonHugePages = result.value
		case 41:
			ram.ShmemHugePages = result.value
		case 42:
			ram.ShmemPmdMapped = result.value
		case 43:
			ram.FileHugePages = result.value
		case 44:
			ram.FilePmdMapped = result.value
		case 45:
			ram.HugePages_Total = result.value
		case 46:
			ram.HugePages_Free = result.value
		case 47:
			ram.HugePages_Rsvd = result.value
		case 48:
			ram.HugePages_Surp = result.value
		case 49:
			ram.Hugepagesize = result.value
		case 50:
			ram.Hugetlb = result.value
		case 51:
			ram.DirectMap4k = result.value
		case 52:
			ram.DirectMap2M = result.value
		case 53:
			ram.DirectMap1G = result.value
		}
	}

	return ram, nil
}

type memResult struct {
	value int
	err   error
}

func GetTotalRam() (int, error) {
	line, err := cmlib.ReadFirstLine(MemoryInfo)

	if err != nil {
		return 0, err
	}

	return parseMemoryVal(line)

}

func GetFreeMem() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 2)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetAvailableMem() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 3)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetBuffers() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 4)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetCached() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 5)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetSwapCached() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 6)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetActive() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 7)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)
}

func GetInactive() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 8)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetUnevictable() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 13)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetMlocked() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 14)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetSwapTotal() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 15)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetSwapFree() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 16)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetZswap() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 17)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetZswapped() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 18)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetDirty() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 19)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetWriteBack() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 20)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetAnonPages() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 21)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetMapped() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 22)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetShmem() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 23)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetKReclaimable() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 24)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetSlab() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 25)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetSReclaimable() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 26)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetSUnreclaim() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 27)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetKernelStack() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 28)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetPageTables() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 29)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetSecPageTables() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 30)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetNfsUnstable() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 31)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetBounce() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 32)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetWritebackTmp() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 33)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetCommitLimit() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 34)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetCommitLimitAS() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 35)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetVmallocTotal() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 36)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetVmallocUsed() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 37)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetVmallocChunk() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 38)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetPercpu() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 39)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetHardwareCorrupted() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 40)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetAnonHugePages() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 41)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetShmemHugePages() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 42)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetShmemPmdMapped() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 43)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetFileHugePages() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 44)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetFilePmdMapped() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 45)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetHugePagesTotal() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 46)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetHugePagesFree() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 47)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetHugePagesRsvd() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 48)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetHugePagesSurp() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 49)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetHugePageSize() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 50)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetHugetLb() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 51)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetDirectMap4k() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 52)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetDirectMap2M() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 53)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}

func GetDirectMap1G() (int, error) {
	line, err := cmlib.ReadLine(MemoryInfo, 54)

	if err != nil {
		return 0, err

	}

	return parseMemoryVal(line)

}
