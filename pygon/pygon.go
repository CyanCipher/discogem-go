package pygon

import (
	//"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func GenImage(prompt string) (bool, error) {

	// WRITE PROMPT TO A TEXT FILE
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file, error := os.Create(filepath.Join(cwd, "Media", "prompt.txt"))
	if error != nil {
		log.Fatal(error)
	}
	defer file.Close()

	_, errx := file.WriteString(prompt)
	if errx != nil {
		log.Fatal(err)
	}

	// CALL THE PYTHON SCRIPT

	cmd := exec.Command("python", filepath.Join(cwd, "pygon", "imagegen.py"))
	fmt.Println("calling the python file")

	_, errr := cmd.CombinedOutput()
	if errr != nil {
		log.Fatal(errr)
	}

	_, er := os.Stat(filepath.Join(cwd, "Media", "prompt.txt"))
	if er == nil {
		return false, er
	} else if os.IsNotExist(er) {
		return true, nil
	} else {
		return true, nil
	}

}
