#!/usr/bin/env python3
"""
Test searching custom fields through the MCP server
"""
import json
import subprocess
import sys

def main():
    search_query = sys.argv[1] if len(sys.argv) > 1 else ""
    
    print(f"Testing search custom fields (query: '{search_query}')...")
    
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
    
    # Search custom fields
    search_req = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/call",
        "params": {
            "name": "search_custom_fields",
            "arguments": {
                "query": search_query
            }
        }
    }
    
    proc = subprocess.Popen(
        ['./start-mcp-server.sh'],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        text=True
    )
    
    proc.stdin.write(json.dumps(init_req) + '\n')
    proc.stdin.write(json.dumps(search_req) + '\n')
    proc.stdin.close()
    
    for line in proc.stdout:
        resp = json.loads(line)
        if resp.get('id') == 2:
            if 'error' in resp:
                print(f"ERROR: {resp['error']}")
            elif 'result' in resp:
                data = json.loads(resp['result']['content'][0]['text'])
                custom_fields = data.get('custom_fields', [])
                print(f"\nFOUND {data.get('count', 0)} MATCHING CUSTOM FIELDS:\n")
                
                for field in custom_fields:
                    print(f"- {field['name']}")
                    print(f"  ID: {field['id']}")
                    print(f"  Type: {field['field_type']}")
                    print(f"  Description: {field.get('description', 'N/A')}")
                    print()

if __name__ == "__main__":
    main()

