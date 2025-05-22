#!/bin/bash

# Check if input file exists
if [ ! -f "$1" ]; then
  echo "Error: $1 not found"
  exit 1
fi

# Compile Java file
javac "$1" 2> compile_errors.txt
if [ $? -eq 0 ]; then
  # Extract class name from filename
  CLASS=$(basename "$1" .java)
  java "$CLASS"
else
  echo "Compilation Failed:"
  cat compile_errors.txt
fi
