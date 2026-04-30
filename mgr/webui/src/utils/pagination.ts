import { useState } from "react";

export const DEFAULT_PAGE_SIZE = 10;
export const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];

interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
}

export function getLastPage(total: number, pageSize: number) {
  return Math.max(1, Math.ceil(Math.max(total, 0) / pageSize));
}

export function getPageAfterDelete(total: number, page: number, pageSize: number) {
  return Math.min(page, getLastPage(total, pageSize));
}

export function useTablePagination(initialPageSize = DEFAULT_PAGE_SIZE) {
  const [pagination, setPagination] = useState<PaginationState>({
    page: 1,
    pageSize: initialPageSize,
    total: 0,
  });

  const resetPage = () => {
    setPagination((current) => ({
      ...current,
      page: 1,
    }));
  };

  const updatePageData = (data: API.PageResult<unknown>) => {
    setPagination({
      page: data.page,
      pageSize: data.pageSize,
      total: data.total,
    });
  };

  const handleTableChange = (page: number, pageSize: number) => {
    setPagination((current) => ({
      ...current,
      page,
      pageSize,
    }));
  };

  return {
    pagination,
    resetPage,
    setPagination,
    updatePageData,
    handleTableChange,
  };
}
