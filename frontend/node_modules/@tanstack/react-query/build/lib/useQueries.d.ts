import type { OmitKeyof, QueryFunction, QueryKey } from '@tanstack/query-core';
import type { DefinedUseQueryResult, UseQueryOptions, UseQueryResult } from './types';
declare type UseQueryOptionsForUseQueries<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey> = OmitKeyof<UseQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'context'>;
declare type MAXIMUM_DEPTH = 20;
declare type GetOptions<T> = T extends {
    queryFnData: infer TQueryFnData;
    error?: infer TError;
    data: infer TData;
} ? UseQueryOptionsForUseQueries<TQueryFnData, TError, TData> : T extends {
    queryFnData: infer TQueryFnData;
    error?: infer TError;
} ? UseQueryOptionsForUseQueries<TQueryFnData, TError> : T extends {
    data: infer TData;
    error?: infer TError;
} ? UseQueryOptionsForUseQueries<unknown, TError, TData> : T extends [infer TQueryFnData, infer TError, infer TData] ? UseQueryOptionsForUseQueries<TQueryFnData, TError, TData> : T extends [infer TQueryFnData, infer TError] ? UseQueryOptionsForUseQueries<TQueryFnData, TError> : T extends [infer TQueryFnData] ? UseQueryOptionsForUseQueries<TQueryFnData> : T extends {
    queryFn?: QueryFunction<infer TQueryFnData, infer TQueryKey>;
    select: (data: any) => infer TData;
} ? UseQueryOptionsForUseQueries<TQueryFnData, unknown, TData, TQueryKey> : T extends {
    queryFn?: QueryFunction<infer TQueryFnData, infer TQueryKey>;
} ? UseQueryOptionsForUseQueries<TQueryFnData, unknown, TQueryFnData, TQueryKey> : UseQueryOptionsForUseQueries;
declare type GetDefinedOrUndefinedQueryResult<T, TData, TError = unknown> = T extends {
    initialData?: infer TInitialData;
} ? unknown extends TInitialData ? UseQueryResult<TData, TError> : TInitialData extends TData ? DefinedUseQueryResult<TData, TError> : TInitialData extends () => infer TInitialDataResult ? unknown extends TInitialDataResult ? UseQueryResult<TData, TError> : TInitialDataResult extends TData ? DefinedUseQueryResult<TData, TError> : UseQueryResult<TData, TError> : UseQueryResult<TData, TError> : UseQueryResult<TData, TError>;
declare type GetUseQueryResult<T> = T extends {
    queryFnData: any;
    error?: infer TError;
    data: infer TData;
} ? GetDefinedOrUndefinedQueryResult<T, TData, TError> : T extends {
    queryFnData: infer TQueryFnData;
    error?: infer TError;
} ? GetDefinedOrUndefinedQueryResult<T, TQueryFnData, TError> : T extends {
    data: infer TData;
    error?: infer TError;
} ? GetDefinedOrUndefinedQueryResult<T, TData, TError> : T extends [any, infer TError, infer TData] ? GetDefinedOrUndefinedQueryResult<T, TData, TError> : T extends [infer TQueryFnData, infer TError] ? GetDefinedOrUndefinedQueryResult<T, TQueryFnData, TError> : T extends [infer TQueryFnData] ? GetDefinedOrUndefinedQueryResult<T, TQueryFnData> : T extends {
    queryFn?: QueryFunction<unknown, any>;
    select: (data: any) => infer TData;
} ? GetDefinedOrUndefinedQueryResult<T, TData> : T extends {
    queryFn?: QueryFunction<infer TQueryFnData, any>;
} ? GetDefinedOrUndefinedQueryResult<T, TQueryFnData> : UseQueryResult;
/**
 * QueriesOptions reducer recursively unwraps function arguments to infer/enforce type param
 */
export declare type QueriesOptions<T extends any[], TResult extends any[] = [], TDepth extends ReadonlyArray<number> = []> = TDepth['length'] extends MAXIMUM_DEPTH ? UseQueryOptionsForUseQueries[] : T extends [] ? [] : T extends [infer Head] ? [...TResult, GetOptions<Head>] : T extends [infer Head, ...infer Tail] ? QueriesOptions<[...Tail], [...TResult, GetOptions<Head>], [...TDepth, 1]> : unknown[] extends T ? T : T extends UseQueryOptionsForUseQueries<infer TQueryFnData, infer TError, infer TData, infer TQueryKey>[] ? UseQueryOptionsForUseQueries<TQueryFnData, TError, TData, TQueryKey>[] : UseQueryOptionsForUseQueries[];
/**
 * QueriesResults reducer recursively maps type param to results
 */
export declare type QueriesResults<T extends any[], TResults extends any[] = [], TDepth extends ReadonlyArray<number> = []> = TDepth['length'] extends MAXIMUM_DEPTH ? UseQueryResult[] : T extends [] ? [] : T extends [infer Head] ? [...TResults, GetUseQueryResult<Head>] : T extends [infer Head, ...infer Tail] ? QueriesResults<[
    ...Tail
], [
    ...TResults,
    GetUseQueryResult<Head>
], [
    ...TDepth,
    1
]> : T extends UseQueryOptionsForUseQueries<infer TQueryFnData, infer TError, infer TData, any>[] ? UseQueryResult<unknown extends TData ? TQueryFnData : TData, TError>[] : UseQueryResult[];
export declare function useQueries<T extends any[]>({ queries, context, }: {
    queries: readonly [...QueriesOptions<T>];
    context?: UseQueryOptions['context'];
}): QueriesResults<T>;
export {};
//# sourceMappingURL=useQueries.d.ts.map