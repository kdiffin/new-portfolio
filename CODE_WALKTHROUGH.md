# Codebase Walkthrough: Go Portfolio Website

Welcome to the codebase! This document provides a detailed explanation of how this Go web application is structured and how the different pieces fit together. This project follows standard Go web development idioms, often seen in resources like "Let's Go" by Alex Edwards.

## Architecture Overview

This project uses a standard multi-tier architecture, separating the application layer (routing, handling HTTP requests) from the business logic and data access layers.

Here is a breakdown of the top-level directories:

### `cmd/web/` - The Application Entry Point
This directory contains the code that runs the web server. It acts as the "wiring" for the application.

*   `main.go`: The starting point of the application. It typically initializes configurations, connects to databases or data stores, sets up the template cache, and starts the HTTP server.
*   `routes.go`: Defines the URL paths (e.g., `/`, `/about`, `/projects`) and maps them to specific handler functions.
*   `handlers.go`: Contains the controller logic. When a user visits a route, the corresponding handler function is executed. It fetches data from the `internal/` packages, processes it, and passes it to the `render` package to generate an HTML response.
*   `middleware.go`: Contains functions that run before or after your handlers. Common uses include logging requests, handling panics, or adding security headers.

### `internal/` - The Core Business Logic
The `internal/` directory is special in Go. Code inside this directory can only be imported by code within the same parent directory (this repository). It hides your core logic from being imported as a library by others.

*   `models/`: Contains the data structures (structs) that represent the core entities of the application (e.g., a `Content`, `Post`, or `Project` struct).
*   `content/`: Likely contains the logic for retrieving and storing content. This could be reading from a database, or reading Markdown files from the disk (given it's a portfolio site). The `store.go` file probably defines an interface or struct for accessing this data.
*   `render/`: Contains logic for parsing and rendering HTML templates. `templates.go` handles building a "template cache" (loading all templates into memory when the app starts) and executing them with data passed from the handlers.

### `ui/` - The User Interface
This directory holds everything related to what the user sees in the browser.

*   `html/`: Contains Go's `html/template` files.
    *   `base.tmpl`: The master layout template (includes the `<html>`, `<head>`, and `<body>` boilerplate).
    *   `pages/`: Individual page templates (e.g., `home.tmpl`, `about.tmpl`). These define the specific content for a route and "plug into" the base template.
    *   `partials/`: Reusable HTML snippets (e.g., `nav.tmpl`, `footer.tmpl`) that can be included in various pages.
*   `static/`: Holds static assets like CSS, JavaScript, and images that are served directly to the browser.
    *   Notice `input.css` and `output.css`. Coupled with `package.json` and `tailwind.config.js` in the root, this indicates the project uses **Tailwind CSS**. `input.css` has the Tailwind directives, and a build process compiles it into `output.css` based on the classes used in the HTML templates.

### `macwrights-website/`
This directory appears to contain static HTML/CSS files, possibly from an older version of the site, a reference design, or parsed content that is being migrated.

## The Lifecycle of an HTTP Request

To understand how everything connects, let's trace a standard HTTP request (e.g., a user visiting the home page):

1.  **Incoming Request**: The user's browser sends a GET request to `/`.
2.  **Routing (`routes.go`)**: The router intercepts the request and matches `/` to the `Home` handler function.
3.  **Middleware (`middleware.go`)**: The request may pass through middleware (e.g., logging the request URL and IP).
4.  **Handling (`handlers.go`)**: The `Home` handler takes over.
5.  **Data Fetching (`internal/content/store.go`)**: The handler might call a function like `store.GetRecentPosts()` to fetch data for the home page.
6.  **Rendering (`internal/render/templates.go`)**: The handler takes the retrieved data and passes it to the rendering engine, specifying the `home.tmpl` page.
7.  **Template Execution (`ui/html/pages/home.tmpl`)**: The rendering engine injects the data into the HTML template, producing a final HTML string.
8.  **Response**: The HTML is sent back to the user's browser with a `200 OK` status.

## Current Known State / Bugs
When running `go run ./cmd/web`, you might encounter an error like:
`could not build template cache: current error: function "safeHTML" not defined`

This happens because the `render` package is trying to parse a template (`micro_index.tmpl`) that uses a custom template function (`safeHTML`), but that function hasn't been mapped in the `template.FuncMap` before parsing. Fixing this is a great first step to dive into `internal/render/templates.go`!
