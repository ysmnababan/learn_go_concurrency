KEY TAKEWAYS:

- Now that we know how to ensure goroutines don’t leak, we can stipulate a conven‐
tion: If a goroutine is responsible for creating a goroutine, it is also responsible for ensur‐
ing it can stop the goroutine.
