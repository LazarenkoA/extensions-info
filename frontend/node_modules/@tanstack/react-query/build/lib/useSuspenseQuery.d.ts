import type { QueryKey } from '@tanstack/query-core';
import type { UseSuspenseQueryOptions, UseSuspenseQueryResult } from './types';
export declare function useSuspenseQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(options: UseSuspenseQueryOptions<TQueryFnData, TError, TData, TQueryKey>): UseSuspenseQueryResult<TData, TError>;
//# sourceMappingURL=useSuspenseQuery.d.ts.map