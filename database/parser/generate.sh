#!/bin/bash

# 默认参数值
DEFAULT_NAME="SQLParser.g4"

ANTLR4_JAR="./antlr-4.13.2-complete.jar"
OUTPUT_DIR="./ast"
DEFAULT_PACKAGE="ast"

if [ ! -f "$ANTLR4_JAR" ]; then
  echo "error: cann't find ANTLR4 JAR file: $ANTLR4_JAR"
  exit 1
fi

if [ -d "$OUTPUT_DIR" ]; then
  echo "warn: output dir $OUTPUT_DIR already exist，deleting..."
  rm -rf "$OUTPUT_DIR"
fi

mkdir -p "$OUTPUT_DIR"

java -Xmx500M -cp "$ANTLR4_JAR:." org.antlr.v4.Tool \
  -Dlanguage=Go \
  -visitor \
  -package "$DEFAULT_PACKAGE" \
  -o "$OUTPUT_DIR" \
  "${DEFAULT_NAME}" SQLLexer.g4 SQLParser.g4

if [ $? -ne 0 ]; then
  echo "error: ANTLR4 faild"
  exit 1
fi

echo "successfully generate $OUTPUT_DIR "