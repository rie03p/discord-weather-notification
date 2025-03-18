#!/bin/bash
set -euo pipefail

# 必要なコマンドの存在確認
for cmd in go terraform zip; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "Error: $cmd がインストールされていません。"
        exit 1
    fi
done

echo "build lambda function..."
pushd lambda > /dev/null
rm -f bootstrap function.zip
GOOS=linux GOARCH=amd64 go build -o bootstrap
zip function.zip bootstrap
popd > /dev/null

echo "deploy lambda function..."
pushd terraform > /dev/null
terraform init

terraform plan -var-file="terraform.tfvars"

read -p "上記の差分を確認しました。適用しますか？ (y/N): " confirm
if [[ "$confirm" =~ ^[Yy]$ ]]; then
    terraform apply -auto-approve -var-file="terraform.tfvars"
    echo "deploy complete."
else
    echo "deploy canceled."
fi
popd > /dev/null
