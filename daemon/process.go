package daemon

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"
)

const (
	goHTTPCLIDirectory = "/.go-http-cli"
	pidFile            = "daemon.pid"
)

var (
	startDaemonTries = 0
)

// EnsureDaemon makes sure a daemon of the right version is up and running. This method will either
// panic or return after the daemon is up and running.
func EnsureDaemon() error {
	version, connErr := Handshake()

	if connErr != nil || version != DaemonMajorVersion {
		KillDaemon()
		startDaemon()
		startDaemonTries++

		if startDaemonTries > 3 {
			panic("Tried to start the daemon more than 3 times.")
		}

		return EnsureDaemon()
	}

	return nil
}

// KillDaemon will kill the process with the PID stored in the daemon.pid file
func KillDaemon() {
	pid, pidError := getDaemonPID()
	if pidError != nil {
		panic(pidError)
	}

	process, processErr := os.FindProcess(pid)
	if processErr != nil {
		panic(processErr)
	}

	process.Kill()
}

// WriteDaemonPID writes the PID of the current process to a file in the go-http-cli process dir
func WriteDaemonPID() error {
	processDir, dirError := getProcessDirectory()
	if dirError != nil {
		return dirError
	}

	if createDirErr := os.MkdirAll(processDir, 0755); createDirErr != nil {
		return createDirErr
	}

	file, fileErr := os.OpenFile(processDir+"/"+pidFile, os.O_WRONLY|os.O_CREATE, 0755)
	if fileErr != nil {
		return fileErr
	}

	file.WriteString(fmt.Sprintf("%d\n", os.Getpid()))

	return nil
}

func getDaemonPID() (int, error) {
	processDir, dirError := getProcessDirectory()
	if dirError != nil {
		return 0, dirError
	}

	buffer, readError := ioutil.ReadFile(processDir + "/" + pidFile)
	if readError != nil {
		return 0, readError
	}

	return strconv.Atoi(strings.TrimSpace(string(buffer)))
}

func getProcessDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + goHTTPCLIDirectory, nil
}

func startDaemon() {
	processDir, dirError := getProcessDirectory()
	if dirError != nil {
		panic(dirError)
	}

	logFilePath := processDir + "/daemon.log"
	logFile, logFileError := os.Create(logFilePath)
	if logFileError != nil {
		panic(logFileError)
	}

	defer logFile.Close()

	cmd := exec.Command("go-http-daemon")
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if waitError := waitToStart(logFilePath); waitError != nil {
		panic(waitError)
	}
}

// waitToStart waits for the daemon to start by checking the log file for a specific string. Make sure
// the log file is empty before calling this or you might check for the wrong thing.
func waitToStart(logFilePath string) error {
	tries := 0
	for {
		tries++
		buffer, readError := ioutil.ReadFile(logFilePath)
		if readError != nil {
			return readError
		}

		logData := string(buffer)
		if strings.Contains(logData, "started and waiting for connections") {
			return nil
		}

		time.Sleep(100 * time.Millisecond)

		if tries >= 50 {
			return errors.New("Waited too long for daemon to start")
		}
	}
}
