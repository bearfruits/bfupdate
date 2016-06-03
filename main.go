package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	var pemPath string
	var host string
	var dir string
	flag.StringVar(&pemPath, "i", "~/.ssh/id_rsa", "through to ssh -i")
	flag.StringVar(&host, "h", "", "user@host")
	flag.StringVar(&dir, "d", "", "working dir")
	flag.Parse()

	err := update(host, pemPath, dir)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func update(host, pemPath, dir string) error {
	cmd := fmt.Sprintf(`
bash -xe -c '
	cd %s
	git checkout master
	git pull
	bin/update
'
`, dir)
	c := exec.Command("ssh", host, "-t", "-t", "-i", pemPath, "-oStrictHostKeyChecking=no", cmd)
	buf := bytes.NewBuffer([]byte{})
	c.Stdout = buf
	c.Stderr = buf

	if err := c.Run(); err != nil {
		return fmt.Errorf("ssh to %s is failed.\n\nError:\n%s\n\nOutput:\n\n%s", host, err, buf.String())
	}
	return nil
}
