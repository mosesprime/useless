package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
    path := flag.String("p", "./", "path to watch for changes")
    //TODO: dotfiles := flag.Bool("d", false, "include dot files")
    //TODO: ignore := flag.String("i", "", "path to some '.ignore' file")
    //TODO: exclude := flag.String("x", "", "additonal files to exclude")
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        output("missing command to execute\n")
    } else {
        run(*path, args)
    }
}

func run(path string, args []string) {
    initStat, err := os.Stat(path)
    if err != nil {
        output("%s\n", err)
        return
    }
    cmd := startProcess(args)
    for {
        stat, err := os.Stat(path)
        if err != nil {
            output("%s\n", err)
            return
        }
        if stat.ModTime() != initStat.ModTime() {
            output("detected changes\n")
            err := stopProcess(cmd)
            if err != nil {
                output("%s\n", err)
                return
            }
            initStat = stat
            cmd = startProcess(args)
            continue
        }
        time.Sleep(1 * time.Second)
    }
}

func startProcess(userCmd []string) *exec.Cmd {
    output("new command: %s\n", userCmd)
    name := userCmd[0]
    rest := userCmd[1:]
    cmd := exec.Command(name, rest...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Start()
    if err != nil {
        output("error starting process: %s\n", err)
        return nil
    }
    output("process started: %v\n", cmd.Process.Pid)
    return cmd
}

func stopProcess(cmd *exec.Cmd) error {
    err := cmd.Process.Kill()
    if err != nil {
        output("failed to kill process: %v\n", cmd.Process.Pid)
        return err
    }
    output("process terminated: %v\n", cmd.Process.Pid)
    return cmd.Wait()
}

func output(f string, args ...any) {
    fmt.Printf("[rerun] "+f, args...)
}
