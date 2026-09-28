#!/usr/bin/env bash

set -euo pipefail

echo "Scanning for *.ap.yaml files..."

find . -type f -name "*.ap.yaml" | while read -r spec_file; do
    dir_path=$(dirname "$spec_file")
    file_name=$(basename "$spec_file" .ap.yaml)
    output_file="${dir_path}/${file_name}.gen.go"

    # Extract the name of the folder containing the file
    package_name=$(basename "$(cd "$dir_path" && pwd)")

    # Fallback to main if it evaluates to root/current directory dots
    if [ "$package_name" = "." ] || [ "$package_name" = "/" ]; then
        package_name="main"
    fi

    echo "Generating: ${output_file} (Package: ${package_name})"
    
    # Run the generator overriding both package and output targets
    oapi-codegen --config=oapi-codegen.yaml -p "$package_name" -o "$output_file" "$spec_file"
done

echo "Done!"
