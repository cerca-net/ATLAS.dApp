package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func main() {
	path := `c:\Users\beatr\Desktop\ATLAS\cercachain-fix-userpage-syntax-errors\cercachain-fix-userpage-syntax-errors\cercaend\lib\components\object\object_widget.dart`

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	text := string(bytes)

	reUpvote := regexp.MustCompile(`await containerCreditsRecord!\s*\n\s*\.reference\s*\n\s*\.update\(\{\s*\n\s*\.\.\.mapToFirestore\(\s*\n\s*\{\s*\n\s*'participation_i':\s*\n\s*FieldValue\s*\n\s*\.increment\(\s*\n\s*0\.01\),\s*\n\s*\},\s*\n\s*\),\s*\n\s*\}\);`)

	// We deduct 1.0 from the user's credits due to interaction (DUs deduction).
	replacement := `// Deduct 1 DU from the user interacting
                                                          await toggleCreditsRecord!
                                                              .reference
                                                              .update({
                                                            ...mapToFirestore(
                                                              {
                                                                'f_x': FieldValue.increment(-1.0),
                                                                'participation_i': FieldValue.increment(0.01),
                                                              },
                                                            ),
                                                          });`

	newText := reUpvote.ReplaceAllString(text, replacement)

	// Actually we need to search and replace containerCreditsRecord -> toggleCreditsRecord

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
