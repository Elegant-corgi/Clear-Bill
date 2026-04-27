import { requestWithFallback } from "@/services/http";

import { billRecordsMock, customersMock, reconciliationTasksMock } from "./mock";
import type { BillRecord, CustomerSnapshot, ReconciliationTask } from "./types";

export async function getBills(): Promise<BillRecord[]> {
  return requestWithFallback("/bills", billRecordsMock);
}

export async function getCustomers(): Promise<CustomerSnapshot[]> {
  return requestWithFallback("/customers", customersMock);
}

export async function getReconciliationTasks(): Promise<ReconciliationTask[]> {
  return requestWithFallback("/reconciliations", reconciliationTasksMock);
}
