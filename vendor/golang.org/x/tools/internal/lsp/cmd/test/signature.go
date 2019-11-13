// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdtest

import (
	"fmt"
	"testing"

	"golang.org/x/tools/internal/lsp/cmd"
	"golang.org/x/tools/internal/lsp/source"
	"golang.org/x/tools/internal/tool"

	"golang.org/x/tools/internal/span"
)

func (r *runner) SignatureHelp(t *testing.T, spn span.Span, expectedSignature *source.SignatureInformation) {
	goldenTag := "-signature"
	if expectedSignature != nil {
		goldenTag = expectedSignature.Label + goldenTag
	}
	uri := spn.URI()
	filename := uri.Filename()
	target := filename + fmt.Sprintf(":%v:%v", spn.Start().Line(), spn.Start().Column())

	app := cmd.New("gopls-test", r.data.Config.Dir, r.data.Config.Env, r.options)
	got := CaptureStdOut(t, func() {
		tool.Run(r.ctx, app, append([]string{"-remote=internal", "signature"}, target))
	})

	expect := string(r.data.Golden(goldenTag, filename, func() ([]byte, error) {
		return []byte(got), nil
	}))

	if expect != got {
		t.Errorf("signature failed failed for %s expected:\n%s\ngot:\n%s", filename, expect, got)
	}
}
