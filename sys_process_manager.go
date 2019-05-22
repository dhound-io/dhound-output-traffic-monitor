package main

import (
	process "github.com/shirou/gopsutil/process"
	//"fmt"
	"time"
)

var (
	hits       = 0
	misses     = 0
	isOutdated = false
	notFound   = false
	isNew      = false
)
type Counter struct {
	hits, misses int32
}
type SysProcessManager struct {
	_pidToProcessInfoMap map[int32]*ProcessInfo
}

type ProcessInfo struct {
	Name, CommandLine string
	Pid               int32
	EventTimeUtcNumber int64
}

func (manager *SysProcessManager) Init() {
	manager._pidToProcessInfoMap = make(map[int32]*ProcessInfo)
}

func (manager *SysProcessManager) Run() {
	go func() {
		for {
			manager._syncProcessInfoOnPids()
			time.Sleep(500 * time.Millisecond)
			//time.Sleep(2000 * time.Millisecond)
		}
	}()
}

func (manager *SysProcessManager) _syncProcessInfoOnPids() bool {
	processes, err := process.Processes()

	if err != nil {
		emitLine(logLevel.important, "could not get processes: %s", err.Error())
		return false
	}

	pids := make([]int32, 0)
	for _, process := range processes {
		pids = append(pids, process.Pid)
	}

	pidsToProcess := make([]int32, 0)

	// sync map
	for _, pid := range pids {
		if _, ok := manager._pidToProcessInfoMap[pid]; !ok {
			manager._pidToProcessInfoMap[pid] = &ProcessInfo{}
			pidsToProcess = append(pidsToProcess, pid)
		}
	}

	/*
	obsoletePids := make([]int32, 0)
	for pid, _ := range manager._pidToProcessInfoMap {
		if ContainsInt32(pids, pid) == false {
			debugJson("ContainsInt32")
			debugJson(pids)
			debugJson(pid)
			// debug("remove pid: %d (%s)", pid, value.Name)
			obsoletePids = append(obsoletePids, pid)
		}
	}

	for _, pid := range obsoletePids {
		delete(manager._pidToProcessInfoMap, pid)
	}
	*/

	if len(pidsToProcess) > 0 {
		// parse name
		for _, process := range processes {
			if ContainsInt32(pidsToProcess, process.Pid) {
				name, _ := process.Name()
				manager._pidToProcessInfoMap[process.Pid].Name = name
				// debug("new pid: %d (%s)", process.Pid, name)
			}
		}
	}

	return true
}

func (manager SysProcessManager) FindProcessInfoByPid(pid int32) *ProcessInfo {
	if pid > 0 {
		currentTime := time.Now().UTC().Unix()
		isOutdated  = false
		notFound = false
		isNew = false
		if processInfo, ok := manager._pidToProcessInfoMap[pid]; ok {
			if processInfo.EventTimeUtcNumber > 0{
				diff := currentTime - processInfo.EventTimeUtcNumber
				if diff > 60 {
					isOutdated  = true
					delete(manager._pidToProcessInfoMap, pid)
				}else{
					hits = hits + 1
					//debugJson(fmt.Sprintf("hits: %d", hits))
					return processInfo
				}
			}
			isNew = true
		}else{
			notFound = true
		}

		if notFound == true || isOutdated == true || isNew == true {
			processInfo := &ProcessInfo{}
			//if len(processInfo.CommandLine) < 1 {
				process, _ := process.NewProcess(pid)

				if process != nil {
					cmdLine, _ := process.Cmdline()
					name, _ := process.Name()
					processInfo.Name = name
					processInfo.CommandLine = cmdLine
				}

				if len(processInfo.CommandLine) < 1 {
					processInfo.CommandLine = "-"
				}
				processInfo.EventTimeUtcNumber = time.Now().UTC().Unix()
			//}
			manager._pidToProcessInfoMap[pid] = processInfo
			return processInfo
		}
		misses = misses + 1
		//debugJson(fmt.Sprintf("misses: %d", misses))
		return nil
	}

	misses = misses + 1
	//debugJson(fmt.Sprintf("misses: %d", misses))
	return nil
}
