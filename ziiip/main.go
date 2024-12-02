package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("missing subcommand")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "zip":
        zb, err := NewZipBuilder(os.Args[2:])
        if err != nil {
            fmt.Printf("failed to parse: %s\n", err)
            os.Exit(1)
        }
        zb.Run()
    case "tar":
        fmt.Println("unimplimented")
        os.Exit(1)
    default:
        fmt.Println("missing subcommand")
        os.Exit(1)
    }
}

type ZipBuilder struct {
    printf func(f string, v ...any)
}

func NewZipBuilder(args []string) (*ZipBuilder, error) {
    set := flag.NewFlagSet("zip", flag.ExitOnError)
    verbose := set.Bool("v", false, "verbose")

    err := set.Parse(args)
    if err != nil {
        return nil, err
    }

    printf := func (f string, v ...any) {}
    if *verbose {
        printf = func(f string, v ...any) { fmt.Printf(f, v...) }
    }

    return &ZipBuilder{
        printf: printf,
    }, nil
}

func (zb *ZipBuilder) Run() {
    zb.printf("hello\n")
}
