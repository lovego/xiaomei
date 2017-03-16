package stack

import (
	"fmt"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/cluster"
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func Logs(svcName string, all bool) error {
	if all {
		return logService(svcName, allLogs)
	}
	if svcName == `` {
		for _, s := range getService() {
			err := logService(s, oneLog)
			if err != nil {
				return err
			}
		}
	}
	return logService(svcName, oneLog)
}

func logService(svcName string, logFunc func(svcName string, ms, ws []string) error) error {
	ms, ws := cluster.GetCluster().List()
	if len(ms) == 0 && len(ws) == 0 {
		return nil
	}
	return logFunc(svcName, ms, ws)
}

func oneLog(svcName string, ms, ws []string) error {
	for _, m := range ms {
		containers, err := getContainers(svcName, m)
		if err != nil {
			return err
		}
		if len(containers) > 0 {
			return printLog(svcName, containers[0])
		}
	}
	for _, w := range ws {
		containers, err := getContainers(svcName, w)
		if err != nil {
			return err
		}
		if len(containers) > 0 {
			return printLog(svcName, containers[0])
		}
	}
	return nil
}

func allLogs(svcName string, ms, ws []string) error {
	for _, m := range ms {
		containers, err := getContainers(svcName, m)
		if err != nil {
			return err
		}
		for _, container := range containers {
			err := printLog(svcName, container)
			if err != nil {
				return err
			}
		}
	}
	for _, w := range ws {
		containers, err := getContainers(svcName, w)
		if err != nil {
			return err
		}
		for _, container := range containers {
			err := printLog(svcName, container)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getService() []string {
	services := []string{}
	for s, _ := range GetStack().Services {
		services = append(services, s)
	}
	return services
}

func getContainers(svcName, addr string) ([]string, error) {
	script := `docker ps -af name=%s -q`
	if svcName != `` {
		script = fmt.Sprintf(script, release.Name()+`_`+svcName)
	} else {
		script = fmt.Sprintf(script, release.Name())
	}
	result, err := cmd.SshRun(cmd.O{Output: true}, addr, script)
	if err != nil {
		return nil, err
	}
	containers := strings.Split(result, "\n")
	return containers, nil
}

func printLog(svcName, container string) error {
	if svcName == `` {
		cmd.Run(cmd.O{}, `echo`, fmt.Sprintf("container %s log:", container))
	} else {
		cmd.Run(cmd.O{}, `echo`, fmt.Sprintf("service %s container %s log:", svcName, container))
	}
	_, err := cluster.Run(cmd.O{}, fmt.Sprintf(`docker logs %s`, container))
	return err
}
