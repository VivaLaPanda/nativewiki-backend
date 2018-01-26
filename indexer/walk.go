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
	"io/ioutil"
	"log"
	"os"

	"github.com/blevesearch/bleve"
	git "gopkg.in/src-d/go-git.v4"
)

func walkForIndexing(path string, index bleve.Index, repo *git.Repository) {

	dirEntries, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, dirEntry := range dirEntries {
		dirEntryPath := path + string(os.PathSeparator) + dirEntry.Name()
		if dirEntry.IsDir() {
			walkForIndexing(dirEntryPath, index, repo)
		} else if pathMatch(dirEntry.Name()) {
			processUpdate(index, repo, dirEntryPath)
		}
	}
}
