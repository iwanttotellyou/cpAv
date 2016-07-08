package main

import (
    "fmt"
    "path/filepath"
    "log"
    "strings"
    "os"
    . "cpAv/file"
    "regexp"
    "errors"
)

type file struct {
    path string
    dir  string
    name string
}

type files []file
type array []string

func pwd(p string) string {
    dir, err := filepath.Abs(p)
    if err != nil {
        log.Fatal(err)
    }
    return strings.Replace(dir, "\\", "/", -1)
}

func (f *files) Files(fileName string) error {
    filepath.Walk(pwd(fileName), func(path string, info os.FileInfo, err error) error {
        dir, name := filepath.Split(path)
        *f = append(*f, file{
            path: path,
            dir: dir,
            name: name,
        })
        return nil
    })
    return nil
}

func (arr array) getArgv(index int) string {
    if index + 1 > len(arr) {
        return ""
    }
    return arr[index]
}

func (f files) Av(match string) (files, error) {
    reg, err := regexp.Compile(match)
    if err != nil {
        return nil, errors.New("regex syntax is wrong")
    }
    tempF := files{}
    for _, i := range f {
        if reg.Match([]byte(i.name)) {
            tempF = append(tempF, i)
        }
    }
    return tempF, nil
}

func Av(origin, target, regex string) error {
    f := files{}
    err := f.Files(origin)
    if err != nil {
        return errors.New("Files error")
    }
    if regex != "" {
        f, err = f.Av(regex)
        if err != nil {
            return errors.New("Av error")
        }
    }
    for _, i := range f {
        if !isExist(target) {
            fmt.Printf("mkdir: %s", target)
            os.Mkdir(target, 0755)
        }
        CopyFile(i.path, filepath.Join(target, i.name))
        fmt.Println(i.path, " : ", filepath.Join(target, i.name))
    }
    return nil
}

func isExist(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil || os.IsExist(err)
}

func main() {
    var arr array = os.Args
    Av(arr.getArgv(1), arr.getArgv(2), arr.getArgv(3))
}
