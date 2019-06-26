// Copyright © Microsoft <wastore@microsoft.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"github.com/Azure/azure-storage-azcopy/azbfs"
	chk "gopkg.in/check.v1"
)

func (s *cmdIntegrationSuite) TestRemoveDirectory(c *chk.C) {
	mockedRPC := interceptor{}
	mockedRPC.init()
	ctx := context.Background()

	// set up the file system
	bfsServiceURL := getBfsSU()
	fsURL, _ := createNewFileSystem(c, bfsServiceURL)
	defer deleteFileSystem(c, fsURL)

	// set up the directory to be deleted
	dirURL := fsURL.NewDirectoryURL(generateName("dir", 0))
	_, err := dirURL.Create(ctx)
	c.Assert(err, chk.IsNil)
	fileURL := dirURL.NewFileURL(generateName("file", 0))
	_, err = fileURL.Create(ctx, azbfs.BlobFSHTTPHeaders{})
	c.Assert(err, chk.IsNil)

	// trying to remove the dir with recursive=false should fail
	raw := getDefaultRemoveRawInput(dirURL.String())
	raw.recursive = false
	runCopyAndVerify(c, raw, func(err error) {
		c.Assert(err, chk.NotNil)
	})

	// removing the dir with recursive=true should succeed
	raw.recursive = true
	runCopyAndVerify(c, raw, func(err error) {
		c.Assert(err, chk.IsNil)

		// make sure the directory does not exist anymore
		_, err = dirURL.GetProperties(ctx)
		c.Assert(err, chk.NotNil)
	})
}

func (s *cmdIntegrationSuite) TestRemoveFile(c *chk.C) {
	mockedRPC := interceptor{}
	mockedRPC.init()
	ctx := context.Background()

	// set up the file system
	bfsServiceURL := getBfsSU()
	fsURL, _ := createNewFileSystem(c, bfsServiceURL)
	defer deleteFileSystem(c, fsURL)

	// set up the file to be deleted
	parentDirURL := fsURL.NewDirectoryURL(generateName("dir", 0))
	_, err := parentDirURL.Create(ctx)
	c.Assert(err, chk.IsNil)
	fileURL := parentDirURL.NewFileURL(generateName("file", 0))
	_, err = fileURL.Create(ctx, azbfs.BlobFSHTTPHeaders{})
	c.Assert(err, chk.IsNil)

	raw := getDefaultRemoveRawInput(fileURL.String())

	runCopyAndVerify(c, raw, func(err error) {
		c.Assert(err, chk.IsNil)

		// make sure the file does not exist anymore
		_, err = fileURL.GetProperties(ctx)
		c.Assert(err, chk.NotNil)
	})
}
