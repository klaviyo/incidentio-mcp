#!/usr/bin/env python3
"""
Test filtering incidents by time range
Example: Get incidents from the last 7 days
"""
import json
import subprocess
from datetime import datetime, timedelta

def main():
    # Calculate time range (last 7 days)
    now = datetime.utcnow()
    seven_days_ago = now - timedelta(days=7)
    
    created_from = seven_days_ago.strftime('%Y-%m-%dT%H:%M:%SZ')
    created_to = now.strftime('%Y-%m-%dT%H:%M:%SZ')
    
    print(f"Filtering incidents from {created_from} to {created_to}")
    print("(Last 7 days)\n")
    
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
    
    # List incidents with time filter
    list_req = {
        "jsonrpc": "2.0",
        "id": 2,
        "method": "tools/call",
        "params": {
            "name": "list_incidents",
            "arguments": {
                "created_at_gte": created_from,
                "created_at_lte": created_to,
                "page_size": 10
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
    proc.stdin.write(json.dumps(list_req) + '\n')
    proc.stdin.close()
    
    for line in proc.stdout:
        resp = json.loads(line)
        if resp.get('id') == 2:
            if 'error' in resp:
                print(f"ERROR: {resp['error']}")
            elif 'result' in resp:
                data = json.loads(resp['result']['content'][0]['text'])
                incidents = data.get('incidents', [])
                print(f"FOUND {len(incidents)} INCIDENTS IN LAST 7 DAYS:\n")
                
                for inc in incidents:
                    print(f"- {inc.get('name')}")
                    print(f"  Status: {inc.get('incident_status', {}).get('name')}")
                    print(f"  Created: {inc.get('created_at')}")
                    print()

if __name__ == "__main__":
    main()

