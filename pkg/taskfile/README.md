# Taskfile

## Directory Structure
Use `Taskfile.yml` with the following structure.
```txt
your_tool_directory/
|--1-1_image_renamer_with_exif/
|--
|--  ...
|--
|--999_converted_images/
|--Taskfile.yml
|--pkg
| |--bin
| | |--cli
| | | |--darwin_arm64
| | | | |--anilist
| | | | |--arithmetic-calculator
| | | | |--
| | | | |--  ...
| | | | |--
| | | |--linux_amd64
| | | | |--anilist
| | | | |--arithmetic-calculator
| | | | |--
| | | | |--  ...
| | | | |--
| | | |--win_amd64
| | | | |--anilist.exe
| | | | |--arithmetic-calculator.exe
| | | | |--
| | | | |--  ...
| | | | |--
| | |--db-server (external dependency)
| | | |--cli
| | | | |--darwin_arm64
| | | | | |--all_cli_entry
| | | | |--linux_amd64
| | | | | |--all_cli_entry
| | | | |--win_amd64
| | | | | |--all_cli_entry.exe
| |--taskfiles
| | |--calculate_for_text.yml
| | |--core.yml
| | |--env.sample.yml
| | |--env.yml
| | |--file_maneuver.yml
| | |--image_convert.yml
| | |--image_rename.yml
| | |--image_synthesize.yml
| | |--iso8601.yml
```

## Installation
### Root taskfile settings
Input `<your-taskfiles-dir>` to use in your working directory.
```yaml
includes:
  calc:
    taskfile: ./<your-taskfiles-dir>/calculate_for_text.yml
    flatten: false
    internal: true

...

```

### Set env vars
Especially, confirm `PACKAGE_ROOT_DIR`.

### Symblic links
Create new symbolic links.
```bash
cd "<your-dir>"
ln -sfn "$HOME/devbox/pkg/taskfile/taskfiles" "taskfiles"
ln -sfn "$HOME/devbox/pkg/bin" "bin"
```
