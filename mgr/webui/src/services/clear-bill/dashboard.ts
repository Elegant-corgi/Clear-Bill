import { requestWithFallback } from "@/services/http";

import { dashboardSummaryMock } from "./mock";
import type { DashboardSummary } from "./types";

export async function getDashboardSummary(): Promise<DashboardSummary> {
  return requestWithFallback("/api/v1/dashboard/summary", dashboardSummaryMock);
}
