package demo

import (
	"io/ioutil"
	"log"
	"os/exec"
	"syscall"

	"github.com/Dataman-Cloud/seckilling/gate/src/cache"
)

func Reset() error {
	// reload nigix
	cmdReload := "docker exec resty nginx -s reload"
	_, _, err := ExecuteSysCommand(cmdReload)
	if err != nil {
		log.Println("reload nginx data has error: ", err)
		return err
	}

	conn := cache.Open()
	defer conn.Close()

	_, err = conn.Do("SET", "wf:1", 0)
	if err != nil {
		log.Println("reset counter has error: ", err)
		return err
	}

	return nil

}

func ExecuteSysCommand(command string) (string, string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("StdoutPipe: ", err.Error())
		return "", "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println("StderrPipe: ", err.Error())
		return "", "", err
	}

	if err := cmd.Start(); err != nil {
		log.Println("Start: ", err.Error())
		return "", "", err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Println("ReadAll stderr: ", err.Error())
		return "", "", err
	}

	if len(bytesErr) != 0 {
		log.Println("stderr is not nil: ", string(bytesErr))
		return "", string(bytesErr), nil
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ReadAll err: ", err.Error())
		return "", "", err
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Wait err: ", err.Error())
		return "", "", err
	}

	return string(bytes), "", nil
}
