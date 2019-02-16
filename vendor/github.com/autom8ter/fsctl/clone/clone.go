package clone

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"log"
	"os"
)

type CloneFunc func(options *git.CloneOptions)


func (c CloneFunc) MemClone() *git.Repository {
	opts := &git.CloneOptions{}
	c(opts)
	log.Printf("Cloning: %s",  opts.URL)
	r, err := git.Clone(memory.NewStorage(), nil, opts)
	ifErr(err)
	return r
}

func (c CloneFunc) Clone(dir string) *git.Repository {
	opts := &git.CloneOptions{}
	c(opts)
	log.Printf("Cloning: %s",  opts.URL)

	r, err := git.PlainClone(dir, false, opts)
	ifErr(err)

	// Clones the repository into the given dir, just as a normal git clone does
	return r
}

func NewAuthCloner() CloneFunc {
	return func(options *git.CloneOptions) {
		options.URL = os.Getenv("CFGURL")
		options.Auth = &http.BasicAuth{
			Username: os.Getenv("GITUSER"),
			Password: os.Getenv("CFGTOKEN"),
		}
	}
}

func ifErr(err error) {
	if err != nil {
		panic(err)
	}
}