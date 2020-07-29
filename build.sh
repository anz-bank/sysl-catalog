#!/bin/bash
git add demo/html || true
git rm -rf demo/html/* || true
git add demo/markdown || true
git rm -rf demo/markdown/* || true
git add docs || true
git rm -rf docs/* || true
make install || true

sysl-catalog -o demo/markdown demo/sizzle.sysl --mermaid
sysl-catalog --type=html --plantuml=https://plantuml.com/plantuml -o demo/html demo/sizzle.sysl --redoc --mermaid
mkdir -p docs
cp -r demo/html/* docs
mkdir demo/html/demo/
mkdir docs/demo/
cp demo/mastercard.yaml demo/html/demo/mastercard.yaml
cp demo/mastercard.yaml docs/demo/mastercard.yaml

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
 sed -i "s/sizzle.sysl/<a href=http:\/\/github.com\/anz-bank\/sysl-catalog>This is an example of sysl catalog deployed to github pages <\/a>/" docs/index.html
elif [[ "$OSTYPE" == "darwin"* ]]; then
 sed -i "" "s/sizzle.sysl/<a href=http:\/\/github.com\/anz-bank\/sysl-catalog>This is an example of sysl catalog deployed to github pages <\/a>/" docs/index.html
fi
