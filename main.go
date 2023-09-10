package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"text/template"

	"github.com/spf13/pflag"
	"github.com/togettoyou/hub-mirror/pkg"
)

var (
	content    = pflag.StringP("content", "", "", "原始镜像，格式为：{ \"platform\": \"\", \"hub-mirror\": [] }")
	maxContent = pflag.IntP("maxContent", "", 11, "原始镜像个数限制")
	repository = pflag.StringP("repository", "", "", "推送仓库地址，为空默认为 hub.docker.com")
	username   = pflag.StringP("username", "", "", "仓库用户名")
	password   = pflag.StringP("password", "", "", "仓库密码")
	// outputPath = pflag.StringP("outputPath", "", "output.sh", "结果输出路径")
)

func main() {
	pflag.Parse()

	fmt.Println("验证原始镜像内容")
	var hubMirrors struct {
		Content  []string `json:"hub-mirror"`
		Platform string   `json:"platform"`
	}

	err := json.Unmarshal([]byte(*content), &hubMirrors)
	if err != nil {
		panic(err)
	}

	if len(hubMirrors.Content) > *maxContent {
		panic("提交的原始镜像个数超出了最大限制")
	}

	fmt.Printf("mirrors: %+v, platform: %+v\n", hubMirrors.Content, hubMirrors.Platform)

	fmt.Println("初始化 Docker 客户端")
	old_cli, err := pkg.NewCli(context.Background(), *repository, *username, *password, os.Stdout)
	if err != nil {
		panic(err)
	}

	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("创建 Docker 客户端失败：%v", err)
	}

	// 镜像摘要引用
	imageDigestReference := "gcr.io/oss-fuzz-base/base-builder@sha256:87ca1e9e19235e731fac8de8d1892ebe8d55caf18e7aa131346fc582a2034fdd"

	// 推送镜像 Manifest
	ctx := context.Background()
	options := types.PushOptions{
		RegistryAuth: old_cli.auth,
	}

	_, err = cli.ImagePush(ctx, imageDigestReference, options)
	if err != nil {
		log.Fatalf("推送 Manifest 失败：%v", err)
	}

	fmt.Println("Manifest 推送成功")

	// outputs := make([]*pkg.Output, 0)
	// mu := sync.Mutex{}
	// wg := sync.WaitGroup{}

	// for _, source := range hubMirrors.Content {
	// 	source := source

	// 	if source == "" {
	// 		continue
	// 	}
	// 	fmt.Println("开始转换镜像", source)
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()

	// 		output, err := cli.PullTagPushImage(context.Background(), source, hubMirrors.Platform)
	// 		if err != nil {
	// 			fmt.Println(source, "转换异常: ", err)
	// 			return
	// 		}

	// 		mu.Lock()
	// 		defer mu.Unlock()
	// 		outputs = append(outputs, output)
	// 	}()
	// }

	// wg.Wait()

	// if len(outputs) == 0 {
	// 	panic("没有转换成功的镜像")
	// }

// 	tmpl, err := template.New("pull_images").Parse(
// 		`{{if .Repository}}# if your repository is private,please login...
// # docker login {{ .Repository }} --username={your username}
// {{end}}
// {{- range .Outputs }}
// docker pull {{ .Target }}
// docker tag {{ .Target }} {{ .Source }}{{ end }}`)
// 	if err != nil {
// 		panic(err)
// 	}
// 	outputFile, err := os.Create(*outputPath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer outputFile.Close()

// 	err = tmpl.Execute(outputFile, map[string]interface{}{
// 		"Outputs":    outputs,
// 		"Repository": *repository,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
}

