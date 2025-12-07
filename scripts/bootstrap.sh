#!/bin/bash
set -e

echo "🚀 Bootstrapping Crowd Unlocked..."

# Check prerequisites
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed."; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js is required but not installed."; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo "❌ Terraform is required but not installed."; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "❌ kubectl is required but not installed."; exit 1; }
command -v flux >/dev/null 2>&1 || { echo "❌ Flux CLI is required but not installed."; exit 1; }

echo "✅ All prerequisites installed"

# Initialize Go modules
echo "📦 Initializing Go modules..."
for service in services/*/; do
    if [ -f "$service/go.mod" ]; then
        echo "  - $(basename $service)"
        (cd "$service" && go mod download)
    fi
done

# Install web dependencies
echo "📦 Installing web dependencies..."
(cd apps/web && npm install)

# Install Flutter dependencies
if command -v flutter >/dev/null 2>&1; then
    echo "📦 Installing Flutter dependencies..."
    (cd apps/mobile && flutter pub get)
else
    echo "⚠️  Flutter not installed, skipping mobile setup"
fi

echo "✅ Bootstrap complete!"
echo ""
echo "Next steps:"
echo "  1. Configure AWS credentials"
echo "  2. Run 'make test' to verify setup"
echo "  3. See docs/SETUP.md for infrastructure deployment"
