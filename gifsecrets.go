package gifsecrets

import (
	"bufio"
	"code.google.com/p/go-uuid/uuid"
	"io"
	"os"
	"os/exec"
	"regexp"
)

// Encode stores a uuid in the comment section of a gif
func Encode(path, id string) error {
	if id == "" {
		id = uuid.New()
	}

	cmd := exec.Command("/usr/local/bin/gifsicle", path, "--no-comments", "-c", "\""+id+"\"", "#0", "#1-")

	outfile, err := os.Create("./out.gif")
	if err != nil {
		return err
	}
	defer outfile.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(outfile)
	if err := cmd.Start(); err != nil {
		return err
	}

	go io.Copy(writer, stdoutPipe)
	cmd.Wait()
	return nil
}

// Decode returns the secret encoded in the gif comment extension
func Decode(path string) (string, error) {
	cmd := exec.Command("/usr/local/bin/gifsicle", "-I", path, "#0")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile("\"(.*?)\"") // this captures everything in between quotes which is pretty ghetto
	raw := reg.Find(out)
	return string(raw[1:len(raw)-1]), nil
}
