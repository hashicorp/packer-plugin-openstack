// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package openstack

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
)

// WaitForVolume waits for the given volume to become available.
func WaitForVolume(ctx context.Context, blockStorageClient *gophercloud.ServiceClient, volumeID string) error {
	maxNumErrors := 10
	numErrors := 0

	for {
		status, err := GetVolumeStatus(ctx, blockStorageClient, volumeID)
		if err != nil {
			errCode, ok := err.(*gophercloud.ErrUnexpectedResponseCode)
			if ok && (errCode.Actual == 500 || errCode.Actual == 404) {
				numErrors++
				if numErrors >= maxNumErrors {
					log.Printf("[ERROR] Maximum number of errors (%d) reached; failing with: %s", numErrors, err)
					return err
				}
				log.Printf("[ERROR] %d error received, will ignore and retry: %s", errCode.Actual, err)
				time.Sleep(2 * time.Second)
				continue
			}

			return err
		}

		if status == "available" {
			return nil
		}

		if status == "error" {
			return errors.New("The status of volume is error")
		}

		log.Printf("Waiting for volume creation status: %s", status)
		time.Sleep(2 * time.Second)
	}
}

// GetVolumeSize returns volume size in gigabytes based on the image min disk
// value if it's not empty.
// Or it calculates needed gigabytes size from the image bytes size.
func GetVolumeSize(ctx context.Context, imageClient *gophercloud.ServiceClient, imageID string) (int, error) {
	sourceImage, err := images.Get(ctx, imageClient, imageID).Extract()
	if err != nil {
		return 0, err
	}

	if sourceImage.MinDiskGigabytes != 0 {
		return sourceImage.MinDiskGigabytes, nil
	}

	volumeSizeMB := sourceImage.SizeBytes / 1024 / 1024
	volumeSizeGB := int(sourceImage.SizeBytes / 1024 / 1024 / 1024)

	// Increment gigabytes size if the initial size can't be divided without
	// remainder.
	if volumeSizeMB%1024 > 0 {
		volumeSizeGB++
	}

	return volumeSizeGB, nil
}

func GetVolumeStatus(ctx context.Context, blockStorageClient *gophercloud.ServiceClient, volumeID string) (string, error) {
	volume, err := volumes.Get(ctx, blockStorageClient, volumeID).Extract()
	if err != nil {
		return "", err
	}

	return volume.Status, nil
}
