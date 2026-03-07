package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func main() {
	path := `c:\Users\beatr\Desktop\ATLAS\cercachain-fix-userpage-syntax-errors\cercachain-fix-userpage-syntax-errors\cercaend\lib\mainpages\newpage\newpage_widget.dart`

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	text := string(bytes)

	re := regexp.MustCompile(`'thread':\s*_model\.choiceChipsThread([a-zA-Z]+)Values,`)

	replacement := `'thread': [
                                                                            ...(_model.choiceChipsThread${1}Values ?? []),
                                                                            ...(currentUserDocument?.userOccupations ?? []),
                                                                            ...(currentUserDocument?.userInterests ?? [])
                                                                        ],`

	newText := re.ReplaceAllString(text, replacement)

	if text == newText {
		fmt.Println("No replacements made.")
		return
	}

	err = ioutil.WriteFile(path, []byte(newText), 0644)
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}

	fmt.Println("Successfully replaced occurrences.")
}
