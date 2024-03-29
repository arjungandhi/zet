//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	return sh.Run("go", "build", "-o", "zet", "./cmd/zet")
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	// install to ~/.local/bin
	homedir, err := os.UserHomeDir()
	if err != nil {
		Clean()
		return err
	}

	err = os.Rename("zet", homedir+"/.local/bin/zet")
	if err != nil {
		Clean()
		return err
	}

	return nil
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.Run("go", "mod", "download")
}

// Test
func Test() error {
	mg.Deps(InstallDeps)
	fmt.Println("Testing...")
	return sh.Run("go", "test", "./...")
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("zet")
}
