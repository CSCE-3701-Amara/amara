# Figma shader source

`CodeComponentId_8ce92017e53431a2f04b3574f4ba7c98f6f55f1e_8358.ts` is the
unmodified shader source returned by the Figma MCP server for welcome frame
`149:119` in file `5alyud0Z3bCoU3yaV0q5BV`. It is excluded from ESLint so the
export can remain verbatim. Do not format or hand-edit the generated WGSL.

The MCP connection returned the shader source but did not expose its referenced
`file://shader-runtime-index.json/` resource or shader runtime templates.
`shared/ui/shader-fill.tsx` supplies the small React canvas host for this fill's
`setup(device, frame)` and `render(device, frame)` exports, using the exact
parameters from the frame. No shader math is recreated in CSS. WebGPU requires
a supported browser and a secure context (localhost is supported).

Reduced motion freezes shader time and the auth-panel animations. Browsers
without a WebGPU adapter retain the hero's orange base surface. The browser
shader check verifies the GPU workload completes without validation errors.

## Maintenance

Keep generated shader source isolated from application code. If the Figma
export changes, replace the generated file as a whole and update the
component wrapper only when the exported `setup` or `render` contract changes.
Do not add secrets, remote shader URLs, or browser storage to this directory.

Validate the welcome screen from `web/` with:

```sh
npm run build
npm run test:e2e
```
