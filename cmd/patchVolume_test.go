/*
Copyright © 2019 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPatchVolumeHalevel runs a series of volume patching test
func TestPatchVolumeHalevel(t *testing.T) {
	volName := genVolName("testVol")
	haLevel := 2
	volCreate(t, volName)
	// Now update halevel to 2
	testPatchVolumeHalevel(t, volName, haLevel)
	volCleanup(t, volName)
}

func TestPatchVolumeResize(t *testing.T) {
	volName := genVolName("testVol")
	var size uint64
	// Setting size to 2GB
	size = 2
	volCreate(t, volName)

	// Now update halevel to 2
	testPatchVolumeResize(t, volName, size)
	volCleanup(t, volName)
}

func TestPatchVolumeShared(t *testing.T) {
	volName := genVolName("testVol")
	shared := true

	volCreate(t, volName)
	// Now update shared to true
	testPatchVolumeShared(t, volName, shared)
	volCleanup(t, volName)
}

func TestPatchVolumeUnsetShared(t *testing.T) {
	volName := genVolName("testVol")
	sharedTrue := true
	sharedFalse := false

	volCreate(t, volName)
	// Now update shared to true
	testPatchVolumeShared(t, volName, sharedTrue)
	//Now unset shared aka to false
	testPatchVolumeShared(t, volName, sharedFalse)
	volCleanup(t, volName)
}

// Helper to create a volume
func volCreate(t *testing.T, volName string) {
	// Create a volume
	testCreateVolume(t, volName, 1)
	// Verify that the volume got created
	assert.True(t, testHasVolume(volName))
}

// Helper function to cleanup volume created
func volCleanup(t *testing.T, volName string) {
	// Delete Volume
	testDeleteVolume(t, volName)
}