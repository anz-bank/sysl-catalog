git add demo/html || true
git rm -rf demo/html/* || true
git add demo/markdown || true
git rm -rf demo/markdown/* || true
git add docs || true
git rm -rf docs/* || true
go run . -o demo/markdown demo/simple2.sysl
go run . --type=html --embed -o demo/html demo/simple2.sysl
mkdir -p docs
cp -r demo/html/* docs/

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  sed -i "s/simple2.sysl/<a href=http:\/\/github.com\/anz-bank\/sysl-catalog>This is an example of sysl catalog deployed to github pages <\/a>/" docs/index.html
elif [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i "" "s/simple2.sysl/<a href=http:\/\/github.com\/anz-bank\/sysl-catalog>This is an example of sysl catalog deployed to github pages <\/a>/" docs/index.html
fi