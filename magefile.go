//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// mg contains helpful utility functions, like Deps

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Run "everything": Tidy
func CI() {
	mg.Deps(Tidy)
}

// Run `go mod tidy` and verifies no changes are required.
func Tidy() error {
	if err := githubDiff(); err != nil {
		return err
	}
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}
	if err := githubDiff(); err != nil {
		return err
	}

	return nil
}

func WatchSCSS() error {
	return sh.RunV("./sass", "--no-source-map", "--watch", "./static/styles/global.scss:./static/public/css/global.css")
}

func BuildSCSS() error {
	return sh.RunV("./sass", "--no-source-map", "./static/styles/global.scss:./static/public/css/global.css")
}

func BuildTempl() error {
	return sh.RunV("templ", "generate")
}

func WatchTempl() error {
	return sh.RunV("templ", "generate", "--watch")
}

func BuildTailwindCSS() error {
	return sh.RunV("./tailwindcss", "-i", "./static/styles/tailwind.css", "-o", "./static/public/css/tailwind.css")
}

func WatchTailwindCSS() error {
	return sh.RunV("./tailwindcss", "-i", "./static/styles/tailwind.css", "-o", "./static/public/css/tailwind.css", "--watch")
}

// Install/upgrade tools used for development.
func InstallDeps() error {
	for _, dep := range []string{
		"github.com/air-verse/air", // For development
		"github.com/gocolly/colly/v2",
		"github.com/a-h/templ/cmd/templ",
		"github.com/spf13/viper",
	} {
		if err := sh.RunV("go", "get", "-u", fmt.Sprintf("%s@latest", dep)); err != nil {
			return err
		}
	}

	return nil
}

// Download TailwindCSS standalone CLI for macOS 64
// TODO(austinlukaschewski): Configure to download different sources based on OS, only macOS arm64 is supported at the moment.
func InstallSassStandaloneCLI() error {
	if err := sh.RunV("curl", "-sLO", "https://github.com/sass/dart-sass/releases/download/1.83.4/dart-sass-1.83.4-macos-arm64.tar.gz | tar xz"); err != nil {
		return err
	}

	if err := sh.RunV("chmod", "+x", "./dart-sass-1.84.0-macos-arm64"); err != nil {
		return err
	}

	return sh.RunV("mv", "./dart-sass-1.84.0-macos-arm64", "./sass")
}

// Download TailwindCSS standalone CLI for macOS 64
// TODO(austinlukaschewski): Configure to download different sources based on OS, only macOS arm64 is supported at the moment.
func InstallTailwindCSSStandaloneCLI() error {
	if err := sh.RunV("curl", "-sLO", "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64"); err != nil {
		return err
	}

	if err := sh.RunV("chmod", "+x", "./tailwindcss-macos-arm64"); err != nil {
		return err
	}

	return sh.RunV("mv", "./tailwindcss-macos-arm64", "./tailwindcss")
}

// Check github to make sure on changes in go.mod & go.sum files.
func githubDiff() error {
	return sh.RunV("git", "diff", "--no-patch", "--exit-code", "--", "go.mod", "go.sum")
}
