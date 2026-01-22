#!/bin/bash

set -e

echo "========================================"
echo "  Search API Test Suite"
echo "========================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Run tests with coverage
echo -e "${BLUE}Running tests with coverage...${NC}"
echo ""

# Test internal/search package
echo "📦 Testing: internal/search"
go test ./internal/search -mod=mod -cover -v 2>&1 | grep -E "(RUN|PASS|FAIL|SKIP|coverage)" | tail -5

echo ""

# Test API package
echo "📦 Testing: cmd/search-api/internal/api"
go test ./cmd/search-api/internal/api -mod=mod -cover -v -timeout=30s 2>&1 | grep -E "(RUN|PASS|FAIL|SKIP|coverage)" | tail -5

echo ""
echo "========================================"
echo "  Test Summary"
echo "========================================"

# Run all tests and capture output
TEST_OUTPUT=$(go test ./internal/search ./cmd/search-api/internal/api -mod=mod -cover 2>&1)
echo "$TEST_OUTPUT"

# Check if tests passed
if echo "$TEST_OUTPUT" | grep -q "FAIL"; then
    echo ""
    echo -e "${RED}❌ Tests FAILED${NC}"
    exit 1
else
    echo ""
    echo -e "${GREEN}✅ All tests PASSED!${NC}"
fi

echo ""
echo "========================================"
echo "  Quick Commands"
echo "========================================"
echo ""
echo "Run specific test:"
echo "  go test ./internal/search -mod=mod -run TestBuildQueryContext"
echo ""
echo "Run with verbose output:"
echo "  go test ./internal/search -mod=mod -v"
echo ""
echo "Generate HTML coverage report:"
echo "  go test ./... -mod=mod -coverprofile=coverage.out"
echo "  go tool cover -html=coverage.out -o coverage.html"
echo "  open coverage.html"
echo ""
