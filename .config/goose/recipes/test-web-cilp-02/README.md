## Run
Run `goose run` in this directory.
```bash
goose run --recipe ./main-recipe.yaml --no-session --max-turns 20 --max-tool-repetitions 2 --debug
```

## Debug
Run all commands from this directory.
```bash
cd .config/goose/recipes/test-web-cilp-02
```

Validate the main recipe and the subrecipe.
```bash
goose recipe validate ./main-recipe.yaml
goose recipe validate ./subrecipes/patch-markdown.yaml
```

Render the main recipe without executing it.
```bash
goose run --recipe ./main-recipe.yaml --render-recipe
```

Run the subrecipe directly. This verifies that `web_clipper__patch_markdown`
works without going through `delegate`.
```bash
goose run \
  --recipe ./subrecipes/patch-markdown.yaml \
  --params target_title="Test Article" \
  --params target_url="https://example.com/test-article" \
  --params out_file_path="./tmp/test-web-clip-subrecipe-direct.md" \
  --no-session \
  --max-turns 10 \
  --max-tool-repetitions 2 \
  --debug
```

Run the main recipe. This verifies that the main recipe calls the subrecipe via
`delegate`, and that the subagent can call `web_clipper__patch_markdown`.
```bash
goose run \
  --recipe ./main-recipe.yaml \
  --no-session \
  --max-turns 20 \
  --max-tool-repetitions 2 \
  --debug
```

Check the generated Markdown.
```bash
sed -n '1,120p' ./tmp/test-web-clip-subrecipe-patch-markdown.md
```

Check that the recipe YAML files do not contain local absolute paths.
Replace '/home/user/devbox' to your repository root directory from home directory.
```bash
rg -n '/home/user/devbox' ./main-recipe.yaml ./subrecipes/patch-markdown.yaml
```

Check that the MCP executable symlink exists and is executable.
```bash
ls -l ./mcps
test -x ./mcps/devbox-mcp-tools
```
