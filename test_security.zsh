#!/usr/bin/env zsh

echo "=== Security Testing Script for Scrypts ===" echo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Start the server
echo "${YELLOW}Starting server...${NC}"
cd "$(dirname "$0")"
rm -f scrypts.db*
JWT_SECRET="a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
MASTER_KEY="q1w2e3r4t5y6u7i8o9p0a1s2d3f4g5h6" \
ALLOWED_ORIGINS="http://localhost:3000" \
./scrypts > /tmp/scrypts_test.log 2>&1 &
SERVER_PID=$!
sleep 2

echo "${GREEN}Server started with PID: $SERVER_PID${NC}"
echo

# Test 1: Username validation
echo "=== Test 1: Username Validation ==="
echo "Testing valid username (alphanumeric with underscore/hyphen)..."
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test_user-123","password":"TestPass123!"}')
if [[ "$RESPONSE" == *"successfully"* ]]; then
  echo "${GREEN}✓ Valid username accepted${NC}"
else
  echo "${RED}✗ Failed: $RESPONSE${NC}"
fi

echo "Testing invalid username (contains @ symbol)..."
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test@user","password":"TestPass123!"}')
if [[ "$RESPONSE" == *"error"* ]] || [[ "$RESPONSE" == *"failed"* ]]; then
  echo "${GREEN}✓ Invalid username rejected${NC}"
else
  echo "${RED}✗ Failed: Invalid username was accepted${NC}"
fi
echo

# Test 2: Security headers
echo "=== Test 2: Security Headers ==="
HEADERS=$(curl -i -s http://localhost:8080/ | head -20)
echo "Checking for security headers..."
if [[ "$HEADERS" == *"X-Frame-Options"* ]]; then
  echo "${GREEN}✓ X-Frame-Options header present${NC}"
else
  echo "${RED}✗ X-Frame-Options header missing${NC}"
fi
if [[ "$HEADERS" == *"Content-Security-Policy"* ]]; then
  echo "${GREEN}✓ Content-Security-Policy header present${NC}"
else
  echo "${RED}✗ Content-Security-Policy header missing${NC}"
fi
if [[ "$HEADERS" == *"X-Content-Type-Options"* ]]; then
  echo "${GREEN}✓ X-Content-Type-Options header present${NC}"
else
  echo "${RED}✗ X-Content-Type-Options header missing${NC}"
fi
echo

# Test 3: Rate limiting
echo "=== Test 3: Rate Limiting ==="
echo "Making 12 rapid registration attempts (limit is 10)..."
SUCCESS_COUNT=0
BLOCKED_COUNT=0
for i in {1..12}; do
  RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"rateuser$i\",\"password\":\"TestPass123!\"}")
  
  if [[ "$RESPONSE" == *"successfully"* ]]; then
    ((SUCCESS_COUNT++))
  elif [[ "$RESPONSE" == *"Rate limit"* ]]; then
    ((BLOCKED_COUNT++))
  fi
done

echo "Successful requests: $SUCCESS_COUNT"
echo "Blocked requests: $BLOCKED_COUNT"

if [[ $SUCCESS_COUNT -le 10 ]] && [[ $BLOCKED_COUNT -ge 2 ]]; then
  echo "${GREEN}✓ Rate limiting working correctly${NC}"
else
  echo "${RED}✗ Rate limiting failed${NC}"
fi
echo

# Test 4: CORS validation
echo "=== Test 4: CORS Policy ==="
echo "Testing request with allowed origin..."
RESPONSE=$(curl -i -s -X OPTIONS http://localhost:8080/register \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST")
if [[ "$RESPONSE" == *"Access-Control-Allow-Origin"* ]]; then
  echo "${GREEN}✓ CORS allows whitelisted origin${NC}"
else
  echo "${YELLOW}⚠ CORS response: Check manually${NC}"
fi
echo

# Test 5: Password complexity
echo "=== Test 5: Password Complexity ==="
echo "Testing weak password..."
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"weakpwduser","password":"weak"}')
if [[ "$RESPONSE" == *"error"* ]] || [[ "$RESPONSE" == *"failed"* ]] || [[ "$RESPONSE" == *"complex"* ]]; then
  echo "${GREEN}✓ Weak password rejected${NC}"
else
  echo "${RED}✗ Weak password was accepted${NC}"
fi
echo

# Cleanup
echo "${YELLOW}Cleaning up...${NC}"
kill $SERVER_PID 2>/dev/null
rm -f scrypts.db*

echo
echo "${GREEN}=== Security testing complete ===${NC}"
echo "Server logs available at: /tmp/scrypts_test.log"
