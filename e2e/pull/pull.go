// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package pull

import (
	"os"
	"os/exec"
	//	"path/filepath"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/sylabs/singularity/e2e/imgbuild"
	"github.com/sylabs/singularity/internal/pkg/test"
)

type testingEnv struct {
	// base env for running tests
	CmdPath     string `split_words:"true"`
	TestDir     string `split_words:"true"`
	RunDisabled bool   `default:"false"`
}

var testenv testingEnv

//func TestEnv(t *Testing.T) {

func imagePull(library string, imagePath string, sourceSpec string, force, unauthenticated bool) ([]byte, error) {
	var argv []string
	argv = append(argv, "pull")
	if force {
		argv = append(argv, "--force")
	}
	if unauthenticated {
		argv = append(argv, "--allow-unauthenticated")
	}
	if library != "" {
		argv = append(argv, "--library", library)
	}
	if imagePath != "" {
		argv = append(argv, imagePath)
	}
	argv = append(argv, sourceSpec)

	return exec.Command(testenv.CmdPath, argv...).CombinedOutput()
}

func testPull(t *testing.T) {
	test.DropPrivilege(t)

	imagePath := "./test_pull.sif"

	// nolint:maligned
	tests := []struct {
		name            string
		sourceSpec      string
		force           bool
		unauthenticated bool
		library         string
		imagePath       string
		success         bool
	}{
		// TODO: Download the Sylabs pub key before running tests
		{"Pull_Library", "library://alpine:3.8", false, true, "", imagePath, true}, // https://cloud.sylabs.io/library
		{"ForceUnauth", "library://alpine:3.8", true, false, "", imagePath, false},
		{"Force", "library://alpine:3.8", true, true, "", imagePath, true},
		{"Unsigned_image", "library://sylabs/tests/unsigned:1.0.0", true, true, "", imagePath, true},
		{"Unsigned_image_fail", "library://sylabs/tests/unsigned:1.0.0", true, false, "", imagePath, false}, // pull a unsigned image; should fail
		{"Pull_Docker", "docker://alpine:3.8", true, false, "", imagePath, true},                            // https://hub.docker.com/
		{"Pull_Shub", "shub://GodloveD/busybox", true, false, "", imagePath, true},                          // https://singularity-hub.org/
		{"PullWithHash", "library://sylabs/tests/signed:sha256.5c439fd262095766693dae95fb81334c3a02a7f0e4dc6291e0648ed4ddc61c6c", true, true, "", imagePath, true},
		{"PullWithoutTransportProtocol", "alpine:3.8", true, true, "", imagePath, true},
		{"PullNonExistent", "library://this_should_not/exist/not_exist", true, false, "", imagePath, false}, // pull a non-existent container
		{"Pull_Library_Latest", "library://alpine:latest", true, true, "", imagePath, true},                 // https://cloud.sylabs.io/library
	}
	defer os.Remove(imagePath)
	for _, tt := range tests {
		t.Run(tt.name, test.WithoutPrivilege(func(t *testing.T) {
			b, err := imagePull(tt.library, tt.imagePath, tt.sourceSpec, tt.force, tt.unauthenticated)
			if tt.success {
				if err != nil {
					t.Log(string(b))
					t.Fatalf("unexpected failure: %v", err)
				}
				//imgbuild.ImageVerify(t, tt.imagePath, false)
				imgbuild.ImageVerify(t, testenv.CmdPath, imagePath, false, testenv.RunDisabled)
			} else {
				if err == nil {
					t.Log(string(b))
					t.Fatalf("unexpected success: command should have failed")
				}
			}
		}))
	}

}

// RunE2ETests is the main func to trigger the test suite
func RunE2ETests(t *testing.T) {
	err := envconfig.Process("E2E", &testenv)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(testenv)

	t.Run("pull", testPull)

}
