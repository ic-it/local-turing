# Local Turing

This is a simple tool to run your tests on your local machine.

![Example of usage](./assets/demo.gif)

<!-- spoiler -->

<details>
  <summary>Open to see the demo</summary>

[![asciicast](https://asciinema.org/a/2Z3VLJZJXx91k32TABBpRsIxA.svg)](https://asciinema.org/a/2Z3VLJZJXx91k32TABBpRsIxA)

</details>

## Installation

```bash
git clone https://github.com/ic-it/local-turing.git
cd local-turing
go install .
```

## Configuration

**Schemas**:
- [tests.json](https://raw.githubusercontent.com/ic-it/local-turing/main/assets/schemas/tests.json)
- [config.json](https://raw.githubusercontent.com/ic-it/local-turing/main/assets/schemas/config.json)

**Test file example**:  
[tests.json gist](https://gist.github.com/ic-it/7c401138b41ffc2b4f3c1105abacdabf)

**Example config file**:
```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/ic-it/local-turing/main/assets/schemas/config.json

# Cloud Turing configuration
cloud-turing:
  name: xchaban
  password: SomePassword
  url: https://www.turing.sk # (optional)

# Local Turing configuration
local-turing:
  tests-file: tests.json
  # (optional) if not specified you must specify the executable for each assignment
  build-commands: 
    - mkdir --parents bin
    - clang -O0 -Wall -Werror -std=gnu11 -g3 -ggdb -fno-omit-frame-pointer main.c -o bin/main
  # (optional) if not specified, you must specify the build-commands for each assignment
  executable: bin/main 
  # (optional) if not specified, you must specify the main-file for each assignment
  main-file: main.c
  assignments:
    - name: teap-uloha-1-1
      dir: 1-1            # (optional) if not specified, the name of the assignment is used
      build-commands:     # (optional) if not specified, the build-commands from the global config are used
        - make
      executable: bin/1-1 # (optional) if not specified, the executable from the global config is used
      main-file: main.c   # (optional) if not specified, the main-file from the global config is used
      push-name: teap-uloha-1-1 # (optional) if not specified, the name of the assignment is used
    - name: teap-uloha-1-2
      dir: 1-2
    - name: teap-uloha-2-1
      dir: 2-1
    - name: teap-uloha-2-2
      dir: 2-2
    - name: teap-uloha-3-1
      dir: 3-1
    - name: teap-uloha-3-2
      dir: 3-2
    - name: teap-uloha-4-1
      dir: 4-1
    - name: teap-uloha-6-1
      dir: 6-1
    - name: teap-uloha-6-2
      dir: 6-2
    - name: teap-uloha-7-1
      dir: 7-1
    - name: teap-uloha-7-2
      dir: 7-2
    - name: teap-uloha-8-1
      dir: 8-1
```

## Usage

**Create or download tests file**:
```bash
curl -o tests.json https://gist.githubusercontent.com/ic-it/7c401138b41ffc2b4f3c1105abacdabf/raw/12eb05df5843543386f3839d90f0825c453f6a30/tests.json
```

**Run local-turing**:
```bash
local-turing test
local-turing -C test | less -S
```

## License

[MIT](LICENSE.txt)