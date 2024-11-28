// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from internal/build/util.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package build

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var DryRunFlag = flag.Bool("n", false, "dry run, don't execute commands")

// TryRun executes the given command and returns an error if error occurs.
func TryRun(cmd *exec.Cmd) error {
	fmt.Println(">>>", strings.Join(cmd.Args, " "))
	if !*DryRunFlag {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
	return nil
}

// TryRunCommand executes the given command and arguments in strings.
func TryRunCommand(cmd string, args ...string) {
	TryRun(exec.Command(cmd, args...))
}

// MustRun executes the given command and exits the host process for
// any error.
func MustRun(cmd *exec.Cmd) {
	if err := TryRun(cmd); err != nil {
		log.Fatal(err)
	}
}

func MustRunCommand(cmd string, args ...string) {
	MustRun(exec.Command(cmd, args...))
}

// GOPATH returns the value that the GOPATH environment
// variable should be set to.
func GOPATH() string {
	if os.Getenv("GOPATH") == "" {
		log.Fatal("GOPATH is not set")
	}
	return os.Getenv("GOPATH")
}

// VERSION returns the content of the VERSION file.
func VERSION() string {
	version, err := os.ReadFile("VERSION")
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes.TrimSpace(version))
}

var warnedAboutGit bool

// RunGit runs a git subcommand and returns its output.
// The command must complete successfully.
func RunGit(args ...string) string {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err == exec.ErrNotFound {
		if !warnedAboutGit {
			log.Println("Warning: can't find 'git' in PATH")
			warnedAboutGit = true
		}
		return ""
	} else if err != nil {
		log.Fatal(strings.Join(cmd.Args, " "), ": ", err, "\n", stderr.String())
	}
	return strings.TrimSpace(stdout.String())
}

// readGitFile returns content of file in .git directory.
func readGitFile(file string) string {
	content, err := os.ReadFile(path.Join(".git", file))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}

// Render renders the given template file into outputFile.
func Render(templateFile, outputFile string, outputPerm os.FileMode, x interface{}) {
	tpl := template.Must(template.ParseFiles(templateFile))
	render(tpl, outputFile, outputPerm, x)
}

// RenderString renders the given template string into outputFile.
func RenderString(templateContent, outputFile string, outputPerm os.FileMode, x interface{}) {
	tpl := template.Must(template.New("").Parse(templateContent))
	render(tpl, outputFile, outputPerm, x)
}

func render(tpl *template.Template, outputFile string, outputPerm os.FileMode, x interface{}) {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		log.Fatal(err)
	}
	out, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, outputPerm)
	if err != nil {
		log.Fatal(err)
	}
	if err := tpl.Execute(out, x); err != nil {
		log.Fatal(err)
	}
	if err := out.Close(); err != nil {
		log.Fatal(err)
	}
}

// CopyFile copies a file.
func CopyFile(dst, src string, mode os.FileMode) {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		log.Fatal(err)
	}
	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		log.Fatal(err)
	}
}

// GoTool returns the command that runs a go tool. This uses go from GOROOT instead of PATH
// so that go commands executed by build use the same version of Go as the 'host' that runs
// build code. e.g.
//
//	/usr/lib/go-1.8/bin/go run build/ci.go ...
//
// runs using go 1.8 and invokes go 1.8 tools from the same GOROOT. This is also important
// because runtime.Version checks on the host should match the tools that are run.
func GoTool(tool string, args ...string) *exec.Cmd {
	args = append([]string{tool}, args...)
	return exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), args...)
}

// ExpandPackages expands a packages list if an input contains "...".
func ExpandPackages(packages []string) []string {
	for _, pkg := range packages {
		if strings.Contains(pkg, "...") {
			var newPkgs []string

			cmd := GoTool("list", packages...)
			out, err := cmd.Output()
			if err != nil {
				log.Fatalf("package listing failed: %v\n%s", err, string(out))
			}

			for _, line := range strings.Split(string(out), "\n") {
				newPkgs = append(newPkgs, strings.TrimSpace(line))
			}
			return newPkgs
		}
	}
	return packages
}

// ExcludePackages excludes packages having patterns from the passed package slice and
// returns a slice including only the remained packages.
func ExcludePackages(packages []string, patterns []string) []string {
	// TODO-Kaia This exclusion code is a naive implementation. Improve this if it hurts build performance.
	packages = ExpandPackages(packages)

	for _, pattern := range patterns {
		var newPkgs []string
		for _, pkg := range packages {
			if !strings.Contains(pkg, pattern) {
				newPkgs = append(newPkgs, pkg)
			}
		}
		packages = newPkgs
	}
	return packages
}

// Version returns the versions defined in the checksumdb.
func Version(csdb *ChecksumDB, version string) (string, error) {
	for _, l := range csdb.allChecksums {
		if !strings.HasPrefix(l, "# version:") {
			continue
		}
		v := strings.Split(l, ":")[1]
		parts := strings.Split(v, " ")
		if len(parts) != 2 {
			log.Print("Erroneous version-string", "v", l)
			continue
		}
		if parts[0] == version {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf("no version found for '%v'", version)
}

// ExtractArchive unpacks a .zip or .tar.gz archive to the destination directory.
func ExtractArchive(archive string, dest string) error {
	ar, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer ar.Close()

	switch {
	case strings.HasSuffix(archive, ".tar.gz"):
		return extractTarball(ar, dest)
	case strings.HasSuffix(archive, ".zip"):
		return extractZip(ar, dest)
	default:
		return fmt.Errorf("unhandled archive type %s", archive)
	}
}

// extractTarball unpacks a .tar.gz file.
func extractTarball(ar io.Reader, dest string) error {
	gzr, err := gzip.NewReader(ar)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		// Move to the next file header.
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		// We only care about regular files, directory modes
		// and special file types are not supported.
		if header.Typeflag == tar.TypeReg {
			armode := header.FileInfo().Mode()
			err := extractFile(header.Name, armode, tr, dest)
			if err != nil {
				return fmt.Errorf("extract %s: %v", header.Name, err)
			}
		}
	}
}

// extractZip unpacks the given .zip file.
func extractZip(ar *os.File, dest string) error {
	info, err := ar.Stat()
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(ar, info.Size())
	if err != nil {
		return err
	}

	for _, zf := range zr.File {
		if !zf.Mode().IsRegular() {
			continue
		}

		data, err := zf.Open()
		if err != nil {
			return err
		}
		err = extractFile(zf.Name, zf.Mode(), data, dest)
		data.Close()
		if err != nil {
			return fmt.Errorf("extract %s: %v", zf.Name, err)
		}
	}
	return nil
}

// extractFile extracts a single file from an archive.
func extractFile(arpath string, armode os.FileMode, data io.Reader, dest string) error {
	// Check that path is inside destination directory.
	target := filepath.Join(dest, filepath.FromSlash(arpath))
	if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("path %q escapes archive destination", target)
	}

	// Remove the preivously-extracted file if it exists
	if err := os.RemoveAll(target); err != nil {
		return err
	}

	// Recreate the destination directory
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	// Copy file data.
	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, armode)
	if err != nil {
		return err
	}
	if _, err = io.Copy(file, data); err != nil {
		file.Close()
		os.Remove(target)
		return err
	}
	return file.Close()
}
