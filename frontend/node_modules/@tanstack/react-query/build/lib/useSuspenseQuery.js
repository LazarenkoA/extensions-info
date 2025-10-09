'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var queryCore = require('@tanstack/query-core');
var useBaseQuery = require('./useBaseQuery.js');

function useSuspenseQuery(options) {
  return useBaseQuery.useBaseQuery({ ...options,
    enabled: true,
    useErrorBoundary: true,
    suspense: true,
    placeholderData: undefined,
    networkMode: 'always',
    onSuccess: undefined,
    onError: undefined,
    onSettled: undefined
  }, queryCore.QueryObserver);
}

exports.useSuspenseQuery = useSuspenseQuery;
//# sourceMappingURL=useSuspenseQuery.js.map
