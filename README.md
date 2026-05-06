# ProfilesCLI

`ProfilesCLI` is a small command-line utility for working with profiles stored as YAML files in the current directory.

The executable is called `mws`. It supports creating, reading, listing and deleting profiles.

Each profile is stored as a separate YAML file. The profile name corresponds to the file name.

For example, profile `test` is stored as:

```text
test.yaml
```

with content:

```yaml
user: example
project: new-project
```

## Requirements

- Go 1.21 or newer

## Installation

### Install with Go

```bash
go install github.com/hamajun504/ProfilesCLI/cmd/mws@latest
```

The binary will be installed into:

```bash
$(go env GOPATH)/bin
```

You can run it directly:

```bash
$(go env GOPATH)/bin/mws help
```

Or add `$(go env GOPATH)/bin` to your `PATH`.

### Build from source

```bash
git clone https://github.com/hamajun504/ProfilesCLI.git
cd ProfilesCLI
go build -o mws ./cmd/mws
```

Run:

```bash
./mws help
```

## Usage

```bash
mws <command> [arguments]
```

Available commands:

```text
mws profile create
mws profile get
mws profile list
mws profile delete
mws profile help
mws help
```

## Commands

### Create profile

Creates a new profile in the current directory.

```bash
mws profile create --name=<name> --user=<user> --project=<project>
```

Example:

```bash
mws profile create --name=test --user=example --project=new-project
```

This creates a file:

```text
test.yaml
```

with content:

```yaml
user: example
project: new-project
```

If a profile with the same name already exists, the utility asks for confirmation before overwriting it.

To overwrite without confirmation, use `--force` or `-f`:

```bash
mws profile create --name=test --user=example --project=new-project --force
```

or:

```bash
mws profile create --name=test --user=example --project=new-project -f
```

### Get profile

Shows profile content by name.

```bash
mws profile get --name=<name>
```

Example:

```bash
mws profile get --name=test
```

Output example:

```text
name: test
user: example
project: new-project
```

### List profiles

Lists profiles in the current directory.

```bash
mws profile list
```

By default, only valid profile files are shown.

Available list modes:

```text
-l    Show only valid profiles
-e    Show valid profiles and profiles with extra YAML fields
-a    Show all YAML profile files, including invalid ones
```

Examples:

```bash
mws profile list
mws profile list -l
mws profile list -e
mws profile list -a
```

### Delete profile

Deletes a profile by name.

```bash
mws profile delete --name=<name>
```

Example:

```bash
mws profile delete --name=test
```

Delete is idempotent: if the profile does not exist, the command still finishes successfully because the desired final state is already reached.

### Help

```bash
mws help
```

or:

```bash
mws profile help
```

## Profile file format

A valid profile YAML file contains exactly two string fields:

```yaml
user: example
project: new-project
```

The file name without `.yaml` is used as the profile name.

Example:

```text
test.yaml
```

corresponds to profile name:

```text
test
```

## Validation rules

Profile name:

- must not be empty;
- must be a safe file name;
- must not contain path separators;
- is used as the YAML file name without extension.

User and project:

- must not be empty;
- must be valid string values;
- must not contain line breaks.

## Examples

Create a profile:

```bash
mws profile create --name=dev --user=alice --project=cloud-platform
```

Show the profile:

```bash
mws profile get --name=dev
```

List profiles:

```bash
mws profile list
```

Overwrite an existing profile:

```bash
mws profile create --name=dev --user=bob --project=new-project --force
```

Delete a profile:

```bash
mws profile delete --name=dev
```

## Notes

All profile operations are performed in the current working directory.

For example:

```bash
cd /tmp/profiles
mws profile create --name=test --user=example --project=new-project
```

creates:

```text
/tmp/profiles/test.yaml
```

not a file near the `mws` binary.