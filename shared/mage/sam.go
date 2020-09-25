package mage

import (
    "github.com/magefile/mage/mg"
    "github.com/magefile/mage/sh"
)

var ExtraArgs []string

type Invoke mg.Namespace

const (
    envFile  = ""
    buildDir = ".aws-sam/"
)

// Build the project
func Build() error {
    return sh.Run("sam", "build")
}

func buildArgs(debug bool) []string {
    args := []string{}
    // Add an env file if specified
    if envFile != "" {
        args = append(args, "--env-vars", envFile)
    }
    if debug {
        // Specify delve as a dependency for debug
        mg.Deps(BuildDelve)
        // Add debug flags
        args = append(args, "-d", "5986", "--debugger-path", getBinFile("delve"), "--debug-args", "-delveAPI=2")
    }
    return args
}

func Sam(cmds []string, debug bool) error {
    // Determine arguments
    args := buildArgs(debug)
    // We always need to build before this
    mg.Deps(Build)
    combined := []string{"local"}
    combined = append(combined, cmds...)
    combined = append(combined, ExtraArgs...)
    combined = append(combined, args...)
    return sh.Run("sam", combined...)
}
