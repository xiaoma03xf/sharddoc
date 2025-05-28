#!/bin/bash

# 默认参数值
DEFAULT_NAME="SQLParser.g4"

# 设置ANTLR4命令和JAR路径
ANTLR4_JAR="./antlr-4.13.2-complete.jar"
OUTPUT_DIR="./ast"

# 检查JAR文件是否存在
if [ ! -f "$ANTLR4_JAR" ]; then
  echo "错误: 找不到ANTLR4 JAR文件: $ANTLR4_JAR"
  exit 1
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 执行ANTLR4命令
java -Xmx500M -cp "$ANTLR4_JAR:." org.antlr.v4.Tool \
  -Dlanguage=Go \
  -visitor \
  -package "$DEFAULT_PACKAGE" \
  -o "$OUTPUT_DIR" \
  "${DEFAULT_NAME}" SQLLexer.g4 SQLParser.g4

# 检查执行结果
if [ $? -ne 0 ]; then
  echo "错误: ANTLR4 生成失败"
  exit 1
fi

echo "生成成功！Go代码已输出到 $OUTPUT_DIR 目录"