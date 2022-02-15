// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build !gogit
// +build !gogit

package git

import (
	"path"
	"strings"
)

// GetTreeEntryByPath get the tree entries according the sub dir
func (t *Tree) GetTreeEntryByPath(relpath string) (*TreeEntry, error) {
	if len(relpath) == 0 {
		return &TreeEntry{
			ptree:     t,
			ID:        t.ID,
			name:      "",
			fullName:  "",
			entryMode: EntryModeTree,
		}, nil
	}

	// FIXME: This should probably use git cat-file --batch to be a bit more efficient
	relpath = path.Clean(relpath)
	parts := strings.Split(relpath, "/")
	var err error
	tree := t
	for i, name := range parts {
		if i == len(parts)-1 {
			entries, err := tree.ListEntries()
			if err != nil {
				return nil, err
			}
			for _, v := range entries {
				if v.Name() == name {
					return v, nil
				}
			}
		} else {
			tree, err = tree.SubTree(name)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, ErrNotExist{"", relpath}
}