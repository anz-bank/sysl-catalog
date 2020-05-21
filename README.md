# sysl-catalog

A markdown/html + Diagram generator for sysl specifications

## Installation

```bash
go get -u -v github.com/anz-bank/sysl-catalog
```

## How to use
1. Set up environment
`export SYSL_PLANTUML=http://www.plantuml.com/plantuml`

2. Run 

```bashs
sysl-catalog -o <output directory> <input.sysl>
```
- You can optionally specify the `--type=html` if you want to generate html instead of markdown, which is useful for use with github pages, which you can see a demo of with this repo [here](https://anz-bank.github.io/sysl-catalog/)

3. That's it (basically!)

    This will generate markdown with integration diagrams + sequence diagrams + data model diagrams as seen in [demo/markdown/README.md](demo/markdown/README.md) or see html generation at [demo/html/index.html](demo/html/index.html).


## Server Mode
sysl-catalog comes with a `serve` mode which will serve on port `:6900` by default

```bash 
sysl-catalog --serve <input.sysl>
```
This will start a server and filewatchers to watch the input file and its directories recursively, and any changes will automatically show:
![example gif](resources/example.gif)

## Requirements
In [demo/markdown/README.md](demo/markdown/README.md) we have an example with a couple of interesting parts:


1. `@package` attribute must be specified:
- This will create a markdown page for `ApplicationPackage` as seen in [demo/markdown/ApplicationPackage/README.md](demo/markdown/ApplicationPackage/README.md).
 Currently the package name is not inferred from the application name (`MobileApp`), so this needs to be added (`ApplicationPackage`).
```
MobileApp:
    @package = "ApplicationPackage"
    Login(input <: Server.Request):
        Server <- Authenticate
        return ok <: MegaDatabase.Empty
```

2. Application names might need to be prefixed to parameter types if the type is defined in another application, since defined parameters are under scope of the application it is defined in:
```diff
MobileApp:
    @package = "ApplicationPackage"
+    Login(input <: Server.Request):
-    Login(input <: Request):
        Server <- Authenticate
        return ok <: MegaDatabase.Empty
```

3. Add `~ignore` to applications/projects that are to be ignored in the markdown creation
```diff
ThisAppShouldntShow[~ignore]:
    NotMySystem:
        ...
# Or ignore only specific endpoints
ThisAppShouldShow[~ignore]:
    NotMySystem[~ignore]:
        ...
```

## CLI options

#### Output default Markdown
`sysl-catalog -o=docs/ filename.sysl`

#### Output default HTML
`sysl-catalog -o=docs/ --type=html filename.sysl`

#### Run with custom templates
- With this the first template will be executed first, then the second
`sysl-catalog --templates=<fileName.tmpl>,<filename.tmpl> filename.sysl`

#### Run in server mode
`sysl-catalog --serve filename.sysl`
![server mode](resources/server.png)

#### Run in server mode without css/rendered images
- good for rendering raw markdown

`sysl-catalog --serve --noCSS filename.sysl`
![server mode raw](resources/standard-template.png)
#### Run server with custom template
`sysl-catalog --serve --templates=<fileName.tmpl>,<filename.tmpl> filename.sysl`

![server mode raw](resources/custom-template.png)

- See templates/ for custom template examples

## Screenshots
![resources/project_view.png](resources/project_view.png)
*project_view*

![resources/package_view.png](resources/package_view.png)
*package_view*

## Docker image

You can create a docker image containing the necessary depdenencies by running:


```make build```

### Run the Image

You can run the image as follows:

`docker run -v $(pwd)/demo:/demo/ anz-bank/sysl-catalog /demo/simple2.sysl --serve`