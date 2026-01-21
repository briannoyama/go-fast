# go-fast

Fast Data Structures for go inspired by work on a game engine. Most of the data structures build off of a slice based data structure `CPage` (Continuous Page) which maintains references(`int`s) to a contiguous block of memory.

There are also fast implementations for Treemaps that use `SPage` (Static Page).

## Organization

```mermaid
flowchart TD
  Queue
  A[CPage]-->|creates|C[RefFactory]
  B[CVisitor]-->|has|A
  C-->|creates|D[Ref]
  D-->|creates|E[RefCached]
  F[Heap]-->|has|A
  F-->|creates|C
  G[FTreeMap]-->|has|K[SPage]
  H[Cache]-->|has|F
  I[CBufPage]-->|has|A
  I-->|creates|C
  J[CBuffer]-->|has*|I
  J-->|creates|C
```

See unit tests for usage.
