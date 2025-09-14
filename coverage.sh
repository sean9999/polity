#!/bin/bash

# Code coverage script for polity project
# Usage: ./coverage.sh [summary|html|check|full]

set -e

MODE=${1:-full}
COVERAGE_FILE="coverage.out"
HTML_FILE="coverage.html"
THRESHOLD=70

generate_coverage() {
    echo "🧪 Running tests with coverage..."
    go test -coverprofile=$COVERAGE_FILE ./...
}

show_summary() {
    echo "📊 Coverage Summary:"
    go tool cover -func=$COVERAGE_FILE
}

generate_html() {
    echo "🌐 Generating HTML coverage report..."
    go tool cover -html=$COVERAGE_FILE -o $HTML_FILE
    echo "HTML report saved to: $HTML_FILE"
}

check_threshold() {
    COVERAGE=$(go tool cover -func=$COVERAGE_FILE | grep total | awk '{print $3}' | sed 's/%//')
    echo "📈 Total coverage: $COVERAGE%"

    if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
        echo "❌ Coverage below $THRESHOLD% threshold"
        exit 1
    else
        echo "✅ Coverage meets $THRESHOLD% threshold"
    fi
}

show_uncovered() {
    echo "🔍 Files with uncovered code paths:"
    go tool cover -func=$COVERAGE_FILE | grep -v "100.0%" | head -20
}

case $MODE in
    summary)
        generate_coverage
        show_summary
        ;;
    html)
        generate_coverage
        generate_html
        ;;
    check)
        generate_coverage
        check_threshold
        ;;
    full)
        generate_coverage
        show_summary
        echo ""
        check_threshold
        echo ""
        show_uncovered
        echo ""
        generate_html
        echo ""
        echo "📋 Available commands:"
        echo "  ./coverage.sh summary  - Show coverage summary only"
        echo "  ./coverage.sh html     - Generate HTML report only"
        echo "  ./coverage.sh check    - Check coverage threshold"
        echo "  ./coverage.sh full     - Full coverage analysis (default)"
        ;;
    *)
        echo "Usage: $0 [summary|html|check|full]"
        exit 1
        ;;
esac