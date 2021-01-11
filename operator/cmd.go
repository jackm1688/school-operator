package operator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	opSdk   = "/usr/local/bin/operator-sdk"
	docker  = "/usr/local/bin/docker"
	kubelet = "/usr/local/bin/kubectl"
	sed     = "/usr/bin/sed"
)

//镜像替换
func replaceImage(srcImage string, dstImage string) error {
	//sed -i "" 's|REPLACE_IMAGE|quay.io/example/memcached-operator:v0.0.1|g' deploy/operator.yaml
	r := fmt.Sprintf(`'s|%s|%s|g'`, srcImage, dstImage)
	return execCmd(sed, "-i", `""`, r, "deploy/operator.yaml")
}

//获取operator.yaml中镜像信息
func ReadImage() (error, string) {
	data, err := ioutil.ReadFile("deploy/operator.yaml")
	if err != nil {
		return err, ""
	}
	r := regexp.MustCompile(`image:\s+(.*)`)
	v := r.FindAllString(string(data), 1)
	s := ""
	if len(v) > 0 {
		s = strings.TrimSpace(strings.Split(v[0], " ")[1])
	}
	return nil, s

}

//执行命令
func execCmd(c string, args ...string) error {
	dir, _ := os.Getwd()
	fmt.Println("--------current dir:", dir, "------------")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(c, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("can't start for the %s cmd,err:%v", opSdk, err)
	}

	time.Sleep(5 * time.Second)
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("can't wait for the %s, err:%v", opSdk, err)
	}
	fmt.Println(stdout.String())
	return nil
}
