# Design System Specification: The Synthetic Intelligence Interface

## 1. Overview & Creative North Star
**Creative North Star: The Digital Curator**
This design system moves away from the "noisy" futurism of Hollywood interfaces toward a sophisticated, editorial high-tech aesthetic. It treats AI management not as a chaotic stream of data, but as a curated, high-fidelity command center. 

The system breaks the "standard dashboard" template through **Intentional Asymmetry**—where dense data modules are balanced by expansive, breathable negative space. We utilize **Tonal Depth** instead of structural lines to create a UI that feels like it’s carved from a single block of dark slate, punctuated by "light-emitters" (our accent tokens) that signal intelligence and activity.

---

## 2. Colors & Surface Philosophy

### The "No-Line" Rule
To achieve a premium, seamless feel, **1px solid borders are prohibited for sectioning.** Boundaries must be defined through background color shifts. A `surface-container-low` section sitting on a `surface` background provides enough contrast for the eye without the "grid-prison" effect of traditional outlines.

### Surface Hierarchy & Nesting
Think of the UI as physical layers of smoked glass and obsidian. 
- **Base Layer:** `surface` (#0b1326) – The infinite void.
- **Sectioning:** `surface-container-low` (#131b2e) – Defines large functional areas.
- **Interaction/Cards:** `surface-container` (#171f33) – The primary staging area for agent data.
- **Floating/Active:** `surface-container-highest` (#2d3449) – Used for modals or high-priority overlays.

### The Glass & Gradient Rule
For streaming text components and tool execution logs, use **Glassmorphism**. Apply `surface-variant` at 40% opacity with a `backdrop-blur` of 12px. Main CTAs should utilize a subtle linear gradient from `primary` (#c0c1ff) to `primary-container` (#8083ff) at a 135-degree angle to provide a "lithium-ion" glow.

---

## 3. Typography: The Intellectual Voice

The system pairs the technical precision of **Inter** with the architectural character of **Space Grotesk**.

- **Display & Headlines (Space Grotesk):** These are the "Editorial" anchors. Use `display-lg` and `headline-md` to establish clear topical hierarchy. The wider tracking and geometric forms of Space Grotesk signal a futuristic, authoritative tone.
- **Body & Labels (Inter):** Inter is used for high-density data. Its high x-height ensures readability in complex logs. 
- **Hierarchy through Scale:** We use an aggressive contrast between `headline-lg` (2rem) and `label-sm` (0.6875rem). This "Large-Small" pairing makes the interface feel like a professional technical manual rather than a generic web app.

---

## 4. Elevation & Depth: Tonal Layering

### The Layering Principle
Depth is achieved by "stacking" tones. 
*Example:* A `surface-container-lowest` (#060e20) card placed inside a `surface-container-low` (#131b2e) section creates a recessed, "etched" look. This is preferred over shadows for a high-tech, integrated feel.

### Ambient Shadows
Shadows are only permitted for "floating" elements (Modals, Tooltips). 
- **Specs:** Blur: 24px–40px, Opacity: 6%, Color: `on-surface` (#dae2fd). This mimics natural ambient occlusion rather than a harsh drop shadow.

### The "Ghost Border" Fallback
If a border is required for accessibility, use a **Ghost Border**: `outline-variant` (#464554) at 20% opacity. Never use 100% opaque borders.

---

## 5. Components

### Primary Action (Button)
- **Style:** Gradient fill (`primary` to `primary-container`). 
- **Radius:** `md` (0.375rem).
- **Interaction:** On hover, increase the `surface-tint` overlay by 10%. No "lifting" shadows; use glow intensity instead.

### The Agent Status Chip
- **Active:** `tertiary-container` background with `tertiary` (#4edea3) text. Use a 4px "pulse" dot.
- **Error:** `error-container` background with `error` (#ffb4ab) text.
- **Idle:** `surface-variant` background with `on-surface-variant` text.

### Streaming Log (Custom Component)
- **Background:** `surface-container-lowest` with a 20% `outline-variant` Ghost Border.
- **Typography:** `body-sm` (Inter) with a 1.6 line-height for readability during high-speed text streaming.
- **Visual:** A subtle `secondary` (#4cd7f6) left-hand vertical accent line to indicate "Execution Mode."

### Cards & Data Lists
- **Rule:** Forbid divider lines. Separate list items using `2.5` (0.5rem) of vertical space and a subtle `surface-container-high` background shift on hover.

### Tool Execution Logs
- **Modular Blocks:** Use nested `surface-container-low` blocks for each tool step.
- **Success Indicators:** Use `tertiary` (#4edea3) sparingly for icons only; do not flood the UI with green.

---

## 6. Do's and Don'ts

### Do
- **Do** use `20` (4.5rem) or `24` (5.5rem) spacing to separate major functional groups to prevent "Data Fatigue."
- **Do** use `secondary` (#4cd7f6) for data visualization and "active" state highlights.
- **Do** lean into `surface-container-lowest` for background areas where you want the user to focus on content.

### Don't
- **Don't** use pure black (#000000). Use the `surface` token to maintain tonal depth.
- **Don't** use standard "Select" dropdowns. Use the `surface-container-highest` overlay pattern with glassmorphism.
- **Don't** crowd the interface. If a module has high data density, it must be surrounded by at least `10` (2.25rem) of white space.
- **Don't** use harsh transitions. All state changes should have a 150ms ease-in-out transition for a premium, liquid feel.