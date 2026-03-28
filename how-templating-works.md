This Go codebase uses the standard library's `html/template` package to generate HTML. 

Here is a step-by-step breakdown of how templating works in your application, drawing from specific examples in your code.

### 1. The Template Files (The View)
Your HTML is broken up into three main parts:
*   **The Base Layout** (base.tmpl): This file contains the outer `<html>`, `<head>`, and `<body>` tags. It acts as the skeleton.
*   **Partials** (e.g., nav.tmpl): Small, reusable snippets like navigation bars or footers.
*   **Pages** (e.g., home.tmpl): The actual unique content for a specific URL route.

If you look at base.tmpl, you'll see things like `{{template "main" .}}` and `{{template "title" .}}`. These act as "hooks" where specific code from a Page template gets injected.

If you look at home.tmpl, it defines those hooks:
```html
{{define "title"}}Home{{end}}

{{define "main"}}
<section class="home-section">
...
{{end}}
```
When Go renders home.tmpl, it knows to inject the content from `{{define "main"}}` directly into the `{{template "main" .}}` slot in base.tmpl.

### 2. The Data Structure (The Model)
When a user visits a page, your Go code needs to pass data into the template. The templates expect a specific Go struct.

If you look at handlers.go, you can see how data is structured for the home page:
```go
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := &models.ViewData{
		Title:       "Home",
		Description: "A minimal publishing site built with Go and Tailwind CSS.",
		Writing:     app.store.List("writing"),
		Micro:       app.store.List("micro"),
		// ...
	}
	app.render(w, r, "home.tmpl", data)
}
```
This `data` struct is what the `.` (dot) represents in your HTML templates. So when home.tmpl says `{{range limit .Writing 9}}`, it is looping over the `Writing` array provided by that `models.ViewData` struct in your Go handler.

### 3. Template Functions (Custom Logic)
Sometimes you need to modify data *while* rendering HTML (like formatting a date or truncating text). Go allows you to register custom "Template Functions".

In home.tmpl, you see a custom function being used:
```html
<time class="home-date" datetime="{{fmtDate .PublishedAt}}">{{fmtDate .PublishedAt}}</time>
```
Go's default templating doesn't know what `fmtDate` is. This means somewhere in your code (likely in main.go or `cmd/web/templates.go`), a `template.FuncMap` is defined that maps the string `"fmtDate"` to a real Go function that formats dates. 

*(Note: The error you are seeing in your terminal right now—`function "safeHTML" not defined`—is because another template is trying to use a function called `safeHTML`, but it hasn't been registered in your Go code yet!)*

### 4. Template Caching (Performance)
Reading files from the hard drive every time a user visits your site is slow. Instead, your app reads all the template files once when the server starts and stores them in memory.

Look at `NewTemplateCache` in templates.go:
```go
func NewTemplateCache(root string, funcs template.FuncMap) (map[string]*template.Template, error) {
    // ...
	for _, page := range pages {
		name := filepath.Base(page) // e.g., "home.tmpl"

        // 1. It parses the base template
		ts, err := template.New(name).Funcs(funcs).ParseFiles(filepath.Join(root, "base.tmpl"))
		
        // 2. It parses and adds all the partials (nav, footer)
		ts, err = ts.ParseGlob(filepath.Join(root, "partials", "*.tmpl"))
		
        // 3. It parses the specific page (e.g. home.tmpl)
		ts, err = ts.ParseFiles(page)

        // 4. It saves the fully compiled group in the cache map
		cache[name] = ts
	}
    // ...
}
```
This creates a map where the key is `"home.tmpl"` and the value is a fully parsed package containing base.tmpl + `nav.tmpl` + `footer.tmpl` + home.tmpl. 

When your handler calls `app.render(w, r, "home.tmpl", data)`, it simply grabs `"home.tmpl"` from this memory cache, injects the `data` struct, and sends the final HTML to the user.