package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	regex := regexp.MustCompile(`^.*func (\([a-zA-Z0-9_ *]+\) )?([a-zA-Z0-9_]+)[\[\(]`)
	f, _ := os.ReadFile(os.Args[1])
	usrMsg := strings.Trim(string(f), "\n")

	cmd := "git --no-pager diff --staged"
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	msg := ""
	funcNames := []string{}
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if line == "" || !strings.Contains(line, "func ") {
			continue
		}
		switch line[0] {
		case '@', '+', '-':
			// get second group from regex
			match := regex.FindStringSubmatch(line)
			if len(match) < 2 {
				continue
			}
			funcName := match[2]
			//if funcName is in funcNames, continue
			if Contains(funcNames, funcName) {
				continue
			}
			funcNames = append(funcNames, funcName)

			if line[0] == '+' && i > 0 && lines[i-1][0] == '-' {
				msg += "->" + funcName
				continue
			}
			msg += "\n- " + funcName
		}
	}

	//if usrMsg has one character, replace it with msg
	if len(usrMsg) == 1 {
		os.WriteFile(os.Args[1], []byte(msg), 0644)
		return
	}
	os.WriteFile(os.Args[1], []byte(usrMsg+"\n"+msg), 0644)

}

func Contains(slice []string, e string) bool {
	for _, n := range slice {
		if e == n {
			return true
		}
	}
	return false
}
