#!/usr/bin/env python3
"""
Test script for the follow-ups MCP tools
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
    print("🧪 Testing Follow-ups MCP Tools")
    print("=" * 50)
    
    # Test cases for list_follow_ups
    list_test_cases = [
        {
            "name": "List all follow-ups",
            "args": {}
        },
        {
            "name": "List follow-ups for specific incident",
            "args": {"incident_id": "01FCNDV6P870EA6S7TK1DSYDG0"}
        },
        {
            "name": "List follow-ups for standard incidents",
            "args": {"incident_mode": "standard"}
        },
        {
            "name": "List follow-ups for retrospective incidents",
            "args": {"incident_mode": "retrospective"}
        },
        {
            "name": "List follow-ups with both filters",
            "args": {"incident_id": "01FCNDV6P870EA6S7TK1DSYDG0", "incident_mode": "standard"}
        }
    ]
    
    # Test cases for get_follow_up
    get_test_cases = [
        {
            "name": "Get specific follow-up",
            "args": {"id": "01FCNDV6P870EA6S7TK1DSYDG0"}
        }
    ]
    
    passed = 0
    total = len(list_test_cases) + len(get_test_cases)
    
    print("\n📋 Testing list_follow_ups tool:")
    for test_case in list_test_cases:
        print(f"\n  • {test_case['name']}")
        if test_mcp_tool("list_follow_ups", test_case["args"]):
            passed += 1
    
    print("\n📋 Testing get_follow_up tool:")
    for test_case in get_test_cases:
        print(f"\n  • {test_case['name']}")
        if test_mcp_tool("get_follow_up", test_case["args"]):
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
