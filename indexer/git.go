//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package indexer

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	gitObj "gopkg.in/src-d/go-git.v4/plumbing/object"
)

func openGitRepo(path string) *git.Repository {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func doGitStuff(repo *git.Repository, path string, wiki *WikiPage) {

	// lookup head
	head, err := repo.Head()
	if err != nil {
		log.Print(err)
	} else {
		// lookup commit object
		headOid := head.Hash()
		commit, err := repo.CommitObject(headOid)
		if err != nil {
			log.Print(err)
		}

		// start diffing backwards
		diffCommit, err := recursiveDiffLookingForFile(repo, commit, path)
		if err != nil {
			log.Print(err)
		} else if diffCommit != nil {
			author := diffCommit.Author
			wiki.ModifiedByName = author.Name
			wiki.ModifiedByEmail = author.Email
			wiki.Modified = author.When
			if wiki.ModifiedByEmail != "" {
				wiki.ModifiedByGravatar = gravatarHashFromEmail(wiki.ModifiedByEmail)
				log.Printf("gravatar hash is: %s", wiki.ModifiedByGravatar)
			}
		} else {
			log.Printf("unable to find commit where file changed")
		}
	}
}

func recursiveDiffLookingForFile(repo *git.Repository, commit *gitObj.Commit, path string) (*gitObj.Commit, error) {
	log.Printf("checking commit %s", commit.ID())
	// if there is a parent, diff against it
	// totally not going to think about branches
	if commit.NumParents() > 0 {
		parent, err := commit.Parent(0)
		if err != nil {
			return nil, err
		}

		parentTree, err := parent.Tree()
		if err != nil {
			return nil, err
		}
		commitTree, err := commit.Tree()
		if err != nil {
			return nil, err
		}
		diffs, err := gitObj.DiffTree(parentTree, commitTree)
		if err != nil {
			return nil, err
		} else {
			found := false
			for i, change := range diffs {
				newFiles, oldFiles, err := change.Files()
				if err != nil {
					return nil, err
				}

				if newFiles.Name == path {
					found = true
					break
				} else if oldFiles.Name == path {
					found = true
					break
				}
			}
			if found {
				return commit, nil
			} else {
				return recursiveDiffLookingForFile(repo, parent, path)
			}
		}
	} else {
		// if there is no parent check to see if this file
		// was in the commit, if so, this is its
		commitTree, err := commit.Tree()
		if err != nil {
			return nil, err
		}

		for i, entry := range commitTree.Entries {
			if entry.Name == path {
				return commit, nil
			}
		}

		return nil, nil
	}
}

func gravatarHashFromEmail(email string) string {
	input := strings.ToLower(strings.TrimSpace(email))
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}
