#!/usr/bin/env python3
"""
Test listing custom fields through the MCP server
"""
import json
import subprocess

def main():
    print("Testing list custom fields...")
    
    # Initialize
    init_req = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "initialize",
        "params": {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "test", "version": "1.0"}
        }
    }
    
    # List custom fields
    list_req = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/call",
        "params": {
            "name": "list_custom_fields",
            "arguments": {}
        }
    }
    
    proc = subprocess.Popen(
        ['./start-mcp-server.sh'],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        text=True
    )
    
    proc.stdin.write(json.dumps(init_req) + '\n')
    proc.stdin.write(json.dumps(list_req) + '\n')
    proc.stdin.close()
    
    for line in proc.stdout:
        resp = json.loads(line)
        if resp.get('id') == 2:
            if 'error' in resp:
                print(f"ERROR: {resp['error']}")
            elif 'result' in resp:
                print("CUSTOM FIELDS:")
                data = json.loads(resp['result']['content'][0]['text'])
                custom_fields = data.get('custom_fields', [])
                print(f"Found {len(custom_fields)} custom fields\n")
                
                for field in custom_fields[:5]:  # Show first 5
                    print(f"- {field['name']}")
                    print(f"  Type: {field['field_type']}")
                    print(f"  Required: {field['required']}")
                    if field.get('options'):
                        opts = [opt['value'] for opt in field['options'][:3]]
                        print(f"  Options: {', '.join(opts)}...")
                    print()
                
                if len(custom_fields) > 5:
                    print(f"... and {len(custom_fields) - 5} more")

if __name__ == "__main__":
    main()

