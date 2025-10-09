import { QueryCache } from './queryCache';
import { MutationCache } from './mutationCache';
import type { OmitKeyof } from '@tanstack/query-core';
import type { CancelOptions, DefaultedQueryObserverOptions } from './types';
import type { Logger } from './logger';
import type { QueryState } from './query';
import type { DefaultOptions, FetchInfiniteQueryOptions, FetchQueryOptions, InfiniteData, InvalidateOptions, InvalidateQueryFilters, MutationKey, MutationObserverOptions, MutationOptions, QueryClientConfig, QueryFunction, QueryKey, QueryObserverOptions, RefetchOptions, RefetchQueryFilters, ResetOptions, ResetQueryFilters, SetDataOptions, WithRequired } from './types';
import type { MutationFilters, QueryFilters, Updater } from './utils';
export declare class QueryClient {
    private queryCache;
    private mutationCache;
    private logger;
    private defaultOptions;
    private queryDefaults;
    private mutationDefaults;
    private mountCount;
    private unsubscribeFocus?;
    private unsubscribeOnline?;
    constructor(config?: QueryClientConfig);
    mount(): void;
    unmount(): void;
    isFetching(filters?: QueryFilters): number;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    isFetching(queryKey?: QueryKey, filters?: OmitKeyof<QueryFilters, 'queryKey'>): number;
    isMutating(filters?: MutationFilters): number;
    getQueryData<TQueryFnData = unknown>(queryKey: QueryKey): TQueryFnData | undefined;
    /**
     * @deprecated This method will accept only queryKey in the next major version.
     */
    getQueryData<TQueryFnData = unknown>(queryKey: QueryKey, filters: OmitKeyof<QueryFilters, 'queryKey'>): TQueryFnData | undefined;
    ensureQueryData<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(options: WithRequired<FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey'>): Promise<TData>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    ensureQueryData<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, options?: OmitKeyof<FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey'>): Promise<TData>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    ensureQueryData<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, queryFn: QueryFunction<TQueryFnData, TQueryKey>, options?: OmitKeyof<FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey' | 'queryFn'>): Promise<TData>;
    getQueriesData<TQueryFnData = unknown>(filters: QueryFilters): [QueryKey, TQueryFnData | undefined][];
    /**
     * @deprecated This method should be used with only one object argument.
     */
    getQueriesData<TQueryFnData = unknown>(queryKey: QueryKey): [QueryKey, TQueryFnData | undefined][];
    setQueryData<TQueryFnData>(queryKey: QueryKey, updater: Updater<TQueryFnData | undefined, TQueryFnData | undefined>, options?: SetDataOptions): TQueryFnData | undefined;
    setQueriesData<TQueryFnData>(filters: QueryFilters, updater: Updater<TQueryFnData | undefined, TQueryFnData | undefined>, options?: SetDataOptions): [QueryKey, TQueryFnData | undefined][];
    /**
     * @deprecated This method should be used with only one object argument.
     */
    setQueriesData<TQueryFnData>(queryKey: QueryKey, updater: Updater<TQueryFnData | undefined, TQueryFnData | undefined>, options?: SetDataOptions): [QueryKey, TQueryFnData | undefined][];
    getQueryState<TQueryFnData = unknown, TError = undefined>(queryKey: QueryKey, 
    /**
     * @deprecated This filters will be removed in the next major version.
     */
    filters?: OmitKeyof<QueryFilters, 'queryKey'>): QueryState<TQueryFnData, TError> | undefined;
    removeQueries(filters?: QueryFilters): void;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    removeQueries(queryKey?: QueryKey, filters?: OmitKeyof<QueryFilters, 'queryKey'>): void;
    resetQueries<TPageData = unknown>(filters?: ResetQueryFilters<TPageData>, options?: ResetOptions): Promise<void>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    resetQueries<TPageData = unknown>(queryKey?: QueryKey, filters?: OmitKeyof<ResetQueryFilters<TPageData>, 'queryKey'>, options?: ResetOptions): Promise<void>;
    cancelQueries(filters?: QueryFilters, options?: CancelOptions): Promise<void>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    cancelQueries(queryKey?: QueryKey, filters?: OmitKeyof<QueryFilters, 'queryKey'>, options?: CancelOptions): Promise<void>;
    invalidateQueries<TPageData = unknown>(filters?: InvalidateQueryFilters<TPageData>, options?: InvalidateOptions): Promise<void>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    invalidateQueries<TPageData = unknown>(queryKey?: QueryKey, filters?: OmitKeyof<InvalidateQueryFilters<TPageData>, 'queryKey'>, options?: InvalidateOptions): Promise<void>;
    refetchQueries<TPageData = unknown>(filters?: RefetchQueryFilters<TPageData>, options?: RefetchOptions): Promise<void>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    refetchQueries<TPageData = unknown>(queryKey?: QueryKey, filters?: OmitKeyof<RefetchQueryFilters<TPageData>, 'queryKey'>, options?: RefetchOptions): Promise<void>;
    fetchQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(options: FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>): Promise<TData>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    fetchQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, options?: OmitKeyof<FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey'>): Promise<TData>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    fetchQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, queryFn: QueryFunction<TQueryFnData, TQueryKey>, options?: OmitKeyof<FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey' | 'queryFn'>): Promise<TData>;
    prefetchQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(options: FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>): Promise<void>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    prefetchQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, options?: OmitKeyof<FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey'>): Promise<void>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    prefetchQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, queryFn: QueryFunction<TQueryFnData, TQueryKey>, options?: OmitKeyof<FetchQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey' | 'queryFn'>): Promise<void>;
    fetchInfiniteQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(options: FetchInfiniteQueryOptions<TQueryFnData, TError, TData, TQueryKey>): Promise<InfiniteData<TData>>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    fetchInfiniteQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, options?: OmitKeyof<FetchInfiniteQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey'>): Promise<InfiniteData<TData>>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    fetchInfiniteQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, queryFn: QueryFunction<TQueryFnData, TQueryKey>, options?: OmitKeyof<FetchInfiniteQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey' | 'queryFn'>): Promise<InfiniteData<TData>>;
    prefetchInfiniteQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(options: FetchInfiniteQueryOptions<TQueryFnData, TError, TData, TQueryKey>): Promise<void>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    prefetchInfiniteQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, options?: OmitKeyof<FetchInfiniteQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey'>): Promise<void>;
    /**
     * @deprecated This method should be used with only one object argument.
     */
    prefetchInfiniteQuery<TQueryFnData = unknown, TError = unknown, TData = TQueryFnData, TQueryKey extends QueryKey = QueryKey>(queryKey: TQueryKey, queryFn: QueryFunction<TQueryFnData, TQueryKey>, options?: OmitKeyof<FetchInfiniteQueryOptions<TQueryFnData, TError, TData, TQueryKey>, 'queryKey' | 'queryFn'>): Promise<void>;
    resumePausedMutations(): Promise<unknown>;
    getQueryCache(): QueryCache;
    getMutationCache(): MutationCache;
    getLogger(): Logger;
    getDefaultOptions(): DefaultOptions;
    setDefaultOptions(options: DefaultOptions): void;
    setQueryDefaults(queryKey: QueryKey, options: QueryObserverOptions<unknown, any, any, any>): void;
    getQueryDefaults(queryKey?: QueryKey): QueryObserverOptions<any, any, any, any, any> | undefined;
    setMutationDefaults(mutationKey: MutationKey, options: MutationObserverOptions<any, any, any, any>): void;
    getMutationDefaults(mutationKey?: MutationKey): MutationObserverOptions<any, any, any, any> | undefined;
    defaultQueryOptions<TQueryFnData, TError, TData, TQueryData, TQueryKey extends QueryKey>(options?: QueryObserverOptions<TQueryFnData, TError, TData, TQueryData, TQueryKey> | DefaultedQueryObserverOptions<TQueryFnData, TError, TData, TQueryData, TQueryKey>): DefaultedQueryObserverOptions<TQueryFnData, TError, TData, TQueryData, TQueryKey>;
    defaultMutationOptions<T extends MutationOptions<any, any, any, any>>(options?: T): T;
    clear(): void;
}
//# sourceMappingURL=queryClient.d.ts.map