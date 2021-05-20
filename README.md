![](imgs/split_list_logo.jpg)

In this Repo you can find an experimental data structure, the SplitList.

In short, SplitList is a sort of **order-statistic** B-Tree/SkipList.

Think of the **many** times you've wished you could access the **something-th** element of a B-Tree faster than /O(n)/.

It supports the following operations:

1. **Find*** ~ O(log n)
2. **Add** ~ O(B)
3. **Delete** ~ O(log n)
4. **Rank** ~ O(log n) // Added
5. **Select** ~ O(log n) // Not yet added, still figuring out, but theoretically it is simple
6. **Range-Query** ~ O(log n) // would returns results unordered though.

It's **very** fast to add due to cache locality, however, the golang code is **very** far from being optimized, thus, even though it has faster lookups and insertions than all other skiplist implementations golang, gearing more towards B-Trees(benchmarks below), due to the fact that it can't be sequentially iterated, it is rather useless at the moment.

Another important point to make is that it is **super** memory efficient.

It could easily support in-order iteration with an aditional pointer cost per data.

* Benchmarks

![](imgs/benchmark.png)
