package mage

import (
    "fmt"
    "github.com/magefile/mage/mg"
    "github.com/magefile/mage/sh"
    "github.com/magefile/mage/target"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "regexp"
)

type lambda struct {
    binDir    string
    toolsFile string
}

type Lambda interface {
    Build() error
    StartApi(debug bool, envFile string) error
    Invoke(debug bool, envFile string) error
    BuildDelve() error
    InstallTools() error
}

// New creates a new lambda helper
func New(binDir, toolsFile string) Lambda {
    l := lambda{
        binDir:    binDir,
        toolsFile: toolsFile,
    }
    return &l
}

// Invokes the lambda function locally
func (l *lambda) Invoke(debug bool, envFile string) error {
    extraArgs := []string{}
    // Add an env file if specified
    if envFile != "" {
        extraArgs = append(extraArgs, "--env-vars", envFile)
    }
    if debug {
        // Specify delve as a dependency for debug
        mg.Deps(l.BuildDelve)
        // Add debug flags
        extraArgs = append(extraArgs, "-d", "5986", "--debugger-path", l.getBinFile("delve"), "--debug-args", "-delveAPI=2")
    }
    // We always need to build before this
    mg.Deps(l.Build)
    return invoke(extraArgs...)
}

// Starts a local API gateway for the lambda function locally
func (l *lambda) StartApi(debug bool, envFile string) error {
    // Start building the arguments to the 'sam' command
    args := []string{"local", "start-api"}
    // Add an env file if specified
    if envFile != "" {
        args = append(args, "--env-vars", envFile)
    }
    if debug {
        // Ensure we have the debugger available
        mg.Deps(l.BuildDelve)
        args = append(args, "-d", "5986", "--debugger-path", l.getBinFile("delve"), "--debug-args", "-delveAPI=2")
    }
    // We always need to build
    mg.Deps(l.Build)
    return sh.Run("sam", args...)
}

// Build the project
func (l *lambda) Build() error {
    return sh.Run("sam", "build")
}

func (l *lambda) BuildDelve() error {
    // Ensure delve is installed
    mg.Deps(l.InstallTools)
    // Find Delve
    pattern := path.Join(
        os.ExpandEnv("$GOPATH"),
        "pkg/mod/github.com/go-delve/delve@*/cmd/dlv")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }
    if len(matches) == 0 {
        return fmt.Errorf("delve not found")
    }
    src := matches[0]
    // State our inputs and outputs
    tgt := l.getBinFile("delve/dlv")
    isNew, err := target.Path(tgt, src)
    if err != nil {
        return err
    }
    if isNew {
        // Build our version of delve
        return sh.RunWith(
            map[string]string{
                "GOARCH": "amd64",
                "GOOS":   "linux",
            },
            "go", "build",
            "-o", tgt,
            src)
    }
    // Nothing to build
    return nil
}

func (l *lambda) InstallTools() error {
    // Read the tools file
    toolsFile, err := ioutil.ReadFile(l.toolsFile)
    if err != nil {
        return err
    }
    // Parse the tools
    re := regexp.MustCompile(`_ \"(.*?)\"`)
    matches := re.FindAllSubmatch(toolsFile, -1)
    // Install
    for _, match := range matches {
        err := sh.Run("go", "install", string(match[1]))
        if err != nil {
            return err
        }
    }
    return nil
}

func invoke(extraArgs ...string) error {
    combined := []string{"local", "invoke", "consumer", "-e", "event.json"}
    combined = append(combined, extraArgs...)
    return sh.Run("sam", combined...)
}

func (l *lambda) getBinFile(file string) string {
    return path.Join(l.binDir, file)
}
