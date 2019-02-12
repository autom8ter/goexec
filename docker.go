package goexec

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/gofunct/goexec/pkg"
	"io"
	"log"
	"os"
	"time"
)

type DkrExecConfigFunc func(config types.ExecConfig)
type DkrImageBuildConfigFunc func(opts types.ImageBuildOptions)
type DkrAuthConfigFunc func(opts types.AuthConfig)

func (c *Command) ListContainers() {
	containers, err := c.dkr.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func (c *Command) DockerVersion() string {
	return c.dkr.ClientVersion()
}

func (c *Command) AttachContainer(ctx context.Context, container string)  error {
	hijak, err := c.dkr.ContainerAttach(ctx, container, types.ContainerAttachOptions{
		Stream:     true,
		Stdin:      true,
		Stdout:     true,
		Stderr:     true,
		Logs:       true,
	})
	if err != nil {
		return err
	}
	defer hijak.Close()
	if _, err := pkg.StdCopy(os.Stdout, os.Stderr, hijak.Reader); err != nil {
		return err
	}
	return nil
}

func (c *Command) CommitContainer(ctx context.Context, container string, author string, comment string) (string, error) {
	id, err := c.dkr.ContainerCommit(ctx, container, types.ContainerCommitOptions{
		Comment:   comment,
		Author:    author,
	})
	return id.ID, err
}

func (c *Command) CreateContainer(ctx context.Context, name string) (string, error) {
	body, err := c.dkr.ContainerCreate(ctx, &container.Config{

	}, &container.HostConfig{

	}, &network.NetworkingConfig{

	}, name,
	)
	if len(body.Warnings) > 0 {
		for _, warn := range body.Warnings {
			log.Println(warn)
		}
	}
	return body.ID, err
}

func (c *Command) DiffContainer(ctx context.Context, name string) {
	chgs, err := c.dkr.ContainerDiff(ctx, name)
	for _, chg := range chgs {
		c.Printf("Kind: %v Path: %s", chg.Kind, chg.Path)
	}
	if err != nil {
		panic(err)
	}
}

func (c *Command) ExecAttachContainer(ctx context.Context, tty, detach bool, execId, name string) error {
	hijak, err := c.dkr.ContainerExecAttach(ctx, execId, types.ExecConfig{
	Detach:	detach,
	Tty: tty,
	})
	if err != nil {
		return err
	}
	defer hijak.Close()
	if _, err := pkg.StdCopy(os.Stdout, os.Stderr, hijak.Reader); err != nil {
		return err
	}
	return nil
}

func (c *Command) CreateExecContainer(ctx context.Context, name string, opts ...ExecConfigFunc) error {
	exconfig := types.ExecConfig{}
	for _, o := range opts {
		o(exconfig)
	}
	id, err := c.dkr.ContainerExecCreate(ctx, name, exconfig)
	if err != nil {
		return err
	}
	log.Printf("ID: %s", id)
	return nil
}

func (c *Command) ExecInspect(ctx context.Context, id string) error {
	res, err := c.dkr.ContainerExecInspect(ctx, id)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Inspection Results: \n%s", res))
	return nil
}

func (c *Command) ExportContainer(ctx context.Context, id string) (io.ReadCloser, error) {
	return c.dkr.ContainerExport(ctx, id)
}

func (c *Command) ShutdownContainer(ctx context.Context, id string) error {
	to := 60 *time.Second
	return c.dkr.ContainerStop(ctx, id, &to)
}


func (c *Command) RemoveContainer(ctx context.Context, id string, vol, links, force bool) error {
	return c.dkr.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
	RemoveLinks: links,
	RemoveVolumes: vol,
	Force: force,
	})
}

func (c *Command) RestartContainer(ctx context.Context, id string) error {
	to := 60 *time.Second
	return c.dkr.ContainerRestart(ctx, id, &to)
}

func (c *Command) CopyFromContainer(ctx context.Context, id string, path string) (io.ReadCloser, error) {
	reader, _, err := c.dkr.CopyFromContainer(ctx, id, path)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (c *Command) CopyToContainer(ctx context.Context, id string, path string, reader io.Reader, overwrite bool) error {
	return c.dkr.CopyToContainer(ctx, id, path, reader, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: overwrite,
	})
}

func (c *Command) CloseDkrClient() error {
	return c.dkr.Close()
}

func (c *Command) DiscUsage(ctx context.Context) error {
	usg, err := c.dkr.DiskUsage(ctx)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("%s", usg))
	return nil
}

func (c *Command) BuildImage(ctx context.Context, reader io.Reader, opts ...DkrImageBuildConfigFunc) (io.ReadCloser, error) {
	cfg := types.ImageBuildOptions{}
	for _, o := range opts {
		o(cfg)
	}
	resp, err := c.dkr.ImageBuild(ctx, reader, cfg)
	if err != nil {
		return resp.Body, err
	}

	return resp.Body, nil
}


