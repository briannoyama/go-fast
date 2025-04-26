# go-fast

Fast Data Structures for go inspired by work on a game engine. Everything here builds off of a slice based data structure `CPage` which maintains references(`int`s) to a contiguous block of memory.

## Organization

```mermaid
lowchart TD
  A[CPage]-->|creates|C[RefFactory]
  B[CVisitor]-->|has|A
  C-->|creates|D[Ref]
  D-->|creates|E[RefCached]
  F[Heap]-->|has|A
  F-->|creates|C
  G[FTreeMap]-->|has|A
  H[Cache]-->|has|F
  I[CBufPage]-->|has|A
  I-->|creates|C
  J[CBuffer]-->|has*|I
  J-->|creates|C
```

See unit tests for usage.
