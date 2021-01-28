package pipeline

import (
	"context"

	"github.com/vomb/x/dir"
	"github.com/vomb/x/git"
	"github.com/vomb/x/proto"
)

func API(ctx context.Context) error {
	repo := "github.com/vomb/baikal"

	gitChange, err := git.Subscribe(ctx, git.Location{
		Git:      repo,
		PathGlob: "api/proto/**/*.proto",
	})
	if err != nil {
		return err
	}

	protoFiles, err := dir.FromGit(ctx, dir.GitLocation{
		Path: "api/proto",
		Git:  gitChange.Git,
		Ref:  gitChange.Ref,
	})
	if err != nil {
		return err
	}

	goFiles, err := proto.ToGo(ctx, proto.ToGoRequest{
		Location: proto.Location{
			Path: protoFiles.Path,
		},
		Package: "api",
	})
	if err != nil {
		return err
	}

	err = dir.ToGit(ctx, dir.ToGitRequest{
		GitLocation: dir.GitLocation{
			Git:  repo,
			Path: "api",
		},
		DeletePattern: "**/*.go",
		FromPath:      goFiles.Path,
	})

	if err != nil {
		return err
	}

	return nil
}
