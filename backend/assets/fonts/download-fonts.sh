#!/bin/bash

# 下载Noto Sans SC字体文件

fonts=(  
  "k3kCo84MPvpLmixcA63oeAL7Iqp5IZJF9bmaG4HFnYw.ttf"  # 300
  "k3kCo84MPvpLmixcA63oeAL7Iqp5IZJF9bmaG9_FnYw.ttf"  # 400
  "k3kCo84MPvpLmixcA63oeAL7Iqp5IZJF9bmaG-3FnYw.ttf"  # 500
  "k3kCo84MPvpLmixcA63oeAL7Iqp5IZJF9bmaGzjCnYw.ttf"  # 700
  "k3kCo84MPvpLmixcA63oeAL7Iqp5IZJF9bmaG3bCnYw.ttf"  # 900
)

for font in "${fonts[@]}"; do
  echo "下载: $font"
  curl -s -O "https://fonts.gstatic.com/s/notosanssc/v40/$font"
done

echo "字体下载完成！"
