#include <immintrin.h>

#define ALIGN(x) x __attribute__((aligned(32)))

// both elements = dot(x,y)
void DP(double* vec1, double* vec2, double* result) {
    __m256d x = _mm256_load_pd(vec1);
    __m256d y = _mm256_load_pd(vec2);
    __m256d xy = _mm256_mul_pd(x, y);

    __m128d xylow  = _mm256_castpd256_pd128(xy);   // (__m128d)cast isn't portable
    __m128d xyhigh = _mm256_extractf128_pd(xy, 1);
    __m128d sum1 =   _mm_add_pd(xylow, xyhigh);

    __m128d swapped = _mm_shuffle_pd(sum1, sum1, 0b01);   // or unpackhi
    __m128d dotproduct = _mm_add_pd(sum1, swapped);
    _mm_storeh_pd(result, dotproduct);
}