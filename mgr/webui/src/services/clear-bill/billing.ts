import { requestWithFallback } from "@/services/http";

import { billRecordsMock, customersMock, reconciliationTasksMock } from "./mock";
import type { BillRecord, CustomerSnapshot, ReconciliationTask } from "./types";

export async function getBills(): Promise<BillRecord[]> {
  return requestWithFallback("/api/v1/bills", billRecordsMock);
}

export async function getCustomers(): Promise<CustomerSnapshot[]> {
  return requestWithFallback("/api/v1/customers", customersMock);
}

export async function getReconciliationTasks(): Promise<ReconciliationTask[]> {
  return requestWithFallback("/api/v1/reconciliations", reconciliationTasksMock);
}
