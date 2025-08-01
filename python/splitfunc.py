import re
import pyperclip

def parse_go_function(func_code):
    # Extract function name
    func_match = re.search(r'func\s+(\w+)\s*\(', func_code)
    if not func_match:
        return None
    
    func_name = func_match.group(1)
    
    # Find switch statement - need to handle nested braces properly
    switch_start = re.search(r'switch\s+e\s*{', func_code)
    if not switch_start:
        return None
    
    # Find the matching closing brace for the switch
    brace_count = 0
    start_pos = switch_start.end() - 1  # Start at the opening brace
    switch_end = None
    
    for i in range(start_pos, len(func_code)):
        if func_code[i] == '{':
            brace_count += 1
        elif func_code[i] == '}':
            brace_count -= 1
            if brace_count == 0:
                switch_end = i
                break
    
    if switch_end is None:
        return None
    
    switch_body = func_code[switch_start.end():switch_end]
    
    # Parse cases by finding case keywords and extracting content between them
    cases = []
    case_starts = list(re.finditer(r'case\s+(\d+):', switch_body))
    
    for i, match in enumerate(case_starts):
        case_num = match.group(1)
        start_pos = match.end()
        
        # Find the end position (next case or end of switch)
        if i + 1 < len(case_starts):
            end_pos = case_starts[i + 1].start()
        else:
            end_pos = len(switch_body)
        
        case_body = switch_body[start_pos:end_pos].strip()
        cases.append((case_num, case_body))
    
    output_functions = []
    
    for case_num, case_body in cases:
        # Skip empty cases
        if not case_body:
            continue
            
        # Check if case ends with return
        lines = [line.strip() for line in case_body.split('\n') if line.strip()]
        
        if lines and lines[-1].startswith('return'):
            return_stmt = lines[-1]
            body_lines = lines[:-1]
        else:
            return_stmt = "return false"
            body_lines = lines
        
        # Build function
        func_def = f"func {func_name}_{case_num}(gb *Gameboy) bool {{\n"
        
        for line in body_lines:
            func_def += f"\t{line}\n"
        
        func_def += f"\t{return_stmt}\n"
        func_def += "}\n"
        
        output_functions.append(func_def)
    
    return '\n'.join(output_functions)

def parse_multiple_functions(code):
    # Split by function boundaries more carefully
    results = []
    
    # Find all function starts
    func_starts = list(re.finditer(r'func\s+\w+\s*\([^)]*\)\s*[^{]*{', code))
    
    for i, match in enumerate(func_starts):
        start_pos = match.start()
        
        # Find the matching closing brace
        brace_count = 0
        brace_start = match.end() - 1  # Position of opening brace
        
        for j in range(brace_start, len(code)):
            if code[j] == '{':
                brace_count += 1
            elif code[j] == '}':
                brace_count -= 1
                if brace_count == 0:
                    func_code = code[start_pos:j+1]
                    result = parse_go_function(func_code)
                    if result:
                        results.append(result)
                    break
    
    return '\n'.join(results) if results else None

# Read input
print("Paste your Go function (press Enter twice when done):")
lines = []
empty_lines = 0
while True:
    line = input()
    if line == "":
        empty_lines += 1
        if empty_lines >= 2:
            break
    else:
        empty_lines = 0
        lines.append(line)

code = '\n'.join(lines)

result = parse_multiple_functions(code)
if result:
    pyperclip.copy(result)
    print("Converted functions copied to clipboard!")
else:
    print("Could not parse any functions")