import { requestWithFallback } from "@/services/http";

import { dashboardSummaryMock } from "./mock";
import type { DashboardSummary } from "./types";

export async function getDashboardSummary(): Promise<DashboardSummary> {
  return requestWithFallback("/dashboard/summary", dashboardSummaryMock);
}
