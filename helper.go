package infowatch_server

import (
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/Pilladian/logger"
)

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello There!</h1><p>This is the main page.</p>")
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func uniqueRandomString(id_len int, files []fs.FileInfo) string {
	cond := true
	filename := randomString(id_len)
	for cond {
		unique := true
		for _, f := range files {
			if filename == f.Name() {
				unique = false
			}
		}
		if !unique {
			filename = randomString(id_len)
		} else {
			cond = false
		}
	}
	return filename
}

func createPath(p string) {
	existent, err := exists(p)
	if err != nil {
		logger.Error(err.Error())
	}

	if !existent {
		split_path := strings.Split(p, "/")
		createPath(strings.Join(split_path[0:len(split_path)-1], "/"))

		err := os.Mkdir(p, 0700)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
