// Copyright (c) 2021 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/decred/politeia/politeiad/cmd/legacypoliteia/gitbe"
)

var (
	// Regular expression that matches the git proposal token from a proposal
	// parent directory.
	expGitProposalToken   = `[0-9a-f]{64}`
	regexGitProposalToken = regexp.MustCompile(expGitProposalToken)
)

// gitProposalToken takes a git repo path and returns the proposal token from
// the path if the path corresponds to the proposal parent directory.
//
// Input:
// "fdd68c87961549750adf29e178128210cb310294080211cf6a35792aa1bb7f63.json"
// Output:
// "fdd68c87961549750adf29e178128210cb310294080211cf6a35792aa1bb7f63", true
//
// Input:
// "mainnet/fdd68c87961549750adf29e178128210cb310294080211cf6a35792aa1bb7f63/1"
// Output:
// "fdd68c87961549750adf29e178128210cb310294080211cf6a35792aa1bb7f63", true
func gitProposalToken(path string) (string, bool) {
	var (
		gitToken = regexGitProposalToken.FindString(path)
		ok       = (gitToken != "")
	)
	return gitToken, ok
}

// gitProposalTokens recursively walks the provided directory and returns an
// inventory of all git proposal tokens found in file paths.
func gitProposalTokens(dirPath string) (map[string]struct{}, error) {
	tokens := make(map[string]struct{}, 256)
	err := filepath.WalkDir(dirPath,
		func(path string, d fs.DirEntry, err error) error {
			// We only care about directories
			if !d.IsDir() {
				return nil
			}

			// Get the token from the path
			token, ok := gitProposalToken(path)
			if ok {
				tokens[token] = struct{}{}
				return nil
			}

			return nil
		})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// latestVersion returns the latest version of a legacy git proposal. The
// version number is parsed from the directory structure.
func latestVersion(gitRepo, token string) (uint64, error) {
	// Compile a list of all directories. The version numbers
	// are the directory name.
	dirs := make(map[string]struct{}, 64)
	err := filepath.WalkDir(filepath.Join(gitRepo, token),
		func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				return nil
			}
			dirs[d.Name()] = struct{}{}
			return nil
		})
	if err != nil {
		return 0, err
	}

	// Parse the version number from the directory name
	versions := make(map[uint64]struct{}, 64)
	for dirname := range dirs {
		u, err := strconv.ParseUint(dirname, 10, 32)
		if err != nil {
			// Not a version directory
			continue
		}
		versions[u] = struct{}{}
	}
	if err != nil {
		return 0, err
	}

	// Find the most recent version
	var latest uint64
	for version := range versions {
		if version > latest {
			latest = version
		}
	}
	if latest == 0 {
		return 0, fmt.Errorf("latest version not found for %v", token)
	}

	return latest, nil
}

var (
	// Regular expersion that matches the token and version number from a
	// proposal directory path.
	expProposalVersion   = `[0-9a-f]{64}\/[0-9]{1,}`
	regexProposalVersion = regexp.MustCompile(expProposalVersion)
)

// proposalVersion parses the version number for the proposal directory path
// and returns it.
func proposalVersion(proposalDir string) (uint32, error) {
	var (
		subPath    = regexProposalVersion.FindString(proposalDir)
		versionStr = filepath.Base(subPath)
	)
	u, err := strconv.ParseUint(versionStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(u), nil
}

// proposalAttachmentFiles returns the filesnames of all proposal attachment
// files. This function does NOT return the file path, just the file name. The
// proposal index file and proposal metadata file are not considered to be
// attachments.
func proposalAttachmentFilenames(proposalDir string) ([]string, error) {
	var (
		notAnAttachmentFile = map[string]struct{}{
			gitbe.IndexFilename:            {},
			gitbe.ProposalMetadataFilename: {},
		}

		payloadDir = payloadDirPath(proposalDir)
		filenames  = make([]string, 0, 64)
	)

	// Walk the payload directory
	err := filepath.WalkDir(payloadDir,
		func(path string, d fs.DirEntry, err error) error {
			// There shouldn't be any nested directories
			// in the payload directory, but check just
			// in case.
			if d.IsDir() {
				return nil
			}

			if _, ok := notAnAttachmentFile[d.Name()]; ok {
				// Not an attachment; skip
				return nil
			}

			// This is an attachment file
			filenames = append(filenames, d.Name())

			return nil
		})
	if err != nil {
		return nil, err
	}

	return filenames, nil
}
