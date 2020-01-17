package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func check(err error) {
	if err != nil {
		if strings.HasPrefix(err.Error(), "exit status") {
			os.Exit(0)
		}
		panic(err)
	}
}

func main() {
	switch os.Args[1] {
	case "run":
		run(os.Args[2:]...)
	case "child":
		child(os.Args[2:]...)
	default:
		panic("no such command")
	}
}

func run(args ...string) {

	fmt.Println("Running", args, os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	check(cmd.Run())
}

func child(args ...string) {
	var env []string
	var mem int
	app := &cli.App{
		Name: "run",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "env",
				Aliases: []string{"e"},
			},
			&cli.IntFlag{
				Name:    "memory",
				Aliases: []string{"m"},
			},
		},
		Action: func(c *cli.Context) error {
			env = c.StringSlice("env")
			args = c.Args().Slice()
			mem = c.Int("memory")
			return nil
		},
	}
	check(app.Run(append([]string{"child"}, args...)))
	fmt.Println("Child", args, os.Getpid())

	if mem > 0 {
		cgMem(mem)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	check(syscall.Sethostname([]byte(fmt.Sprintf("cfs-%d", rand.Int31n(10000)))))
	check(syscall.Chroot("rootfs"))
	check(syscall.Chdir("/"))
	check(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer check(syscall.Unmount("proc", 0))

	check(cmd.Run())
}

func cgMem(maxMemSize int) {
	size := fmt.Sprintf("%d", maxMemSize)

	cgroups := "/sys/fs/cgroup/"
	memory := filepath.Join(cgroups, "memory")
	_ = os.Mkdir(filepath.Join(memory, "cfs"), 0755)
	check(ioutil.WriteFile(filepath.Join(memory, "cfs/memory.swappiness"), []byte("0"), 0700))
	check(ioutil.WriteFile(filepath.Join(memory, "cfs/memory.limit_in_bytes"), []byte(size), 0700))

	// Writing the pid on which to apply the cgroup to
	check(ioutil.WriteFile(filepath.Join(memory, "cfs/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}
