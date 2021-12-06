package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/tinkerbell/hub/actions/archive2disk/v1/pkg/archive"
)

const mountAction = "/mountAction"

func main() {
	fmt.Printf("Archive2Disk - Archive streamer\n------------------------\n")
	blockDevice := os.Getenv("DEST_DISK")
	filesystemType := os.Getenv("FS_TYPE")
	path := os.Getenv("DEST_PATH")
	archiveURL := os.Getenv("ARCHIVE_URL")
	archiveType := os.Getenv("ARCHIVE_TYPE")
	//optional checksum to validate tarfile, must be of the format
	//checksum name:checsum
	//ex:
	//sha256:shasum
	//sha512:shasum	
	archiveChecksum := os.Getenv("TARFILE_CHECKSUM")	
	if blockDevice == "" {
		log.Fatalf("No Block Device speified with Environment Variable [DEST_DISK]")
	}

	// Create the /mountAction mountpoint (no folders exist previously in scratch container)
	err := os.Mkdir(mountAction, os.ModeDir)
	if err != nil {
		log.Fatalf("Error creating the action Mountpoint [%s]", mountAction)
	}

	// Mount the block device to the /mountAction point
	err = syscall.Mount(blockDevice, mountAction, filesystemType, 0, "")
	if err != nil {
		log.Fatalf("Mounting [%s] -> [%s] error [%v]", blockDevice, mountAction, err)
	}
	log.Infof("Mounted [%s] -> [%s]", blockDevice, mountAction)

	// Write the image to disk
	err = archive.Write(archiveURL, archiveType, filepath.Join(mountAction, path), archiveChecksum)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully unpacked [%s] to [%s] on device [%s]", archiveURL, path, blockDevice)
}
