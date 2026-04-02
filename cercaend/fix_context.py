import re
import os

def fix_context_warnings():
    output_files = {}
    with open('analyze_output.txt', 'r', encoding='utf-8') as f:
        lines = f.readlines()
        
    for line in lines:
        if 'use_build_context_synchronously' in line:
            # Parse line:  info - lib\mainpages\orderpage\orderpage_widget.dart:608:18 - ...
            parts = line.split(' - ')
            if len(parts) >= 3:
                file_info = parts[1].strip()
                if ':' in file_info:
                    filepath, line_num, col_num = file_info.rsplit(':', 2)
                    line_num = int(line_num) - 1 # 0-indexed
                    if filepath not in output_files:
                        output_files[filepath] = set()
                    output_files[filepath].add(line_num)
    
    for filepath, line_nums in output_files.items():
        if not os.path.exists(filepath):
            continue
            
        with open(filepath, 'r', encoding='utf-8') as f:
            lines = f.readlines()
            
        # Process from bottom to top to not mess up line numbers
        for ln in sorted(list(line_nums), reverse=True):
            if ln < len(lines):
                target_line = lines[ln]
                indent = len(target_line) - len(target_line.lstrip())
                indent_str = ' ' * indent
                
                # Check if we already have a mounted check
                already_has_mounted = False
                for i in range(1, 4):
                    if ln - i >= 0:
                        prev_line = lines[ln - i].strip()
                        if 'mounted' in prev_line:
                            already_has_mounted = True
                            break
                        if prev_line and prev_line != '}' and prev_line != '{':
                            break # stop checking if we hit a real line without mounted
                
                if not already_has_mounted:
                    # Insert context.mounted
                    lines.insert(ln, indent_str + 'if (!context.mounted) return;\n')
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.writelines(lines)
            
    print(f"Fixed files: {len(output_files)}")

if __name__ == '__main__':
    fix_context_warnings()
