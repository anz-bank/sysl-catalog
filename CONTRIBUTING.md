## Project Layout
```
.
├── README.md
├── resources                 README images
├── Makefile
├── build.sh
├── .github/workflows         GitHub Actions configurations
├── main.go                   command entry
├── pkg
│   ├── catalog               core logic to generate catalog files and serve service
│   ├── catalogdiagrams       functions to generate diagrams
│   └── watcher               functions to watch for file changes in server mode
├── templates                 pre-defined custom template examples, used in flag --templates
├── demo
│   ├── html                  demo on generating HTML files from sysl files
│   ├── markdown              demo on generating Markdown files from sysl files
│   ├── protos                demo on generating Markdown files from proto files
│   ├── simple.yaml
│   └── simple2.sysl
├── docs                      duplicated with demo/html, as GitHub Pages publishing source
├── scripts                   Docker build PlantUML dependancy, can be removed once PlantUML is removed
├── java                      Docker build PlantUML dependancy, can be removed once PlantUML is removed
└── tests
```

## Generator Template
pkg/catalog/template.go:
```
ProjectTemplate
	# Optional
    MacroPackage
        NewPackageTemplate
    MacroPackage
        NewPackageTemplate
    MacroPackage
```
demo/simple2.sysl:
```
simple2[~project]:
    FirstDivision:
        # you can specify packages to include
        ApplicationPackage
    SecondDivision:
        MegaDatabase
    ThirdDivision:
        ServerPackage
        foo
        simpleredoc
```
