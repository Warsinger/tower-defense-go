# Shared web footer pattern

Tower Defense #23 and math-visualization #74 use the same landing-page link design.
For future pages, copy the small HTML/CSS pattern below into their own static
entry point. Each repository owns its copy; no shared JavaScript, CDN stylesheet,
fetch or new package is required to show the link. Coordinate changes through
these issue references rather than tying independent deployments together.

## Content and behavior

- Place the footer after the app and instructions, outside the canvas/iframe.
- Use “More apps by Matt Lee.” and “Explore all projects” consistently.
- Link to `https://eelwehttam.pages.dev/`, not a preview or repository URL.
- Open a new tab to preserve the active app's state. Announce that behavior to
  screen readers and use `noopener noreferrer`.
- Keep the link underlined, provide a visible keyboard focus outline, and allow
  text to wrap at narrow widths. The arrow is decorative.
- Keep the footer in normal document flow, never an overlay. `flex: none` keeps
  it visible on full-height flex layouts such as Math Viz, whose canvas should
  continue using `flex: 1; min-height: 0` for the remaining space.

```html
<footer class="project-footer">
  More apps by Matt Lee.
  <a href="https://eelwehttam.pages.dev/" target="_blank" rel="noopener noreferrer" aria-label="Explore all projects (opens in a new tab)">Explore all projects <span aria-hidden="true">↗</span></a>
</footer>
```

The initial dark-page style uses muted text, a subtle divider and a light blue
link. Keep these rules consistent between the two apps; future light themes can
adapt the palette while retaining the content, placement and keyboard behavior.

```css
.project-footer { box-sizing: border-box; flex: none; width: 100%; max-width: 900px; margin: 16px auto 0; padding: 16px; border-top: 1px solid #405163; color: #b9c8d4; font: 14px/1.6 system-ui, sans-serif; text-align: center; }
.project-footer a { color: #9adfff; text-underline-offset: 3px; }
.project-footer a:hover { text-decoration-thickness: 2px; }
.project-footer a:focus-visible { outline: 2px solid currentColor; outline-offset: 4px; border-radius: 2px; }
```

Before publishing, inspect the footer at desktop and narrow widths, reach the
link using the keyboard, and verify it opens the landing page without replacing
the active app. Check both local and packaged output. Native apps and internal
runtime documents do not need this website navigation footer.
