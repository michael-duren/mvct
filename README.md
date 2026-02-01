# MVC Bubbletea (mvct)

> [!WARNING]
> This project is a work in progress and may not be stable for production use.

A framework to simplify building complex TUI applications in Go using the Bubbletea library,
based on the Model-View-Controller (MVC) design pattern.

## Table of Contents

- [About](#about)
- [Features](#features)
- [Installation](#installation)

## About

MVCT is a framework designed to help developers build complex TUI applications in Go using the Bubbletea library.
Bubble tea is a powerful TUI framework, but as applications grow in complexity, managing state and interactions can become challenging.
If you're simply building small TUI apps, you might not need this framework. This is more suited for larger applications with multiple views, complex state management, and user interactions.

## Features

- **MVC Architecture**: Separates concerns into Models, Views, and Controllers for better organization and maintainability.
- **State Management**: Provides a structured way to manage application state across different views and components.
- **Routing**: Built-in support for navigating between different views in the TUI.
- **Event Handling**: Simplifies handling user input and events in a structured manner.
- **Middleware Support**: Allows for adding middleware to handle cross-cutting concerns like logging, authentication, etc.
- **Extensible**: Easily extendable to add custom functionality as needed.
- **Integration with Bubbletea**: Leverages the power of the Bubbletea library while providing additional structure and organization.
- **Default Layouts and Components**: Provides pre-built layouts and components examples to speed up development.

## Installation

You can install the `mvct` package using go or download from the releases page.

```bash
go install github.com/michael-duren/bubbletea-mvc@latest
```

After installation, scaffold a new project:

```bash
mvct new my-tui-app
```

Check out the [documentation](documentation/framework.md) for more information.