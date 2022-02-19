package env

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

var envVarsDockerfile = `
FROM "{{.BaseImage}}"

ENV root_dir="{{.RootDir}}"
ENV start_cmd="{{.StartCmd}}"

WORKDIR "{{.RootDir}}"
`

type EnvVars struct {
	RootDir   string
	StartCmd  string
	BaseImage string
}

func getEnvVarsDockerfile(baseImage string, conf *EnvConfig) (string, error) {
	tmpl, err := template.New("Dockerfile").Parse(envVarsDockerfile)
	if err != nil {
		return "", fmt.Errorf("cannot assemble internal Dockerfile template: %v", err)
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, EnvVars{
		BaseImage: baseImage,
		StartCmd:  conf.StartCmd,
		RootDir:   conf.RootDir,
	})
	if err != nil {
		return "", fmt.Errorf("cannot customize internal Dockerfile: %v", err)
	}

	return buf.String(), nil
}

func BuildEnv(ctx context.Context, client *docker.Client, conf *EnvConfig, dir string) (string, error) {
	outputbuf := bytes.NewBuffer(nil)

	imageName := registryPath + "/" + conf.ID
	imageNameWithoutEnvs := imageName + ":no-envs"

	// Build user's env based on a devbook image
	err := client.BuildImage(docker.BuildImageOptions{
		Name:         imageNameWithoutEnvs,
		OutputStream: outputbuf,
		Context:      ctx,
		ContextDir:   dir,
	})

	// fmt.Println(outputbuf)

	if err != nil {
		return "", fmt.Errorf("cannot build custom env Docker image: %v", err)
	}

	outputbuf.Reset()

	dockerfile, err := getEnvVarsDockerfile(imageNameWithoutEnvs, conf)
	if err != nil {
		return "", fmt.Errorf("cannot assemble internal Dockerfile: %v", err)
	}

	inputbuf := bytes.NewBuffer(nil)

	t := time.Now()
	tr := tar.NewWriter(inputbuf)
	dockerfileBytes := []byte(dockerfile + "\n")
	size := int64(len(dockerfileBytes))
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: size, ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write(dockerfileBytes)
	tr.Close()

	// Build image based on the user's image, injecting Docker env vars so the tinit can access them.
	err = client.BuildImage(docker.BuildImageOptions{
		Name:         imageName,
		OutputStream: outputbuf,
		InputStream:  inputbuf,
		Context:      ctx,
	})

	// fmt.Println(outputbuf)

	if err != nil {
		return "", fmt.Errorf("cannot inject env vars to custom env Docker image: %v", err)
	}

	return imageName, nil
}