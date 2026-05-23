package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := filepath.Walk("c:\\Users\\beatr\\Desktop\\ATLAS.dApp\\cercaend\\lib", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".dart") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			strContent := string(content)
			if strings.Contains(strContent, "firebase_auth/auth_util.dart") {
				strContent = strings.ReplaceAll(strContent, "auth/firebase_auth/auth_util.dart", "auth/auth_util.dart")
				strContent = strings.ReplaceAll(strContent, "../auth/firebase_auth/auth_util.dart", "../auth/auth_util.dart")
				strContent = strings.ReplaceAll(strContent, "firebase_auth/auth_util.dart", "auth_util.dart")
				ioutil.WriteFile(path, []byte(strContent), 0644)
				fmt.Println("Updated:", path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
}
