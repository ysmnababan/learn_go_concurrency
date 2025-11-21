KEY TAKEWAYS:

- Now that we know how to ensure goroutines don’t leak, we can stipulate a conven‐
tion: If a goroutine is responsible for creating a goroutine, it is also responsible for ensur‐
ing it can stop the goroutine.
- Errors should be considered first-class citizens when constructing values to return
from goroutines. If your goroutine can produce errors, those errors should be tightly
coupled with your result type, and passed along through the same lines of communi‐
cation—just like regular synchronous functions.