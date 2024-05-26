# docker-list-context -- list files taken by Docker in a build context


`docker-list-context` lists files which are sent by the `docker build` command
to the Docker daemon as the *build context*.

You can reduce the build context (better performance, avoid leaking secrets) by ignoring
files in a `.dockerignore` file.
The official documentation for `.dockerignore` files is here: https://docs.docker.com/build/building/context/#dockerignore-files

Tip: you can even ignore `Dockerfile` and `.dockerignore`, unless those files are referenced in the `Dockerfile`.

## Usage

List context files for `docker build .` (using `Dockerfile.dockerignore` or `.dockerignore`):

```console
$ docker-list-context
```

List context files for `docker build -f Prj1 .` (using `Prj1.dockerignore` or `.dockerignore`):

```console
$ docker-list-context -f Prj1 .
```

## Install

```console
$ go install github.com/dolmen-go/docker-list-context@latest
```

## License

Copyright 2021-2024 Olivier Mengué

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
