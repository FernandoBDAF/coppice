# IDE Configuration Guide

-> IMPORTANT: Never add fictional dates, version numbers, or metrics. Only include real, verified information. If information is not available, mark it as "To be determined" or remove the section.

## Primary Purpose and Main Goals

### Primary Purpose

This guide provides comprehensive IDE setup instructions, recommended extensions, and configuration settings for the Profile Service Microservices project, ensuring a consistent and efficient development environment.

### Main Goals

1. Standardize IDE setup
2. Optimize development workflow
3. Ensure code quality
4. Facilitate debugging
5. Enhance productivity

## Recommended IDEs

### 1. Visual Studio Code

#### Installation

- Download from [VS Code website](https://code.visualstudio.com/)
- Install for your platform
- Launch and verify installation

#### Essential Extensions

1. **Go Development**

   - Go (by Go Team at Google)
   - Go Test Explorer
   - Go Doc
   - Go Outliner

2. **JavaScript/TypeScript**

   - ESLint
   - Prettier
   - JavaScript and TypeScript Nightly
   - npm Intellisense

3. **Docker**

   - Docker
   - Remote - Containers
   - Kubernetes

4. **Testing**

   - Test Explorer UI
   - Coverage Gutters
   - Jest Runner

5. **Git Integration**
   - GitLens
   - Git History
   - Git Graph

#### Recommended Settings

```json
{
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "go.formatTool": "gofmt",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v"],
  "files.autoSave": "afterDelay",
  "files.autoSaveDelay": 1000,
  "editor.rulers": [80, 100],
  "editor.tabSize": 4,
  "files.trimTrailingWhitespace": true
}
```

### 2. GoLand

#### Installation

- Download from [JetBrains website](https://www.jetbrains.com/go/)
- Install for your platform
- Launch and verify installation

#### Essential Plugins

1. **Go Development**

   - Go
   - Go Tools
   - Go Test Explorer

2. **Testing**

   - Coverage
   - Test Runner

3. **Docker**

   - Docker
   - Kubernetes

4. **Version Control**
   - Git Integration
   - GitHub

#### Recommended Settings

1. **Editor Settings**

   - Enable "Format on Save"
   - Set tab size to 4
   - Enable "Show line numbers"
   - Enable "Show whitespaces"

2. **Go Settings**
   - Enable "Format on Save"
   - Set "Go fmt" as formatter
   - Enable "Optimize imports"
   - Set "golangci-lint" as linter

## Debugging Configuration

### 1. VS Code Debugging

#### Go Debugging

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Package",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${fileDirname}",
      "env": {},
      "args": []
    }
  ]
}
```

#### JavaScript Debugging

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "type": "node",
      "request": "launch",
      "name": "Debug Program",
      "skipFiles": ["<node_internals>/**"],
      "program": "${workspaceFolder}/index.js"
    }
  ]
}
```

### 2. GoLand Debugging

#### Debug Configuration

1. **Run Configuration**

   - Set working directory
   - Configure environment variables
   - Set program arguments
   - Configure build tags

2. **Debug Tools**
   - Breakpoints
   - Watches
   - Variables
   - Call Stack

## Code Navigation

### 1. VS Code Navigation

#### Keyboard Shortcuts

- `F12`: Go to Definition
- `Alt+F12`: Peek Definition
- `Shift+F12`: Find References
- `Ctrl+Space`: Trigger Suggestions
- `F2`: Rename Symbol

#### Features

- Code Outline
- Symbol Search
- File Search
- Command Palette

### 2. GoLand Navigation

#### Keyboard Shortcuts

- `Ctrl+B`: Go to Declaration
- `Alt+F7`: Find Usages
- `Ctrl+Alt+Left`: Back
- `Ctrl+Alt+Right`: Forward
- `Ctrl+Shift+F`: Find in Path

#### Features

- Structure View
- TODO Tool Window
- Terminal Integration
- Database Tools

## Performance Optimization

### 1. VS Code Performance

- Disable unnecessary extensions
- Use workspace-specific settings
- Optimize file watching
- Configure search exclusions

### 2. GoLand Performance

- Increase memory allocation
- Disable unnecessary plugins
- Configure indexing
- Optimize file system watching

## Notes

- Keep extensions updated
- Regular IDE updates
- Backup settings
- Share configurations
- Document custom settings

## Version History

### Current Version

- Version: To be determined
- Date: To be determined
- Changes:
  - Initial IDE configuration guide
  - VS Code setup documented
  - GoLand setup documented
  - Debugging configuration added
