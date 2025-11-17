#!/usr/bin/env python3
"""
Quick test to verify tool names are registered correctly
"""
import json
import subprocess

print("Testing MCP server tool registration...\n")

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

# List tools
list_req = {
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
}

proc = subprocess.Popen(
    ['./bin/mcp-server'],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
    text=True
)

proc.stdin.write(json.dumps(init_req) + '\n')
proc.stdin.write(json.dumps(list_req) + '\n')
proc.stdin.close()

for line in proc.stdout:
    resp = json.loads(line)
    if resp.get('id') == 2 and 'result' in resp:
        tools = resp['result']['tools']
        print(f"Total tools registered: {len(tools)}\n")
        
        # Show all tool names
        print("All tool names:")
        for tool in tools:
            print(f"  - {tool['name']}")
        
        # Highlight custom field tools
        print("\nCustom field tools:")
        custom_tools = [t for t in tools if 'custom' in t['name'].lower()]
        for tool in custom_tools:
            print(f"  ✓ {tool['name']}")
        
        if not custom_tools:
            print("  ❌ NO CUSTOM FIELD TOOLS FOUND!")
        break

proc.terminate()

