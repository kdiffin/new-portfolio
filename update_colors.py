import os
import glob

replacements = [
    ('bg-[#111111]', 'bg-primary'),
    ('light:bg-white', 'light:bg-primary-light'),
    ('light:bg-[#f9f9f9]', 'light:bg-secondary-light'),
    ('bg-[#1a1a1a]', 'bg-secondary'),
    ('text-[#e0e0e0]', 'text-body'),
    ('light:text-[#111111]', 'light:text-body-light'),
    ('text-[#777]', 'text-muted'),
    ('light:text-[#888]', 'light:text-muted-light'),
    ('text-[#ccc]', 'text-link'),
    ('light:text-[#333]', 'light:text-link-light'),
    ('text-white', 'text-heading'),
    ('light:text-[#000]', 'light:text-heading-light'),
    ('light:text-black', 'light:text-heading-light'),
    ('text-[#aaa]', 'text-nav'),
    ('light:text-[#555]', 'light:text-nav-light'),
    ('border-[#333]', 'border-border'),
    ('light:border-[#ddd]', 'light:border-border-light'),
    ('border-[#444]', 'border-border-subtle'),
    ('hover:border-[#888]', 'hover:border-border-subtle-hover'),
    ('light:border-[#ccc]', 'light:border-border-subtle-light'),
    ('light:border-t-[#eee]', 'light:border-t-border-light'),
]

files = glob.glob('ui/html/**/*.tmpl', recursive=True) + ['ui/static/css/input.css']

for file in files:
    with open(file, 'r') as f:
        content = f.read()
        
    for old, new in replacements:
        content = content.replace(old, new)
        
    with open(file, 'w') as f:
        f.write(content)
        
print("Replacements done!")
