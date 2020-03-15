package shared

import (
    "fmt"
    "github.com/magefile/mage/mg"
    "github.com/magefile/mage/sh"
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
    Invoke() error
    InvokeDebug() error
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
func (l *lambda) Invoke() error {
    mg.Deps(l.Build)
    return invoke()
}

// Invokes the lambda function locally, running the debug server
func (l *lambda) InvokeDebug() error {
    mg.Deps(l.Build, l.BuildDelve)
    return invoke("-d", "5986", "--debugger-path", l.getBinFile("delve"), "--debug-args", "-delveAPI=2")
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
    // Build our version of delve
    return sh.RunWith(
        map[string]string{
            "GOARCH": "amd64",
            "GOOS":   "linux",
        },
        "go", "build",
        "-o", l.getBinFile("delve/dlv"),
        matches[0])
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
