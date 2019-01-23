## Building from source

This project uses [Gradle](https://www.gradle.org) as the build system. To run it, the only thing
you'll need is to have Java installed. If you have a Go workspace [correctly setup](https://golang.org/doc/code.html)
you can run the following to run a complete build.

```bash
cd $GOPATH
mkdir -p src/github.com/visola
cd src/github.com/visola
git clone https://github.com/visola/go-http-cli.git
cd go-http-cli
./gradlew updateDependencies updateLinter build
```

There are also tasks to build specific packages for the `amd64` arch. You can build all the packages
using:

```bash
./gradlew buildPackages
```

or for one specific platform:

```bash
./gradlew packageDarwinAMD64
```
