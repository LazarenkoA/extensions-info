import type { UseQueryOptions, UseSuspenseQueryOptions, UseSuspenseQueryResult } from './types';
import type { QueryFunction } from '@tanstack/query-core';
declare type MAXIMUM_DEPTH = 20;
declare type GetSuspenseOptions<T> = T extends {
    queryFnData: infer TQueryFnData;
    error?: infer TError;
    data: infer TData;
} ? UseSuspenseQueryOptions<TQueryFnData, TError, TData> : T extends {
    queryFnData: infer TQueryFnData;
    error?: infer TError;
} ? UseSuspenseQueryOptions<TQueryFnData, TError> : T extends {
    data: infer TData;
    error?: infer TError;
} ? UseSuspenseQueryOptions<unknown, TError, TData> : T extends [infer TQueryFnData, infer TError, infer TData] ? UseSuspenseQueryOptions<TQueryFnData, TError, TData> : T extends [infer TQueryFnData, infer TError] ? UseSuspenseQueryOptions<TQueryFnData, TError> : T extends [infer TQueryFnData] ? UseSuspenseQueryOptions<TQueryFnData> : T extends {
    queryFn?: QueryFunction<infer TQueryFnData, infer TQueryKey>;
    select?: (data: any) => infer TData;
} ? UseSuspenseQueryOptions<TQueryFnData, unknown, TData, TQueryKey> : T extends {
    queryFn?: QueryFunction<infer TQueryFnData, infer TQueryKey>;
} ? UseSuspenseQueryOptions<TQueryFnData, unknown, TQueryFnData, TQueryKey> : UseSuspenseQueryOptions;
declare type GetSuspenseResults<T> = T extends {
    queryFnData: any;
    error?: infer TError;
    data: infer TData;
} ? UseSuspenseQueryResult<TData, TError> : T extends {
    queryFnData: infer TQueryFnData;
    error?: infer TError;
} ? UseSuspenseQueryResult<TQueryFnData, TError> : T extends {
    data: infer TData;
    error?: infer TError;
} ? UseSuspenseQueryResult<TData, TError> : T extends [any, infer TError, infer TData] ? UseSuspenseQueryResult<TData, TError> : T extends [infer TQueryFnData, infer TError] ? UseSuspenseQueryResult<TQueryFnData, TError> : T extends [infer TQueryFnData] ? UseSuspenseQueryResult<TQueryFnData> : T extends {
    queryFn?: QueryFunction<infer TQueryFnData, any>;
    select?: (data: any) => infer TData;
} ? UseSuspenseQueryResult<unknown extends TData ? TQueryFnData : TData> : T extends {
    queryFn?: QueryFunction<infer TQueryFnData, any>;
} ? UseSuspenseQueryResult<TQueryFnData> : UseSuspenseQueryResult;
/**
 * SuspenseQueriesOptions reducer recursively unwraps function arguments to infer/enforce type param
 */
export declare type SuspenseQueriesOptions<T extends Array<any>, TResult extends Array<any> = [], TDepth extends ReadonlyArray<number> = []> = TDepth['length'] extends MAXIMUM_DEPTH ? Array<UseSuspenseQueryOptions> : T extends [] ? [] : T extends [infer Head] ? [...TResult, GetSuspenseOptions<Head>] : T extends [infer Head, ...infer Tail] ? SuspenseQueriesOptions<[
    ...Tail
], [
    ...TResult,
    GetSuspenseOptions<Head>
], [
    ...TDepth,
    1
]> : Array<unknown> extends T ? T : T extends Array<UseSuspenseQueryOptions<infer TQueryFnData, infer TError, infer TData, infer TQueryKey>> ? Array<UseSuspenseQueryOptions<TQueryFnData, TError, TData, TQueryKey>> : Array<UseSuspenseQueryOptions>;
/**
 * SuspenseQueriesResults reducer recursively maps type param to results
 */
export declare type SuspenseQueriesResults<T extends Array<any>, TResult extends Array<any> = [], TDepth extends ReadonlyArray<number> = []> = TDepth['length'] extends MAXIMUM_DEPTH ? Array<UseSuspenseQueryResult> : T extends [] ? [] : T extends [infer Head] ? [...TResult, GetSuspenseResults<Head>] : T extends [infer Head, ...infer Tail] ? SuspenseQueriesResults<[
    ...Tail
], [
    ...TResult,
    GetSuspenseResults<Head>
], [
    ...TDepth,
    1
]> : T extends Array<UseSuspenseQueryOptions<infer TQueryFnData, infer TError, infer TData, any>> ? Array<UseSuspenseQueryResult<unknown extends TData ? TQueryFnData : TData, TError>> : Array<UseSuspenseQueryResult>;
export declare function useSuspenseQueries<T extends any[]>({ queries, context, }: {
    queries: readonly [...SuspenseQueriesOptions<T>];
    context?: UseQueryOptions['context'];
}): SuspenseQueriesResults<T>;
export {};
//# sourceMappingURL=useSuspenseQueries.d.ts.map