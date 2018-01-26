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
	"log"
	"path/filepath"
	"regexp"

	"github.com/VivaLaPanda/nativewiki-backend/wikiroot"
	bleveHttp "github.com/blevesearch/bleve/http"
)

// A function which starts an indexing API given a directory to store the indexer
// data and a list of wikiroots
// index is a string to a path where the index should be stored
func StartIndexer(indexPath string, root wikiroot.WikiRoot) {
	if indexPath == "" {
		indexPath = root.Name + "wiki.bleve"
	}

	dir := root.Directory
	pathFilter := root.Pathfilter

	if dir == "" {
		log.Fatalf("must specify a directory to watch")
	}

	if pathFilter == "" {
		pathFilter = `\.md$`
	}

	// cleanup the dir
	dir = filepath.Clean(*dir)

	var err error
	if pathFilter != "" {
		pathRegexp, err = regexp.Compile(pathFilter)
		if err != nil {
			log.Fatal(err)
		}
	}

	// open the index
	index := openIndex(indexPath)

	// open the git repo
	repo := openGitRepo(dir)

	// create a router to serve static files
	router := staticFileRouter()

	// add the API
	bleveHttp.RegisterIndexName(root.Name+"wiki", index)
	searchHandler := bleveHttp.NewSearchHandler(root.Name + "wiki")
	router.Handle("/api/search/"+root.Name, searchHandler).Methods("POST")

	// start a watcher on the directory
	watcher := startWatching(dir, index, repo)
	defer watcher.Close()

	// walk the directory to ensure current
	walkForIndexing(dir, index, repo)
}

func pathMatch(path string) bool {
	if pathRegexp != nil {
		return pathRegexp.MatchString(path)
	}
	return true
}
