package env

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/devbookhq/devbookctl/cmd/err"

	docker "github.com/fsouza/go-dockerclient"
)

var registryPath = "us-central1-docker.pkg.dev/devbookhq/devbook-runner-templates"

var baseImageName = "base"
var baseImagePath = registryPath + "/" + "base"
var baseImageAlias = "devbook"

var envVarsDockerfile = `
FROM {{.BaseImage}}

ENV root_dir="{{.RootDir}}"
ENV start_cmd="{{.StartCmd}}"
`

type EnvVars struct {
	RootDir   string
	StartCmd  string
	BaseImage string
}

func getEnvVarsDockerfile(baseImage string, conf *EnvConfig) string {
	tmpl, tmplErr := template.New("Dockerfile").Parse(envVarsDockerfile)
	err.Check(tmplErr)

	buf := new(bytes.Buffer)

	envVars := EnvVars{
		BaseImage: baseImage,
		StartCmd:  conf.StartCmd,
		RootDir:   conf.RootDir,
	}

	exeErr := tmpl.Execute(buf, envVars)
	err.Check(exeErr)

	return buf.String()
}

func BuildEnv(client *docker.Client, conf *EnvConfig) string {
	var buf bytes.Buffer

	auth := docker.AuthConfiguration{}

	// Pull base image path from registry
	pullOpts := docker.PullImageOptions{}

	pullErr := client.PullImage(pullOpts, auth)
	err.Check(pullErr)

	tagOpts := docker.TagImageOptions{}

	// Retag base image from registry as "devbook" image
	tagErr := client.TagImage(baseImagePath, tagOpts)
	err.Check(tagErr)

	imageName := registryPath + "/" + conf.Id

	imageNameWithoutEnvs := imageName + ":no-envs"

	buildOpts := docker.BuildImageOptions{
		Name:         imageNameWithoutEnvs,
		OutputStream: &buf,
	}

	// Build user's env based on a devbook image
	buildErr := client.BuildImage(buildOpts)
	err.Check(buildErr)
	fmt.Println(buf)

	buf.Reset()

	envVarsBuildOpts := docker.BuildImageOptions{
		Name:         imageName,
		OutputStream: &buf,
		Dockerfile:   getEnvVarsDockerfile(imageNameWithoutEnvs, conf),
	}

	// Build image based on the user's image, injecting Docker env vars so the tinit can access them.
	envVarsBuildErr := client.BuildImage(envVarsBuildOpts)
	err.Check(envVarsBuildErr)
	fmt.Println(buf)

	return imageName
}
