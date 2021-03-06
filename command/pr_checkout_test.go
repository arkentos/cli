package command

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/cli/cli/context"
	"github.com/cli/cli/utils"
)

func TestPRCheckout_sameRepo(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	ctx.SetRemotes(map[string]string{
		"origin": "OWNER/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequest": {
		"number": 123,
		"headRefName": "feature",
		"headRepositoryOwner": {
			"login": "hubot"
		},
		"headRepository": {
			"name": "REPO",
			"defaultBranchRef": {
				"name": "master"
			}
		},
		"isCrossRepository": false,
		"maintainerCanModify": false
	} } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git show-ref --verify --quiet refs/heads/pr/123/feature":
			return &errorStub{"exit status: 1"}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout 123`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 4)
	eq(t, strings.Join(ranCommands[0], " "), "git fetch origin +refs/heads/feature:refs/remotes/origin/feature")
	eq(t, strings.Join(ranCommands[1], " "), "git checkout -b pr/123/feature --no-track origin/feature")
	eq(t, strings.Join(ranCommands[2], " "), "git config branch.pr/123/feature.remote origin")
	eq(t, strings.Join(ranCommands[3], " "), "git config branch.pr/123/feature.merge refs/heads/feature")
}

func TestPRCheckout_urlArg(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	ctx.SetRemotes(map[string]string{
		"origin": "OWNER/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequest": {
		"number": 123,
		"headRefName": "feature",
		"headRepositoryOwner": {
			"login": "hubot"
		},
		"headRepository": {
			"name": "REPO",
			"defaultBranchRef": {
				"name": "master"
			}
		},
		"isCrossRepository": false,
		"maintainerCanModify": false
	} } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git show-ref --verify --quiet refs/heads/pr/123/feature":
			return &errorStub{"exit status: 1"}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout https://github.com/OWNER/REPO/pull/123/files`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 4)
	eq(t, strings.Join(ranCommands[1], " "), "git checkout -b pr/123/feature --no-track origin/feature")
}

func TestPRCheckout_branchArg(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	ctx.SetRemotes(map[string]string{
		"origin": "OWNER/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequests": { "nodes": [
		{ "number": 123,
		  "headRefName": "feature",
		  "headRepositoryOwner": {
		  	"login": "hubot"
		  },
		  "headRepository": {
		  	"name": "REPO",
		  	"defaultBranchRef": {
		  		"name": "master"
		  	}
		  },
		  "isCrossRepository": true,
		  "maintainerCanModify": false }
	] } } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git show-ref --verify --quiet refs/heads/pr/123/feature":
			return &errorStub{"exit status: 1"}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout hubot:feature`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 5)
	eq(t, strings.Join(ranCommands[1], " "), "git fetch origin refs/pull/123/head:pr/123/feature")
}

func TestPRCheckout_existingBranch(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	ctx.SetRemotes(map[string]string{
		"origin": "OWNER/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequest": {
		"number": 123,
		"headRefName": "feature",
		"headRepositoryOwner": {
			"login": "hubot"
		},
		"headRepository": {
			"name": "REPO",
			"defaultBranchRef": {
				"name": "master"
			}
		},
		"isCrossRepository": false,
		"maintainerCanModify": false
	} } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git show-ref --verify --quiet refs/heads/pr/123/feature":
			return &outputStub{}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout 123`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 3)
	eq(t, strings.Join(ranCommands[0], " "), "git fetch origin +refs/heads/feature:refs/remotes/origin/feature")
	eq(t, strings.Join(ranCommands[1], " "), "git checkout pr/123/feature")
	eq(t, strings.Join(ranCommands[2], " "), "git merge --ff-only refs/remotes/origin/feature")
}

func TestPRCheckout_differentRepo_remoteExists(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	ctx.SetRemotes(map[string]string{
		"origin":     "OWNER/REPO",
		"robot-fork": "hubot/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequest": {
		"number": 123,
		"headRefName": "feature",
		"headRepositoryOwner": {
			"login": "hubot"
		},
		"headRepository": {
			"name": "REPO",
			"defaultBranchRef": {
				"name": "master"
			}
		},
		"isCrossRepository": true,
		"maintainerCanModify": false
	} } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git show-ref --verify --quiet refs/heads/pr/123/feature":
			return &errorStub{"exit status: 1"}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout 123`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 4)
	eq(t, strings.Join(ranCommands[0], " "), "git fetch robot-fork +refs/heads/feature:refs/remotes/robot-fork/feature")
	eq(t, strings.Join(ranCommands[1], " "), "git checkout -b pr/123/feature --no-track robot-fork/feature")
	eq(t, strings.Join(ranCommands[2], " "), "git config branch.pr/123/feature.remote robot-fork")
	eq(t, strings.Join(ranCommands[3], " "), "git config branch.pr/123/feature.merge refs/heads/feature")
}

func TestPRCheckout_differentRepo(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	ctx.SetRemotes(map[string]string{
		"origin": "OWNER/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequest": {
		"number": 123,
		"headRefName": "feature",
		"headRepositoryOwner": {
			"login": "hubot"
		},
		"headRepository": {
			"name": "REPO",
			"defaultBranchRef": {
				"name": "master"
			}
		},
		"isCrossRepository": true,
		"maintainerCanModify": false
	} } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git config branch.pr/123/feature.merge":
			return &errorStub{"exit status 1"}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout 123`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 4)
	eq(t, strings.Join(ranCommands[0], " "), "git fetch origin refs/pull/123/head:pr/123/feature")
	eq(t, strings.Join(ranCommands[1], " "), "git checkout pr/123/feature")
	eq(t, strings.Join(ranCommands[2], " "), "git config branch.pr/123/feature.remote origin")
	eq(t, strings.Join(ranCommands[3], " "), "git config branch.pr/123/feature.merge refs/pull/123/head")
}

func TestPRCheckout_differentRepo_existingBranch(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	ctx.SetRemotes(map[string]string{
		"origin": "OWNER/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequest": {
		"number": 123,
		"headRefName": "feature",
		"headRepositoryOwner": {
			"login": "hubot"
		},
		"headRepository": {
			"name": "REPO",
			"defaultBranchRef": {
				"name": "master"
			}
		},
		"isCrossRepository": true,
		"maintainerCanModify": false
	} } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git config branch.pr/123/feature.merge":
			return &outputStub{[]byte("refs/heads/feature\n")}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout 123`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 2)
	eq(t, strings.Join(ranCommands[0], " "), "git fetch origin refs/pull/123/head:pr/123/feature")
	eq(t, strings.Join(ranCommands[1], " "), "git checkout pr/123/feature")
}

func TestPRCheckout_differentRepo_currentBranch(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("pr/123/feature")
	ctx.SetRemotes(map[string]string{
		"origin": "OWNER/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequest": {
		"number": 123,
		"headRefName": "feature",
		"headRepositoryOwner": {
			"login": "hubot"
		},
		"headRepository": {
			"name": "REPO",
			"defaultBranchRef": {
				"name": "master"
			}
		},
		"isCrossRepository": true,
		"maintainerCanModify": false
	} } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git config branch.pr/123/feature.merge":
			return &outputStub{[]byte("refs/heads/feature\n")}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout 123`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 2)
	eq(t, strings.Join(ranCommands[0], " "), "git fetch origin refs/pull/123/head")
	eq(t, strings.Join(ranCommands[1], " "), "git merge --ff-only FETCH_HEAD")
}

func TestPRCheckout_maintainerCanModify(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	ctx.SetRemotes(map[string]string{
		"origin": "OWNER/REPO",
	})
	initContext = func() context.Context {
		return ctx
	}
	http := initFakeHTTP()

	http.StubResponse(200, bytes.NewBufferString(`
	{ "data": { "repository": { "pullRequest": {
		"number": 123,
		"headRefName": "feature",
		"headRepositoryOwner": {
			"login": "hubot"
		},
		"headRepository": {
			"name": "REPO",
			"defaultBranchRef": {
				"name": "master"
			}
		},
		"isCrossRepository": true,
		"maintainerCanModify": true
	} } } }
	`))

	ranCommands := [][]string{}
	restoreCmd := utils.SetPrepareCmd(func(cmd *exec.Cmd) utils.Runnable {
		switch strings.Join(cmd.Args, " ") {
		case "git config branch.pr/123/feature.merge":
			return &errorStub{"exit status 1"}
		default:
			ranCommands = append(ranCommands, cmd.Args)
			return &outputStub{}
		}
	})
	defer restoreCmd()

	output, err := RunCommand(prCheckoutCmd, `pr checkout 123`)
	eq(t, err, nil)
	eq(t, output.String(), "")

	eq(t, len(ranCommands), 4)
	eq(t, strings.Join(ranCommands[0], " "), "git fetch origin refs/pull/123/head:pr/123/feature")
	eq(t, strings.Join(ranCommands[1], " "), "git checkout pr/123/feature")
	eq(t, strings.Join(ranCommands[2], " "), "git config branch.pr/123/feature.remote https://github.com/hubot/REPO.git")
	eq(t, strings.Join(ranCommands[3], " "), "git config branch.pr/123/feature.merge refs/heads/feature")
}
