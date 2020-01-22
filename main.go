package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

	// Executing the same binary again, but replacing run with child in order to utilize the new namespaces
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{

		Cloneflags:
		//  Creates a new namespace  ($ man 2 clone)
		//  for hostname and more
			syscall.CLONE_NEWUTS |

		//  Creates a new namespace  ($ man 2 clone)
		//  for process ids
			syscall.CLONE_NEWPID |

		//  Creates a new namespace  ($ man 2 clone)
		//  for mounting things
			syscall.CLONE_NEWNS,
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

	// Sets the hostname of in the New UTS for the process
	check(syscall.Sethostname([]byte(fmt.Sprintf("cfs-%d", rand.Int31n(10000)))))

	// Sets the root filesystem to the dir rootfs
	check(syscall.Chroot("rootfs"))

	// Sets the dir in which the process starts from
	check(syscall.Chdir("/"))

	// Mounts a news proc folder in order to gain access to process information within the container
	check(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer check(syscall.Unmount("proc", 0))

	check(cmd.Run())
}

// Set the memory cgroup for the process given the number of bytes
func cgMem(maxMemSize int) {
	size := fmt.Sprintf("%d", maxMemSize)

	memory := "/sys/fs/cgroup/memory/cfs"
	// Creates the cgroup dir if it does not exist
	_ = os.Mkdir(memory, 0755)

	// Disables swap for the container
	check(ioutil.WriteFile(filepath.Join(memory, "memory.swappiness"), []byte("0"), 0700))

	// Writes the maximum amount of memory allowed to be used
	check(ioutil.WriteFile(filepath.Join(memory, "memory.limit_in_bytes"), []byte(size), 0700))

	// Writing the pid of the current process to apply the cgroup to it
	check(ioutil.WriteFile(filepath.Join(memory, "cgroup.procs"), []byte(fmt.Sprintf("%d", os.Getpid())), 0700))
}
