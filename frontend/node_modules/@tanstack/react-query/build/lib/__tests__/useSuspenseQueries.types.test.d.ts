export declare const queryKey: readonly ["key"];
export declare const queryFn: () => Promise<{
    text: string;
}>;
export declare const select: (data: Awaited<ReturnType<typeof queryFn>>) => string;
//# sourceMappingURL=useSuspenseQueries.types.test.d.ts.map