// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package library

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-vela/types/raw"
)

func TestLibrary_Build_Duration(t *testing.T) {
	// setup types
	unfinished := testBuild()
	unfinished.SetFinished(0)

	// setup tests
	tests := []struct {
		build *Build
		want  string
	}{
		{
			build: testBuild(),
			want:  "1s",
		},
		{
			build: unfinished,
			want:  time.Since(time.Unix(unfinished.GetStarted(), 0)).Round(time.Second).String(),
		},
		{
			build: new(Build),
			want:  "...",
		},
	}

	// run tests
	for _, test := range tests {
		got := test.build.Duration()

		if got != test.want {
			t.Errorf("Duration is %v, want %v", got, test.want)
		}
	}
}

func TestLibrary_Build_Environment(t *testing.T) {
	// setup types
	_comment := testBuild()
	_comment.SetEvent("comment")
	_comment.SetRef("refs/pulls/1/head")

	_deploy := testBuild()
	_deploy.SetEvent("deployment")
	_deploy.SetDeploy("production")
	_deploy.SetDeployPayload(map[string]string{
		"foo": "test1",
		"bar": "test2",
	})

	_deployTag := testBuild()
	_deployTag.SetRef("refs/tags/v0.1.0")
	_deployTag.SetEvent("deployment")
	_deployTag.SetDeploy("production")
	_deployTag.SetDeployPayload(map[string]string{
		"foo": "test1",
		"bar": "test2",
	})

	_pull := testBuild()
	_pull.SetEvent("pull_request")
	_pull.SetRef("refs/pulls/1/head")

	_tag := testBuild()
	_tag.SetEvent("tag")
	_tag.SetRef("refs/tags/v0.1.0")

	// setup tests
	tests := []struct {
		build *Build
		want  map[string]string
	}{
		{
			build: testBuild(),
			want: map[string]string{
				"VELA_BUILD_AUTHOR":       "OctoKitty",
				"VELA_BUILD_AUTHOR_EMAIL": "OctoKitty@github.com",
				"VELA_BUILD_BASE_REF":     "",
				"VELA_BUILD_BRANCH":       "master",
				"VELA_BUILD_CHANNEL":      "TODO",
				"VELA_BUILD_CLONE":        "https://github.com/github/octocat.git",
				"VELA_BUILD_COMMIT":       "48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_CREATED":      "1563474076",
				"VELA_BUILD_DISTRIBUTION": "linux",
				"VELA_BUILD_ENQUEUED":     "1563474077",
				"VELA_BUILD_EVENT":        "push",
				"VELA_BUILD_HOST":         "example.company.com",
				"VELA_BUILD_LINK":         "https://example.company.com/github/octocat/1",
				"VELA_BUILD_MESSAGE":      "First commit...",
				"VELA_BUILD_NUMBER":       "1",
				"VELA_BUILD_PARENT":       "1",
				"VELA_BUILD_REF":          "refs/heads/master",
				"VELA_BUILD_RUNTIME":      "docker",
				"VELA_BUILD_SENDER":       "OctoKitty",
				"VELA_BUILD_STARTED":      "1563474078",
				"VELA_BUILD_SOURCE":       "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_STATUS":       "running",
				"VELA_BUILD_TITLE":        "push received from https://github.com/github/octocat",
				"VELA_BUILD_WORKSPACE":    "TODO",
				"BUILD_AUTHOR":            "OctoKitty",
				"BUILD_AUTHOR_EMAIL":      "OctoKitty@github.com",
				"BUILD_BASE_REF":          "",
				"BUILD_BRANCH":            "master",
				"BUILD_CHANNEL":           "TODO",
				"BUILD_CLONE":             "https://github.com/github/octocat.git",
				"BUILD_COMMIT":            "48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_CREATED":           "1563474076",
				"BUILD_ENQUEUED":          "1563474077",
				"BUILD_EVENT":             "push",
				"BUILD_HOST":              "example.company.com",
				"BUILD_LINK":              "https://example.company.com/github/octocat/1",
				"BUILD_MESSAGE":           "First commit...",
				"BUILD_NUMBER":            "1",
				"BUILD_PARENT":            "1",
				"BUILD_REF":               "refs/heads/master",
				"BUILD_SENDER":            "OctoKitty",
				"BUILD_STARTED":           "1563474078",
				"BUILD_SOURCE":            "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_STATUS":            "running",
				"BUILD_TITLE":             "push received from https://github.com/github/octocat",
				"BUILD_WORKSPACE":         "TODO",
			},
		},
		{
			build: _comment,
			want: map[string]string{
				"VELA_BUILD_AUTHOR":         "OctoKitty",
				"VELA_BUILD_AUTHOR_EMAIL":   "OctoKitty@github.com",
				"VELA_BUILD_BASE_REF":       "",
				"VELA_BUILD_BRANCH":         "master",
				"VELA_BUILD_CHANNEL":        "TODO",
				"VELA_BUILD_CLONE":          "https://github.com/github/octocat.git",
				"VELA_BUILD_COMMIT":         "48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_CREATED":        "1563474076",
				"VELA_BUILD_DISTRIBUTION":   "linux",
				"VELA_BUILD_ENQUEUED":       "1563474077",
				"VELA_BUILD_EVENT":          "comment",
				"VELA_BUILD_HOST":           "example.company.com",
				"VELA_BUILD_LINK":           "https://example.company.com/github/octocat/1",
				"VELA_BUILD_MESSAGE":        "First commit...",
				"VELA_BUILD_NUMBER":         "1",
				"VELA_BUILD_PARENT":         "1",
				"VELA_BUILD_PULL_REQUEST":   "1",
				"VELA_BUILD_REF":            "refs/pulls/1/head",
				"VELA_BUILD_RUNTIME":        "docker",
				"VELA_BUILD_SENDER":         "OctoKitty",
				"VELA_BUILD_STARTED":        "1563474078",
				"VELA_BUILD_SOURCE":         "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_STATUS":         "running",
				"VELA_BUILD_TITLE":          "push received from https://github.com/github/octocat",
				"VELA_BUILD_WORKSPACE":      "TODO",
				"VELA_PULL_REQUEST":         "1",
				"BUILD_AUTHOR":              "OctoKitty",
				"BUILD_AUTHOR_EMAIL":        "OctoKitty@github.com",
				"BUILD_BASE_REF":            "",
				"BUILD_BRANCH":              "master",
				"BUILD_CHANNEL":             "TODO",
				"BUILD_CLONE":               "https://github.com/github/octocat.git",
				"BUILD_COMMIT":              "48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_CREATED":             "1563474076",
				"BUILD_ENQUEUED":            "1563474077",
				"BUILD_EVENT":               "comment",
				"BUILD_HOST":                "example.company.com",
				"BUILD_LINK":                "https://example.company.com/github/octocat/1",
				"BUILD_MESSAGE":             "First commit...",
				"BUILD_NUMBER":              "1",
				"BUILD_PARENT":              "1",
				"BUILD_PULL_REQUEST_NUMBER": "1",
				"BUILD_REF":                 "refs/pulls/1/head",
				"BUILD_SENDER":              "OctoKitty",
				"BUILD_STARTED":             "1563474078",
				"BUILD_SOURCE":              "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_STATUS":              "running",
				"BUILD_TITLE":               "push received from https://github.com/github/octocat",
				"BUILD_WORKSPACE":           "TODO",
			},
		},
		{
			build: _deploy,
			want: map[string]string{
				"VELA_BUILD_AUTHOR":        "OctoKitty",
				"VELA_BUILD_AUTHOR_EMAIL":  "OctoKitty@github.com",
				"VELA_BUILD_BASE_REF":      "",
				"VELA_BUILD_BRANCH":        "master",
				"VELA_BUILD_CHANNEL":       "TODO",
				"VELA_BUILD_CLONE":         "https://github.com/github/octocat.git",
				"VELA_BUILD_COMMIT":        "48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_CREATED":       "1563474076",
				"VELA_BUILD_DISTRIBUTION":  "linux",
				"VELA_BUILD_ENQUEUED":      "1563474077",
				"VELA_BUILD_EVENT":         "deployment",
				"VELA_BUILD_HOST":          "example.company.com",
				"VELA_BUILD_LINK":          "https://example.company.com/github/octocat/1",
				"VELA_BUILD_MESSAGE":       "First commit...",
				"VELA_BUILD_NUMBER":        "1",
				"VELA_BUILD_PARENT":        "1",
				"VELA_BUILD_REF":           "refs/heads/master",
				"VELA_BUILD_RUNTIME":       "docker",
				"VELA_BUILD_SENDER":        "OctoKitty",
				"VELA_BUILD_STARTED":       "1563474078",
				"VELA_BUILD_SOURCE":        "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_STATUS":        "running",
				"VELA_BUILD_TARGET":        "production",
				"VELA_BUILD_TITLE":         "push received from https://github.com/github/octocat",
				"VELA_BUILD_WORKSPACE":     "TODO",
				"VELA_DEPLOYMENT":          "production",
				"BUILD_TARGET":             "production",
				"BUILD_AUTHOR":             "OctoKitty",
				"BUILD_AUTHOR_EMAIL":       "OctoKitty@github.com",
				"BUILD_BASE_REF":           "",
				"BUILD_BRANCH":             "master",
				"BUILD_CHANNEL":            "TODO",
				"BUILD_CLONE":              "https://github.com/github/octocat.git",
				"BUILD_COMMIT":             "48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_CREATED":            "1563474076",
				"BUILD_ENQUEUED":           "1563474077",
				"BUILD_EVENT":              "deployment",
				"BUILD_HOST":               "example.company.com",
				"BUILD_LINK":               "https://example.company.com/github/octocat/1",
				"BUILD_MESSAGE":            "First commit...",
				"BUILD_NUMBER":             "1",
				"BUILD_PARENT":             "1",
				"BUILD_REF":                "refs/heads/master",
				"BUILD_SENDER":             "OctoKitty",
				"BUILD_STARTED":            "1563474078",
				"BUILD_SOURCE":             "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_STATUS":             "running",
				"BUILD_TITLE":              "push received from https://github.com/github/octocat",
				"BUILD_WORKSPACE":          "TODO",
				"DEPLOYMENT_PARAMETER_FOO": "test1",
				"DEPLOYMENT_PARAMETER_BAR": "test2",
			},
		},
		{
			build: _deployTag,
			want: map[string]string{
				"VELA_BUILD_AUTHOR":        "OctoKitty",
				"VELA_BUILD_AUTHOR_EMAIL":  "OctoKitty@github.com",
				"VELA_BUILD_BASE_REF":      "",
				"VELA_BUILD_BRANCH":        "master",
				"VELA_BUILD_CHANNEL":       "TODO",
				"VELA_BUILD_CLONE":         "https://github.com/github/octocat.git",
				"VELA_BUILD_COMMIT":        "48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_CREATED":       "1563474076",
				"VELA_BUILD_DISTRIBUTION":  "linux",
				"VELA_BUILD_ENQUEUED":      "1563474077",
				"VELA_BUILD_EVENT":         "deployment",
				"VELA_BUILD_HOST":          "example.company.com",
				"VELA_BUILD_LINK":          "https://example.company.com/github/octocat/1",
				"VELA_BUILD_MESSAGE":       "First commit...",
				"VELA_BUILD_NUMBER":        "1",
				"VELA_BUILD_PARENT":        "1",
				"VELA_BUILD_REF":           "refs/tags/v0.1.0",
				"VELA_BUILD_RUNTIME":       "docker",
				"VELA_BUILD_SENDER":        "OctoKitty",
				"VELA_BUILD_STARTED":       "1563474078",
				"VELA_BUILD_SOURCE":        "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_STATUS":        "running",
				"VELA_BUILD_TAG":           "v0.1.0",
				"VELA_BUILD_TARGET":        "production",
				"VELA_BUILD_TITLE":         "push received from https://github.com/github/octocat",
				"VELA_BUILD_WORKSPACE":     "TODO",
				"VELA_DEPLOYMENT":          "production",
				"BUILD_AUTHOR":             "OctoKitty",
				"BUILD_AUTHOR_EMAIL":       "OctoKitty@github.com",
				"BUILD_BASE_REF":           "",
				"BUILD_BRANCH":             "master",
				"BUILD_CHANNEL":            "TODO",
				"BUILD_CLONE":              "https://github.com/github/octocat.git",
				"BUILD_COMMIT":             "48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_CREATED":            "1563474076",
				"BUILD_ENQUEUED":           "1563474077",
				"BUILD_EVENT":              "deployment",
				"BUILD_HOST":               "example.company.com",
				"BUILD_LINK":               "https://example.company.com/github/octocat/1",
				"BUILD_MESSAGE":            "First commit...",
				"BUILD_NUMBER":             "1",
				"BUILD_PARENT":             "1",
				"BUILD_REF":                "refs/tags/v0.1.0",
				"BUILD_SENDER":             "OctoKitty",
				"BUILD_STARTED":            "1563474078",
				"BUILD_SOURCE":             "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_STATUS":             "running",
				"BUILD_TAG":                "v0.1.0",
				"BUILD_TARGET":             "production",
				"BUILD_TITLE":              "push received from https://github.com/github/octocat",
				"BUILD_WORKSPACE":          "TODO",
				"DEPLOYMENT_PARAMETER_FOO": "test1",
				"DEPLOYMENT_PARAMETER_BAR": "test2",
			},
		},
		{
			build: _pull,
			want: map[string]string{
				"VELA_BUILD_AUTHOR":         "OctoKitty",
				"VELA_BUILD_AUTHOR_EMAIL":   "OctoKitty@github.com",
				"VELA_BUILD_BASE_REF":       "",
				"VELA_BUILD_BRANCH":         "master",
				"VELA_BUILD_CHANNEL":        "TODO",
				"VELA_BUILD_CLONE":          "https://github.com/github/octocat.git",
				"VELA_BUILD_COMMIT":         "48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_CREATED":        "1563474076",
				"VELA_BUILD_DISTRIBUTION":   "linux",
				"VELA_BUILD_ENQUEUED":       "1563474077",
				"VELA_BUILD_EVENT":          "pull_request",
				"VELA_BUILD_HOST":           "example.company.com",
				"VELA_BUILD_LINK":           "https://example.company.com/github/octocat/1",
				"VELA_BUILD_MESSAGE":        "First commit...",
				"VELA_BUILD_NUMBER":         "1",
				"VELA_BUILD_PARENT":         "1",
				"VELA_BUILD_PULL_REQUEST":   "1",
				"VELA_BUILD_REF":            "refs/pulls/1/head",
				"VELA_BUILD_RUNTIME":        "docker",
				"VELA_BUILD_SENDER":         "OctoKitty",
				"VELA_BUILD_STARTED":        "1563474078",
				"VELA_BUILD_SOURCE":         "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_STATUS":         "running",
				"VELA_BUILD_TITLE":          "push received from https://github.com/github/octocat",
				"VELA_BUILD_WORKSPACE":      "TODO",
				"VELA_PULL_REQUEST":         "1",
				"VELA_PULL_REQUEST_SOURCE":  "changes",
				"VELA_PULL_REQUEST_TARGET":  "",
				"BUILD_AUTHOR":              "OctoKitty",
				"BUILD_AUTHOR_EMAIL":        "OctoKitty@github.com",
				"BUILD_BASE_REF":            "",
				"BUILD_BRANCH":              "master",
				"BUILD_CHANNEL":             "TODO",
				"BUILD_CLONE":               "https://github.com/github/octocat.git",
				"BUILD_COMMIT":              "48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_CREATED":             "1563474076",
				"BUILD_ENQUEUED":            "1563474077",
				"BUILD_EVENT":               "pull_request",
				"BUILD_HOST":                "example.company.com",
				"BUILD_LINK":                "https://example.company.com/github/octocat/1",
				"BUILD_MESSAGE":             "First commit...",
				"BUILD_NUMBER":              "1",
				"BUILD_PARENT":              "1",
				"BUILD_PULL_REQUEST_NUMBER": "1",
				"BUILD_REF":                 "refs/pulls/1/head",
				"BUILD_SENDER":              "OctoKitty",
				"BUILD_STARTED":             "1563474078",
				"BUILD_SOURCE":              "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_STATUS":              "running",
				"BUILD_TITLE":               "push received from https://github.com/github/octocat",
				"BUILD_WORKSPACE":           "TODO",
			},
		},
		{
			build: _tag,
			want: map[string]string{
				"VELA_BUILD_AUTHOR":       "OctoKitty",
				"VELA_BUILD_AUTHOR_EMAIL": "OctoKitty@github.com",
				"VELA_BUILD_BASE_REF":     "",
				"VELA_BUILD_BRANCH":       "master",
				"VELA_BUILD_CHANNEL":      "TODO",
				"VELA_BUILD_CLONE":        "https://github.com/github/octocat.git",
				"VELA_BUILD_COMMIT":       "48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_CREATED":      "1563474076",
				"VELA_BUILD_DISTRIBUTION": "linux",
				"VELA_BUILD_ENQUEUED":     "1563474077",
				"VELA_BUILD_EVENT":        "tag",
				"VELA_BUILD_HOST":         "example.company.com",
				"VELA_BUILD_LINK":         "https://example.company.com/github/octocat/1",
				"VELA_BUILD_MESSAGE":      "First commit...",
				"VELA_BUILD_NUMBER":       "1",
				"VELA_BUILD_PARENT":       "1",
				"VELA_BUILD_REF":          "refs/tags/v0.1.0",
				"VELA_BUILD_RUNTIME":      "docker",
				"VELA_BUILD_SENDER":       "OctoKitty",
				"VELA_BUILD_STARTED":      "1563474078",
				"VELA_BUILD_SOURCE":       "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"VELA_BUILD_STATUS":       "running",
				"VELA_BUILD_TAG":          "v0.1.0",
				"VELA_BUILD_TITLE":        "push received from https://github.com/github/octocat",
				"VELA_BUILD_WORKSPACE":    "TODO",
				"BUILD_AUTHOR":            "OctoKitty",
				"BUILD_AUTHOR_EMAIL":      "OctoKitty@github.com",
				"BUILD_BASE_REF":          "",
				"BUILD_BRANCH":            "master",
				"BUILD_CHANNEL":           "TODO",
				"BUILD_CLONE":             "https://github.com/github/octocat.git",
				"BUILD_COMMIT":            "48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_CREATED":           "1563474076",
				"BUILD_ENQUEUED":          "1563474077",
				"BUILD_EVENT":             "tag",
				"BUILD_HOST":              "example.company.com",
				"BUILD_LINK":              "https://example.company.com/github/octocat/1",
				"BUILD_MESSAGE":           "First commit...",
				"BUILD_NUMBER":            "1",
				"BUILD_PARENT":            "1",
				"BUILD_REF":               "refs/tags/v0.1.0",
				"BUILD_SENDER":            "OctoKitty",
				"BUILD_STARTED":           "1563474078",
				"BUILD_SOURCE":            "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163",
				"BUILD_STATUS":            "running",
				"BUILD_TAG":               "v0.1.0",
				"BUILD_TITLE":             "push received from https://github.com/github/octocat",
				"BUILD_WORKSPACE":         "TODO",
			},
		},
	}

	// run test
	for _, test := range tests {
		got := test.build.Environment("TODO", "TODO")

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Environment is %v, want %v", got, test.want)
		}
	}
}

func TestLibrary_Build_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		build *Build
		want  *Build
	}{
		{
			build: testBuild(),
			want:  testBuild(),
		},
		{
			build: new(Build),
			want:  new(Build),
		},
	}

	// run tests
	for _, test := range tests {
		if test.build.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.build.GetID(), test.want.GetID())
		}

		if test.build.GetRepoID() != test.want.GetRepoID() {
			t.Errorf("GetRepoID is %v, want %v", test.build.GetRepoID(), test.want.GetRepoID())
		}

		if test.build.GetNumber() != test.want.GetNumber() {
			t.Errorf("GetNumber is %v, want %v", test.build.GetNumber(), test.want.GetNumber())
		}

		if test.build.GetParent() != test.want.GetParent() {
			t.Errorf("GetParent is %v, want %v", test.build.GetParent(), test.want.GetParent())
		}

		if test.build.GetEvent() != test.want.GetEvent() {
			t.Errorf("GetEvent is %v, want %v", test.build.GetEvent(), test.want.GetEvent())
		}

		if test.build.GetStatus() != test.want.GetStatus() {
			t.Errorf("GetStatus is %v, want %v", test.build.GetStatus(), test.want.GetStatus())
		}

		if test.build.GetError() != test.want.GetError() {
			t.Errorf("GetError is %v, want %v", test.build.GetError(), test.want.GetError())
		}

		if test.build.GetEnqueued() != test.want.GetEnqueued() {
			t.Errorf("GetEnqueued is %v, want %v", test.build.GetEnqueued(), test.want.GetEnqueued())
		}

		if test.build.GetCreated() != test.want.GetCreated() {
			t.Errorf("GetCreated is %v, want %v", test.build.GetCreated(), test.want.GetCreated())
		}

		if test.build.GetStarted() != test.want.GetStarted() {
			t.Errorf("GetStarted is %v, want %v", test.build.GetStarted(), test.want.GetStarted())
		}

		if test.build.GetFinished() != test.want.GetFinished() {
			t.Errorf("GetFinished is %v, want %v", test.build.GetFinished(), test.want.GetFinished())
		}

		if test.build.GetDeploy() != test.want.GetDeploy() {
			t.Errorf("GetDeploy is %v, want %v", test.build.GetDeploy(), test.want.GetDeploy())
		}

		if !reflect.DeepEqual(test.build.GetDeployPayload(), test.want.GetDeployPayload()) {
			t.Errorf("GetDeployPayload is %v, want %v", test.build.GetDeployPayload(), test.want.GetDeployPayload())
		}

		if test.build.GetClone() != test.want.GetClone() {
			t.Errorf("GetClone is %v, want %v", test.build.GetClone(), test.want.GetClone())
		}

		if test.build.GetSource() != test.want.GetSource() {
			t.Errorf("GetSource is %v, want %v", test.build.GetSource(), test.want.GetSource())
		}

		if test.build.GetTitle() != test.want.GetTitle() {
			t.Errorf("GetTitle is %v, want %v", test.build.GetTitle(), test.want.GetTitle())
		}

		if test.build.GetMessage() != test.want.GetMessage() {
			t.Errorf("GetMessage is %v, want %v", test.build.GetMessage(), test.want.GetMessage())
		}

		if test.build.GetCommit() != test.want.GetCommit() {
			t.Errorf("GetCommit is %v, want %v", test.build.GetCommit(), test.want.GetCommit())
		}

		if test.build.GetSender() != test.want.GetSender() {
			t.Errorf("GetSender is %v, want %v", test.build.GetSender(), test.want.GetSender())
		}

		if test.build.GetAuthor() != test.want.GetAuthor() {
			t.Errorf("GetAuthor is %v, want %v", test.build.GetAuthor(), test.want.GetAuthor())
		}

		if test.build.GetEmail() != test.want.GetEmail() {
			t.Errorf("GetEmail is %v, want %v", test.build.GetEmail(), test.want.GetEmail())
		}

		if test.build.GetLink() != test.want.GetLink() {
			t.Errorf("GetLink is %v, want %v", test.build.GetLink(), test.want.GetLink())
		}

		if test.build.GetBranch() != test.want.GetBranch() {
			t.Errorf("GetBranch is %v, want %v", test.build.GetBranch(), test.want.GetBranch())
		}

		if test.build.GetRef() != test.want.GetRef() {
			t.Errorf("GetRef is %v, want %v", test.build.GetRef(), test.want.GetRef())
		}

		if test.build.GetBaseRef() != test.want.GetBaseRef() {
			t.Errorf("GetBaseRef is %v, want %v", test.build.GetBaseRef(), test.want.GetBaseRef())
		}

		if test.build.GetHeadRef() != test.want.GetHeadRef() {
			t.Errorf("GetHeadRef is %v, want %v", test.build.GetHeadRef(), test.want.GetHeadRef())
		}

		if test.build.GetHost() != test.want.GetHost() {
			t.Errorf("GetHost is %v, want %v", test.build.GetHost(), test.want.GetHost())
		}

		if test.build.GetRuntime() != test.want.GetRuntime() {
			t.Errorf("GetRuntime is %v, want %v", test.build.GetRuntime(), test.want.GetRuntime())
		}

		if test.build.GetDistribution() != test.want.GetDistribution() {
			t.Errorf("GetDistribution is %v, want %v", test.build.GetDistribution(), test.want.GetDistribution())
		}
	}
}

func TestLibrary_Build_Setters(t *testing.T) {
	// setup types
	var b *Build

	// setup tests
	tests := []struct {
		build *Build
		want  *Build
	}{
		{
			build: testBuild(),
			want:  testBuild(),
		},
		{
			build: b,
			want:  new(Build),
		},
	}

	// run tests
	for _, test := range tests {
		test.build.SetID(test.want.GetID())
		test.build.SetRepoID(test.want.GetRepoID())
		test.build.SetNumber(test.want.GetNumber())
		test.build.SetParent(test.want.GetParent())
		test.build.SetEvent(test.want.GetEvent())
		test.build.SetStatus(test.want.GetStatus())
		test.build.SetError(test.want.GetError())
		test.build.SetEnqueued(test.want.GetEnqueued())
		test.build.SetCreated(test.want.GetCreated())
		test.build.SetStarted(test.want.GetStarted())
		test.build.SetFinished(test.want.GetFinished())
		test.build.SetDeploy(test.want.GetDeploy())
		test.build.SetDeployPayload(test.want.GetDeployPayload())
		test.build.SetClone(test.want.GetClone())
		test.build.SetSource(test.want.GetSource())
		test.build.SetTitle(test.want.GetTitle())
		test.build.SetMessage(test.want.GetMessage())
		test.build.SetCommit(test.want.GetCommit())
		test.build.SetSender(test.want.GetSender())
		test.build.SetAuthor(test.want.GetAuthor())
		test.build.SetEmail(test.want.GetEmail())
		test.build.SetLink(test.want.GetLink())
		test.build.SetBranch(test.want.GetBranch())
		test.build.SetRef(test.want.GetRef())
		test.build.SetBaseRef(test.want.GetBaseRef())
		test.build.SetHeadRef(test.want.GetHeadRef())
		test.build.SetHost(test.want.GetHost())
		test.build.SetRuntime(test.want.GetRuntime())
		test.build.SetDistribution(test.want.GetDistribution())

		if test.build.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.build.GetID(), test.want.GetID())
		}

		if test.build.GetRepoID() != test.want.GetRepoID() {
			t.Errorf("SetRepoID is %v, want %v", test.build.GetRepoID(), test.want.GetRepoID())
		}

		if test.build.GetNumber() != test.want.GetNumber() {
			t.Errorf("SetNumber is %v, want %v", test.build.GetNumber(), test.want.GetNumber())
		}

		if test.build.GetParent() != test.want.GetParent() {
			t.Errorf("SetParent is %v, want %v", test.build.GetParent(), test.want.GetParent())
		}

		if test.build.GetEvent() != test.want.GetEvent() {
			t.Errorf("SetEvent is %v, want %v", test.build.GetEvent(), test.want.GetEvent())
		}

		if test.build.GetStatus() != test.want.GetStatus() {
			t.Errorf("SetStatus is %v, want %v", test.build.GetStatus(), test.want.GetStatus())
		}

		if test.build.GetError() != test.want.GetError() {
			t.Errorf("SetError is %v, want %v", test.build.GetError(), test.want.GetError())
		}

		if test.build.GetEnqueued() != test.want.GetEnqueued() {
			t.Errorf("SetEnqueued is %v, want %v", test.build.GetEnqueued(), test.want.GetEnqueued())
		}

		if test.build.GetCreated() != test.want.GetCreated() {
			t.Errorf("SetCreated is %v, want %v", test.build.GetCreated(), test.want.GetCreated())
		}

		if test.build.GetStarted() != test.want.GetStarted() {
			t.Errorf("SetStarted is %v, want %v", test.build.GetStarted(), test.want.GetStarted())
		}

		if test.build.GetFinished() != test.want.GetFinished() {
			t.Errorf("SetFinished is %v, want %v", test.build.GetFinished(), test.want.GetFinished())
		}

		if test.build.GetDeploy() != test.want.GetDeploy() {
			t.Errorf("SetDeploy is %v, want %v", test.build.GetDeploy(), test.want.GetDeploy())
		}

		if !reflect.DeepEqual(test.build.GetDeployPayload(), test.want.GetDeployPayload()) {
			t.Errorf("GetDeployPayload is %v, want %v", test.build.GetDeployPayload(), test.want.GetDeployPayload())
		}

		if test.build.GetClone() != test.want.GetClone() {
			t.Errorf("SetClone is %v, want %v", test.build.GetClone(), test.want.GetClone())
		}

		if test.build.GetSource() != test.want.GetSource() {
			t.Errorf("SetSource is %v, want %v", test.build.GetSource(), test.want.GetSource())
		}

		if test.build.GetTitle() != test.want.GetTitle() {
			t.Errorf("SetTitle is %v, want %v", test.build.GetTitle(), test.want.GetTitle())
		}

		if test.build.GetMessage() != test.want.GetMessage() {
			t.Errorf("SetMessage is %v, want %v", test.build.GetMessage(), test.want.GetMessage())
		}

		if test.build.GetCommit() != test.want.GetCommit() {
			t.Errorf("SetCommit is %v, want %v", test.build.GetCommit(), test.want.GetCommit())
		}

		if test.build.GetSender() != test.want.GetSender() {
			t.Errorf("SetSender is %v, want %v", test.build.GetSender(), test.want.GetSender())
		}

		if test.build.GetAuthor() != test.want.GetAuthor() {
			t.Errorf("SetAuthor is %v, want %v", test.build.GetAuthor(), test.want.GetAuthor())
		}

		if test.build.GetEmail() != test.want.GetEmail() {
			t.Errorf("SetEmail is %v, want %v", test.build.GetEmail(), test.want.GetEmail())
		}

		if test.build.GetLink() != test.want.GetLink() {
			t.Errorf("SetLink is %v, want %v", test.build.GetLink(), test.want.GetLink())
		}

		if test.build.GetBranch() != test.want.GetBranch() {
			t.Errorf("SetBranch is %v, want %v", test.build.GetBranch(), test.want.GetBranch())
		}

		if test.build.GetRef() != test.want.GetRef() {
			t.Errorf("SetRef is %v, want %v", test.build.GetRef(), test.want.GetRef())
		}

		if test.build.GetBaseRef() != test.want.GetBaseRef() {
			t.Errorf("SetBaseRef is %v, want %v", test.build.GetBaseRef(), test.want.GetBaseRef())
		}

		if test.build.GetHeadRef() != test.want.GetHeadRef() {
			t.Errorf("SetHeadRef is %v, want %v", test.build.GetHeadRef(), test.want.GetHeadRef())
		}

		if test.build.GetHost() != test.want.GetHost() {
			t.Errorf("SetHost is %v, want %v", test.build.GetHost(), test.want.GetHost())
		}

		if test.build.GetRuntime() != test.want.GetRuntime() {
			t.Errorf("SetRuntime is %v, want %v", test.build.GetRuntime(), test.want.GetRuntime())
		}

		if test.build.GetDistribution() != test.want.GetDistribution() {
			t.Errorf("SetDistribution is %v, want %v", test.build.GetDistribution(), test.want.GetDistribution())
		}
	}
}

func TestLibrary_Build_String(t *testing.T) {
	// setup types
	b := testBuild()

	want := fmt.Sprintf(`{
  Author: %s,
  BaseRef: %s,
  Branch: %s,
  Clone: %s,
  Commit: %s,
  Created: %d,
  Deploy: %s,
  DeployPayload: %s,
  Distribution: %s,
  Email: %s,
  Enqueued: %d,
  Error: %s,
  Event: %s,
  Finished: %d,
  HeadRef: %s,
  Host: %s,
  ID: %d,
  Link: %s,
  Message: %s,
  Number: %d,
  Parent: %d,
  Ref: %s,
  RepoID: %d,
  Runtime: %s,
  Sender: %s,
  Source: %s,
  Started: %d,
  Status: %s,
  Title: %s,
}`,
		b.GetAuthor(),
		b.GetBaseRef(),
		b.GetBranch(),
		b.GetClone(),
		b.GetCommit(),
		b.GetCreated(),
		b.GetDeploy(),
		b.GetDeployPayload(),
		b.GetDistribution(),
		b.GetEmail(),
		b.GetEnqueued(),
		b.GetError(),
		b.GetEvent(),
		b.GetFinished(),
		b.GetHeadRef(),
		b.GetHost(),
		b.GetID(),
		b.GetLink(),
		b.GetMessage(),
		b.GetNumber(),
		b.GetParent(),
		b.GetRef(),
		b.GetRepoID(),
		b.GetRuntime(),
		b.GetSender(),
		b.GetSource(),
		b.GetStarted(),
		b.GetStatus(),
		b.GetTitle(),
	)

	// run test
	got := b.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testBuild is a test helper function to create a Build
// type with all fields set to a fake value.
func testBuild() *Build {
	b := new(Build)

	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)
	b.SetParent(1)
	b.SetEvent("push")
	b.SetStatus("running")
	b.SetError("")
	b.SetEnqueued(1563474077)
	b.SetCreated(1563474076)
	b.SetStarted(1563474078)
	b.SetFinished(1563474079)
	b.SetDeploy("")
	b.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
	b.SetClone("https://github.com/github/octocat.git")
	b.SetSource("https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163")
	b.SetTitle("push received from https://github.com/github/octocat")
	b.SetMessage("First commit...")
	b.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	b.SetSender("OctoKitty")
	b.SetAuthor("OctoKitty")
	b.SetEmail("OctoKitty@github.com")
	b.SetLink("https://example.company.com/github/octocat/1")
	b.SetBranch("master")
	b.SetRef("refs/heads/master")
	b.SetBaseRef("")
	b.SetHeadRef("changes")
	b.SetHost("example.company.com")
	b.SetRuntime("docker")
	b.SetDistribution("linux")

	return b
}
