# Static asset performance baseline

Baseline commit: `e49867fe67d9f0a3c965a1d70ce06b104e21aa51`.

| Asset | Baseline | Generated delivery fallback | Change |
|---|---:|---:|---:|
| Before desktop | 1,579,472 B | 63,286 B | -96.0% |
| Before mobile | 1,579,472 B | 28,902 B | -98.2% |
| After desktop | 1,721,676 B | 71,364 B | -95.9% |
| After mobile | 1,721,676 B | 33,866 B | -98.0% |
| Resume GLB | 28,128,092 B | 3,858,836 B | -86.3% |
| favicon | 1,485,199 B | 2,725 B | -99.8% |
| Apple touch icon | 1,485,199 B source | 18,840 B | -98.7% |

The production web image removes the three original source files after the
frontend build, while Git retains them for migration/rollback. The two desktop
X-Ray delivery files total 134,650 B, below the 800 KiB budget. Mobile initially
requests only the 28,902 B Before fallback. The GLB is below the 6 MiB hard
limit and its camera, animation, eye, and sticker-node contract is covered by a
backend test.

The excluded originals total 32,914,439 B. The generated WebP, optimized GLB,
favicon, and touch icon total 4,077,819 B, reducing the production image's
static-file payload by 28,836,620 B (27.5 MiB) before layer compression.

Lighthouse/RUM and Cloudflare cache-hit results require a deployed origin and
custom domain. Record three identical Lighthouse runs and use the median after
production deployment; do not claim those values from a local build.
