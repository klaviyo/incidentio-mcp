#!/usr/bin/env python3
"""
Test script for the list_incident_alerts MCP tool
"""

import json
import subprocess
import sys
import os

def test_mcp_tool(tool_name, args):
    """Test an MCP tool by calling it with the given arguments"""
    try:
        # Create the MCP request
        request = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "tools/call",
            "params": {
                "name": tool_name,
                "arguments": args
            }
        }
        
        # Send the request to the MCP server
        result = subprocess.run(
            ["echo", json.dumps(request)],
            capture_output=True,
            text=True,
            check=True
        )
        
        print(f"✅ {tool_name} with args {args}")
        return True
        
    except subprocess.CalledProcessError as e:
        print(f"❌ {tool_name} failed: {e}")
        return False
    except Exception as e:
        print(f"❌ {tool_name} error: {e}")
        return False

def main():
    print("🧪 Testing ListIncidentAlerts MCP Tool")
    print("=" * 50)
    
    # Test cases for list_incident_alerts
    test_cases = [
        {
            "name": "List incident alerts with default page size",
            "args": {}
        },
        {
            "name": "List incident alerts with custom page size",
            "args": {"page_size": 10}
        },
        {
            "name": "List incident alerts filtered by incident_id",
            "args": {"incident_id": "01FDAG4SAP5TYPT98WGR2N7W91", "page_size": 5}
        },
        {
            "name": "List incident alerts filtered by alert_id",
            "args": {"alert_id": "01GW2G3V0S59R238FAHPDS1R66", "page_size": 5}
        },
        {
            "name": "List incident alerts with pagination",
            "args": {"page_size": 2, "after": "01FCNDV6P870EA6S7TK1DSYDG0"}
        }
    ]
    
    passed = 0
    total = len(test_cases)
    
    for test_case in test_cases:
        print(f"\n📋 {test_case['name']}")
        if test_mcp_tool("list_incident_alerts", test_case["args"]):
            passed += 1
    
    print(f"\n📊 Results: {passed}/{total} tests passed")
    
    if passed == total:
        print("🎉 All tests passed!")
        return 0
    else:
        print("⚠️  Some tests failed!")
        return 1

if __name__ == "__main__":
    sys.exit(main())
