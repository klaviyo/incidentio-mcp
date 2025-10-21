#!/bin/bash
# Test what filter parameters the incident.io API actually supports

if [ -z "$INCIDENT_IO_API_KEY" ]; then
    echo "Error: INCIDENT_IO_API_KEY not set"
    exit 1
fi

API_BASE="https://api.incident.io/v2"
HEADERS="Authorization: Bearer $INCIDENT_IO_API_KEY"

echo "Testing incident.io API filter parameters..."
echo "============================================"

echo -e "\n1. Basic query (baseline):"
curl -s -H "$HEADERS" "$API_BASE/incidents?page_size=1" | jq -r 'if .type then "ERROR: \(.type) - \(.errors[0].message)" else "SUCCESS - \(.pagination_meta.total_count // 0) total incidents" end'

echo -e "\n2. With status filter:"
curl -s -H "$HEADERS" "$API_BASE/incidents?page_size=1&status=active" | jq -r 'if .type then "ERROR: \(.type) - \(.errors[0].message)" else "SUCCESS - status filter works" end'

echo -e "\n3. With created_at[from] filter:"
curl -s -H "$HEADERS" "$API_BASE/incidents?page_size=1&created_at%5Bfrom%5D=2024-10-01T00:00:00Z" | jq -r 'if .type then "ERROR: \(.type) - \(.errors[0].message)" else "SUCCESS - date filter works" end'

echo -e "\n4. Testing invalid parameter to see error format:"
curl -s -H "$HEADERS" "$API_BASE/incidents?page_size=1&invalid_param=test" | jq -r 'if .type then "ERROR: \(.type) - \(.errors[0].message)" else "SUCCESS or ignored" end'

echo -e "\n5. Full response with created_at filter (to inspect):"
curl -s -H "$HEADERS" "$API_BASE/incidents?page_size=1&created_at%5Bfrom%5D=2024-10-01T00:00:00Z"

