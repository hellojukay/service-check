package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
)

type StringArray []string

func (s *StringArray) String() string {
	return "[" + strings.Join(*s, ",") + "]"
}

func (s *StringArray) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s StringArray) Index(i int) (string, error) {
	if i >= len(s) {
		return "", fmt.Errorf("index out of range")
	}
	var arr []string
	arr = s
	return arr[i], nil
}

func (s StringArray) IsBoolFlag() bool {
	return false
}

// Service 系统服务
type Service interface {
	IsRunning() bool
	IsEnabled() bool
}

// Unit systemd Unit
type Unit struct {
	Name  string
	exsit bool
}

// IsRunning systemd check Unit is active
func (unit Unit) IsRunning() bool {
	if !unit.exsit {
		return false
	}
	return run([]string{"systemctl", "is-active", unit.Name})
}

// Check systemd Unit is enabled
func (unit Unit) IsEnabled() bool {
	if !unit.exsit {
		return false
	}
	return run([]string{"systemctl", "is-enabled", unit.Name})
}

func (unit Unit) Exsit() bool {
	return unit.exsit
}

// NewUnit create a Unit
func NewService(name string) Service {
	var u = Unit{
		Name: name,
	}
	u.exsit = run([]string{"systemctl", "status", name})
	return u
}

func run(cmds []string) bool {
	var cmd = exec.Command("sudo", cmds...)
	if cmd.Run() != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
}

var (
	RUNNING  = color.GreenString("active")
	FAILED   = color.RedString("inactive")
	ENABLEED = color.GreenString("enabled")
	DISABLED = color.RedString("diabled")
)
var allService StringArray

func init() {
	flag.Var(&allService, "s", "all services you  want to check")
	flag.Parse()
}

var printFmt = "%-20s%-20s%-20s\n"

func main() {
	if len([]string(allService)) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	for _, service := range []string(allService) {
		var unit = NewService(service)
		if unit.IsEnabled() && unit.IsRunning() {
			fmt.Printf(printFmt, service, color.GreenString("active"), color.GreenString("enabled"))
			continue
		}
		if unit.IsEnabled() && (!unit.IsRunning()) {
			fmt.Printf(printFmt, service, color.RedString("inactive"), color.GreenString("enabled"))
			continue
		}
		if !unit.IsEnabled() && unit.IsRunning() {
			fmt.Printf(printFmt, service, color.GreenString("active"), color.RedString("disabled"))
			continue
		}
		if !unit.IsEnabled() && !unit.IsRunning() {
			fmt.Printf(printFmt, service, color.RedString("inactive"), color.RedString("disabled"))
			continue
		}
	}
}
