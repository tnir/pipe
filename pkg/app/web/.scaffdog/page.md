---
name: "page"
questions:
  name: "Please enter page name"
root: "src/pages"
output: "**/*"
ignore: []
---

# `{{ inputs.name }}.tsx`

```tsx
import { memo, FC, useEffect } from "react";
import { useParams } from "react-router-dom";

export const {{ inputs.name | pascal }}Page: FC = memo(() => {
  return <div>hello</div>;
});
```
