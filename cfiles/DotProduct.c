#include <immintrin.h>

#define ALIGN(x) x __attribute__((aligned(32)))

// Dot product is the sum of the products of the corresponding entries of the two sequences of numbers
// a product is simply put the result of a multiplication. The dot product of two tuples is simply
// t1.x * t2.x + t1.y * t2.y + t1.z * t2.z + t1[3] * t2[3]
// void DotProduct(double* arg1, double* arg2, double *result) {
//     __m256d vec1 = _mm256_load_pd(arg1);
//     __m256d vec2 = _mm256_load_pd(arg2);
//     __m256d xy = _mm256_mul_pd( vec1, vec2 );
//     __m256d temp = _mm256_hadd_pd( xy, xy );
//     __m128d hi128 = _mm256_extractf128_pd( temp, 1 );
//     __m128d dotproduct = _mm_add_pd( (__m128d)temp, hi128 );    
//     _mm256_storeu_pd(result, dotproduct);
// }

void DotProduct(double* vec1, double* vec2, double* result) {
    __m256d x = _mm256_load_pd(vec1);
    __m256d y = _mm256_load_pd(vec2);
    __m256d xy = _mm256_mul_pd(x, y);

    __m128d xylow  = _mm256_castpd256_pd128(xy);   // (__m128d)cast isn't portable
    __m128d xyhigh = _mm256_extractf128_pd(xy, 1);
    __m128d sum1 =   _mm_add_pd(xylow, xyhigh);

    __m128d swapped = _mm_shuffle_pd(sum1, sum1, 0b01);   // or unpackhi
    __m128d dotproduct = _mm_add_pd(sum1, swapped);
    _mm_storeh_pd(result, dotproduct);
    //_mm_store_pd(result, dotproduct[0]);
    //return dotproduct; __m128d
}
