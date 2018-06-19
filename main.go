package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/mattes/go-expand-tilde.v1"
)

func setup() *log.Logger {
	return log.New(os.Stdout, "[bearsync] ", log.Lshortfile)
}

func check(logger *log.Logger, err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func main() {
	logger := setup()
	bearConfigName := "net.shinyfrog.bear"

	localPath := "~/Library/Containers/"
	localBearConfig, err := tilde.Expand(path.Join(localPath, bearConfigName))
	check(logger, err)

	dropboxPath := flag.String("dropboxPath", "~/Dropbox/Personal/CONFIGS/Bear", "path to where you want to store your Bear config")
	cloudBearConfig := path.Join(*dropboxPath, bearConfigName)

	if strings.HasPrefix(*dropboxPath, "~") {
		cloudBearConfig, err = tilde.Expand(cloudBearConfig)
	} else if !filepath.IsAbs(cloudBearConfig) {
		cloudBearConfig, err = filepath.Abs(cloudBearConfig)
	}
	check(logger, err)

	// Check that the dropboxPath exists and contains the config folder
	if _, err = os.Stat(cloudBearConfig); os.IsNotExist(err) {
		logger.Fatalf("Make sure the path %s exits.\n%v", cloudBearConfig, err)
	}

	logger.Printf("Moving %s to %s...", localBearConfig, path.Join(localBearConfig, ".local"))
	// Set old config folder as local
	err = os.Rename(localBearConfig, fmt.Sprintf("%s.%s", localBearConfig, "local"))
	if os.IsNotExist(err) {
		logger.Printf("Folder in path %s does not exist. Continuing regardless", localBearConfig)
	} else {
		check(logger, err)
	}

	logger.Printf("Creating symlink from %s to %s...", cloudBearConfig, localBearConfig)
	// Create symlink from cloudBearConfig to localBearConfig
	err = os.Symlink(cloudBearConfig, localBearConfig)
	check(logger, err)

	logger.Println("Succesfully finished.")
}
