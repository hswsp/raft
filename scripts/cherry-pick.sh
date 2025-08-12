#!/usr/bin/env bash

set -eo pipefail

# Auto cherry-pick script for raft repository
# Usage: ./scripts/cherry-pick.sh <commit_hash> <source_branch> <target_branch>

if [ $# -ne 3 ]; then
    echo "Usage: $0 <commit_hash> <source_branch> <target_branch>"
    echo "Example: $0 57e6240431f035cc2d4fe3da14c1f2496575f30f feature-branch-2 feature-branch-1"
    exit 1
fi

COMMIT_HASH=$1
SOURCE_BRANCH=$2
TARGET_BRANCH=$3

echo "Auto cherry-pick case 2: Cherry-picking commit $COMMIT_HASH from $SOURCE_BRANCH to $TARGET_BRANCH"

# Store current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Verify the commit exists in the source branch
if ! git rev-list "$SOURCE_BRANCH" | grep -q "^$COMMIT_HASH"; then
    echo "Error: Commit $COMMIT_HASH not found in branch $SOURCE_BRANCH"
    exit 1
fi

# Verify target branch exists
if ! git show-ref --verify --quiet "refs/heads/$TARGET_BRANCH"; then
    echo "Error: Target branch $TARGET_BRANCH does not exist"
    exit 1
fi

echo "Switching to target branch: $TARGET_BRANCH"
git checkout "$TARGET_BRANCH"

echo "Cherry-picking commit: $COMMIT_HASH"
if git cherry-pick "$COMMIT_HASH"; then
    echo "Cherry-pick successful!"
    echo "Cherry-picked commit details:"
    git log --oneline -1
else
    echo "Cherry-pick failed. Please resolve conflicts manually."
    git status
    exit 1
fi

# Return to original branch
echo "Returning to original branch: $CURRENT_BRANCH"
git checkout "$CURRENT_BRANCH"

echo "Cherry-pick operation completed successfully."