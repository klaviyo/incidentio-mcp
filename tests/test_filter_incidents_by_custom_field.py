#!/usr/bin/env python3
"""
Test filtering incidents by custom field (e.g., affected team)
Example workflow: Find custom field ID, then filter incidents
"""
import json
import subprocess
import sys

def main():
    # Get custom field name from command line or use default
    field_name = sys.argv[1] if len(sys.argv) > 1 else "team"
    field_value = sys.argv[2] if len(sys.argv) > 2 else "Engineering"
    
    print(f"Searching for custom field: '{field_name}'")
    print(f"Will filter incidents where {field_name} = '{field_value}'\n")
    
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
    
    # Step 1: Search for the custom field
    search_req = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/call",
        "params": {
            "name": "search_custom_fields",
            "arguments": {"query": field_name}
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
    
    custom_field_id = None
    
    # Read responses for steps 1-2
    for i in range(2):
        line = proc.stdout.readline()
        if not line:
            break
        resp = json.loads(line)
        
        if resp.get('id') == 2:
            if 'result' in resp:
                data = json.loads(resp['result']['content'][0]['text'])
                fields = data.get('custom_fields', [])
                if fields:
                    custom_field_id = fields[0]['id']
                    print(f"✓ Found custom field: {fields[0]['name']} (ID: {custom_field_id})\n")
                else:
                    print(f"✗ No custom field found matching '{field_name}'")
                    print("  Use list_custom_fields to see all available fields")
                    proc.stdin.close()
                    return
    
    if not custom_field_id:
        print("Could not find custom field")
        proc.stdin.close()
        return
    
    # Step 2: Filter incidents by this custom field
    print(f"Filtering incidents where custom field = '{field_value}'...\n")
    
    filter_req = {
        "jsonrpc": "2.0",
        "id": 3,
        "method": "tools/call",
        "params": {
            "name": "list_incidents",
            "arguments": {
                "custom_field_id": custom_field_id,
                "custom_field_value": field_value,
                "page_size": 10
            }
        }
    }
    
    proc.stdin.write(json.dumps(filter_req) + '\n')
    proc.stdin.close()
    
    # Read response for step 3
    for line in proc.stdout:
        resp = json.loads(line)
        if resp.get('id') == 3:
            if 'error' in resp:
                print(f"ERROR: {resp['error']}")
            elif 'result' in resp:
                data = json.loads(resp['result']['content'][0]['text'])
                incidents = data.get('incidents', [])
                print(f"FOUND {len(incidents)} INCIDENTS:\n")
                
                for inc in incidents:
                    print(f"- {inc.get('name')}")
                    print(f"  Status: {inc.get('incident_status', {}).get('name')}")
                    print(f"  Created: {inc.get('created_at')}")
                    print()

if __name__ == "__main__":
    main()

