#!/usr/bin/env bash

set -eo pipefail

# Test script for auto cherry-pick functionality
# Tests the cherry-pick.sh script

echo "Testing auto cherry-pick functionality..."

# Store original branch
ORIGINAL_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Test 1: Verify the cherry-pick script exists and is executable
if [ ! -x "./scripts/cherry-pick.sh" ]; then
    echo "FAIL: cherry-pick.sh script not found or not executable"
    exit 1
fi
echo "PASS: cherry-pick.sh script exists and is executable"

# Test 2: Test script usage message
set +e  # Allow command to fail
USAGE_OUTPUT=$(./scripts/cherry-pick.sh 2>&1)
USAGE_EXIT_CODE=$?
set -e  # Resume exit on error
if [[ $USAGE_EXIT_CODE -eq 1 ]] && echo "$USAGE_OUTPUT" | grep -q "Usage:"; then
    echo "PASS: cherry-pick.sh shows usage when called with wrong arguments"
else
    echo "FAIL: cherry-pick.sh should show usage message"
    exit 1
fi

# Test 3: Verify that feature branches exist
if ! git show-ref --verify --quiet "refs/heads/feature-branch-1"; then
    echo "FAIL: feature-branch-1 does not exist"
    exit 1
fi
if ! git show-ref --verify --quiet "refs/heads/feature-branch-2"; then
    echo "FAIL: feature-branch-2 does not exist"
    exit 1
fi
echo "PASS: Both feature branches exist"

# Test 4: Verify that the cherry-pick was successful
git checkout feature-branch-1 >/dev/null 2>&1
if git log --oneline | grep -q "Add cherry-pick test section to README"; then
    echo "PASS: Cherry-picked commit found in feature-branch-1"
else
    echo "FAIL: Cherry-picked commit not found in feature-branch-1"
    exit 1
fi

# Test 5: Verify the commit exists in both branches with different hashes (as expected after cherry-pick)
git checkout feature-branch-2 >/dev/null 2>&1
ORIGINAL_COMMIT=$(git log --oneline | grep "Add cherry-pick test section to README" | cut -d' ' -f1)

git checkout feature-branch-1 >/dev/null 2>&1
CHERRYPICKED_COMMIT=$(git log --oneline | grep "Add cherry-pick test section to README" | cut -d' ' -f1)

if [ "$ORIGINAL_COMMIT" != "$CHERRYPICKED_COMMIT" ]; then
    echo "PASS: Cherry-picked commit has different hash than original (as expected)"
else
    echo "FAIL: Cherry-picked commit should have different hash than original"
    exit 1
fi

# Return to original branch
git checkout "$ORIGINAL_BRANCH" >/dev/null 2>&1

echo ""
echo "All cherry-pick tests passed successfully!"
echo "Original commit: $ORIGINAL_COMMIT (feature-branch-2)"
echo "Cherry-picked commit: $CHERRYPICKED_COMMIT (feature-branch-1)"