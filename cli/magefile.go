//+build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Install() error {

	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	return sh.Run("go", "install", "./...")
}

func Remove() error {

	return sh.Run("go", "clean", "-i", "github.com/chayev/yurl")
}
