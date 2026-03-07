import re
import sys

path = r'c:\Users\beatr\Desktop\ATLAS\cercachain-fix-userpage-syntax-errors\cercachain-fix-userpage-syntax-errors\cercaend\lib\mainpages\newpage\newpage_widget.dart'

with open(path, 'r', encoding='utf-8') as f:
    text = f.read()

pattern = re.compile(r"'thread':\s*_model\.choiceChipsThread([a-zA-Z]+)Values,")
replacement = r"'thread': [\n                                                                            ...(_model.choiceChipsThread\g<1>Values ?? []),\n                                                                            ...(currentUserDocument?.userOccupations ?? []),\n                                                                            ...(currentUserDocument?.userInterests ?? [])\n                                                                        ],"

new_text, count = pattern.subn(replacement, text)

with open(path, 'w', encoding='utf-8') as f:
    f.write(new_text)

print(f'Replaced {count} occurrences.')
